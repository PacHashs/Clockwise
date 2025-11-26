package lexer

// Position tracks the line/column of a token or rune.
type Position struct {
    Line   int
    Column int
}

func (p *Position) Advance(r rune) {
    if r == '\n' {
        p.Line++
        p.Column = 0
        return
    }
    p.Column++
}

func NewPosition() *Position { return &Position{Line: 1, Column: 0} }
