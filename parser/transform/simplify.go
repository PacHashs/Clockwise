package transform

import pr "codeberg.org/clockwise-lang/clockwise/parser"

// Simplify is a tiny example transform that could be used to lower complex
// constructs to simpler ones. For now it is a noop but provides a hook for
// future passes.
func Simplify(prog *pr.Program) *pr.Program {
    return prog
}
