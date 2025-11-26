package lexutils

import lex "codeberg.org/clockwise-lang/clockwise/lexer"

// MergePositions returns a position describing the start of p1 and end of p2
// (useful for diagnostics). This file holds small position utilities implemented
// in a separate subpackage to avoid circular dependencies.
func MergePositions(p1, p2 *lex.Position) *lex.Position {
    if p1 == nil {
        if p2 == nil {
            return lex.NewPosition()
        }
        return &lex.Position{Line: p2.Line, Column: p2.Column}
    }
    if p2 == nil {
        return &lex.Position{Line: p1.Line, Column: p1.Column}
    }
    return &lex.Position{Line: p1.Line, Column: p2.Column}
}
