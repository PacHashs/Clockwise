package codegen

import "strings"

// Quote wraps and escapes a string as a Go string literal.
func Quote(s string) string {
    // inline safe escaping
    s = strings.ReplaceAll(s, "\\", "\\\\")
    s = strings.ReplaceAll(s, "\"", "\\\"")
    s = strings.ReplaceAll(s, "\n", "\\n")
    return "\"" + s + "\""
}
