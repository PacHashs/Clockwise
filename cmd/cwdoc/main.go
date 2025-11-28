package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Documentation represents the extracted information for a single source file.
type Documentation struct {
	FilePath    string        `json:"filePath"`
	PackageName string        `json:"package"`
	Functions   []FunctionDoc `json:"functions"`
	Types       []TypeDoc     `json:"types"`
	Constants   []ConstantDoc `json:"constants"`
}

type FunctionDoc struct {
	Name       string  `json:"name"`
	Signature  string  `json:"signature"`
	Params     []Param `json:"params"`
	ReturnType string  `json:"returnType"`
	Comment    string  `json:"comment"`
}

type TypeDoc struct {
	Name    string  `json:"name"`
	Fields  []Field `json:"fields"`
	Comment string  `json:"comment"`
}

type Field struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type ConstantDoc struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	Comment string `json:"comment"`
}

type Param struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type docSet struct {
	GeneratedAt time.Time       `json:"generatedAt"`
	SourceRoot  string          `json:"sourceRoot"`
	Documents   []Documentation `json:"documents"`
}

func main() {
	root := flag.String("dir", ".", "Root directory to scan for Clockwise sources")
	outputDir := flag.String("output", "docs/cwdoc", "Directory for generated docs")
	format := flag.String("format", "markdown", "Output format: markdown, html, or json")
	extensions := flag.String("extensions", ".cw", "Comma-separated source extensions to include")
	combined := flag.Bool("single", false, "Emit a single combined file instead of per-source files")
	flag.Parse()

	exts := normalizeExtensions(*extensions)
	files, err := collectSources(*root, exts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cwdoc: failed to walk %s: %v\n", *root, err)
		os.Exit(1)
	}
	if len(files) == 0 {
		fmt.Printf("cwdoc: no %v files found under %s\n", exts, *root)
		return
	}

	var docs []Documentation
	for _, abs := range files {
		display := abs
		if rel, err := filepath.Rel(*root, abs); err == nil {
			display = rel
		}
		doc, err := parseFile(abs, display)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cwdoc: skipping %s: %v\n", abs, err)
			continue
		}
		docs = append(docs, doc)
	}
	if len(docs) == 0 {
		fmt.Println("cwdoc: no documentation artifacts were produced")
		return
	}

	if err := os.MkdirAll(*outputDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "cwdoc: unable to create %s: %v\n", *outputDir, err)
		os.Exit(1)
	}

	set := docSet{
		GeneratedAt: time.Now().UTC(),
		SourceRoot:  *root,
		Documents:   docs,
	}

	if *combined {
		outfile := filepath.Join(*outputDir, fmt.Sprintf("cwdoc.%s", formatExtension(*format)))
		if err := writeCombined(set, outfile, *format); err != nil {
			fmt.Fprintf(os.Stderr, "cwdoc: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("cwdoc: wrote %s\n", outfile)
		return
	}

	for _, doc := range docs {
		outfile := filepath.Join(*outputDir, fmt.Sprintf("%s.%s", trimExt(doc.FilePath), formatExtension(*format)))
		if err := writeSingle(doc, outfile, *format); err != nil {
			fmt.Fprintf(os.Stderr, "cwdoc: %v\n", err)
			continue
		}
		fmt.Printf("cwdoc: wrote %s\n", outfile)
	}
}

func collectSources(root string, extensions []string) ([]string, error) {
	var files []string
	extMap := make(map[string]struct{}, len(extensions))
	for _, ext := range extensions {
		extMap[strings.ToLower(ext)] = struct{}{}
	}
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(d.Name()))
		if _, ok := extMap[ext]; ok {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	return files, nil
}

func parseFile(absPath, display string) (Documentation, error) {
	f, err := os.Open(absPath)
	if err != nil {
		return Documentation{}, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var (
		doc         = Documentation{FilePath: display}
		commentBuf  []string
		inStruct    bool
		currentType *TypeDoc
		inConst     bool
	)

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "//") {
			commentBuf = append(commentBuf, strings.TrimSpace(strings.TrimPrefix(trimmed, "//")))
			continue
		}
		if trimmed == "" {
			commentBuf = nil
			continue
		}
		if strings.HasPrefix(trimmed, "package ") {
			fields := strings.Fields(trimmed)
			if len(fields) > 1 {
				doc.PackageName = fields[1]
			}
			commentBuf = nil
			continue
		}
		if inStruct {
			if strings.HasPrefix(trimmed, "}") {
				inStruct = false
				currentType = nil
				commentBuf = nil
				continue
			}
			parts := strings.Fields(trimmed)
			if len(parts) > 0 {
				field := Field{Name: parts[0]}
				if len(parts) > 1 {
					field.Type = strings.Join(parts[1:], " ")
				}
				if currentType != nil {
					currentType.Fields = append(currentType.Fields, field)
				}
			}
			continue
		}
		if inConst {
			if strings.HasPrefix(trimmed, ")") {
				inConst = false
				commentBuf = nil
				continue
			}
			if constDoc, ok := parseConstant(trimmed); ok {
				constDoc.Comment = strings.Join(commentBuf, "\n")
				doc.Constants = append(doc.Constants, constDoc)
			}
			commentBuf = nil
			continue
		}
		if strings.HasPrefix(trimmed, "const (") {
			inConst = true
			continue
		}
		if fn, ok := parseFunction(trimmed); ok {
			fn.Comment = strings.Join(commentBuf, "\n")
			doc.Functions = append(doc.Functions, fn)
			commentBuf = nil
			continue
		}
		if td, ok := parseType(trimmed); ok {
			td.Comment = strings.Join(commentBuf, "\n")
			doc.Types = append(doc.Types, td)
			if strings.Contains(trimmed, "struct") && strings.HasSuffix(trimmed, "{") {
				inStruct = true
				currentType = &doc.Types[len(doc.Types)-1]
			}
			commentBuf = nil
			continue
		}
		if constDoc, ok := parseConstant(trimmed); ok {
			constDoc.Comment = strings.Join(commentBuf, "\n")
			doc.Constants = append(doc.Constants, constDoc)
			commentBuf = nil
			continue
		}
		commentBuf = nil
	}

	return doc, scanner.Err()
}

func parseFunction(line string) (FunctionDoc, bool) {
	line = strings.TrimSpace(line)
	var keyword string
	switch {
	case strings.HasPrefix(line, "fn "):
		keyword = "fn"
	case strings.HasPrefix(line, "func "):
		keyword = "func"
	default:
		return FunctionDoc{}, false
	}

	rest := strings.TrimSpace(line[len(keyword):])
	parenStart := strings.Index(rest, "(")
	parenEnd := strings.Index(rest, ")")
	if parenStart == -1 || parenEnd == -1 || parenEnd < parenStart {
		return FunctionDoc{}, false
	}

	name := strings.TrimSpace(rest[:parenStart])
	if name == "" {
		return FunctionDoc{}, false
	}

	fn := FunctionDoc{Name: name, Signature: line}
	paramsChunk := strings.TrimSpace(rest[parenStart+1 : parenEnd])
	if paramsChunk != "" {
		parts := strings.Split(paramsChunk, ",")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			pieces := strings.Fields(part)
			param := Param{Name: pieces[0]}
			if len(pieces) > 1 {
				param.Type = strings.Join(pieces[1:], " ")
			}
			fn.Params = append(fn.Params, param)
		}
	}

	returnSegment := strings.TrimSpace(rest[parenEnd+1:])
	if strings.HasPrefix(returnSegment, "->") {
		returnSegment = strings.TrimSpace(strings.TrimPrefix(returnSegment, "->"))
	}
	returnSegment = strings.TrimSuffix(returnSegment, "{")
	returnSegment = strings.TrimSpace(returnSegment)
	if returnSegment != "" && returnSegment != "{" {
		fn.ReturnType = returnSegment
	}

	return fn, true
}

func parseType(line string) (TypeDoc, bool) {
	if !strings.HasPrefix(line, "type ") {
		return TypeDoc{}, false
	}
	parts := strings.Fields(line)
	if len(parts) < 3 {
		return TypeDoc{}, false
	}
	return TypeDoc{Name: parts[1]}, true
}

func parseConstant(line string) (ConstantDoc, bool) {
	switch {
	case strings.HasPrefix(line, "const "):
		line = strings.TrimSpace(line[len("const "):])
	case strings.HasPrefix(line, "let "):
		line = strings.TrimSpace(line[len("let "):])
	default:
		return ConstantDoc{}, false
	}
	if line == "" {
		return ConstantDoc{}, false
	}
	parts := strings.SplitN(line, "=", 2)
	name := strings.TrimSpace(parts[0])
	value := ""
	if len(parts) == 2 {
		value = strings.TrimSpace(parts[1])
	}
	return ConstantDoc{Name: name, Value: value}, true
}

func writeSingle(doc Documentation, outfile, format string) error {
	data, err := render(docSet{Documents: []Documentation{doc}}, format)
	if err != nil {
		return err
	}
	return os.WriteFile(outfile, data, 0o644)
}

func writeCombined(set docSet, outfile, format string) error {
	data, err := render(set, format)
	if err != nil {
		return err
	}
	return os.WriteFile(outfile, data, 0o644)
}

func render(set docSet, format string) ([]byte, error) {
	switch strings.ToLower(format) {
	case "markdown", "md":
		return []byte(generateMarkdown(set)), nil
	case "html", "htm":
		html, err := generateHTML(set)
		if err != nil {
			return nil, err
		}
		return []byte(html), nil
	case "json":
		return json.MarshalIndent(set, "", "  ")
	default:
		return nil, fmt.Errorf("unknown format %q", format)
	}
}

func generateMarkdown(set docSet) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("# Clockwise Documentation\n\nGenerated %s from %s\n\n", set.GeneratedAt.Format(time.RFC3339), set.SourceRoot))
	for _, doc := range set.Documents {
		writeMarkdownDoc(&b, doc)
	}
	return b.String()
}

func writeMarkdownDoc(b *strings.Builder, doc Documentation) {
	b.WriteString(fmt.Sprintf("## %s\n\n", doc.FilePath))
	if doc.PackageName != "" {
		b.WriteString(fmt.Sprintf("Package: `%s`\n\n", doc.PackageName))
	}
	if len(doc.Functions) > 0 {
		b.WriteString("### Functions\n\n")
		for _, fn := range doc.Functions {
			b.WriteString(fmt.Sprintf("- **%s** — %s\n", fn.Name, sanitize(fn.Comment)))
			if len(fn.Params) > 0 {
				b.WriteString("  - Params:\n")
				for _, p := range fn.Params {
					if p.Type != "" {
						b.WriteString(fmt.Sprintf("    - `%s %s`\n", p.Name, p.Type))
					} else {
						b.WriteString(fmt.Sprintf("    - `%s`\n", p.Name))
					}
				}
			}
			if fn.ReturnType != "" {
				b.WriteString(fmt.Sprintf("  - Returns: `%s`\n", fn.ReturnType))
			}
		}
		b.WriteString("\n")
	}
	if len(doc.Types) > 0 {
		b.WriteString("### Types\n\n")
		for _, t := range doc.Types {
			b.WriteString(fmt.Sprintf("- **%s** — %s\n", t.Name, sanitize(t.Comment)))
			if len(t.Fields) > 0 {
				b.WriteString("  - Fields:\n")
				for _, f := range t.Fields {
					if f.Type != "" {
						b.WriteString(fmt.Sprintf("    - `%s %s`\n", f.Name, f.Type))
					} else {
						b.WriteString(fmt.Sprintf("    - `%s`\n", f.Name))
					}
				}
			}
		}
		b.WriteString("\n")
	}
	if len(doc.Constants) > 0 {
		b.WriteString("### Constants\n\n")
		for _, c := range doc.Constants {
			if c.Value != "" {
				b.WriteString(fmt.Sprintf("- `%s = %s` — %s\n", c.Name, c.Value, sanitize(c.Comment)))
			} else {
				b.WriteString(fmt.Sprintf("- `%s` — %s\n", c.Name, sanitize(c.Comment)))
			}
		}
		b.WriteString("\n")
	}
}

func generateHTML(set docSet) (string, error) {
	tmpl := `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <title>Clockwise Documentation</title>
  <style>
  body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; margin: 2rem; }
  code { background: #f3f3f3; padding: 0 .25rem; }
  h2 { margin-top: 2rem; }
  ul { line-height: 1.6; }
  </style>
</head>
<body>
  <h1>Clockwise Documentation</h1>
  <p>Generated {{ .GeneratedAt }} from {{ .SourceRoot }}</p>
  {{ range .Documents }}
    <h2>{{ .FilePath }}</h2>
    {{ if .PackageName }}<p>Package <code>{{ .PackageName }}</code></p>{{ end }}
    {{ if .Functions }}
      <h3>Functions</h3>
      <ul>
      {{ range .Functions }}
        <li><strong>{{ .Name }}</strong> — {{ .Comment }}
          {{ if .Params }}
            <div>Params:
              <ul>
                {{ range .Params }}<li><code>{{ .Name }} {{ .Type }}</code></li>{{ end }}
              </ul>
            </div>
          {{ end }}
          {{ if .ReturnType }}<div>Returns: <code>{{ .ReturnType }}</code></div>{{ end }}
        </li>
      {{ end }}
      </ul>
    {{ end }}
    {{ if .Types }}
      <h3>Types</h3>
      <ul>
        {{ range .Types }}
          <li><strong>{{ .Name }}</strong> — {{ .Comment }}
            {{ if .Fields }}
              <div>Fields:
                <ul>
                {{ range .Fields }}<li><code>{{ .Name }} {{ .Type }}</code></li>{{ end }}
                </ul>
              </div>
            {{ end }}
          </li>
        {{ end }}
      </ul>
    {{ end }}
    {{ if .Constants }}
      <h3>Constants</h3>
      <ul>
        {{ range .Constants }}<li><code>{{ .Name }}</code> — {{ .Value }} {{ .Comment }}</li>{{ end }}
      </ul>
    {{ end }}
  {{ end }}
</body>
</html>`
	var b strings.Builder
	if err := template.Must(template.New("cwdoc").Parse(tmpl)).Execute(&b, set); err != nil {
		return "", err
	}
	return b.String(), nil
}

func sanitize(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "(no summary)"
	}
	return s
}

func normalizeExtensions(raw string) []string {
	parts := strings.Split(raw, ",")
	var out []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if !strings.HasPrefix(p, ".") {
			p = "." + p
		}
		out = append(out, p)
	}
	if len(out) == 0 {
		return []string{".cw"}
	}
	return out
}

func trimExt(path string) string {
	base := filepath.Base(path)
	return strings.TrimSuffix(base, filepath.Ext(base))
}

func formatExtension(format string) string {
	switch strings.ToLower(format) {
	case "markdown", "md":
		return "md"
	case "html", "htm":
		return "html"
	case "json":
		return "json"
	default:
		return format
	}
}

