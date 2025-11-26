package parser

// Import represents a simple import directive in the source.
type Import struct {
    Path string
}

// AddImport appends an import path to the program's import list.
func AddImport(p *Program, path string) {
    if p == nil {
        return
    }
    for _, ex := range p.Imports {
        if ex == path {
            return
        }
    }
    p.Imports = append(p.Imports, path)
}
