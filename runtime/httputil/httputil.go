package runtimelib

import "net/url"

// JoinURL joins base and path into a single URL string (naive helper).
func JoinURL(base, path string) string {
    b, err := url.Parse(base)
    if err != nil {
        return ""
    }
    p, err := url.Parse(path)
    if err != nil {
        return ""
    }
    return b.ResolveReference(p).String()
}
