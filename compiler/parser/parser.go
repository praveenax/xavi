package parser

import (
	"fmt"

	"xavi/compiler/ast"
	"xavi/compiler/lexer"
)

type Parser struct {
	tokens []lexer.Token
	pos    int
}

func New(tokens []lexer.Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Imports:   make([]*ast.Import, 0),
		Functions: make([]*ast.Function, 0),
	}

	for p.current().Type != lexer.EOF {
		p.skipNewlines()
		if p.current().Type == lexer.EOF {
			break
		}

		switch p.current().Type {
		case lexer.IMPORT:
			program.Imports = append(program.Imports, p.parseImport())
		case lexer.FN:
			program.Functions = append(program.Functions, p.ParseFunction())
		default:
			panic(fmt.Sprintf(
				"unsupported top-level statement %q at line %d, col %d",
				p.current().Literal,
				p.current().Line,
				p.current().Col,
			))
		}
		p.skipNewlines()
	}

	return program
}

func (p *Parser) parseImport() *ast.Import {
	p.expect(lexer.IMPORT)
	name := ""
	if p.current().Type == lexer.IDENT {
		name = p.advance().Literal
		p.expect(lexer.LT)
	}
	path := p.expect(lexer.STRING)
	return &ast.Import{
		Name: name,
		Path: path.Literal,
	}
}

func (p *Parser) ParseFunction() *ast.Function {
	p.expect(lexer.FN)
	name := p.expect(lexer.IDENT).Literal
	p.expect(lexer.LPAREN)
	params := p.parseParams()
	p.expect(lexer.RPAREN)

	returnType := ""
	if p.match(lexer.ARROW) {
		returnType = p.expect(lexer.IDENT).Literal
	}

	p.expect(lexer.COLON)
	body := p.parseBlock()

	return &ast.Function{
		Name:       name,
		Params:     params,
		ReturnType: returnType,
		Body:       body,
	}
}

func (p *Parser) parseParams() []ast.Param {
	params := make([]ast.Param, 0)

	if p.current().Type == lexer.RPAREN {
		return params
	}

	for {
		name := p.expect(lexer.IDENT).Literal
		paramType := ""
		if p.match(lexer.COLON) {
			paramType = p.expect(lexer.IDENT).Literal
		}

		params = append(params, ast.Param{
			Name: name,
			Type: paramType,
		})

		if !p.match(lexer.COMMA) {
			break
		}
	}

	return params
}

func (p *Parser) parseBlock() []ast.Stmt {
	p.expect(lexer.NEWLINE)
	p.expect(lexer.INDENT)

	body := make([]ast.Stmt, 0)

	for {
		p.skipNewlines()
		if p.current().Type == lexer.DEDENT || p.current().Type == lexer.EOF {
			break
		}

		body = append(body, p.parseStatement())
		if p.current().Type == lexer.NEWLINE {
			p.advance()
		}
	}

	p.expect(lexer.DEDENT)
	return body
}

func (p *Parser) parseStatement() ast.Stmt {
	switch p.current().Type {
	case lexer.LET:
		p.advance()
		name := p.expect(lexer.IDENT).Literal
		p.expect(lexer.ASSIGN)
		return &ast.LetStmt{
			Name:  name,
			Value: p.parseExpr(),
		}
	case lexer.RETURN:
		p.advance()
		return &ast.ReturnStmt{
			Value: p.parseExpr(),
		}
	default:
		panic(fmt.Sprintf(
			"unsupported statement %q at line %d, col %d",
			p.current().Literal,
			p.current().Line,
			p.current().Col,
		))
	}
}

func (p *Parser) parseExpr() ast.Expr {
	return p.parseAddSub()
}

func (p *Parser) parseAddSub() ast.Expr {
	left := p.parseMulDiv()

	for {
		op := p.current().Type
		if op != lexer.PLUS && op != lexer.MINUS {
			return left
		}

		operator := p.advance().Literal
		right := p.parseMulDiv()
		left = &ast.BinaryExpr{
			Op:    operator,
			Left:  left,
			Right: right,
		}
	}
}

func (p *Parser) parseMulDiv() ast.Expr {
	left := p.parseCall()

	for {
		op := p.current().Type
		if op != lexer.STAR && op != lexer.SLASH {
			return left
		}

		operator := p.advance().Literal
		right := p.parsePrimary()
		left = &ast.BinaryExpr{
			Op:    operator,
			Left:  left,
			Right: right,
		}
	}
}

func (p *Parser) parseCall() ast.Expr {
	expr := p.parsePrimary()

	for p.current().Type == lexer.LPAREN {
		ident, ok := expr.(*ast.Ident)
		if !ok {
			panic(fmt.Sprintf(
				"unsupported call target at line %d, col %d",
				p.current().Line,
				p.current().Col,
			))
		}

		p.expect(lexer.LPAREN)
		args := p.parseArgs()
		p.expect(lexer.RPAREN)

		expr = &ast.CallExpr{
			Callee: ident.Name,
			Args:   args,
		}
	}

	return expr
}

func (p *Parser) parseArgs() []ast.Expr {
	args := make([]ast.Expr, 0)
	if p.current().Type == lexer.RPAREN {
		return args
	}

	for {
		args = append(args, p.parseExpr())
		if !p.match(lexer.COMMA) {
			break
		}
	}

	return args
}

func (p *Parser) parsePrimary() ast.Expr {
	token := p.current()

	switch token.Type {
	case lexer.IDENT:
		p.advance()
		return &ast.Ident{Name: token.Literal}
	case lexer.NUMBER:
		p.advance()
		return &ast.NumberLiteral{Value: lexer.ParseNumber(token)}
	case lexer.LPAREN:
		p.advance()
		expr := p.parseExpr()
		p.expect(lexer.RPAREN)
		return expr
	default:
		panic(fmt.Sprintf(
			"unsupported expression %q at line %d, col %d",
			token.Literal,
			token.Line,
			token.Col,
		))
	}
}

func (p *Parser) expect(expected lexer.TokenType) lexer.Token {
	token := p.current()
	if token.Type != expected {
		panic(fmt.Sprintf(
			"expected %s, got %s at line %d, col %d",
			expected,
			token.Type,
			token.Line,
			token.Col,
		))
	}

	p.pos++
	return token
}

func (p *Parser) match(expected lexer.TokenType) bool {
	if p.current().Type != expected {
		return false
	}

	p.pos++
	return true
}

func (p *Parser) skipNewlines() {
	for p.current().Type == lexer.NEWLINE {
		p.pos++
	}
}

func (p *Parser) advance() lexer.Token {
	token := p.current()
	p.pos++
	return token
}

func (p *Parser) current() lexer.Token {
	if p.pos >= len(p.tokens) {
		return lexer.Token{Type: lexer.EOF}
	}

	return p.tokens[p.pos]
}
