package runtimelib

import (
    "os"
)

// CreateTemp creates an empty temp file and returns its full path (empty on error).
func CreateTemp(prefix string) string {
    f, err := os.CreateTemp("", prefix)
    if err != nil {
        return ""
    }
    name := f.Name()
    _ = f.Close()
    return name
}

// WriteTemp writes data into a new temp file and returns its full path (empty on error).
func WriteTemp(prefix, data string) string {
    f, err := os.CreateTemp("", prefix)
    if err != nil {
        return ""
    }
    if _, err := f.Write([]byte(data)); err != nil {
        _ = f.Close()
        _ = os.Remove(f.Name())
        return ""
    }
    if err := f.Close(); err != nil {
        _ = os.Remove(f.Name())
        return ""
    }
    return f.Name()
}

// Remove deletes the given path. Returns 0 on success, 1 on failure.
func Remove(path string) int {
    if err := os.Remove(path); err != nil {
        return 1
    }
    return 0
}
