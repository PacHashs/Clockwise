package parser

// SymbolTable is a minimal symbol table tracking functions and variables.
type SymbolTable struct {
    Funcs map[string]*Function
    Vars  map[string]string // var name -> type
}

func NewSymbolTable() *SymbolTable {
    return &SymbolTable{Funcs: map[string]*Function{}, Vars: map[string]string{}}
}

func (s *SymbolTable) RegisterFunc(f *Function) {
    s.Funcs[f.Name] = f
}

func (s *SymbolTable) LookupFunc(name string) (*Function, bool) {
    f, ok := s.Funcs[name]
    return f, ok
}

func (s *SymbolTable) RegisterVar(name, typ string) {
    s.Vars[name] = typ
}

func (s *SymbolTable) LookupVar(name string) (string, bool) {
    t, ok := s.Vars[name]
    return t, ok
}
