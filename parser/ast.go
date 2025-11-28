package parser

// AST definitions for a minimal Clockwise language

type Program struct {
    Functions []*Function
    Imports   []string
}

func (p *Program) Children() []Node {
    out := make([]Node, 0, len(p.Functions))
    for _, fn := range p.Functions {
        out = append(out, fn)
    }
    return out
}

// Node is a generic AST node used by helpers and visitors.
type Node interface{
    Children() []Node
}

type Function struct {
    Name string
    ReturnType string
    Body *BlockStatement
}

func (f *Function) Children() []Node {
    if f.Body == nil {
        return nil
    }
    return []Node{f.Body}
}

type Statement interface{}
type Expression interface{}

type BlockStatement struct {
    Statements []Statement
}

func (b *BlockStatement) Children() []Node {
    out := make([]Node, 0, len(b.Statements))
    for _, s := range b.Statements {
        if n, ok := s.(Node); ok {
            out = append(out, n)
        }
    }
    return out
}

type ReturnStatement struct {
    Value Expression
}

func (r *ReturnStatement) Children() []Node {
    if n, ok := r.Value.(Node); ok {
        return []Node{n}
    }
    return nil
}

type ExpressionStatement struct {
    Expr Expression
}

func (e *ExpressionStatement) Children() []Node {
    if n, ok := e.Expr.(Node); ok {
        return []Node{n}
    }
    return nil
}

type VarStatement struct {
    Name  string
    Type  string
    Value Expression
}

func (v *VarStatement) Children() []Node {
    if n, ok := v.Value.(Node); ok {
        return []Node{n}
    }
    return nil
}

type IntegerLiteral struct {
    Value string
}

func (i *IntegerLiteral) Children() []Node { return nil }

type StringLiteral struct {
    Value string
}

func (s *StringLiteral) Children() []Node { return nil }

type Identifier struct {
    Value string
}

func (id *Identifier) Children() []Node { return nil }

type CallExpression struct {
    Function Expression
    Args     []Expression
}

func (c *CallExpression) Children() []Node {
    out := []Node{}
    if n, ok := c.Function.(Node); ok {
        out = append(out, n)
    }
    for _, a := range c.Args {
        if n, ok := a.(Node); ok {
            out = append(out, n)
        }
    }
    return out
}

type InfixExpression struct {
    Left     Expression
    Operator string
    Right    Expression
}

func (in *InfixExpression) Children() []Node {
    out := []Node{}
    if n, ok := in.Left.(Node); ok {
        out = append(out, n)
    }
    if n, ok := in.Right.(Node); ok {
        out = append(out, n)
    }
    return out
}
