package runtimelib

import "regexp"

var emailRe = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

// IsEmail returns true when s looks like a simple email address.
func IsEmail(s string) bool {
    return emailRe.MatchString(s)
}
