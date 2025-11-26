# Contributing to Clockwise

Thank you for helping improve Clockwise.

How to contribute
- Open issues for bugs, feature requests, or design discussions.
- Keep changes small and focused; add unit tests for new behavior.
- Run `go test ./...` and ensure builds succeed before opening PRs.

Repository layout
- `cmd/` — CLI tools
- `lexer/`, `parser/`, `codegen/`, `checker/` — compiler subsystems
- `runtime/` — Go-based runtime helpers merged into generated modules
- `docs/` — documentation

See `docs/LANGUAGE_SPEC.md` for the language overview and `docs/USAGE.md` for build and quick-start instructions.

Donations (optional)
--------------------
If you find Clockwise useful, you can support the project via Liberapay. The project account is `Listedroot`.

Donate link: https://liberapay.com/Listedroot/donate

If you publish a website or release notes, linking to the Liberapay page is sufficient; embedding payment widgets is optional.
