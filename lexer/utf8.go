package lexer

// RuneAt returns the rune at rune-index i in s, or 0 if out of range.
func RuneAt(s string, i int) rune {
    var idx int
    for _, r := range s {
        if idx == i {
            return r
        }
        idx++
    }
    return 0
}

// RuneSlice returns substring of s from rune index 'start' to 'end' (end exclusive).
// If end<0 it's treated as len(runes).
func RuneSlice(s string, start, end int) string {
    var out []rune
    var idx int
    for _, r := range s {
        if idx >= start && (end < 0 || idx < end) {
            out = append(out, r)
        }
        idx++
    }
    return string(out)
}
