package lexer

// TokenType identifies the kind of token
type TokenType string

// Token represents a lexical token
type Token struct {
    Type TokenType
    Lit  string
}

const (
    ILLEGAL TokenType = "ILLEGAL"
    EOF     TokenType = "EOF"

    IDENT  TokenType = "IDENT"  // add, foobar, x, y, ...
    INT    TokenType = "INT"    // 1343456
    STRING TokenType = "STRING" // "foobar"

    ASSIGN   TokenType = "="
    PLUS     TokenType = "+"
    MINUS    TokenType = "-"
    BANG     TokenType = "!"
    ASTERISK TokenType = "*"
    SLASH    TokenType = "/"
    LT       TokenType = "<"
    GT       TokenType = ">"
    EQ       TokenType = "=="
    NOT_EQ   TokenType = "!="

    COMMA     TokenType = ","
    SEMICOLON TokenType = ";"

    LPAREN   TokenType = "("
    RPAREN   TokenType = ")"
    LBRACE   TokenType = "{"
    RBRACE   TokenType = "}"
    LBRACKET TokenType = "["
    RBRACKET TokenType = "]"
    DOT      TokenType = "."
    COLON    TokenType = ":"
    ARROW    TokenType = "->"

    FUNCTION TokenType = "FUNCTION"
    VAR      TokenType = "VAR"
    RETURN   TokenType = "RETURN"
    IF       TokenType = "IF"
    ELSE     TokenType = "ELSE"
    WHILE    TokenType = "WHILE"
    TRUE     TokenType = "TRUE"
    FALSE    TokenType = "FALSE"
)

var keywords = map[string]TokenType{
    "fn":     FUNCTION,
    "var":    VAR,
    "return": RETURN,
    "if":     IF,
    "else":   ELSE,
    "while":  WHILE,
    "true":   TRUE,
    "false":  FALSE,
}

// LookupIdent checks if an identifier is a reserved keyword
func LookupIdent(ident string) TokenType {
    if tok, ok := keywords[ident]; ok {
        return tok
    }
    return IDENT
}
