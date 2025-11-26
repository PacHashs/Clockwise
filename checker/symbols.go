package checker

import "codeberg.org/clockwise-lang/clockwise/parser"

// SymbolTable tracks top-level functions and simple variable types per function
type SymbolTable struct {
    Funcs map[string]*parser.Function
    Vars  map[string]map[string]string // func -> var name -> type
}

func NewSymbolTable() *SymbolTable {
    return &SymbolTable{Funcs: map[string]*parser.Function{}, Vars: map[string]map[string]string{}}
}

func (s *SymbolTable) RegisterFunc(f *parser.Function) {
    s.Funcs[f.Name] = f
    if _, ok := s.Vars[f.Name]; !ok {
        s.Vars[f.Name] = map[string]string{}
    }
}

func (s *SymbolTable) RegisterVar(funcName, name, typ string) {
    if _, ok := s.Vars[funcName]; !ok {
        s.Vars[funcName] = map[string]string{}
    }
    s.Vars[funcName][name] = typ
}

func (s *SymbolTable) LookupFunc(name string) (*parser.Function, bool) {
    f, ok := s.Funcs[name]
    return f, ok
}

func (s *SymbolTable) LookupVar(funcName, name string) (string, bool) {
    if vars, ok := s.Vars[funcName]; ok {
        t, ok := vars[name]
        return t, ok
    }
    return "", false
}
