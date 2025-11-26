package parser

// Module groups multiple AST files into a single compilation unit.
type Module struct {
    Name  string
    Files []Node
}

func NewModule(name string) *Module { return &Module{Name: name} }

func (m *Module) AddFile(n Node) { m.Files = append(m.Files, n) }

// TotalNodes returns total node count across all files in the module.
func (m *Module) TotalNodes() int {
    total := 0
    for _, f := range m.Files {
        total += CountNodes(f)
    }
    return total
}
