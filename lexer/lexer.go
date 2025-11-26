package lexer

import (
    "unicode"
)

type Lexer struct {
    input        string
    position     int
    readPosition int
    ch           rune
}

func New(input string) *Lexer {
    l := &Lexer{input: input}
    l.readChar()
    return l
}

func (l *Lexer) readChar() {
    if l.readPosition >= len(l.input) {
        l.ch = 0
    } else {
        l.ch = rune(l.input[l.readPosition])
    }
    l.position = l.readPosition
    l.readPosition++
}

func (l *Lexer) peekChar() rune {
    if l.readPosition >= len(l.input) {
        return 0
    }
    return rune(l.input[l.readPosition])
}

func (l *Lexer) Tokenize() []Token {
    var tokens []Token
    for {
        l.skipWhitespace()
        var tok Token
        switch l.ch {
        case '=':
            if l.peekChar() == '=' {
                ch := l.ch
                l.readChar()
                lit := string(ch) + string(l.ch)
                tok = Token{Type: EQ, Lit: lit}
            } else {
                tok = Token{Type: ASSIGN, Lit: string(l.ch)}
            }
        case '+':
            tok = Token{Type: PLUS, Lit: string(l.ch)}
        case '-':
            if l.peekChar() == '>' {
                ch := l.ch
                l.readChar()
                lit := string(ch) + string(l.ch)
                tok = Token{Type: ARROW, Lit: lit}
            } else {
                tok = Token{Type: MINUS, Lit: string(l.ch)}
            }
        case '!':
            if l.peekChar() == '=' {
                ch := l.ch
                l.readChar()
                lit := string(ch) + string(l.ch)
                tok = Token{Type: NOT_EQ, Lit: lit}
            } else {
                tok = Token{Type: BANG, Lit: string(l.ch)}
            }
        case '/':
            if l.peekChar() == '/' {
                // line comment
                l.readChar()
                l.readChar()
                for l.ch != '\n' && l.ch != 0 {
                    l.readChar()
                }
                l.readChar()
                continue
            }
            tok = Token{Type: SLASH, Lit: string(l.ch)}
        case '*':
            tok = Token{Type: ASTERISK, Lit: string(l.ch)}
        case '<':
            tok = Token{Type: LT, Lit: string(l.ch)}
        case '>':
            tok = Token{Type: GT, Lit: string(l.ch)}
        case ';':
            tok = Token{Type: SEMICOLON, Lit: string(l.ch)}
        case ',':
            tok = Token{Type: COMMA, Lit: string(l.ch)}
        case '(':
            tok = Token{Type: LPAREN, Lit: string(l.ch)}
        case ')':
            tok = Token{Type: RPAREN, Lit: string(l.ch)}
        case '{':
            tok = Token{Type: LBRACE, Lit: string(l.ch)}
        case '}':
            tok = Token{Type: RBRACE, Lit: string(l.ch)}
        case ':':
            tok = Token{Type: COLON, Lit: string(l.ch)}
        case '[':
            tok = Token{Type: LBRACKET, Lit: string(l.ch)}
        case ']':
            tok = Token{Type: RBRACKET, Lit: string(l.ch)}
        case '.':
            tok = Token{Type: DOT, Lit: string(l.ch)}
        case '"':
            tok.Type = STRING
            tok.Lit = l.readString()
        case 0:
            tok.Lit = ""
            tok.Type = EOF
            tokens = append(tokens, tok)
            return tokens
        default:
            if isLetter(l.ch) {
                lit := l.readIdentifier()
                tok.Type = LookupIdent(lit)
                tok.Lit = lit
                tokens = append(tokens, tok)
                continue
            } else if isDigit(l.ch) {
                lit := l.readNumber()
                tok.Type = INT
                tok.Lit = lit
                tokens = append(tokens, tok)
                continue
            } else {
                tok = Token{Type: ILLEGAL, Lit: string(l.ch)}
            }
        }
        tokens = append(tokens, tok)
        l.readChar()
    }
}

func (l *Lexer) skipWhitespace() {
    for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
        l.readChar()
    }
}

func (l *Lexer) readIdentifier() string {
    pos := l.position
    for isLetter(l.ch) || isDigit(l.ch) {
        l.readChar()
    }
    return l.input[pos:l.position]
}

func (l *Lexer) readNumber() string {
    pos := l.position
    for isDigit(l.ch) {
        l.readChar()
    }
    return l.input[pos:l.position]
}

func (l *Lexer) readString() string {
    // assume starting quote was consumed
    l.readChar()
    pos := l.position
    for l.ch != '"' && l.ch != 0 {
        l.readChar()
    }
    s := l.input[pos:l.position]
    // closing quote will be consumed by caller loop's readChar
    return s
}

func isLetter(ch rune) bool {
    return unicode.IsLetter(ch) || ch == '_'
}

func isDigit(ch rune) bool {
    return '0' <= ch && ch <= '9'
}
