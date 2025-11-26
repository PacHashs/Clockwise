package runtimelib

import "strings"

// Split splits s by sep and returns first two parts joined with a comma if more than one.
// (keeps API simple for small runtime)
func SplitFirstTwo(s, sep string) string {
    parts := strings.SplitN(s, sep, 2)
    if len(parts) == 1 {
        return parts[0]
    }
    return parts[0] + "," + parts[1]
}

// Trim trims spaces from s
func Trim(s string) string {
    return strings.TrimSpace(s)
}

// Slice returns substring of s starting at 'start' (inclusive) and ending at 'end' (exclusive).
// If end is negative, the substring extends to the end of the string.
func Slice(s string, start int, end int) string {
    rs := []rune(s)
    l := len(rs)
    if start < 0 {
        start = 0
    }
    if end < 0 || end > l {
        end = l
    }
    if start >= l || start >= end {
        return ""
    }
    return string(rs[start:end])
}
