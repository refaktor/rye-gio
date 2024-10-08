package main

import (
	/*RYEGEN: BEGIN IMPORTS*/
	"rye-gio/ryegen_bindings/gioui_org"
	/*RYEGEN: END IMPORTS*/

	"github.com/refaktor/rye/env"
	"github.com/refaktor/rye/evaldo"
	"github.com/refaktor/rye/runner"
)

func main() {
	runner.DoMain(func(ps *env.ProgramState) {
		/*RYEGEN: BEGIN BUILTINS*/
		evaldo.RegisterBuiltinsInContext(gioui_org.Builtins, ps, "gio")
		/*RYEGEN: END BUILTINS*/
	})
}
