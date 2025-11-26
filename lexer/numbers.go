package lexer

import "unicode"

// readNumber scans a numeric literal (supports integers and simple floats)
// starting at pos in input and returns the literal and new position.
func readNumber(input string, pos int) (string, int) {
    i := pos
    seenDot := false
    for i < len(input) {
        r := rune(input[i])
        if unicode.IsDigit(r) {
            i++
            continue
        }
        if r == '.' && !seenDot {
            seenDot = true
            i++
            continue
        }
        break
    }
    return input[pos:i], i
}
