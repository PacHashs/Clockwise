package runtimelib

import "regexp"

// RegexMatch returns true if pattern matches s
func RegexMatch(pattern, s string) bool {
    re, err := regexp.Compile(pattern)
    if err != nil {
        return false
    }
    return re.MatchString(s)
}

// RegexReplaceAll replaces all occurrences of pattern in s with repl
func RegexReplaceAll(pattern, s, repl string) string {
    re, err := regexp.Compile(pattern)
    if err != nil {
        return s
    }
    return re.ReplaceAllString(s, repl)
}
