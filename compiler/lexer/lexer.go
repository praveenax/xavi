package lexer

type TokenType string

const (
    // Keywords
    FN      TokenType = "FN"
    RECORD  TokenType = "RECORD"
    RETURN  TokenType = "RETURN"
    LET     TokenType = "LET"
    AGENT   TokenType = "AGENT"
    ON      TokenType = "ON"
    EVENT   TokenType = "EVENT"

    // Symbols
    COLON   TokenType = ":"
    ARROW   TokenType = "->"
    LPAREN  TokenType = "("
    RPAREN  TokenType = ")"

    // Literals & identifiers
    IDENT   TokenType = "IDENT"
    NUMBER  TokenType = "NUMBER"
    STRING  TokenType = "STRING"
    INDENT  TokenType = "INDENT"
    DEDENT  TokenType = "DEDENT"
    NEWLINE TokenType = "NEWLINE"

    EOF     TokenType = "EOF"
)

type Token struct {
    Type    TokenType
    Literal string
    Line    int
    Col     int
}
