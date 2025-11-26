package runtimelib

import "path/filepath"

// JoinPaths joins multiple path elements using the OS separator
func JoinPaths(parts ...string) string {
    return filepath.Join(parts...)
}

// BaseName returns the last element of path
func BaseName(p string) string {
    return filepath.Base(p)
}

// DirName returns all but the last element of path
func DirName(p string) string {
    return filepath.Dir(p)
}
