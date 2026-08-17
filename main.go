package main

import (
	"fmt"

	"os"
	"xavi/compiler/bytecode"
	"xavi/vm/exec"
	"xavi/vm/loader"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <file.xavi>")
		os.Exit(1)
	}

	fileName := os.Args[1]

	// 1. Load the entry file and any imported source files.
	program, err := loader.LoadProgram(fileName)
	if err != nil {
		fmt.Println("Error loading Xavi file:", err)
		os.Exit(1)
	}

	if len(program.Functions) == 0 {
		fmt.Println("Error: no function found")
		os.Exit(1)
	}

	// 2. Generate Xavi bytecode
	generator := bytecode.NewGenerator()
	compiledProgram, err := generator.GenerateProgram(program)
	if err != nil {
		fmt.Println("Error generating bytecode:", err)
		os.Exit(1)
	}

	// 3. Execute bytecode using Xavi VM.
	vm := exec.NewInterpreter(compiledProgram)

	result := vm.Run("main")

	fmt.Println("Result:", result)
}
