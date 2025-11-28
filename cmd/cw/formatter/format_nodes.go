package formatter

import (
	"codeberg.org/clockwise-lang/clockwise/ast"
	"fmt"
)

func (f *Formatter) formatExprStmt(s *ast.ExprStmt) error {
	return f.formatNode(s.X)
}

func (f *Formatter) formatIfStmt(s *ast.IfStmt) error {
	f.printf("if ")
	if err := f.formatNode(s.Cond); err != nil {
		return err
	}
	f.printf(" ")

	// Format the "then" block
	if err := f.formatNode(s.Body); err != nil {
		return err
	}

	// Format the "else" clause if it exists
	if s.Else != nil {
		f.printf(" else ")
		switch stmt := s.Else.(type) {
		case *ast.IfStmt:
			return f.formatIfStmt(stmt)
		case *ast.BlockStmt:
			return f.formatNode(stmt)
		default:
			return fmt.Errorf("unexpected else statement type: %T", s.Else)
		}
	}

	return nil
}

func (f *Formatter) formatForStmt(s *ast.ForStmt) error {
	f.printf("for ")

	// Format the initialization
	if s.Init != nil {
		if err := f.formatNode(s.Init); err != nil {
			return err
		}
	}

	f.printf("; ")

	// Format the condition
	if s.Cond != nil {
		if err := f.formatNode(s.Cond); err != nil {
			return err
		}
	}

	f.printf("; ")

	// Format the post statement
	if s.Post != nil {
		if err := f.formatNode(s.Post); err != nil {
			return err
		}
	}

	f.printf(" ")

	// Format the loop body
	return f.formatNode(s.Body)
}

func (f *Formatter) formatReturnStmt(s *ast.ReturnStmt) error {
	f.printf("return")
	if s.Result != nil {
		f.printf(" ")
		return f.formatNode(s.Result)
	}
	return nil
}

func (f *Formatter) formatVarDecl(decl *ast.VarDecl) error {
	f.printf("var %s", decl.Name.Name)
	if decl.Type != nil {
		f.printf(" %s", decl.Type)
	}
	if decl.Value != nil {
		f.printf(" = ")
		return f.formatNode(decl.Value)
	}
	return nil
}

func (f *Formatter) formatAssignStmt(s *ast.AssignStmt) error {
	// Format left-hand side
	for i, lhs := range s.Lhs {
		if i > 0 {
			f.printf(", ")
		}
		if err := f.formatNode(lhs); err != nil {
			return err
		}
	}

	// Format operator
	f.printf(" %s ", s.Tok.String())

	// Format right-hand side
	for i, rhs := range s.Rhs {
		if i > 0 {
			f.printf(", ")
		}
		if err := f.formatNode(rhs); err != nil {
			return err
		}
	}

	return nil
}

func (f *Formatter) formatBinaryExpr(expr *ast.BinaryExpr) error {
	if err := f.formatNode(expr.X); err != nil {
		return err
	}

	f.printf(" %s ", expr.Op.String())

	return f.formatNode(expr.Y)
}

func (f *Formatter) formatIdent(id *ast.Ident) error {
	f.printf("%s", id.Name)
	return nil
}

func (f *Formatter) formatBasicLit(lit *ast.BasicLit) error {
	f.printf("%s", lit.Value)
	return nil
}
