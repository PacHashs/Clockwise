# Clockwise Language Specification (short)

This is a short, intentionally minimal spec for the Clockwise language.

1. Lexical grammar
- Identifiers: `[A-Za-z_][A-Za-z0-9_]*`
- Integers: decimal sequences (e.g. `123`)
- Strings: double-quoted `"..."` with `\"` supported for an embedded quote

2. Top-level
- Functions: `fn <name>(<params>) -> <type> { ... }`
- Imports: `import "<filename>.cw";`
- The entry point is `fn main() -> int` which returns an integer exit code.

3. Types
- Builtins: `int`, `string`
- Future additions may include floats, arrays, structs, and pointers.

4. Runtime helpers
- Runtime helpers are implemented as Go functions under the `runtime/` folder.
  During compilation these helpers are merged into the temporary Go module used
  for `go build` so generated code can call them directly (for example, `Print`,
  `HttpGet`, `SHA256Hex`). See `docs/USAGE.md` for examples and available helpers.

5. Statements
- `import "<filename>.cw";` - Import another Clockwise file
- `var <name>: <type> = <expr>;`
- `return <expr>;`
- Expression statements (function calls and side-effecting expressions)

6. Expressions
- Literals: integers and strings
- Identifiers and function calls
- Infix expressions with standard precedence (e.g., `*`/`/` before `+`/`-`, comparisons last)

7. Interop and conventions
- The compiler maps language-level print/call expressions to runtime helpers
  (for example a `Print` runtime helper). Generated code calls exported Go
  functions found in `runtime/`.

This spec is intentionally compact to keep the implementation small and auditable. See `docs/USAGE.md` for usage examples.
