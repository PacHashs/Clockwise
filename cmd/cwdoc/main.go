package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Define command-line flags
	dir := flag.String("dir", ".", "Directory containing .cw files to document")
	output := flag.String("output", "docs", "Output directory for generated documentation")
	format := flag.String("format", "markdown", "Output format (markdown, html, or json)")
	help := flag.Bool("help", false, "Show help message")

	flag.Parse()

	if *help {
		showHelp()
		return
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(*output, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	// Find all .cw files in the directory
	files, err := filepath.Glob(filepath.Join(*dir, "*.cw"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding .cw files: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Println("No .cw files found in", *dir)
		return
	}

	// Process each file
	for _, file := range files {
		if err := processFile(file, *output, *format); err != nil {
			fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", file, err)
		}
	}

	fmt.Printf("Documentation generated in %s\n", *output)
}

func processFile(inputFile, outputDir, format string) error {
	// Read the source file
	content, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return err
	}

	// Extract documentation and code structure
	docs := extractDocumentation(string(content))

	// Generate output based on format
	outputFile := filepath.Join(outputDir, strings.TrimSuffix(filepath.Base(inputFile), ".cw")+".md")
	return ioutil.WriteFile(outputFile, []byte(generateMarkdown(docs)), 0644)
}

type Documentation struct {
	FilePath    string
	PackageName string
	Imports     []string
	Functions   []FunctionDoc
	Types       []TypeDoc
	Constants   []ConstantDoc
}

type FunctionDoc struct {
	Name       string
	Params     []Param
	ReturnType string
	Comment    string
}

type TypeDoc struct {
	Name    string
	Fields  []Field
	Comment string
}

type Field struct {
	Name string
	Type string
}

type ConstantDoc struct {
	Name  string
	Value string
}

type Param struct {
	Name string
	Type string
}

func extractDocumentation(content string) Documentation {
	// This is a simplified implementation
	// A real implementation would parse the source code properly
	var doc Documentation
	lines := strings.Split(content, "\n")

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "//") {
			// Handle comments
		} else if strings.HasPrefix(line, "func ") {
			// Parse function
			if i > 0 && strings.HasPrefix(lines[i-1], "//") {
				// Found a function with a comment
			}
		} else if strings.HasPrefix(line, "type ") && strings.Contains(line, "struct") {
			// Parse struct
		} else if strings.HasPrefix(line, "const ") || strings.HasPrefix(line, "let ") {
			// Parse constant
		}
	}

	return doc
}

func generateMarkdown(doc Documentation) string {
	var builder strings.Builder

	// Title
	builder.WriteString(fmt.Sprintf("# %s\n\n", doc.PackageName))

	// Functions
	if len(doc.Functions) > 0 {
		builder.WriteString("## Functions\n\n")
		for _, fn := range doc.Functions {
			builder.WriteString(fmt.Sprintf("### `%s`\n\n", fn.Name))
			if fn.Comment != "" {
				builder.WriteString(fmt.Sprintf("%s\n\n", fn.Comment))
			}
			builder.WriteString("**Parameters:**\n\n")
			for _, param := range fn.Params {
				builder.WriteString(fmt.Sprintf("- `%s %s`\n", param.Name, param.Type))
			}
			if fn.ReturnType != "" {
				builder.WriteString(fmt.Sprintf("\n**Returns:** `%s`\n\n", fn.ReturnType))
			}
		}
	}

	// Types
	if len(doc.Types) > 0 {
		builder.WriteString("## Types\n\n")
		for _, t := range doc.Types {
			builder.WriteString(fmt.Sprintf("### `%s`\n\n", t.Name))
			if t.Comment != "" {
				builder.WriteString(fmt.Sprintf("%s\n\n", t.Comment))
			}
			if len(t.Fields) > 0 {
				builder.WriteString("**Fields:**\n\n")
				for _, field := range t.Fields {
					builder.WriteString(fmt.Sprintf("- `%s %s`\n", field.Name, field.Type))
				}
			}
			builder.WriteString("\n")
		}
	}

	// Constants
	if len(doc.Constants) > 0 {
		builder.WriteString("## Constants\n\n")
		for _, c := range doc.Constants {
			builder.WriteString(fmt.Sprintf("- `%s = %s`\n", c.Name, c.Value))
		}
	}

	return builder.String()
}

func showHelp() {
	helpText := `cwdoc - ClockWise Documentation Generator

Usage:
  cwdoc [flags]

Flags:
  -dir string
        Directory containing .cw files to document (default ".")
  -output string
        Output directory for generated documentation (default "docs")
  -format string
        Output format (markdown, html, or json) (default "markdown")
  -help
        Show this help message
`
	fmt.Print(helpText)
}
