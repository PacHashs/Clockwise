package runtimelib

import (
    "bytes"
    "io"
    "net/http"
    "os"
)

// Simple HTTP GET returning body as string (empty on error)
func HttpGet(url string) string {
    resp, err := http.Get(url)
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

// HttpGetStatus performs GET and returns body and status code (0 on error)
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

// HttpPost sends a POST with content-type text/plain and returns body or empty on error
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

// DownloadFile downloads url to path; returns error or nil
func DownloadFile(url, path string) error {
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    out, err := os.Create(path)
    if err != nil {
        return err
    }
    defer out.Close()
    _, err = io.Copy(out, resp.Body)
    return err
}
