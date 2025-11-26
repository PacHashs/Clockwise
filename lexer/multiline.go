package lexer

// ReadBlockString reads a block string starting at pos using backtick delimiters.
// It returns the string content (without delimiters) and the index after the closing backtick.
// If no closing backtick is found, it returns the content up to the end and the end index.
func ReadBlockString(input string, pos int) (string, int) {
    if pos >= len(input) || input[pos] != '`' {
        return "", pos
    }
    i := pos + 1
    start := i
    for i < len(input) {
        if input[i] == '`' {
            return input[start:i], i + 1
        }
        i++
    }
    // unterminated, return remainder
    return input[start:], i
}
