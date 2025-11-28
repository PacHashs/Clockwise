# Clockwise — Introduction

Clockwise is my “tiny but honest” systems programming language. It borrows the
predictability of C, keeps the runtime in plain Go, and aims to be readable
enough that you can onboard in an evening. The guiding principles are:

- Predictability and performance: tiny runtime, zero magic.
- Readability: familiar C-ish syntax plus a lightweight standard library.
- Friendly tooling: the repo ships `cwc`, `cwfmt`, `cwdoc`, and a Windows GUI
  installer so you can share binaries without asking folks to compile from
  source.

This folder contains the developer-facing docs: language overview, usage tips,
and contribution notes. If you’re new, start with:

- `LANGUAGE_SPEC.md` for the grammar and semantics.
- `USAGE.md` for everyday commands (compiling, running, formatting, documenting).
- `CONTRIBUTING.md` if you plan to send patches.

For a guided install experience on Windows, grab the packaged installer from a
release or run `./build_installer.ps1` to build one yourself. On other
platforms, `cwc` is a single Go build away.
