package lexer

import (
    "bufio"
    "strings"
)

// StatefulLexer is a small scaffold that demonstrates a more capable lexer
// implementation while remaining independent of the simple lexer in lexer.go.
type StatefulLexer struct {
    r  *bufio.Reader
    pos *Position
    src string
}

func NewStateful(src string) *StatefulLexer {
    return &StatefulLexer{r: bufio.NewReader(strings.NewReader(src)), pos: NewPosition(), src: src}
}

// PeekRune returns the next rune without consuming it.
func (s *StatefulLexer) PeekRune() (rune, error) {
    ch, _, err := s.r.ReadRune()
    if err != nil {
        return 0, err
    }
    _ = s.r.UnreadRune()
    return ch, nil
}

// NextRune returns the next rune and advances position.
func (s *StatefulLexer) NextRune() (rune, error) {
    ch, _, err := s.r.ReadRune()
    if err != nil {
        return 0, err
    }
    s.pos.Advance(ch)
    return ch, nil
}
