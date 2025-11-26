package main

import (
    "flag"
    "fmt"
    "io/ioutil"
    "os"
    "strings"
)

// cwfmt is a tiny formatter: trims trailing spaces and ensures unix newlines.
func main() {
    in := flag.String("in", "", "input .cw file")
    out := flag.String("out", "", "output file (defaults to overwrite input)")
    flag.Parse()
    if *in == "" {
        fmt.Fprintln(os.Stderr, "-in required")
        os.Exit(2)
    }
    data, err := ioutil.ReadFile(*in)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    s := string(data)
    // normalize newlines
    s = strings.ReplaceAll(s, "\r\n", "\n")
    // trim trailing spaces on each line
    lines := strings.Split(s, "\n")
    for i, l := range lines {
        lines[i] = strings.TrimRight(l, " \t")
    }
    outS := strings.Join(lines, "\n")
    if *out == "" {
        *out = *in
    }
    if err := ioutil.WriteFile(*out, []byte(outS), 0644); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
