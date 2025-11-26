package checker

import (
    "fmt"
    "codeberg.org/clockwise-lang/clockwise/parser"
)

// Very small type system: only 'int' and 'string' for now
type Type string

const (
    TInt    Type = "int"
    TString Type = "string"
)

func typeFromIdent(s string) (Type, error) {
    switch s {
    case "int":
        return TInt, nil
    case "string":
        return TString, nil
    default:
        return "", fmt.Errorf("unknown type: %s", s)
    }
}

// CheckProgram runs basic type checking and returns error if any
func CheckProgram(p *parser.Program) error {
    for _, fn := range p.Functions {
        if err := checkFunction(fn); err != nil {
            return fmt.Errorf("in function %s: %w", fn.Name, err)
        }
    }
    return nil
}

func checkFunction(fn *parser.Function) error {
    // currently only verify return statements match declared return type
    retT, err := typeFromIdent(fn.ReturnType)
    if err != nil {
        return err
    }
    for _, s := range fn.Body.Statements {
        switch st := s.(type) {
        case *parser.ReturnStatement:
            t, err := inferExprType(st.Value)
            if err != nil {
                return err
            }
            if t != retT {
                return fmt.Errorf("return type mismatch: want %s, got %s", retT, t)
            }
        case *parser.VarStatement:
            declared, err := typeFromIdent(st.Type)
            if err != nil {
                return err
            }
            t, err := inferExprType(st.Value)
            if err != nil {
                return err
            }
            if t != declared {
                return fmt.Errorf("variable %s: type mismatch: declared %s but assigned %s", st.Name, declared, t)
            }
        }
    }
    return nil
}

func inferExprType(e parser.Expression) (Type, error) {
    switch ex := e.(type) {
    case *parser.IntegerLiteral:
        return TInt, nil
    case *parser.StringLiteral:
        return TString, nil
    case *parser.Identifier:
        // identifiers assumed to be functions or variables; for now assume int
        return TInt, nil
    case *parser.InfixExpression:
        lt, err := inferExprType(ex.Left)
        if err != nil {
            return "", err
        }
        rt, err := inferExprType(ex.Right)
        if err != nil {
            return "", err
        }
        if lt != rt {
            return "", fmt.Errorf("type mismatch in infix: left %s right %s", lt, rt)
        }
        // arithmetic operators produce int
        return lt, nil
    case *parser.CallExpression:
        // special-case 'print' returns int 0 for now
        if id, ok := ex.Function.(*parser.Identifier); ok {
            if id.Value == "print" {
                return TInt, nil
            }
        }
        return TInt, nil
    default:
        return "", fmt.Errorf("unknown expression type")
    }
}
