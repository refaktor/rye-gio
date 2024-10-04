// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"go/format"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unsafe"

	_ "embed"

	"gioui.org/cpu"
)

/*
#cgo pkg-config: vulkan

#include <stdlib.h>
#include <vulkan/vulkan.h>
*/
import "C"

type descriptorType uint8

type descriptor struct {
	binding int
	_type   descriptorType
	count   int
}

// program describes parameters constant to a program.
type program struct {
	hasControlBarriers bool
	workgroupSize      [3]int
	memorySize         int
}

const (
	descriptorTypeBuffer descriptorType = iota
	descriptorTypeImage
)

var (
	descSetLayout = flag.String("layout", "", "Descriptor set layout in <binding>:[<num>]<type> format\nFor example: '0:ssbo,1:[10]ssbo'")
	arch          = flag.String("arch", "amd64", "GOARCH prefix (amd64, arm64, arm,386)")
	objcopy       = flag.String("objcopy", "objcopy", "The objcopy binary suitable for the architecture")
)

const (
	supportConstraints      = "linux && (arm64 || arm || amd64)"
	supportConstraints116   = "// +build linux\n// +build arm64 arm amd64"
	nosupportConstraints116 = "// +build !linux !arm64,!arm,!amd64"
)

//go:embed support.c.inc
var supportc []byte

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "specify compute program filename\n")
		os.Exit(1)
	}
	layout, err := parseDescriptorSetLayout(*descSetLayout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse layout: %v\n", err)
		os.Exit(1)
	}

	file := flag.Arg(0)
	glsl, err := os.ReadFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}
	spirv, err := glslToSPIRV(file, glsl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", file, err)
		os.Exit(2)
	}
	prog, err := parseProgramParams(spirv)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse SPIR-V parameters: %v\n", err)
		os.Exit(1)
	}
	obj, err := compileProgram(spirv, layout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build SPIR-V program: %v\n", err)
		os.Exit(2)
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get current dir: %v\n", err)
		os.Exit(2)
	}
	pkg := filepath.Base(cwd)

	name := filepath.Base(file)
	name = name[:len(name)-len(filepath.Ext(name))]

	sum := sha256.Sum256(glsl)
	hexSum := hex.EncodeToString(sum[:])
	header := generateHeader(name, layout)
	impl := generateImpl(name, prog)
	goImpl := generateGo(pkg, name, hexSum, layout)
	goFallbackImpl := generateGoFallback(pkg, name, layout)

	syso := fmt.Sprintf("%s_linux_%s.syso", name, *arch)
	if err := os.WriteFile(syso, obj, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}
	// Rename symbols to avoid name clashes when linking multiple programs.
	var renames []string
	symbols := []string{
		"coroutine_begin",
		"coroutine_await",
		"coroutine_destroy",
		"coroutine_begin.resume",
		"coroutine_begin.destroy",
		"coroutine_begin.cleanup",
	}
	for _, sym := range symbols {
		rename := fmt.Sprintf("%s=%s_%s", sym, name, sym)
		renames = append(renames, "--redefine-sym", rename)
	}
	renames = append(renames, syso)
	objcopy := exec.Command(*objcopy, renames...)
	objcopy.Stderr = os.Stderr
	if err := objcopy.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v: %v\n", objcopy, err)
		os.Exit(2)
	}
	files := []struct {
		name    string
		content []byte
	}{
		{name + "_abi.h", header},
		{name + "_abi.c", impl},
		{name + "_abi.go", goImpl},
		{name + "_abi_nosupport.go", goFallbackImpl},
		{"support.c", supportc},
		{"runtime.h", cpu.RuntimeH},
		{"abi.h", cpu.ABIH},
	}
	for _, f := range files {
		if err := os.WriteFile(f.name, f.content, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(2)
		}
	}
}

func glslToSPIRV(file string, glsl []byte) ([]byte, error) {
	output := filepath.Join(os.TempDir(), filepath.Base(file)+".spv")
	defer os.Remove(output)
	cmd := exec.Command(
		"glslangValidator",
		"-V", "-o", output,
		"--stdin",
		"-I"+filepath.Dir(file),
		"-S", "comp",
	)
	cmd.Stdin = bytes.NewBuffer(glsl)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return os.ReadFile(output)
}

func parseProgramParams(spirv []byte) (program, error) {
	var stdout bytes.Buffer
	disassemble := exec.Command("spirv-dis", "--no-indent", "--no-header")
	disassemble.Stdin = bytes.NewBuffer(spirv)
	disassemble.Stdout = &stdout

	var p program
	// TODO: extract from SwiftShader.
	p.memorySize = 1e5
	if err := disassemble.Run(); err != nil {
		return p, err
	}
	for {
		line, err := stdout.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				return p, err
			}
			break
		}
		switch {
		case strings.HasPrefix(line, "OpExecutionMode"):
			s := &p.workgroupSize
			i := strings.Index(line, "LocalSize ")
			if i == -1 {
				return p, fmt.Errorf("unknown execution mode: %s", line)
			}
			sub := line[i:]
			if _, err := fmt.Sscanf(sub, "LocalSize %d %d %d", &s[0], &s[1], &s[2]); err != nil {
				return p, err
			}
		case strings.Index(line, "OpControlBarrier") != -1:
			p.hasControlBarriers = true
		}
	}
	return p, nil
}

func generateImpl(name string, p program) []byte {
	var b printer
	b.printf("// Code generated by gioui.org/cpu/cmd/compile DO NOT EDIT.\n\n")
	b.printf("//go:build %s\n", supportConstraints)
	b.printf("%s\n\n", supportConstraints116)
	b.printf("#include <stdint.h>\n")
	b.printf("#include <stddef.h>\n")
	b.printf("#include \"abi.h\"\n")
	b.printf("#include \"runtime.h\"\n")
	b.printf("#include \"%s_abi.h\"\n\n", name)
	b.printf("const struct program_info %s_program_info = {\n", name)
	c := 0
	if p.hasControlBarriers {
		c = 1
	}
	b.printf("\t.has_cbarriers = %d,\n", c)
	b.printf("\t.min_memory_size = %d,\n", p.memorySize)
	b.printf("\t.desc_set_size = sizeof(struct %s_descriptor_set_layout),\n", name)
	b.printf("\t.workgroup_size_x = %d,\n", p.workgroupSize[0])
	b.printf("\t.workgroup_size_y = %d,\n", p.workgroupSize[1])
	b.printf("\t.workgroup_size_z = %d,\n", p.workgroupSize[2])
	b.printf("\t.begin = %s_coroutine_begin,\n", name)
	b.printf("\t.await = %s_coroutine_await,\n", name)
	b.printf("\t.destroy = %s_coroutine_destroy,\n", name)
	b.printf("};\n")
	return b.Bytes()
}

func generateGoFallback(pkg, name string, layout []descriptor) []byte {
	var b printer
	b.printf("// Code generated by gioui.org/cpu/cmd/compile DO NOT EDIT.\n\n")

	b.printf("//go:build !(%s)\n", supportConstraints)
	b.printf("%s\n\n", nosupportConstraints116)

	b.printf("package %s\n\n", pkg)

	b.printf("import \"gioui.org/cpu\"\n")

	b.printf("var %sProgramInfo *cpu.ProgramInfo\n\n", strings.Title(name))
	b.printf("type %sDescriptorSetLayout struct{}\n\n", strings.Title(name))
	b.printf("const %sHash = \"\"\n\n", strings.Title(name))
	for _, desc := range layout {
		var _type string
		switch desc._type {
		case descriptorTypeBuffer:
			_type = "Buffer"
		case descriptorTypeImage:
			_type = "Image"
		}
		args := ""
		if desc.count > 1 {
			args = "index int"
		}
		b.printf("func (l *%sDescriptorSetLayout) Binding%d(%s) *cpu.%sDescriptor {\n", strings.Title(name), desc.binding, args, _type)
		b.printf("\tpanic(\"unsupported\")\n")
		b.printf("}\n\n")
	}
	src, err := format.Source(b.Bytes())
	if err != nil {
		panic(err)
	}
	return src
}

func generateGo(pkg, name, hexSum string, layout []descriptor) []byte {
	var b printer
	b.printf("// Code generated by gioui.org/cpu/cmd/compile DO NOT EDIT.\n\n")
	b.printf("//go:build %s\n", supportConstraints)
	b.printf("%s\n\n", supportConstraints116)
	b.printf("package %s\n\n", pkg)

	b.printf("import \"gioui.org/cpu\"\n")
	b.printf("import \"unsafe\"\n\n")

	b.printf("/*\n")
	b.printf("#cgo LDFLAGS: -lm\n\n")
	b.printf("#include <stdint.h>\n")
	b.printf("#include <stdlib.h>\n")
	b.printf("#include \"abi.h\"\n")
	b.printf("#include \"runtime.h\"\n")
	b.printf("#include \"%s_abi.h\"\n", name)
	b.printf("*/\n")
	b.printf("import \"C\"\n\n")

	b.printf("var %sProgramInfo = (*cpu.ProgramInfo)(unsafe.Pointer(&C.%s_program_info))\n\n", strings.Title(name), name)
	b.printf("type %sDescriptorSetLayout = C.struct_%s_descriptor_set_layout\n\n", strings.Title(name), name)
	b.printf("const %sHash = %q\n\n", strings.Title(name), hexSum)
	for _, desc := range layout {
		var _type string
		switch desc._type {
		case descriptorTypeBuffer:
			_type = "Buffer"
		case descriptorTypeImage:
			_type = "Image"
		}
		args := ""
		index := ""
		if desc.count > 1 {
			args = "index int"
			index = "[index]"
		}
		b.printf("func (l *%sDescriptorSetLayout) Binding%d(%s) *cpu.%sDescriptor {\n", strings.Title(name), desc.binding, args, _type)
		b.printf("\treturn (*cpu.%sDescriptor)(unsafe.Pointer(&l.binding%d%s))\n", _type, desc.binding, index)
		b.printf("}\n\n")
	}
	src, err := format.Source(b.Bytes())
	if err != nil {
		panic(err)
	}
	return src
}

func generateHeader(name string, layout []descriptor) []byte {
	var b printer

	b.printf("// Code generated by gioui.org/cpu/cmd/compile DO NOT EDIT.\n\n")
	// Descriptor set layout.
	b.printf("struct %s_descriptor_set_layout {\n", name)
	for _, desc := range layout {
		var _type string
		switch desc._type {
		case descriptorTypeBuffer:
			_type = "buffer"
		case descriptorTypeImage:
			_type = "image"
		}
		b.printf("\tstruct %s_descriptor binding%d", _type, desc.binding)
		if desc.count > 1 {
			b.printf("[%d]", desc.count)
		}
		b.printf(";\n")
	}
	b.printf("};\n\n")

	// Program routines.
	b.printf("extern coroutine %s_coroutine_begin(struct program_data *data,\n", name)
	b.printf("\tint32_t workgroupX, int32_t workgroupY, int32_t workgroupZ,\n")
	b.printf("\tvoid *workgroupMemory,\n")
	b.printf("\tint32_t firstSubgroup,\n")
	b.printf("\tint32_t subgroupCount) ATTR_HIDDEN;\n\n")
	b.printf("extern bool %s_coroutine_await(coroutine r, yield_result *res) ATTR_HIDDEN;\n", name)
	b.printf("extern void %s_coroutine_destroy(coroutine r) ATTR_HIDDEN;\n\n", name)

	b.printf("extern const struct program_info %s_program_info ATTR_HIDDEN;\n", name)

	return b.Bytes()
}

func parseDescriptorSetLayout(layoutDef string) ([]descriptor, error) {
	descriptors := strings.Split(layoutDef, ",")
	var layout []descriptor
	for _, def := range descriptors {
		var desc descriptor
		var typeName string
		if _, err := fmt.Sscanf(def, "%d:[%d]%s", &desc.binding, &desc.count, &typeName); err != nil {
			desc.count = 1
			if _, err := fmt.Sscanf(def, "%d:%s", &desc.binding, &typeName); err != nil {
				return nil, err
			}
		}
		switch typeName {
		case "buffer":
			desc._type = descriptorTypeBuffer
		case "image":
			desc._type = descriptorTypeImage
		default:
			return nil, fmt.Errorf("unknown descriptor type: %s", typeName)
		}
		layout = append(layout, desc)
	}
	return layout, nil
}

func vkErr(code C.VkResult) error {
	if code == C.VK_SUCCESS {
		return nil
	}
	return fmt.Errorf("error code: %d", code)
}

// compileProgram uses the SwiftShader vulkan implementation to compile a SPIR-V
// compute program to assembly. It returns the object file in ELF format.
func compileProgram(spirv []byte, layout []descriptor) ([]byte, error) {
	// Create vulkan instance.
	instInf := C.VkInstanceCreateInfo{
		sType: C.VK_STRUCTURE_TYPE_INSTANCE_CREATE_INFO,
	}
	var inst C.VkInstance
	if err := vkErr(C.vkCreateInstance(&instInf, nil, &inst)); err != nil {
		return nil, err
	}
	defer C.vkDestroyInstance(inst, nil)

	// Find a vulkan device.
	var (
		ndev C.uint32_t = 1
		pdev C.VkPhysicalDevice
	)
	if err := vkErr(C.vkEnumeratePhysicalDevices(inst, &ndev, &pdev)); err != nil {
		return nil, err
	}

	// Find queue with compute ability.
	var nqueues C.uint32_t
	C.vkGetPhysicalDeviceQueueFamilyProperties(pdev, &nqueues, nil)
	qprops := make([]C.VkQueueFamilyProperties, nqueues)
	C.vkGetPhysicalDeviceQueueFamilyProperties(pdev, &nqueues, &qprops[0])
	qidx := -1
	for i, prop := range qprops {
		if prop.queueFlags&C.VK_QUEUE_COMPUTE_BIT != 0 {
			qidx = i
			break
		}
	}
	if qidx == -1 {
		return nil, errors.New("no compute queues available")
	}

	// Create device and queue.
	// Place queue create info in the C heap because Cgo calls cannot
	// pass pointer to pointer of Go memory.
	qinfMem := C.malloc(C.size_t(unsafe.Sizeof(C.VkDeviceQueueCreateInfo{})))
	defer C.free(qinfMem)
	qinf := (*C.VkDeviceQueueCreateInfo)(qinfMem)
	*qinf = C.VkDeviceQueueCreateInfo{
		sType:            C.VK_STRUCTURE_TYPE_DEVICE_QUEUE_CREATE_INFO,
		queueFamilyIndex: C.uint32_t(qidx),
		queueCount:       1,
	}
	devInf := C.VkDeviceCreateInfo{
		sType:                C.VK_STRUCTURE_TYPE_DEVICE_CREATE_INFO,
		queueCreateInfoCount: 1,
		pQueueCreateInfos:    qinf,
	}
	var dev C.VkDevice
	if err := vkErr(C.vkCreateDevice(pdev, &devInf, nil, &dev)); err != nil {
		return nil, err
	}
	defer C.vkDestroyDevice(dev, nil)

	spirvMem := C.malloc(C.size_t(int(unsafe.Sizeof(spirv[0])) * len(spirv)))
	defer C.free(spirvMem)
	spirvBuf := (*(*[1 << 30]byte)(spirvMem))[:len(spirv):len(spirv)]
	copy(spirvBuf, spirv)
	modInf := C.VkShaderModuleCreateInfo{
		sType:    C.VK_STRUCTURE_TYPE_SHADER_MODULE_CREATE_INFO,
		codeSize: C.size_t(len(spirvBuf)),
		pCode:    (*C.uint32_t)(unsafe.Pointer(&spirvBuf[0])),
	}
	var mod C.VkShaderModule
	if err := vkErr(C.vkCreateShaderModule(dev, &modInf, nil, &mod)); err != nil {
		return nil, err
	}
	defer C.vkDestroyShaderModule(dev, mod, nil)

	// Create descriptor set layout.
	nbindings := len(layout)
	var bindings []C.VkDescriptorSetLayoutBinding
	bindingsMem := C.malloc(C.size_t(int(unsafe.Sizeof(bindings[0])) * nbindings))
	defer C.free(bindingsMem)
	bindings = unsafe.Slice((*C.VkDescriptorSetLayoutBinding)(bindingsMem), nbindings)
	for i, desc := range layout {
		binding := C.VkDescriptorSetLayoutBinding{
			binding:         C.uint32_t(desc.binding),
			descriptorCount: C.uint32_t(desc.count),
			stageFlags:      C.VK_SHADER_STAGE_COMPUTE_BIT,
		}
		switch desc._type {
		case descriptorTypeBuffer:
			binding.descriptorType = C.VK_DESCRIPTOR_TYPE_STORAGE_BUFFER
		case descriptorTypeImage:
			binding.descriptorType = C.VK_DESCRIPTOR_TYPE_STORAGE_IMAGE
		default:
			panic("unhandled descriptor type")
		}
		bindings[i] = binding
	}
	descSetLayoutInf := C.VkDescriptorSetLayoutCreateInfo{
		sType:        C.VK_STRUCTURE_TYPE_DESCRIPTOR_SET_LAYOUT_CREATE_INFO,
		bindingCount: C.uint32_t(len(bindings)),
		pBindings:    (*C.VkDescriptorSetLayoutBinding)(unsafe.Pointer(&bindings[0])),
	}
	var descSetLayout C.VkDescriptorSetLayout
	if err := vkErr(C.vkCreateDescriptorSetLayout(dev, &descSetLayoutInf, nil, &descSetLayout)); err != nil {
		return nil, err
	}
	defer C.vkDestroyDescriptorSetLayout(dev, descSetLayout, nil)

	// Create pipeline layout.
	descSetLayoutMem := C.malloc(C.size_t(unsafe.Sizeof(descSetLayout)))
	defer C.free(descSetLayoutMem)
	descSetLayoutBuf := (*C.VkDescriptorSetLayout)(descSetLayoutMem)
	*descSetLayoutBuf = descSetLayout
	pipeLayoutInf := C.VkPipelineLayoutCreateInfo{
		sType:          C.VK_STRUCTURE_TYPE_PIPELINE_LAYOUT_CREATE_INFO,
		setLayoutCount: 1,
		pSetLayouts:    descSetLayoutBuf,
	}
	var pipeLayout C.VkPipelineLayout
	if err := vkErr(C.vkCreatePipelineLayout(dev, &pipeLayoutInf, nil, &pipeLayout)); err != nil {
		return nil, err
	}
	defer C.vkDestroyPipelineLayout(dev, pipeLayout, nil)
	mainName := C.CString("main")
	defer C.free(unsafe.Pointer(mainName))
	pipeInf := C.VkComputePipelineCreateInfo{
		sType: C.VK_STRUCTURE_TYPE_COMPUTE_PIPELINE_CREATE_INFO,
		stage: C.VkPipelineShaderStageCreateInfo{
			sType:  C.VK_STRUCTURE_TYPE_PIPELINE_SHADER_STAGE_CREATE_INFO,
			stage:  C.VK_SHADER_STAGE_COMPUTE_BIT,
			module: mod,
			pName:  mainName,
		},
		layout: pipeLayout,
	}

	var pipe C.VkPipeline
	var nilPipeCache C.VkPipelineCache
	if err := vkErr(C.vkCreateComputePipelines(dev, nilPipeCache, 1, &pipeInf, nil, &pipe)); err != nil {
		return nil, err
	}
	const objFile = "reactor_jit_llvm_0000_ComputeProgram.o"
	defer os.Remove(objFile)
	return os.ReadFile(objFile)
}

type printer struct {
	bytes.Buffer
}

func (p *printer) printf(format string, args ...interface{}) {
	fmt.Fprintf(&p.Buffer, format, args...)
}
