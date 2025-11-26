package runtimelib

import "time"

// ParseISO parses an RFC3339 timestamp, returns empty on error
func ParseISO(s string) string {
    t, err := time.Parse(time.RFC3339, s)
    if err != nil {
        return ""
    }
    return t.Format(time.RFC3339)
}

// FormatISO formats Unix seconds into RFC3339 string
func FormatISO(sec int64) string {
    return time.Unix(sec, 0).Format(time.RFC3339)
}
