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
# Easiest: run the bootstrap installer script (per-user, no admin)
./installer.ps1
```

## Basic Commands

### Building Programs
```bash
# Build a single file program
cwc build program.cw -o program

# Build multi-file program
cwc build main.cw utils.cw helpers.cw -o myapp

# Build with debug information
cwc build -d program.cw -o debug_program
```

### Running Programs
```bash
# Compile and run single file in one step
cwc run program.cw

# Compile and run multi-file program
cwc run main.cw utils.cw

# Pass command-line arguments (use -- to separate files from args)
cwc run program.cw -- arg1 value1
cwc run main.cw utils.cw -- --program-flag value
```

### Multi-File Projects

Clockwise supports organizing code across multiple files using import statements:

#### Import Syntax
```cw
// Import another Clockwise file
import "utils.cw";

fn main() -> int {
    var result: string = greet("World");
    Print(result + "\n");
    return 0;
}
```

#### Multi-File Compilation
```bash
# Compile multiple files together
cwc build main.cw utils.cw helpers.cw -o application

# Run multiple files together
cwc run main.cw utils.cw helpers.cw

# All functions from all files are merged into a single binary
# No external dependencies are created
```

#### Rules and Limitations
- Import statements use relative file paths with `.cw` extension
- All functions from imported files are available globally
- Duplicate function names across files are detected as errors
- No circular imports allowed
- All files are compiled into a single executable

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
- `Gzip(s string) string` / `Gunzip(b64 string) string` — gzip to base64 and back
- `HMACSHA256(key, data string) string` — hex-encoded HMAC-SHA256
- `CRC32Hex(s string) string` — hex-encoded CRC32
- `URLEncode(s string) string` / `URLDecode(s string) string` — percent-encode/decode
- `CreateTemp(prefix string) string` / `WriteTemp(prefix, data string) string` / `Remove(path string) int` — temp file helpers

More helpers live under `runtime/` (examples: `environs`, `metriclib`, `complib`, `urllib`, `uuidlib`, `stringx`). Call them directly from your `.cw` code once compiled.

Inline example — simple "Hello" program

```cw
fn main() -> int {
    Print("Hello Clockwise!\n")
    return 0
}
```

See `docs/LANGUAGE_SPEC.md` for a brief language overview and `docs/CONTRIBUTING.md` for contribution notes.
