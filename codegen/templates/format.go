package codegen

import "strings"

// TrimTrailingWhitespace normalizes generated code by trimming trailing spaces
// on each line. Small helper for code formatting before writing to disk.
func TrimTrailingWhitespace(src string) string {
    lines := strings.Split(src, "\n")
    for i, l := range lines {
        lines[i] = strings.TrimRight(l, " \t")
    }
    return strings.Join(lines, "\n")
}
