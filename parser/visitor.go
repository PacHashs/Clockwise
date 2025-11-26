package parser

// Visitor allows walking the AST nodes.
type Visitor interface {
    VisitNode(n Node)
}

// Walk performs a simple preorder traversal calling v.VisitNode on each Node.
func Walk(n Node, v Visitor) {
    if n == nil || v == nil {
        return
    }
    v.VisitNode(n)
    for _, c := range n.Children() {
        Walk(c, v)
    }
}

