package lexer

import "unicode"

// IsIdentStart reports if r is a valid identifier start rune.
func IsIdentStart(r rune) bool {
    return unicode.IsLetter(r) || r == '_'
}

// IsIdentPart reports if r is valid inside an identifier.
func IsIdentPart(r rune) bool {
    return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}
