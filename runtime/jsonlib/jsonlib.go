package runtimelib

import (
    "encoding/json"
)

// JSONEscape returns a JSON string encoding of s (including quotes)
func JSONEscape(s string) string {
    b, _ := json.Marshal(s)
    return string(b)
}
