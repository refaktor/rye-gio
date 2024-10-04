// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"runtime"
	"unsafe"

	"gioui.org/cpu"
	"gioui.org/cpu/example"
)

func main() {
	buffer := cpu.NewBuffer(1000)
	defer buffer.Free()
	const (
		imgWidth  = 17
		imgHeight = 19
	)
	img := cpu.NewImageRGBA(imgWidth, imgHeight)
	defer img.Free()

	var descSet example.ExampleDescriptorSetLayout
	*descSet.Binding0() = buffer
	*descSet.Binding1() = img

	ctx := cpu.NewDispatchContext()
	defer ctx.Free()
	bufData := buffer.Data()
	intBuf := unsafe.Slice((*int32)(unsafe.Pointer(&bufData[0])), len(bufData)/4)
	nthreads := runtime.NumCPU()
	ctx.Prepare(nthreads, example.ExampleProgramInfo, unsafe.Pointer(&descSet), 4, 3, 2)
	done := make(chan struct{})
	for i := 0; i < nthreads; i++ {
		i := i
		go func() {
			thread := cpu.NewThreadContext()
			defer thread.Free()
			ctx.Dispatch(i, thread)
			done <- struct{}{}
		}()
	}
	for i := 0; i < nthreads; i++ {
		<-done
	}
	fmt.Println(img.Data())
	rgba := &image.RGBA{
		Pix:    img.Data(),
		Stride: imgWidth * 4,
		Rect:   image.Rect(0, 0, imgWidth, imgHeight),
	}
	dump, err := os.Create("dump.png")
	if err != nil {
		panic(err)
	}
	if err := png.Encode(dump, rgba); err != nil {
		panic(err)
	}

	for i, v := range intBuf {
		if v == 0 {
			continue
		}
		fmt.Printf("buf[%d]: %d\n", i, v)
	}
}
