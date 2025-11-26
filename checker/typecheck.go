package checker

import (
    "fmt"
    "codeberg.org/clockwise-lang/clockwise/parser"
)

// RunChecks performs additional checker passes: symbol collection, simple
// sanity checks and lint hooks. It delegates to the main type-checker and
// aggregates diagnostics.
func RunChecks(p *parser.Program) error {
    if p == nil {
        return fmt.Errorf("nil program")
    }
    // First, run the existing strict type checker in this package (if present).
    if err := CheckProgram(p); err != nil {
        return err
    }

    diags := Diagnostics{}
    st := NewSymbolTable()
    // register functions and gather variable info
    for _, f := range p.Functions {
        if _, exists := st.Funcs[f.Name]; exists {
            diags = append(diags, Diagnostic{Msg: fmt.Sprintf("duplicate function %s", f.Name), Node: f})
        }
        st.RegisterFunc(f)
        if f.Body != nil {
            for _, s := range f.Body.Statements {
                if vs, ok := s.(*parser.VarStatement); ok {
                    if vs.Type != "" && vs.Type != "int" && vs.Type != "string" {
                        diags = append(diags, Diagnostic{Msg: fmt.Sprintf("unsupported var type '%s' for %s", vs.Type, vs.Name), Node: f})
                    }
                    st.RegisterVar(f.Name, vs.Name, vs.Type)
                }
            }
        }
    }
    if _, ok := st.LookupFunc("main"); !ok {
        diags = append(diags, Diagnostic{Msg: "no main function found"})
    }
    // Print lint warnings (non-fatal)
    for _, w := range LintProgram(p) {
        diags = append(diags, Diagnostic{Msg: "lint: " + w})
    }
    return ReportError(diags)
}
