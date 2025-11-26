package runtimelib

import (
    "bytes"
    "io"
    "net/http"
)

// HttpGetStatus performs GET and returns response body as string (empty on error) and status code
func HttpGetStatus(url string) (string, int) {
    resp, err := http.Get(url)
    if err != nil {
        return "", 0
    }
    defer resp.Body.Close()
    b, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", resp.StatusCode
    }
    return string(b), resp.StatusCode
}

// HttpPost sends a POST with content-type text/plain and returns response body or empty on error
func HttpPost(url, body string) string {
    resp, err := http.Post(url, "text/plain", bytes.NewBufferString(body))
    if err != nil {
        return ""
    }
    defer resp.Body.Close()
    b, err := io.ReadAll(resp.Body)
    if err != nil {
        return ""
    }
    return string(b)
}
