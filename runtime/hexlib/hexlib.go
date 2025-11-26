package runtimelib

import "encoding/hex"

// HexEncode encodes bytes (as string) to hex string
func HexEncode(s string) string {
    return hex.EncodeToString([]byte(s))
}

// HexDecode decodes a hex string, returns empty string on error
func HexDecode(h string) string {
    b, err := hex.DecodeString(h)
    if err != nil {
        return ""
    }
    return string(b)
}
