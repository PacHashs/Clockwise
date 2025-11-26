package codegen

import (
    "fmt"
    "codeberg.org/clockwise-lang/clockwise/parser"
)

// GenFunction generates Go source for a single parser.Function.
func GenFunction(f *parser.Function) string {
    if f == nil {
        return ""
    }
    ret := "int"
    if f.ReturnType != "" {
        ret = f.ReturnType
    }
    var sb string
    sb += fmt.Sprintf("func %s() %s {\n", f.Name, ret)
    if f.Body != nil {
        for _, s := range f.Body.Statements {
            switch st := s.(type) {
            case *parser.ReturnStatement:
                sb += "    return 0\n"
            case *parser.VarStatement:
                sb += fmt.Sprintf("    var %s %s = 0\n", st.Name, mapType(st.Type))
            case *parser.ExpressionStatement:
                sb += "    // expression\n"
            default:
                sb += "    // stmt\n"
            }
        }
    }
    sb += "}\n\n"
    return sb
}
