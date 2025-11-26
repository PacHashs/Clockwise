package runtimelib

import (
    "io"
    "os"
)

// CopyFile copies source to destination, returns error on failure
func CopyFile(src, dst string) error {
    in, err := os.Open(src)
    if err != nil {
        return err
    }
    defer in.Close()
    out, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer out.Close()
    _, err = io.Copy(out, in)
    if err != nil {
        return err
    }
    return out.Sync()
}

// FileSize returns file size in bytes or -1 on error
func FileSize(path string) int64 {
    fi, err := os.Stat(path)
    if err != nil {
        return -1
    }
    return fi.Size()
}
