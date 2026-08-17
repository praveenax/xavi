package bytecode

import (
	"fmt"

	"xavi/compiler/ast"
	"xavi/vm/opcode"
)

type CompiledProgram struct {
	Functions     []*CompiledFunction
	FunctionIndex map[string]uint8
}

type CompiledFunction struct {
	Name      string
	Params    []string
	Bytecode  []uint8
	Consts    []interface{}
	FrameSize int
}

// Generator converts Xavi AST nodes into Xavi VM bytecode.
type Generator struct {
	Vars          map[string]int
	Consts        []interface{}
	Bytecode      []uint8
	FunctionIndex map[string]uint8
}

// NewGenerator creates a fresh bytecode generator.
func NewGenerator() *Generator {
	return &Generator{
		Vars:          make(map[string]int),
		Consts:        make([]interface{}, 0),
		Bytecode:      make([]uint8, 0),
		FunctionIndex: make(map[string]uint8),
	}
}

// addConst adds a value to the constant pool
// and returns its index.
func (g *Generator) addConst(value interface{}) uint8 {
	for i, existing := range g.Consts {
		if existing == value {
			return uint8(i)
		}
	}

	g.Consts = append(g.Consts, value)
	return uint8(len(g.Consts) - 1)
}

func (g *Generator) GenerateProgram(program *ast.Program) (*CompiledProgram, error) {
	compiled := &CompiledProgram{
		Functions:     make([]*CompiledFunction, 0, len(program.Functions)),
		FunctionIndex: make(map[string]uint8, len(program.Functions)),
	}

	for i, fn := range program.Functions {
		if _, exists := compiled.FunctionIndex[fn.Name]; exists {
			return nil, fmt.Errorf("duplicate function: %s", fn.Name)
		}
		compiled.FunctionIndex[fn.Name] = uint8(i)
	}

	g.FunctionIndex = compiled.FunctionIndex

	for _, fn := range program.Functions {
		compiledFn := g.generateFunction(fn)
		compiled.Functions = append(compiled.Functions, compiledFn)
	}

	return compiled, nil
}

// generateFunction converts a function AST into bytecode.
func (g *Generator) generateFunction(fn *ast.Function) *CompiledFunction {
	g.Vars = make(map[string]int)
	g.Consts = make([]interface{}, 0)
	g.Bytecode = make([]uint8, 0)

	params := make([]string, 0, len(fn.Params))
	for index, param := range fn.Params {
		g.Vars[param.Name] = index
		params = append(params, param.Name)
	}

	for _, statement := range fn.Body {
		switch stmt := statement.(type) {
		case *ast.LetStmt:
			g.emitExpr(stmt.Value)

			index, exists := g.Vars[stmt.Name]
			if !exists {
				index = len(g.Vars)
				g.Vars[stmt.Name] = index
			}

			g.Bytecode = append(g.Bytecode, opcode.STORE_VAR, uint8(index))

		case *ast.ReturnStmt:
			g.emitExpr(stmt.Value)
			g.Bytecode = append(g.Bytecode, opcode.RETURN)

		case *ast.ExprStmt:
			g.emitExpr(stmt.Value)
			g.Bytecode = append(g.Bytecode, opcode.POP)

		default:
			panic(fmt.Sprintf("unsupported statement type %T", stmt))
		}
	}

	return &CompiledFunction{
		Name:      fn.Name,
		Params:    params,
		Bytecode:  append([]uint8(nil), g.Bytecode...),
		Consts:    append([]interface{}(nil), g.Consts...),
		FrameSize: len(g.Vars),
	}
}

// emitExpr converts an expression into bytecode.
func (g *Generator) emitExpr(expression ast.Expr) {
	switch expr := expression.(type) {
	case *ast.NumberLiteral:
		constIndex := g.addConst(expr.Value)
		g.Bytecode = append(g.Bytecode, opcode.LOAD_CONST, constIndex)

	case *ast.StringLiteral:
		constIndex := g.addConst(expr.Value)
		g.Bytecode = append(g.Bytecode, opcode.LOAD_CONST, constIndex)

	case *ast.Ident:
		varIndex, exists := g.Vars[expr.Name]
		if !exists {
			panic(fmt.Sprintf("undefined variable: %s", expr.Name))
		}

		g.Bytecode = append(g.Bytecode, opcode.LOAD_VAR, uint8(varIndex))

	case *ast.BinaryExpr:
		g.emitExpr(expr.Left)
		g.emitExpr(expr.Right)

		switch expr.Op {
		case "+":
			g.Bytecode = append(g.Bytecode, opcode.ADD)
		case "-":
			g.Bytecode = append(g.Bytecode, opcode.SUB)
		case "*":
			g.Bytecode = append(g.Bytecode, opcode.MUL)
		case "/":
			g.Bytecode = append(g.Bytecode, opcode.DIV)
		default:
			panic(fmt.Sprintf("unsupported operator: %s", expr.Op))
		}

	case *ast.CallExpr:
		for _, arg := range expr.Args {
			g.emitExpr(arg)
		}

		if builtinIndex, ok := opcode.BuiltinIndex(expr.Callee); ok {
			g.Bytecode = append(g.Bytecode, opcode.CALL_BUILTIN, builtinIndex, uint8(len(expr.Args)))
			return
		}

		functionIndex, exists := g.FunctionIndex[expr.Callee]
		if !exists {
			panic(fmt.Sprintf("undefined function: %s", expr.Callee))
		}

		g.Bytecode = append(g.Bytecode, opcode.CALL, functionIndex, uint8(len(expr.Args)))

	default:
		panic(fmt.Sprintf("unsupported expression type %T", expr))
	}
}
