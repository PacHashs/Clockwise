package codegen

import "strings"

// EscapeString returns a Go-safe quoted string body (not surrounding quotes).
func EscapeString(s string) string {
    s = strings.ReplaceAll(s, "\\", "\\\\")
    s = strings.ReplaceAll(s, "\"", "\\\"")
    s = strings.ReplaceAll(s, "\n", "\\n")
    return s
}
