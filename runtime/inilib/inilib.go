package runtimelib

import "strings"

// ParseINI finds the value for key in a simple INI-like content where lines
// look like `key=value`. Returns empty string if key not found.
func ParseINI(content, key string) string {
    lines := strings.Split(content, "\n")
    for _, l := range lines {
        l = strings.TrimSpace(l)
        if l == "" || strings.HasPrefix(l, "#") || strings.HasPrefix(l, ";") {
            continue
        }
        parts := strings.SplitN(l, "=", 2)
        if len(parts) != 2 {
            continue
        }
        if strings.TrimSpace(parts[0]) == key {
            return strings.TrimSpace(parts[1])
        }
    }
    return ""
}
