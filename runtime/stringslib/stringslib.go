package runtimelib

import "strings"

// Concat joins two strings
func Concat(a, b string) string {
    return a + b
}

// ToUpper returns the uppercase string
func ToUpper(s string) string {
    return strings.ToUpper(s)
}
