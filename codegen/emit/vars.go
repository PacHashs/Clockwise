package codegen

import "fmt"

// EmitVar returns Go source for a variable declaration: `var name type = value`.
func EmitVar(name, typ, value string) string {
    if typ == "" {
        return fmt.Sprintf("var %s = %s\n", name, value)
    }
    return fmt.Sprintf("var %s %s = %s\n", name, typ, value)
}
