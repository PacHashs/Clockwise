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
	tokens := l.Tokenize()

	// 3. Parsing
	p := parser.New(tokens)
	program, err := p.ParseProgram()
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	// 4. Semantic analysis
	if err := checker.CheckProgram(program); err != nil {
		return fmt.Errorf("semantic error: %w", err)
	}

	// 5. Code generation
	goCode := codegen.Generate(program)

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
// legacy error aggregation helpers kept for compatibility but unused
func (c *Compiler) addLexerErrors(_ interface{})           {}
func (c *Compiler) addParserErrors(_ interface{})          {}
func (c *Compiler) addSemanticErrors(_ interface{})        {}
func (c *Compiler) addGenerationErrors(_ interface{})      {}

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
