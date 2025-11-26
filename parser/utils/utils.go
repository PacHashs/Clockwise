package parser

// IsExported returns true when the name looks exported (starts with uppercase).
func IsExported(name string) bool {
    if name == "" {
        return false
    }
    c := name[0]
    return c >= 'A' && c <= 'Z'
}
