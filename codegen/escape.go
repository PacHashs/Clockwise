package codegen

import "strings"

// EscapeString returns a Go-safe string body (without surrounding quotes).
func EscapeString(s string) string {
    s = strings.ReplaceAll(s, "\\", "\\\\")
    s = strings.ReplaceAll(s, "\"", "\\\"")
    s = strings.ReplaceAll(s, "\n", "\\n")
    return s
}
