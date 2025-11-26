package lexer

import "unicode"

// isWhitespace returns true for typical ASCII whitespace runes.
func isWhitespace(r rune) bool {
    return r == ' ' || r == '\t' || r == '\n' || r == '\r' || unicode.IsSpace(r)
}

