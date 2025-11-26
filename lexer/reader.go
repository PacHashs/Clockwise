package lexer

import "io"

// Reader is a tiny rune reader wrapper used by more advanced lexers.
type Reader struct {
    r    io.RuneReader
    buff []rune
}

func NewReader(rr io.RuneReader) *Reader {
    return &Reader{r: rr}
}

// Next returns the next rune or io.EOF.
func (r *Reader) Next() (rune, error) {
    if len(r.buff) > 0 {
        ch := r.buff[0]
        r.buff = r.buff[1:]
        return ch, nil
    }
    ch, _, err := r.r.ReadRune()
    return ch, err
}

// Unread pushes back a rune so Next will return it again.
func (r *Reader) Unread(ch rune) {
    r.buff = append([]rune{ch}, r.buff...)
}
