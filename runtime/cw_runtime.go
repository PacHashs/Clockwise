package runtimelib

import "fmt"

// cwPrint is the runtime helper used by generated code for printing.
func cwPrint(s string) int {
    fmt.Print(s)
    return 0
}
