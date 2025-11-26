package lexer

import "fmt"

// TokenStream provides utilities for stepping through a token slice with lookahead
// and basic expectations (useful for parsers/tests).
type TokenStream struct {
    tokens []Token
    pos    int
}

func NewTokenStream(tokens []Token) *TokenStream { return &TokenStream{tokens: tokens, pos: 0} }

func (ts *TokenStream) Peek() Token { return ts.cur() }

func (ts *TokenStream) Next() Token {
    t := ts.cur()
    if ts.pos < len(ts.tokens) {
        ts.pos++
    }
    return t
}

func (ts *TokenStream) Backup() {
    if ts.pos > 0 {
        ts.pos--
    }
}

func (ts *TokenStream) cur() Token {
    if ts.pos >= len(ts.tokens) {
        return Token{Type: EOF, Lit: ""}
    }
    return ts.tokens[ts.pos]
}

// Expect ensures the next token matches tt or returns an error.
func (ts *TokenStream) Expect(tt TokenType) (Token, error) {
    t := ts.Next()
    if t.Type != tt {
        return t, fmt.Errorf("expected %s, got %s (%s)", tt, t.Type, t.Lit)
    }
    return t, nil
}
