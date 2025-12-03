package parser

import (
    "xavi/compiler/lexer"
    "xavi/compiler/ast"
)

type Parser struct {
    tokens []lexer.Token
    pos    int
}

func New(tokens []lexer.Token) *Parser {
    return &Parser{tokens: tokens}
}

func (p *Parser) ParseFunction() *ast.Function {
    p.expect(lexer.FN)

    name := p.expect(lexer.IDENT).Literal
    p.expect(lexer.LPAREN)

    params := p.parseParams()

    p.expect(lexer.RPAREN)
    p.expect(lexer.ARROW)

    returnType := p.expect(lexer.IDENT).Literal

    body := p.parseBlock()

    return &ast.Function{
        Name:       name,
        Params:     params,
        ReturnType: returnType,
        Body:       body,
    }
}
