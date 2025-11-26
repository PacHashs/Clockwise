package codegen

// OptPass is a minimal optimization pass placeholder. In a fuller compiler
// this would perform dead-code removal, constant folding, etc.
func OptPass(src string) string {
    // noop for now
    return src
}
