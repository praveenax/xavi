package loader

import (
	"fmt"
	"os"
	"path/filepath"

	"xavi/compiler/ast"
	"xavi/compiler/lexer"
	"xavi/compiler/parser"
)

func LoadProgram(entryPath string) (*ast.Program, error) {
	visited := make(map[string]bool)
	return loadProgramFile(entryPath, visited)
}

func loadProgramFile(filePath string, visited map[string]bool) (*ast.Program, error) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, fmt.Errorf("resolve path %q: %w", filePath, err)
	}

	if visited[absPath] {
		return &ast.Program{
			Imports:   make([]*ast.Import, 0),
			Functions: make([]*ast.Function, 0),
		}, nil
	}
	visited[absPath] = true

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("read %q: %w", absPath, err)
	}

	program := parseProgram(string(data))
	merged := &ast.Program{
		Imports:   append([]*ast.Import(nil), program.Imports...),
		Functions: make([]*ast.Function, 0),
	}

	baseDir := filepath.Dir(absPath)
	for _, imp := range program.Imports {
		importPath := filepath.Join(baseDir, imp.Path)
		importedProgram, err := loadProgramFile(importPath, visited)
		if err != nil {
			return nil, err
		}

		functions := importedProgram.Functions
		if imp.Name != "" {
			functions, err = selectImportedFunctions(importedProgram.Functions, imp.Name, imp.Path)
			if err != nil {
				return nil, err
			}
		}

		merged.Functions = append(merged.Functions, functions...)
	}

	merged.Functions = append(merged.Functions, program.Functions...)
	return merged, nil
}

func parseProgram(src string) *ast.Program {
	lex := lexer.New(src)
	tokens := make([]lexer.Token, 0)

	for {
		token := lex.Next()
		tokens = append(tokens, token)
		if token.Type == lexer.EOF {
			break
		}
	}

	p := parser.New(tokens)
	return p.ParseProgram()
}

func selectImportedFunctions(functions []*ast.Function, name string, importPath string) ([]*ast.Function, error) {
	for _, fn := range functions {
		if fn.Name == name {
			return []*ast.Function{fn}, nil
		}
	}

	return nil, fmt.Errorf("import %q not found in %q", name, importPath)
}
