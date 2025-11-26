package runtimelib

import "os"

// GetEnv returns the value of the environment variable named by the key.
// If the variable is not present the empty string is returned.
func GetEnv(key string) string {
    return os.Getenv(key)
}

// SetEnv sets the value of the environment variable named by the key.
// Returns nil on success, or an error on failure.
func SetEnv(key, value string) error {
    return os.Setenv(key, value)
}
