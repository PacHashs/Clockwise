package runtimelib

import "net/url"

// URLEncode percent-encodes the given string for use in a URL query.
func URLEncode(s string) string {
    return url.QueryEscape(s)
}

// URLDecode decodes a percent-encoded string; returns empty string on error.
func URLDecode(s string) string {
    v, err := url.QueryUnescape(s)
    if err != nil {
        return ""
    }
    return v
}
