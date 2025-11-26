package checks

import (
    "fmt"
    pr "codeberg.org/clockwise-lang/clockwise/parser"
)

// ValidateProgram runs a few basic sanity checks on the parsed program.
// Returns an error describing the first problem found.
func ValidateProgram(p *pr.Program) error {
    if p == nil {
        return fmt.Errorf("nil program")
    }
    // ensure there is at least one function named main
    found := false
    for _, f := range p.Functions {
        if f.Name == "main" {
            found = true
            break
        }
    }
    if !found {
        return fmt.Errorf("no main function found")
    }
    return nil
}
