package codegen

// RuntimeStub returns a small comment stub explaining runtime helpers link.
func RuntimeStub() string {
    return "// Runtime helpers are provided by the 'runtime' folder and merged into\n" +
        "// the generated temporary module during build. Functions declared there\n" +
        "// are callable directly from generated code.\n\n"
}
