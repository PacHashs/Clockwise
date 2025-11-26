package runtimelib

import "fmt"

// Print writes a string to stdout. Returns 0 for compatibility with earlier API.
func Print(s string) int {
    fmt.Print(s)
    return 0
}

// Sconcat concatenates strings
func Sconcat(a, b string) string {
    return a + b
}
