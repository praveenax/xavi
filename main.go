package main

import (
	"fmt"
	"os"

	"xavi/compiler/bytecode"
	"xavi/compiler/lexer"
	"xavi/compiler/parser"
	"xavi/vm/exec"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <file.xavi>")
		os.Exit(1)
	}

	fileName := os.Args[1]

	data, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println("Error reading Xavi file:", err)
		os.Exit(1)
	}

	src := string(data)

	// 1. Lexical analysis
	lex := lexer.New(src)

	var tokens []lexer.Token

	for {
		token := lex.Next()
		tokens = append(tokens, token)

		if token.Type == lexer.EOF {
			break
		}
	}

	// 2. Parse tokens into AST
	p := parser.New(tokens)
	program := p.ParseProgram()

	if len(program.Functions) == 0 {
		fmt.Println("Error: no function found")
		os.Exit(1)
	}

	// For now, execute the first function.
	fn := program.Functions[0]

	// 3. Generate Xavi bytecode
	generator := bytecode.NewGenerator()

	bc, constants := generator.Generate(fn)

	// 4. Allocate enough local variable slots.
	frameSize := len(generator.Vars)

	// 5. Execute bytecode using Xavi VM.
	vm := exec.NewInterpreter(
		bc,
		constants,
		frameSize,
	)

	result := vm.Run()

	fmt.Println("Result:", result)
}