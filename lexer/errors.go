package lexer

import "fmt"

// LexError represents an error found during lexing.
type LexError struct{
    Pos Position
    Msg string
}

func (e *LexError) Error() string {
    return fmt.Sprintf("lex error at %d:%d: %s", e.Pos.Line, e.Pos.Column, e.Msg)
}
