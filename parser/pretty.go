package parser

import "strings"

// PrettyPrint produces a simple human-readable representation of the AST.
func PrettyPrint(p *Program) string {
    if p == nil {
        return "<nil>"
    }
    var sb strings.Builder
    for _, f := range p.Functions {
        sb.WriteString("fn ")
        sb.WriteString(f.Name)
        sb.WriteString("() -> ")
        if f.ReturnType == "" {
            sb.WriteString("int")
        } else {
            sb.WriteString(f.ReturnType)
        }
        sb.WriteString(" {\n")
        if f.Body != nil {
            for _, s := range f.Body.Statements {
                switch st := s.(type) {
                case *ReturnStatement:
                    sb.WriteString("  return ...\n")
                case *VarStatement:
                    sb.WriteString("  var ")
                    sb.WriteString(st.Name)
                    sb.WriteString(" : ")
                    sb.WriteString(st.Type)
                    sb.WriteString("\n")
                default:
                    sb.WriteString("  stmt\n")
                }
            }
        }
        sb.WriteString("}\n")
    }
    return sb.String()
}
