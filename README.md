# Clockwise — tiny, practical, and chill
<p align="center">
  <img src="assets/logo.svg" alt="Clockwise logo" width="420" />
  
</p>
Thanks for dropping by. Clockwise is a small systems language that stays out of your way. The compiler is written in Go (that's just an implementation detail). After you install Clockwise, you do not need Go to use it or to run the programs you build. Outputs are normal, standalone native binaries. Go is the oven; Clockwise is the cake.

What you get:
- Minimal language with C-ish vibes and a tiny standard "runtime" (plain Go helpers).
- Simple tool: `cwc` to build, run, and format.
- Everything lives here: lexer → parser → checker → codegen → runtime.
- 30+ built-in runtime helpers for common tasks (HTTP, crypto, compression, file ops, etc).
- Multi-file project support with import statements.

## Installation

**Windows (one command, no admin):**
```powershell
./installer.ps1
```

**Everyone else:**
```bash
git clone https://codeberg.org/clockwise-lang/clockwise.git
cd clockwise
go build -o bin/cwc ./cmd/cw
```

## Quick Start

Single file:
```cw
fn main() -> int {
    Print("Hello Clockwise!\n")
    return 0
}
```

Multi-file:
```cw
// main.cw
import "utils.cw";

fn main() -> int {
    var message: string = greet("World");
    Print(message + "\n");
    return 0
}
```

```cw
// utils.cw
fn greet(name: string) -> string {
    return "Hello, " + name + "!";
}
```

## Usage

```bash
# Single file
cwc build program.cw -o program
cwc run program.cw

# Multi-file
cwc build main.cw utils.cw -o myapp
cwc run main.cw utils.cw

# Formatting
cwc fmt -w *.cw
```

## Language Features

- **C-like syntax** - `var name: type = value` declarations
- **Type-safe** - Explicit types with `int`, `string` 
- **Functions** - `fn name() -> type { }` syntax
- **Built-in runtime helpers** - HttpGet(), SHA256Hex(), Print(), etc.
- **Import statements** - `import "filename.cw"` for multi-file projects
- **Native binaries** - No runtime, no VM, instant startup
- **Zero dependencies** - Everything compiled into single executable

## Multi-File Projects

Clockwise supports organizing code across multiple files:

- Import other files: `import "utils.cw"`
- Compile multiple files: `cwc build main.cw utils.cw helpers.cw`
- All functions merged into single binary
- No external dependencies or package managers
- Duplicate function names are detected and reported

Repo home (Codeberg):
- https://codeberg.org/clockwise-lang/clockwise

Support Clockwise development:
<a href="https://liberapay.com/root"><img alt="Donate using Liberapay" src="https://liberapay.com/assets/widgets/donate.svg"></a>
