package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/your-org/clockwise/cmd/cw/formatter"
	"github.com/your-org/clockwise/cmd/cw/compiler"
	"github.com/your-org/clockwise/internal/logger"
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

	startTime := time.Now()

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
		fmt.Fprintf(os.Stderr, "Usage: cwc build [input.cw] [-o output]\n")
		fs.PrintDefaults()
	}

	// Parse flags and args
	if err := fs.Parse(os.Args[2:]); err != nil {
		log.Fatalf("Error parsing flags: %v", err)
	}

	// Get input file
	args := fs.Args()
	if len(args) == 0 {
		fs.Usage()
		os.Exit(1)
	}
	inputFile := args[0]

	// Set default output file if not specified
	if *outputFile == "" {
		*outputFile = strings.TrimSuffix(inputFile, filepath.Ext(inputFile))
	}

	// Run the compiler
	compiler := NewCompiler(inputFile, *outputFile)
	compiler.Verbose = *verbose

	if err := compiler.Compile(); err != nil {
		log.Fatalf("Compilation failed: %v", err)
	}
}

func runCmd() {
	fs := flag.NewFlagSet("run", flag.ExitOnError)
	verbose := fs.Bool("v", false, "Enable verbose output")
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: cwc run [input.cw] [args...]\n")
		fs.PrintDefaults()
	}

	// Parse flags and args
	if err := fs.Parse(os.Args[2:]); err != nil {
		log.Fatalf("Error parsing flags: %v", err)
	}

	if len(fs.Args()) == 0 {
		fs.Usage()
		os.Exit(1)
	}

	inputFile := fs.Args()[0]
	tempExe := filepath.Join(os.TempDir(), "cwc-run-*"+filepath.Ext(inputFile))

	// Compile to temp file
	compiler := NewCompiler(inputFile, tempExe)
	compiler.Verbose = *verbose

	if err := compiler.Compile(); err != nil {
		log.Fatalf("Compilation failed: %v", err)
	}

	// Run the compiled program
	cmd := exec.Command(tempExe, fs.Args()[1:]...)
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
	diff := fs.Bool("d", false, "display diffs instead of rewriting files")
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: cwc fmt [flags] [path ...]\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(os.Args[2:]); err != nil {
		log.Fatalf("Error parsing flags: %v", err)
	}

	args := fs.Args()
	if len(args) == 0 {
		// Default to current directory
		args = []string{"."}
	}

	var files []string
	for _, arg := range args {
		fileInfo, err := os.Stat(arg)
		if err != nil {
			log.Printf("Error reading %s: %v", arg, err)
			continue
		}

		if fileInfo.IsDir() {
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
		if err := processFile(file, *write, *diff); err != nil {
			log.Printf("Error processing %s: %v", file, err)
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

func processFile(filename string, write, showDiff bool) error {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	res, err := formatter.Format(src)
	if err != nil {
		return fmt.Errorf("formatting %s: %v", filename, err)
	}

	if bytes.Equal(src, res) {
		return nil // No changes
	}

	if showDiff {
		diff, err := diff(src, res, filename)
		if err != nil {
			return fmt.Errorf("computing diff for %s: %v", filename, err)
		}
		fmt.Printf("diff %s cw/fmt/%s\n", filename, filename)
		fmt.Printf("--- %s\n", filename)
		fmt.Printf("+++ cw/fmt/%s\n", filename)
		os.Stdout.Write(diff)
		return nil
	}

	if !write {
		_, err = os.Stdout.Write(res)
		return err
	}

	// Create a backup
	if err := backupFile(filename, src); err != nil {
		return fmt.Errorf("creating backup: %v", err)
	}

	// Write formatted content
	return ioutil.WriteFile(filename, res, 0644)
}

func diff(b1, b2 []byte, filename string) (data []byte, err error) {
	f1, err := writeTempFile("", "cwc-fmt-", b1)
	if err != nil {
		return nil, err
	}
	defer os.Remove(f1)

	f2, err := writeTempFile("", "cwc-fmt-", b2)
	if err != nil {
		return nil, err
	}
	defer os.Remove(f2)

	cmd := exec.Command("diff", "-u", f1, f2)
	data, err = cmd.CombinedOutput()
	if len(data) > 0 {
		// diff exits with a non-zero status when the files don't match.
		// Ignore that failure as long as we get output.
		err = nil
	}
	return data, err
}

func writeTempFile(dir, prefix string, data []byte) (string, error) {
	file, err := ioutil.TempFile(dir, prefix)
	if err != nil {
		return "", err
	}
	_, err = file.Write(data)
	if err1 := file.Close(); err == nil {
		err = err1
	}
	if err != nil {
		os.Remove(file.Name())
		return "", err
	}
	return file.Name(), nil
}

func backupFile(filename string, data []byte) error {
	err := ioutil.WriteFile(filename+".bak", data, 0644)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %v", err)
	}
	return nil
}
    fs := flag.NewFlagSet("build", flag.ExitOnError)
    in := fs.String("in", "", "input .cw file")
    out := fs.String("out", "a.out", "output exe")
    keep := fs.Bool("keep", false, "keep temp build files")
    fs.Parse(args)
    if *in == "" {
        fmt.Fprintln(os.Stderr, "-in required for build")
        os.Exit(2)
    }
    if err := doBuild(*in, *out, *keep); err != nil {
        fmt.Fprintln(os.Stderr, "build failed:", err)
        os.Exit(1)
    }
}

func runCmd(args []string) {
    fs := flag.NewFlagSet("run", flag.ExitOnError)
    in := fs.String("in", "", "input .cw file")
    fs.Parse(args)
    if *in == "" {
        fmt.Fprintln(os.Stderr, "-in required for run")
        os.Exit(2)
    }
    // build to a temp file and run
    tmp := filepath.Join(os.TempDir(), "cw-run-"+filepath.Base(*in)+".exe")
    if err := doBuild(*in, tmp, false); err != nil {
        fmt.Fprintln(os.Stderr, "build failed:", err)
        os.Exit(1)
    }
    cmd := exec.Command(tmp)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin
    if err := cmd.Run(); err != nil {
        fmt.Fprintln(os.Stderr, "program exited with error:", err)
        os.Exit(1)
    }
}

func fmtCmd(args []string) {
    fs := flag.NewFlagSet("fmt", flag.ExitOnError)
    in := fs.String("in", "", "input .cw file")
    write := fs.Bool("w", false, "write changes")
    fs.Parse(args)
    if *in == "" {
        fmt.Fprintln(os.Stderr, "-in required for fmt")
        os.Exit(2)
    }
    out, err := ioutil.ReadFile(*in)
    if err != nil {
        fmt.Fprintln(os.Stderr, "read error:", err)
        os.Exit(1)
    }
    s := string(out)
    s = strings.ReplaceAll(s, "\r\n", "\n")
    lines := strings.Split(s, "\n")
    for i, l := range lines {
        lines[i] = strings.TrimRight(l, " \t")
    }
    res := strings.Join(lines, "\n")
    if *write {
        if err := ioutil.WriteFile(*in, []byte(res), 0644); err != nil {
            fmt.Fprintln(os.Stderr, "write error:", err)
            os.Exit(1)
        }
        fmt.Println("wrote", *in)
    } else {
        fmt.Print(res)
    }
}

func initCmd(args []string) {
    fs := flag.NewFlagSet("init", flag.ExitOnError)
    name := fs.String("name", "clockwise-project", "project directory name")
    fs.Parse(args)
    if err := os.MkdirAll(*name, 0755); err != nil {
        fmt.Fprintln(os.Stderr, "failed to create dir:", err)
        os.Exit(1)
    }
    // create simple main.cw
    example := "fn main() -> int {\n    Print(\"Hello Clockwise!\\n\")\n    return 0\n}\n"
    if err := ioutil.WriteFile(filepath.Join(*name, "main.cw"), []byte(example), 0644); err != nil {
        fmt.Fprintln(os.Stderr, "failed to write example:", err)
        os.Exit(1)
    }
    fmt.Println("created project:", *name)
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
