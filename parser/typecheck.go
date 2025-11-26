package parser

import "fmt"

// CheckTypes performs a very small type-check pass validating declared types
// used in variable declarations and function return types. It's intentionally
// conservative and only recognizes "int" and "string" for now.
func CheckTypes(p *Program) error {
    if p == nil {
        return fmt.Errorf("nil program")
    }
    for _, f := range p.Functions {
        if f.ReturnType != "" && f.ReturnType != "int" && f.ReturnType != "string" {
            return fmt.Errorf("unsupported return type '%s' in function %s", f.ReturnType, f.Name)
        }
        // walk body for var statements
        if f.Body != nil {
            for _, st := range f.Body.Statements {
                if vs, ok := st.(*VarStatement); ok {
                    if vs.Type != "" && vs.Type != "int" && vs.Type != "string" {
                        return fmt.Errorf("unsupported var type '%s' for %s in function %s", vs.Type, vs.Name, f.Name)
                    }
                }
            }
        }
    }
    return nil
}
