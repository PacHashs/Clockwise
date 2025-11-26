package parser

// CountNodes returns the total number of nodes in the subtree rooted at n.
func CountNodes(n Node) int {
    if n == nil {
        return 0
    }
    c := 1
    for _, ch := range n.Children() {
        c += CountNodes(ch)
    }
    return c
}
