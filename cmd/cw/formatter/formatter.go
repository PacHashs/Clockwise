package formatter

import (
	"bytes"
	"fmt"
	"go/token"
	"io"
	"strings"

	"codeberg.org/clockwise-lang/clockwise/parser"
	"codeberg.org/clockwise-lang/clockwise/lexer"
	"codeberg.org/clockwise-lang/clockwise/parser"
)

type Formatter struct {
	src       []byte
	fset      *token.FileSet
	node      ast.Node
	indent    int
	lineStart bool
	output    *bytes.Buffer
}

func Format(src []byte) ([]byte, error) {
	// Create a new lexer and parse the source
	l := lexer.New(string(src))
	p := parser.New(l)
	file, err := p.ParseFile()
	if err != nil {
		return nil, fmt.Errorf("parsing error: %v", err)
	}

	f := &Formatter{
		src:       src,
		fset:      token.NewFileSet(),
		node:      file,
		indent:    0,
		lineStart: true,
		output:    &bytes.Buffer{},
	}

	// Format the AST
	if err := f.formatNode(file); err != nil {
		return nil, err
	}

	return f.output.Bytes(), nil
}

func (f *Formatter) formatNode(node ast.Node) error {
	switch n := node.(type) {
	case *ast.File:
		return f.formatFile(n)
	case *ast.FuncDecl:
		return f.formatFuncDecl(n)
	case *ast.BlockStmt:
		return f.formatBlockStmt(n)
	case *ast.ExprStmt:
		return f.formatExprStmt(n)
	case *ast.IfStmt:
		return f.formatIfStmt(n)
	case *ast.ForStmt:
		return f.formatForStmt(n)
	case *ast.ReturnStmt:
		return f.formatReturnStmt(n)
	case *ast.VarDecl:
		return f.formatVarDecl(n)
	case *ast.AssignStmt:
		return f.formatAssignStmt(n)
	case *ast.BinaryExpr:
		return f.formatBinaryExpr(n)
	case *ast.Ident:
		return f.formatIdent(n)
	case *ast.BasicLit:
		return f.formatBasicLit(n)
	default:
		return fmt.Errorf("unsupported node type: %T", node)
	}
}

func (f *Formatter) formatFile(file *ast.File) error {
	// Format package declaration
	f.printf("package %s\n\n", file.Name.Name)

	// Format imports
	if len(file.Imports) > 0 {
		f.printf("import (\n")
		f.indent++
		for _, imp := range file.Imports {
			f.printf("\"%s\"\n", imp.Path)
		}
		f.indent--
		f.printf(")\n\n")
	}

	// Format declarations
	for _, decl := range file.Decls {
		if err := f.formatNode(decl); err != nil {
			return err
		}
		f.printf("\n")
	}

	return nil
}

func (f *Formatter) formatFuncDecl(fn *ast.FuncDecl) error {
	// Function signature
	f.printf("func %s(", fn.Name.Name)

	// Parameters
	for i, param := range fn.Type.Params.List {
		if i > 0 {
			f.printf(", ")
		}
		for j, name := range param.Names {
			if j > 0 {
				f.printf(", ")
			}
			f.printf("%s ", name.Name)
		}
		f.printf("%s", param.Type)
	}

	// Return type
	if fn.Type.Results != nil {
		f.printf(") ")
		if len(fn.Type.Results.List) > 1 {
			f.printf("(")
		}
		for i, result := range fn.Type.Results.List {
			if i > 0 {
				f.printf(", ")
			}
			f.printf("%s", result.Type)
		}
		if len(fn.Type.Results.List) > 1 {
			f.printf(")")
		}
	} else {
		f.printf(")")
	}

	// Function body
	return f.formatNode(fn.Body)
}

func (f *Formatter) formatBlockStmt(block *ast.BlockStmt) error {
	f.printf(" {\n")
	f.indent++

	for _, stmt := range block.List {
		f.printIndent()
		if err := f.formatNode(stmt); err != nil {
			return err
		}
		f.printf("\n")
	}

	f.indent--
	f.printIndent()
	f.printf("}")

	return nil
}

// Helper methods for formatting
func (f *Formatter) printf(format string, args ...interface{}) {
	if f.lineStart {
		f.printIndent()
	}
	fmt.Fprintf(f.output, format, args...)
}

func (f *Formatter) printIndent() {
	if f.lineStart {
		f.lineStart = false
		for i := 0; i < f.indent; i++ {
			f.output.WriteString("\t")
		}
	}
}

// ... (implement other format methods for different AST nodes)

// FormatFile formats a single file
func FormatFile(filename string) ([]byte, error) {
	src, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	res, err := Format(src)
	if err != nil {
		return nil, fmt.Errorf("formatting %s: %v", filename, err)
	}

	if bytes.Equal(src, res) {
		return nil, nil // No changes
	}

	return res, nil
}

// Process processes the input and writes the formatted output to w.
// The results of processing the source are the same as calling Format.
func Process(filename string, src []byte, w io.Writer) error {
	res, err := Format(src)
	if err != nil {
		return err
	}

	if src != nil && bytes.Equal(src, res) {
		return nil // No changes
	}

	_, err = w.Write(res)
	return err
}
