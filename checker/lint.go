package checker

import (
    "fmt"
    "codeberg.org/clockwise-lang/clockwise/parser"
)

// LintProgram performs small lint checks and returns a list of warnings.
func LintProgram(p *parser.Program) []string {
    var out []string
    if p == nil {
        return out
    }
    for _, f := range p.Functions {
        if f.Name == "main" && f.ReturnType != "int" {
            out = append(out, "main should return int")
        }
    }
    return out
}

// PrintLint prints lint warnings to stdout (used by CLI tools optionally).
func PrintLint(p *parser.Program) {
    for _, w := range LintProgram(p) {
        fmt.Println("lint:", w)
    }
}
