package runtimelib

import (
    "bytes"
    "compress/gzip"
    "encoding/base64"
    "io"
)

// Gzip encodes the input as gzip and returns base64 text (empty on error).
func Gzip(s string) string {
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

// Gunzip decodes base64(gzip(data)) back to string (empty on error).
func Gunzip(b64 string) string {
    raw, err := base64.StdEncoding.DecodeString(b64)
    if err != nil {
        return ""
    }
    zr, err := gzip.NewReader(bytes.NewReader(raw))
    if err != nil {
        return ""
    }
    defer zr.Close()
    out, err := io.ReadAll(zr)
    if err != nil {
        return ""
    }
    return string(out)
}
