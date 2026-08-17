package bytecode

import (
	"fmt"

	"xavi/compiler/ast"
	"xavi/vm/exec"
)

// Generator converts Xavi AST nodes into Xavi VM bytecode.
type Generator struct {
	Vars     map[string]int
	Consts   []float64
	Bytecode []uint8
}

// NewGenerator creates a fresh bytecode generator.
func NewGenerator() *Generator {
	return &Generator{
		Vars:     make(map[string]int),
		Consts:   make([]float64, 0),
		Bytecode: make([]uint8, 0),
	}
}

// addConst adds a numeric value to the constant pool
// and returns its index.
func (g *Generator) addConst(value float64) uint8 {
	for i, existing := range g.Consts {
		if existing == value {
			return uint8(i)
		}
	}

	g.Consts = append(g.Consts, value)
	return uint8(len(g.Consts) - 1)
}

// Generate converts a function AST into bytecode.
func (g *Generator) Generate(fn *ast.Function) ([]uint8, []float64) {
	g.Vars = make(map[string]int)
	g.Consts = make([]float64, 0)
	g.Bytecode = make([]uint8, 0)

	for _, statement := range fn.Body {
		switch stmt := statement.(type) {
		case *ast.LetStmt:
			g.emitExpr(stmt.Value)

			index, exists := g.Vars[stmt.Name]
			if !exists {
				index = len(g.Vars)
				g.Vars[stmt.Name] = index
			}

			g.Bytecode = append(g.Bytecode, exec.OP_STORE_VAR, uint8(index))

		case *ast.ReturnStmt:
			g.emitExpr(stmt.Value)
			g.Bytecode = append(g.Bytecode, exec.OP_RETURN)

		default:
			panic(fmt.Sprintf("unsupported statement type %T", stmt))
		}
	}

	return g.Bytecode, g.Consts
}

// emitExpr converts an expression into bytecode.
func (g *Generator) emitExpr(expression ast.Expr) {
	switch expr := expression.(type) {
	case *ast.NumberLiteral:
		constIndex := g.addConst(expr.Value)
		g.Bytecode = append(g.Bytecode, exec.OP_LOAD_CONST, constIndex)

	case *ast.Ident:
		varIndex, exists := g.Vars[expr.Name]
		if !exists {
			panic(fmt.Sprintf("undefined variable: %s", expr.Name))
		}

		g.Bytecode = append(g.Bytecode, exec.OP_LOAD_VAR, uint8(varIndex))

	case *ast.BinaryExpr:
		g.emitExpr(expr.Left)
		g.emitExpr(expr.Right)

		switch expr.Op {
		case "+":
			g.Bytecode = append(g.Bytecode, exec.OP_ADD)
		case "-":
			g.Bytecode = append(g.Bytecode, exec.OP_SUB)
		case "*":
			g.Bytecode = append(g.Bytecode, exec.OP_MUL)
		case "/":
			g.Bytecode = append(g.Bytecode, exec.OP_DIV)
		default:
			panic(fmt.Sprintf("unsupported operator: %s", expr.Op))
		}

	default:
		panic(fmt.Sprintf("unsupported expression type %T", expr))
	}
}
