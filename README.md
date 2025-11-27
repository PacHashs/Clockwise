
# Clockwise.  a tiny, practical systems language
<p align="center">
	<img src="assets/logo.svg" alt="Clockwise logo" width="420" />
</p>
Clockwise is a pragmatic, small systems language focused on readability,
fast iteration, and practical interoperability. It targets a compact surface
area so you can ship tiny, self-contained native binaries while leveraging
the Go toolchain for compilation and cross-compilation.

The language is intentionally simple and predictable: a C-like syntax, first-
class functions, and a small runtime implemented as plain Go helpers you can
inspect, extend, or replace.

This repository contains a working compiler (written in Go), a collection of
Go-based runtime helpers, developer tools (formatter and an installer builder),
and documentation. The compiler emits Go source from `.cw` files and uses
`go build` to produce native binaries, so you don't need a C toolchain.

Why Clockwise
- Small & predictable: minimal semantics and a tiny runtime.
- Fast iteration: generate Go and reuse the Go toolchain for builds and
	cross-compiles.
- Easy interop: runtime helpers are plain Go functions that are callable from
	generated code.

What you'll find here
- `cmd/` — CLI tools: `cwc` (compiler & toolchain), `installer` (installer builder)
  - `cwc build` - Compile ClockWise programs
  - `cwc run` - Compile and run in one step
  - `cwc fmt` - Format ClockWise source code
  - `cwc doc` - Generate documentation
- `lexer/`, `parser/`, `checker/`, `codegen/` — compiler subsystems
- `runtime/` — Go-based runtime helpers (printing, fs, strings, crypto, net, math, time, and many small utilities)
- `docs/` — language spec, usage, and contributing notes

Quick start (you only need Go):

Windows PowerShell:
```powershell
# Build the compiler and tools
cd C:\path\to\ClockWise
go build -o bin\cwc.exe ./cmd/cw

# Or use the installer for system-wide installation
.\build_installer.ps1
bin\installer.exe --bootstrap
```

Linux/macOS:
```bash
# Build the compiler and tools
cd /path/to/ClockWise
go build -o bin/cwc ./cmd/cw

# Install system-wide (requires sudo)
sudo cp bin/cwc /usr/local/bin/
```

Basic Usage

```bash
# Build a program
cwc build program.cw -o program

# Run directly without explicit build
cwc run program.cw --arg1 value1

# Format your code
cwc fmt -w *.cw  # Format all .cw files in current directory

# Show documentation
cwc doc package.cw
```

Example `program.cw`:
```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, ClockWise!")
}
```

Note: the `examples/` folder has been removed from this repository. See
`docs/USAGE.md` for inline snippets and quick-start examples, or create your
own small `.cw` files to experiment.

## Installation

### Windows
1. Download the latest installer from the releases page
2. Run the installer with administrative privileges
3. The installer will add `cwc` to your system PATH

### Linux/macOS
```bash
# Download and install
curl -L https://codeberg.org/clockwise-lang/clockwise/releases/latest/download/install.sh | sh

# Or build from source
git clone https://codeberg.org/clockwise-lang/clockwise
cd clockwise
go build -o /usr/local/bin/cwc ./cmd/cw
```

### Building the Installer (Developers)
From the repo root:

Windows PowerShell:
```powershell
./build_installer.ps1
# Produces `bin\installer.exe`
```

Linux/macOS:
```bash
./build_installer.sh
# Produces `bin/installer`
```

The installer can bootstrap the entire toolchain with the `--bootstrap` flag.


Donations
---------
If you like Clockwise and want to support its development, you can donate via Liberapay. Donations go to the project account `Listedroot`.


<a href="https://liberapay.com/Listedroot/donate" target="_blank" rel="noopener noreferrer">
	<img src="https://liberapay.com/assets/widgets/donate.svg" alt="Donate to Listedroot on Liberapay" />
</a>

License & rights
-----------------
This project is licensed under the MIT License. Copyright (c) 2025 BloodyHell-Industries-Inc.

Getting involved
----------------
See `docs/CONTRIBUTING.md` for how to file issues and contribute code. If you want new runtime helpers,
add them under `runtime/` (they're plain Go files and will be merged into generated modules).

Questions or want a feature fast? Open an issue or drop a PR . I try to keep changes small and testable.
