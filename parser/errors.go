package parser

import (
    "fmt"
    "codeberg.org/clockwise-lang/clockwise/lexer"
)

// ParseError represents a parser error with optional position information.
type ParseError struct {
    Msg string
    Pos *lexer.Position
}

func (e *ParseError) Error() string {
    if e == nil {
        return "parser: <nil>"
    }
    if e.Pos == nil {
        return e.Msg
    }
    return fmt.Sprintf("%s at line %d col %d", e.Msg, e.Pos.Line, e.Pos.Column)
}

