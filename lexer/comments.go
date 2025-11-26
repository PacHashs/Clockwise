package lexer

// ConsumeBlockComment consumes a C-style block comment from input starting
// at the given position. It returns the index after the closing '*/' or the
// original pos if no block comment was found. This helper is intentionally
// simple and does not handle nested comments.
func ConsumeBlockComment(input string, pos int) int {
    if pos+1 >= len(input) || input[pos] != '/' || input[pos+1] != '*' {
        return pos
    }
    i := pos + 2
    for i+1 < len(input) {
        if input[i] == '*' && input[i+1] == '/' {
            return i + 2
        }
        i++
    }
    return i
}
