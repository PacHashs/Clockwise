# ClockWise Usage Guide

## Quick Start

### Prerequisites
- Go 1.20 or later
- Git (for source installation)

### Installation

#### From Source

Windows PowerShell:
```powershell
git clone https://codeberg.org/clockwise-lang/clockwise.git
cd clockwise
go build -o bin\cwc.exe ./cmd/cw
```

Linux/macOS:
```bash
git clone https://codeberg.org/clockwise-lang/clockwise.git
cd clockwise
go build -o bin/cwc ./cmd/cw
```

#### Using Installer (Windows)
```powershell
# Build the installer
.\build_installer.ps1

# Run the installer (requires admin)
bin\installer.exe --bootstrap
```

## Basic Commands

### Building Programs
```bash
# Build a program
cwc build program.cw -o program

# Build with debug information
cwc build -d program.cw -o debug_program
```

### Running Programs
```bash
# Compile and run in one step
cwc run program.cw

# Pass command-line arguments
cwc run program.cw --arg1 value1
```

### Code Formatting
```bash
# Format a file and print to stdout
cwc fmt program.cw

# Format and overwrite files
cwc fmt -w *.cw

# Show diffs of what would change
cwc fmt -d program.cw
```

### Documentation
```bash
# Generate documentation for a package
cwc doc package.cw

# Write docs to a directory
cwc doc -o docs/ package/
```

Runtime and libraries

Runtime helpers are simple Go files under `runtime/`. During compilation the
compiler merges these helpers into the temporary Go module used for `go build`,
so functions implemented in `runtime/` are callable directly from your Clockwise
programs.

Common helpers (examples)
- `Print(s string) int` — write to stdout
- `Sconcat(a, b string) string` — string concat helper
- `NowISO()` / `SleepMs(ms int)` — time helpers
- `HttpGet(url string) string` — simple HTTP helper
- `SHA256Hex(s string) string` — cryptographic helper

More helpers live under `runtime/` (examples: `environs`, `metriclib`, `complib`, `urllib`, `uuidlib`, `stringx`). Call them directly from your `.cw` code once compiled.

Inline example — simple "Hello" program

```cw
fn main() -> int {
    Print("Hello Clockwise!\n")
    return 0
}
```

See `docs/LANGUAGE_SPEC.md` for a brief language overview and `docs/CONTRIBUTING.md` for contribution notes.
