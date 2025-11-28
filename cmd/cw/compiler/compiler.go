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
	InputFiles []string
	OutputFile string
	Verbose    bool
	Optimize   bool

	// Internal state
	errors   []error
	warnings []string
}

// NewCompiler creates a new instance of the ClockWise compiler
func NewCompiler(inputFiles []string, outputFile string) *Compiler {
	return &Compiler{
		InputFiles: inputFiles,
		OutputFile: outputFile,
		Verbose:    false,
		Optimize:   true,
	}
}

// NewSingleFileCompiler creates a compiler for a single input file (backward compatibility)
func NewSingleFileCompiler(inputFile, outputFile string) *Compiler {
	return NewCompiler([]string{inputFile}, outputFile)
}

// Compile compiles the input files to the output file
func (c *Compiler) Compile() error {
	if len(c.InputFiles) == 0 {
		return fmt.Errorf("no input files specified")
	}

	// Create a module to hold all files
	module := parser.NewModule("main")

	// Process each input file
	for _, inputFile := range c.InputFiles {
		if c.Verbose {
			fmt.Printf("Processing file: %s\n", inputFile)
		}

		// 1. Read input file
		src, err := ioutil.ReadFile(inputFile)
		if err != nil {
			return fmt.Errorf("failed to read input file %s: %w", inputFile, err)
		}

		// 2. Lexical analysis
		l := lexer.New(string(src))
		tokens := l.Tokenize()

		// 3. Parsing
		p := parser.New(tokens)
		program, err := p.ParseProgram()
		if err != nil {
			return fmt.Errorf("parse error in %s: %w", inputFile, err)
		}

		// 4. Add file to module
		module.AddFile(program)
	}

	// 5. Resolve imports and create unified program
	unifiedProgram, err := c.resolveImports(module)
	if err != nil {
		return fmt.Errorf("import resolution failed: %w", err)
	}

	// 6. Semantic analysis
	if err := checker.CheckProgram(unifiedProgram); err != nil {
		return fmt.Errorf("semantic error: %w", err)
	}

	// 7. Code generation
	goCode := codegen.Generate(unifiedProgram)

	// 8. Format the generated Go code
	formatted, err := format.Source([]byte(goCode))
	if err != nil {
		return fmt.Errorf("failed to format generated code: %w", err)
	}

	// 9. Write output file
	if err := os.MkdirAll(filepath.Dir(c.OutputFile), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	if err := ioutil.WriteFile(c.OutputFile, formatted, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	if c.Verbose {
		fmt.Printf("Successfully compiled %d files to %s\n", len(c.InputFiles), c.OutputFile)
	}

	return nil
}

// resolveImports resolves imports across all files in the module and creates a unified program
func (c *Compiler) resolveImports(module *parser.Module) (*parser.Program, error) {
	// Start with an empty program
	unifiedProgram := &parser.Program{
		Imports:   []string{},
		Functions: []*parser.Function{},
	}
	
	// Track all function names to detect duplicates
	functionNames := make(map[string]bool)
	
	// Process each file in the module
	for _, fileNode := range module.Files {
		program, ok := fileNode.(*parser.Program)
		if !ok {
			return nil, fmt.Errorf("invalid file node type in module")
		}
		
		// Check for duplicate function names
		for _, fn := range program.Functions {
			if functionNames[fn.Name] {
				return nil, fmt.Errorf("duplicate function name '%s' found in multiple files", fn.Name)
			}
			functionNames[fn.Name] = true
		}
		
		// Merge imports
		for _, imp := range program.Imports {
			// Check if import already exists
			found := false
			for _, existing := range unifiedProgram.Imports {
				if existing == imp {
					found = true
					break
				}
			}
			if !found {
				unifiedProgram.Imports = append(unifiedProgram.Imports, imp)
			}
		}
		
		// Merge functions
		unifiedProgram.Functions = append(unifiedProgram.Functions, program.Functions...)
	}
	
	return unifiedProgram, nil
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
