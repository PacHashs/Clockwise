package lexer

import "unicode/utf8"

// RemoveBOM removes a UTF-8 BOM if present at the start of s.
func RemoveBOM(s string) string {
    if len(s) >= 3 && s[0] == 0xEF && s[1] == 0xBB && s[2] == 0xBF {
        return s[3:]
    }
    return s
}

// NormalizeSpaces collapses runs of whitespace into a single ASCII space
// and trims leading/trailing whitespace.
func NormalizeSpaces(s string) string {
    var out []rune
    lastSpace := false
    for len(s) > 0 {
        r, size := utf8.DecodeRuneInString(s)
        s = s[size:]
        if r == utf8.RuneError && size == 1 {
            // preserve invalid bytes as-is
            out = append(out, r)
            lastSpace = false
            continue
        }
        if isWhitespace(r) {
            if !lastSpace {
                out = append(out, ' ')
                lastSpace = true
            }
            continue
        }
        out = append(out, r)
        lastSpace = false
    }
    // trim leading/trailing
    if len(out) > 0 && out[0] == ' ' {
        out = out[1:]
    }
    if len(out) > 0 && out[len(out)-1] == ' ' {
        out = out[:len(out)-1]
    }
    return string(out)
}

// RuneCount returns the number of runes in s.
func RuneCount(s string) int {
    return utf8.RuneCountInString(s)
}
