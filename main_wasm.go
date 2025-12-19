//go:build js && wasm

package main

import (
	"syscall/js"
)

func main() {
	c := make(chan struct{})
	js.Global().Set("evaluate", js.FuncOf(evaluate))
	<-c
}

func evaluate(_ js.Value, args []js.Value) interface{} {
	if len(args) != 2 {
		return "Invalid number of arguments"
	}
	yamlDoc := args[0].String()
	jsonPath := args[1].String()
	return playgroundEvaluate(yamlDoc, jsonPath)
}
