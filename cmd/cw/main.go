package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	cwcompiler "codeberg.org/clockwise-lang/clockwise/cmd/cw/compiler"
	"codeberg.org/clockwise-lang/clockwise/checker"
	"codeberg.org/clockwise-lang/clockwise/codegen"
	"codeberg.org/clockwise-lang/clockwise/lexer"
	"codeberg.org/clockwise-lang/clockwise/parser"
)

const (
	appName        = "cwc"
	appDescription = `ClockWise Compiler - A pragmatic systems language

Usage:
  cwc [command] [flags]
  cwc [input.cw] [flags]
  cwc build [input.cw] [-o output]
  cwc run [input.cw] [args...]
  cwc fmt [input.cw]
  cwc --help | -h
  cwc --version | -v

Examples:
  cwc build program.cw -o program
  cwc run program.cw arg1 arg2
  cwc fmt program.cw
  cwc --help`
)

var (
	// These will be set during build with -ldflags
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Handle no arguments case
	if len(os.Args) == 1 {
		printUsage(nil)
		os.Exit(0)
	}

	// Handle subcommands
	switch os.Args[1] {
	case "build":
		buildCmd()
	case "run":
		runCmd()
	case "fmt":
		fmtCmd()
	case "--help", "-h":
		printUsage(nil)
		os.Exit(0)
	case "--version", "-v":
		printVersion()
		os.Exit(0)
	default:
		// If no command specified, default to build
		buildCmd()
	}
}

func buildCmd() {
	fs := flag.NewFlagSet("build", flag.ExitOnError)
	outputFile := fs.String("o", "", "Output file (default: input filename without extension)")
	verbose := fs.Bool("v", false, "Enable verbose output")
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: cwc build [input1.cw input2.cw ...] [-o output]\n")
		fmt.Fprintf(os.Stderr, "  If multiple input files are provided, they will be compiled together.\n")
		fs.PrintDefaults()
	}

	// Parse flags and args
	if err := fs.Parse(os.Args[2:]); err != nil {
		log.Fatalf("Error parsing flags: %v", err)
	}

	// Get input files
	args := fs.Args()
	if len(args) == 0 {
		fs.Usage()
		os.Exit(1)
	}
	inputFiles := args

	// Set default output file if not specified
	if *outputFile == "" {
		if len(inputFiles) == 1 {
			*outputFile = strings.TrimSuffix(inputFiles[0], filepath.Ext(inputFiles[0]))
		} else {
			*outputFile = "program" // default for multi-file compilation
		}
	}

	// Run the compiler
	comp := cwcompiler.NewCompiler(inputFiles, *outputFile)
	comp.Verbose = *verbose

	if err := comp.Compile(); err != nil {
		log.Fatalf("Compilation failed: %v", err)
	}
}

func runCmd() {
	fs := flag.NewFlagSet("run", flag.ExitOnError)
	verbose := fs.Bool("v", false, "Enable verbose output")
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: cwc run [input1.cw input2.cw ...] [-- program args...]\n")
		fmt.Fprintf(os.Stderr, "  If multiple input files are provided, they will be compiled together.\n")
		fmt.Fprintf(os.Stderr, "  Use '--' to separate Clockwise files from program arguments.\n")
		fs.PrintDefaults()
	}

	// Parse flags and args
	if err := fs.Parse(os.Args[2:]); err != nil {
		log.Fatalf("Error parsing flags: %v", err)
	}

	args := fs.Args()
	if len(args) == 0 {
		fs.Usage()
		os.Exit(1)
	}

	// Separate input files from program arguments
	var inputFiles []string
	var programArgs []string
	seenDoubleDash := false
	
	for _, arg := range args {
		if arg == "--" {
			seenDoubleDash = true
			continue
		}
		if seenDoubleDash {
			programArgs = append(programArgs, arg)
		} else {
			inputFiles = append(inputFiles, arg)
		}
	}

	tempExe := filepath.Join(os.TempDir(), "cwc-run-*")

	// Compile to temp file
	comp := cwcompiler.NewCompiler(inputFiles, tempExe)
	comp.Verbose = *verbose

	if err := comp.Compile(); err != nil {
		log.Fatalf("Compilation failed: %v", err)
	}

	// Run the compiled program
	cmd := exec.Command(tempExe, programArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("Program exited with error: %v", err)
	}

	// Clean up
	os.Remove(tempExe)
}

func fmtCmd() {
	fs := flag.NewFlagSet("fmt", flag.ExitOnError)
	write := fs.Bool("w", false, "write result to (source) file instead of stdout")
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: cwc fmt [flags] [path ...]\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(os.Args[2:]); err != nil {
		log.Fatalf("Error parsing flags: %v", err)
	}

	args := fs.Args()
	if len(args) == 0 {
		args = []string{"."}
	}

	var files []string
	for _, arg := range args {
		info, err := os.Stat(arg)
		if err != nil {
			log.Printf("Error reading %s: %v", arg, err)
			continue
		}
		if info.IsDir() {
			dirFiles, err := findCWFiles(arg)
			if err != nil {
				log.Printf("Error finding .cw files in %s: %v", arg, err)
				continue
			}
			files = append(files, dirFiles...)
		} else if strings.HasSuffix(arg, ".cw") {
			files = append(files, arg)
		}
	}

	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			log.Printf("Error reading %s: %v", file, err)
			continue
		}
		// very simple formatter: normalize line endings and trim trailing spaces
		s := strings.ReplaceAll(string(data), "\r\n", "\n")
		lines := strings.Split(s, "\n")
		for i, l := range lines {
			lines[i] = strings.TrimRight(l, " \t")
		}
		out := strings.Join(lines, "\n")
		if *write {
			if err := ioutil.WriteFile(file, []byte(out), 0644); err != nil {
				log.Printf("Error writing %s: %v", file, err)
			}
		} else {
			fmt.Print(out)
		}
	}
}

func findCWFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".cw") {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

func printUsage(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
	fmt.Fprintf(os.Stderr, "%s\n", appDescription)
}

func printVersion() {
	fmt.Printf("%s version %s (commit %s, built %s)\n", appName, version, commit, date)
}

func doBuild(inPath, outPath string, keep bool) error {
	src, err := ioutil.ReadFile(inPath)
	if err != nil {
		return err
	}
	tokens := lexer.New(string(src)).Tokenize()
	p := parser.New(tokens)
	program, perr := p.ParseProgram()
	if perr != nil {
		return perr
	}
	if err := checker.CheckProgram(program); err != nil {
		return err
	}
	gosrc := codegen.Generate(program)
	tmpDir, err := ioutil.TempDir("", "clockwise-build-")
	if err != nil {
		return err
	}
	if !keep {
		defer os.RemoveAll(tmpDir)
	}
	mainFile := filepath.Join(tmpDir, "main.go")
	if err := ioutil.WriteFile(mainFile, []byte(gosrc), 0644); err != nil {
		return err
	}
	// merge runtime helpers into tmpDir (copy all files under runtime/* into tmpDir)
	if info, err := os.Stat("runtime"); err == nil && info.IsDir() {
		entries, err := os.ReadDir("runtime")
		if err != nil {
			return err
		}
		for _, e := range entries {
			srcPath := filepath.Join("runtime", e.Name())
			if err := copyDir(srcPath, tmpDir); err != nil {
				return err
			}
		}
	}
	if err := ioutil.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module cwtemp\n\n go 1.20\n"), 0644); err != nil {
		return err
	}
	buildCmd := exec.Command("go", "build", "-o", outPath)
	buildCmd.Dir = tmpDir
	buildCmd.Stderr = os.Stderr
	buildCmd.Stdout = os.Stdout
	return buildCmd.Run()
}

// copyDir recursively copies a directory from src to dst (dst will be created).
func copyDir(src string, dst string) error {
    entries, err := os.ReadDir(src)
    if err != nil {
        return err
    }
    if err := os.MkdirAll(dst, 0755); err != nil {
        return err
    }
    for _, e := range entries {
        srcPath := filepath.Join(src, e.Name())
        dstPath := filepath.Join(dst, e.Name())
        if e.IsDir() {
            // Flatten .go files from subdirectories into dst root so helpers
            // defined in runtime subfolders become available to the generated
            // `package main` code. Non-.go files are copied into subdirs.
            subEntries, err := os.ReadDir(srcPath)
            if err != nil {
                return err
            }
            for _, se := range subEntries {
                ssrc := filepath.Join(srcPath, se.Name())
                if se.IsDir() {
                    // recurse into deeper dirs
                    if err := copyDir(ssrc, dst); err != nil {
                        return err
                    }
                    continue
                }
                // If it's a Go source file, copy it into dst root (flatten).
                if filepath.Ext(se.Name()) == ".go" {
                    b, err := os.ReadFile(ssrc)
                    if err != nil {
                        return err
                    }
                    // Rewrite package declaration to `package main` so runtime
                    // helper files (which may live in a different package in-repo)
                    // become part of the generated `package main` module.
                    // This is a simple heuristic replacing the first `package` line.
                    out := b
                    // attempt to find the package declaration at the top of the file
                    lines := strings.SplitN(string(b), "\n", 4)
                    if len(lines) > 0 && strings.HasPrefix(strings.TrimSpace(lines[0]), "package ") {
                        lines[0] = "package main"
                        out = []byte(strings.Join(lines, "\n"))
                    }
                    // write into dst root
                    if err := os.WriteFile(filepath.Join(dst, se.Name()), out, 0644); err != nil {
                        return err
                    }
                } else {
                    // keep non-go files under their respective subdirectory
                    if err := os.MkdirAll(dstPath, 0755); err != nil {
                        return err
                    }
                    b, err := os.ReadFile(ssrc)
                    if err != nil {
                        return err
                    }
                    if err := os.WriteFile(filepath.Join(dstPath, se.Name()), b, 0644); err != nil {
                        return err
                    }
                }
            }
        } else {
            b, err := os.ReadFile(srcPath)
            if err != nil {
                return err
            }
            if err := os.WriteFile(dstPath, b, 0644); err != nil {
                return err
            }
        }
    }
    return nil
}
