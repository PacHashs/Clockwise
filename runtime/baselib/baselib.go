package runtimelib

import "encoding/base64"

// Base64Encode encodes s to base64 (no newline)
func Base64Encode(s string) string {
    return base64.StdEncoding.EncodeToString([]byte(s))
}

// Base64Decode decodes a base64 string, returns empty string on error
func Base64Decode(s string) string {
    b, err := base64.StdEncoding.DecodeString(s)
    if err != nil {
        return ""
    }
    return string(b)
}
