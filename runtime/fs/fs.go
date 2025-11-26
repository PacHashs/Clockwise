package runtimelib

import (
    "os"
)

// ReadFile returns file contents as string
func ReadFile(path string) (string, error) {
    b, err := os.ReadFile(path)
    if err != nil {
        return "", err
    }
    return string(b), nil
}

// WriteFile writes content to path
func WriteFile(path, content string) error {
    return os.WriteFile(path, []byte(content), 0644)
}

// Exists checks file existence
func Exists(path string) bool {
    _, err := os.Stat(path)
    return err == nil
}
