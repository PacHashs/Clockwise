package compiler

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"codeberg.org/clockwise-lang/clockwise/checker"
	"codeberg.org/clockwise-lang/clockwise/codegen"
	"codeberg.org/clockwise-lang/clockwise/lexer"
	"codeberg.org/clockwise-lang/clockwise/parser"
)

type Compiler struct {
	// Configuration
	InputFile  string
	OutputFile string
	Verbose    bool
	Optimize   bool

	// Internal state
	errors   []error
	warnings []string
}

// NewCompiler creates a new instance of the ClockWise compiler
func NewCompiler(inputFile, outputFile string) *Compiler {
	return &Compiler{
		InputFile:  inputFile,
		OutputFile: outputFile,
		Verbose:    false,
		Optimize:   true,
	}
}

// Compile compiles the input file to the output file
func (c *Compiler) Compile() error {
	// 1. Read input file
	src, err := ioutil.ReadFile(c.InputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// 2. Lexical analysis
	l := lexer.New(string(src))
	tokens, lexErrs := l.Lex()
	if len(lexErrs) > 0 {
		c.addLexerErrors(lexErrs)
		return c.compileError()
	}

	// 3. Parsing
	p := parser.New(tokens)
	ast, parseErrs := p.Parse()
	if len(parseErrs) > 0 {
		c.addParserErrors(parseErrs)
		return c.compileError()
	}

	// 4. Semantic analysis
	chk := checker.New()
	semErrs := chk.Check(ast)
	if len(semErrs) > 0 {
		c.addSemanticErrors(semErrs)
		return c.compileError()
	}

	// 5. Code generation
	gen := codegen.New()
	goCode, genErrs := gen.Generate(ast)
	if len(genErrs) > 0 {
		c.addGenerationErrors(genErrs)
		return c.compileError()
	}

	// 6. Format the generated Go code
	formatted, err := format.Source([]byte(goCode))
	if err != nil {
		return fmt.Errorf("failed to format generated code: %w", err)
	}

	// 7. Write output file
	if err := os.MkdirAll(filepath.Dir(c.OutputFile), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	if err := ioutil.WriteFile(c.OutputFile, formatted, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

// addLexerErrors adds lexer errors to the compiler's error list
func (c *Compiler) addLexerErrors(errs []lexer.Error) {
	for _, err := range errs {
		c.errors = append(c.errors, fmt.Errorf("lexer error at %d:%d: %s",
			err.Line, err.Column, err.Msg))
	}
}

// addParserErrors adds parser errors to the compiler's error list
func (c *Compiler) addParserErrors(errs []parser.Error) {
	for _, err := range errs {
		c.errors = append(c.errors, fmt.Errorf("parse error at %d:%d: %s",
			err.Token.Line, err.Token.Column, err.Msg))
	}
}

// addSemanticErrors adds semantic analysis errors to the compiler's error list
func (c *Compiler) addSemanticErrors(errs []checker.Error) {
	for _, err := range errs {
		c.errors = append(c.errors, fmt.Errorf("semantic error at %d:%d: %s",
			err.Node.Pos().Line, err.Node.Pos().Column, err.Msg))
	}
}

// addGenerationErrors adds code generation errors to the compiler's error list
func (c *Compiler) addGenerationErrors(errs []error) {
	for _, err := range errs {
		c.errors = append(c.errors, fmt.Errorf("code generation error: %w", err))
	}
}

// compileError returns a formatted compilation error with all accumulated errors
func (c *Compiler) compileError() error {
	var sb strings.Builder
	sb.WriteString("compilation failed with ")
	sb.WriteString(fmt.Sprintf("%d errors", len(c.errors)))
	if len(c.warnings) > 0 {
		sb.WriteString(fmt.Sprintf(" and %d warnings", len(c.warnings)))
	}
	sb.WriteString(":\n")

	for i, err := range c.errors {
		sb.WriteString(fmt.Sprintf("  %d. %s\n", i+1, err.Error()))
	}

	for i, warn := range c.warnings {
		sb.WriteString(fmt.Sprintf("  warning %d: %s\n", i+1, warn))
	}

	return fmt.Errorf(sb.String())
}
