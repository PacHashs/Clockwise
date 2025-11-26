package runtimelib

import (
    "bytes"
    "compress/gzip"
    "encoding/base64"
)

// GzipBase64 compresses input string with gzip and returns base64-encoded bytes.
// Returns empty string on error.
func GzipBase64(s string) string {
    var buf bytes.Buffer
    zw := gzip.NewWriter(&buf)
    if _, err := zw.Write([]byte(s)); err != nil {
        return ""
    }
    if err := zw.Close(); err != nil {
        return ""
    }
    return base64.StdEncoding.EncodeToString(buf.Bytes())
}
