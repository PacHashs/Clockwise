package runtimelib

import (
    "encoding/json"
    "fmt"
)

// Expr is a tiny JSON-friendly expression representation used by the optimizer.
// Supported shapes:
// {"kind":"int","int":123}
// {"kind":"string","str":"hello"}
// {"kind":"binary","op":"+","left":{...},"right":{...}}
type Expr struct {
    Kind  string `json:"kind"`
    Int   int    `json:"int,omitempty"`
    Str   string `json:"str,omitempty"`
    Op    string `json:"op,omitempty"`
    Left  *Expr  `json:"left,omitempty"`
    Right *Expr  `json:"right,omitempty"`
}

// OptimizeExpr takes a JSON-encoded Expr, performs simple constant folding
// (integer arithmetic and string concatenation for +), and returns the
// optimized JSON-encoded Expr or an error.
func OptimizeExpr(in string) (string, error) {
    var e Expr
    if err := json.Unmarshal([]byte(in), &e); err != nil {
        return "", fmt.Errorf("invalid expr json: %w", err)
    }
    out := fold(&e)
    b, err := json.Marshal(out)
    if err != nil {
        return "", err
    }
    return string(b), nil
}

func fold(e *Expr) *Expr {
    if e == nil {
        return nil
    }
    switch e.Kind {
    case "int", "string":
        return e
    case "binary":
        l := fold(e.Left)
        r := fold(e.Right)
        if l == nil || r == nil {
            return &Expr{Kind: "binary", Op: e.Op, Left: l, Right: r}
        }
        // integer arithmetic folding
        if l.Kind == "int" && r.Kind == "int" {
            switch e.Op {
            case "+":
                return &Expr{Kind: "int", Int: l.Int + r.Int}
            case "-":
                return &Expr{Kind: "int", Int: l.Int - r.Int}
            case "*":
                return &Expr{Kind: "int", Int: l.Int * r.Int}
            case "/":
                if r.Int == 0 {
                    // avoid divide-by-zero folding; leave as-is
                    return &Expr{Kind: "binary", Op: e.Op, Left: l, Right: r}
                }
                return &Expr{Kind: "int", Int: l.Int / r.Int}
            }
        }
        // string concatenation folding
        if e.Op == "+" && l.Kind == "string" && r.Kind == "string" {
            return &Expr{Kind: "string", Str: l.Str + r.Str}
        }
        // otherwise return folded children
        return &Expr{Kind: "binary", Op: e.Op, Left: l, Right: r}
    default:
        return e
    }
}
