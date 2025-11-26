package runtimelib

import (
    "flag"
    "fmt"
    "os"
    "path/filepath"
)

// gen_stdlib writes a small C stdlib into runtime/cw_stdio.c. This allows
// using Go to author and regenerate the C runtime sources deterministically.
func main() {
    out := flag.String("out", "runtime/cw_stdio.c", "output C runtime path")
    flag.Parse()

    data := `// Clockwise runtime (generated)
#include <stdio.h>

void cw_print(const char *s) {
    fputs(s, stdout);
}
`

    if err := os.MkdirAll(filepath.Dir(*out), 0755); err != nil {
        fmt.Fprintln(os.Stderr, "failed to create runtime dir:", err)
        os.Exit(1)
    }
    if err := os.WriteFile(*out, []byte(data), 0644); err != nil {
        fmt.Fprintln(os.Stderr, "failed to write runtime file:", err)
        os.Exit(1)
    }
    fmt.Println("wrote", *out)
}
