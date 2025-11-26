package codegen

import (
    "io/ioutil"
    "path/filepath"
)

// WriteMain writes a main source file into dir with the provided content.
// This helper is used by the CLI build step to create temporary files.
func WriteMain(dir, content string) error {
    mainPath := filepath.Join(dir, "main.go")
    return ioutil.WriteFile(mainPath, []byte(content), 0644)
}
