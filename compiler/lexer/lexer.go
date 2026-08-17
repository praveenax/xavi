package lexer

import (
	"strconv"
	"strings"
	"unicode"
)

type TokenType string

const (
	FN     TokenType = "FN"
	RETURN TokenType = "RETURN"
	LET    TokenType = "LET"
	RECORD TokenType = "RECORD"
	AGENT  TokenType = "AGENT"
	ON     TokenType = "ON"
	EVENT  TokenType = "EVENT"

	COLON  TokenType = ":"
	ARROW  TokenType = "->"
	LPAREN TokenType = "("
	RPAREN TokenType = ")"
	COMMA  TokenType = ","
	ASSIGN TokenType = "="
	PLUS   TokenType = "+"
	MINUS  TokenType = "-"
	STAR   TokenType = "*"
	SLASH  TokenType = "/"

	IDENT   TokenType = "IDENT"
	NUMBER  TokenType = "NUMBER"
	STRING  TokenType = "STRING"
	INDENT  TokenType = "INDENT"
	DEDENT  TokenType = "DEDENT"
	NEWLINE TokenType = "NEWLINE"

	EOF TokenType = "EOF"
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Col     int
}

type Lexer struct {
	tokens []Token
	pos    int
}

func New(src string) *Lexer {
	return &Lexer{tokens: tokenize(src)}
}

func (l *Lexer) Next() Token {
	if l.pos >= len(l.tokens) {
		return Token{Type: EOF}
	}

	token := l.tokens[l.pos]
	l.pos++
	return token
}

func tokenize(src string) []Token {
	src = strings.ReplaceAll(src, "\r\n", "\n")

	lines := strings.Split(src, "\n")
	tokens := make([]Token, 0, len(lines)*4)
	indentStack := []int{0}

	for lineIndex, rawLine := range lines {
		lineNumber := lineIndex + 1

		if strings.TrimSpace(rawLine) == "" {
			tokens = append(tokens, Token{Type: NEWLINE, Line: lineNumber, Col: 1})
			continue
		}

		indentWidth := countIndent(rawLine)
		currentIndent := indentStack[len(indentStack)-1]

		switch {
		case indentWidth > currentIndent:
			indentStack = append(indentStack, indentWidth)
			tokens = append(tokens, Token{Type: INDENT, Line: lineNumber, Col: 1})
		case indentWidth < currentIndent:
			for len(indentStack) > 1 && indentWidth < indentStack[len(indentStack)-1] {
				indentStack = indentStack[:len(indentStack)-1]
				tokens = append(tokens, Token{Type: DEDENT, Line: lineNumber, Col: 1})
			}
		}

		line := strings.TrimLeft(rawLine, " \t")
		col := indentWidth + 1

		for i := 0; i < len(line); {
			ch := line[i]

			if ch == ' ' || ch == '\t' {
				i++
				col++
				continue
			}

			if isLetter(ch) {
				start := i
				startCol := col
				for i < len(line) && (isLetter(line[i]) || isDigit(line[i])) {
					i++
					col++
				}

				literal := line[start:i]
				tokens = append(tokens, Token{
					Type:    keywordType(literal),
					Literal: literal,
					Line:    lineNumber,
					Col:     startCol,
				})
				continue
			}

			if isDigit(ch) {
				start := i
				startCol := col
				dotSeen := false
				for i < len(line) {
					if line[i] == '.' && !dotSeen {
						dotSeen = true
						i++
						col++
						continue
					}
					if !isDigit(line[i]) {
						break
					}
					i++
					col++
				}

				tokens = append(tokens, Token{
					Type:    NUMBER,
					Literal: line[start:i],
					Line:    lineNumber,
					Col:     startCol,
				})
				continue
			}

			switch ch {
			case ':':
				tokens = append(tokens, Token{Type: COLON, Literal: ":", Line: lineNumber, Col: col})
			case '(':
				tokens = append(tokens, Token{Type: LPAREN, Literal: "(", Line: lineNumber, Col: col})
			case ')':
				tokens = append(tokens, Token{Type: RPAREN, Literal: ")", Line: lineNumber, Col: col})
			case ',':
				tokens = append(tokens, Token{Type: COMMA, Literal: ",", Line: lineNumber, Col: col})
			case '=':
				tokens = append(tokens, Token{Type: ASSIGN, Literal: "=", Line: lineNumber, Col: col})
			case '+':
				tokens = append(tokens, Token{Type: PLUS, Literal: "+", Line: lineNumber, Col: col})
			case '*':
				tokens = append(tokens, Token{Type: STAR, Literal: "*", Line: lineNumber, Col: col})
			case '/':
				tokens = append(tokens, Token{Type: SLASH, Literal: "/", Line: lineNumber, Col: col})
			case '-':
				if i+1 < len(line) && line[i+1] == '>' {
					tokens = append(tokens, Token{Type: ARROW, Literal: "->", Line: lineNumber, Col: col})
					i++
					col++
				} else {
					tokens = append(tokens, Token{Type: MINUS, Literal: "-", Line: lineNumber, Col: col})
				}
			case '"':
				startCol := col
				i++
				col++
				start := i
				for i < len(line) && line[i] != '"' {
					i++
					col++
				}
				literal := line[start:i]
				tokens = append(tokens, Token{
					Type:    STRING,
					Literal: literal,
					Line:    lineNumber,
					Col:     startCol,
				})
			}

			i++
			col++
		}

		tokens = append(tokens, Token{Type: NEWLINE, Line: lineNumber, Col: col})
	}

	for len(indentStack) > 1 {
		indentStack = indentStack[:len(indentStack)-1]
		tokens = append(tokens, Token{Type: DEDENT})
	}

	tokens = append(tokens, Token{Type: EOF, Line: len(lines) + 1, Col: 1})
	return tokens
}

func countIndent(line string) int {
	count := 0
	for _, ch := range line {
		if ch == ' ' {
			count++
			continue
		}
		if ch == '\t' {
			count += 4
			continue
		}
		break
	}
	return count
}

func keywordType(literal string) TokenType {
	switch literal {
	case "fn":
		return FN
	case "return":
		return RETURN
	case "let":
		return LET
	case "record":
		return RECORD
	case "agent":
		return AGENT
	case "on":
		return ON
	case "event":
		return EVENT
	default:
		return IDENT
	}
}

func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || ch == '_'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func ParseNumber(token Token) float64 {
	value, _ := strconv.ParseFloat(token.Literal, 64)
	return value
}
