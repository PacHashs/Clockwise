package checker

import (
    "fmt"
    "codeberg.org/clockwise-lang/clockwise/parser"
)

// Diagnostic represents a single checker warning or error with optional node context.
type Diagnostic struct {
    Msg  string
    Node parser.Node
}

func (d Diagnostic) Error() string {
    return d.Msg
}

// ReportError returns a formatted error for a diagnostics list.
func ReportError(d Diagnostics) error {
    if len(d) == 0 {
        return nil
    }
    // return first error for simplicity
    return fmt.Errorf(d[0].Msg)
}

// Diagnostics is a slice of Diagnostic
type Diagnostics []Diagnostic
