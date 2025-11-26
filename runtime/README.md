# Runtime

This folder contains Go-based runtime helpers used by generated Clockwise programs.
Each library is a small Go file (or directory) providing plain functions the compiler
can call from generated code.

Included libraries (high level):
- `cwlib` — core low-level helpers
- `fs` — file read/write and existence checks
- `stringslib` — string utilities
- `jsonlib` — JSON helpers
- `cryptolib` — hashing helpers (SHA256)
- `mathlib` — math helpers (ints and floats)
- `netlib` — HTTP helpers (GET, POST, download, status)
- `timelib` — time helpers (now, sleep)
- `randlib` — random helpers
- `osenv` — environment helpers

Recently added helper libraries:
- `baselib` — base64 encode/decode
- `pathlib` — path joining and basename/dirname helpers
- `regexlib` — regexp match/replace helpers
- `urllib` — simple HTTP helpers (GET with status, POST)
- `statlib` — small integer stats (sum, mean)
- `uuidlib` — UUIDv4 generator
- `hexlib` — hex encode/decode
- `fileutil` — file copy and size helpers
- `stringx` — small string helpers (split/trim)
- `timeparse` — parsing/formatting RFC3339 timestamps

To extend the runtime, add small Go files with `package main` and export simple
functions that the code generator will call. Keep dependencies small to avoid
pulling heavy transitive packages into user binaries.
