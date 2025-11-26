package lexer

// LexTokens returns a slice of token type+lit strings for quick inspection.
func LexTokens(input string) []string {
    l := New(input)
    toks := l.Tokenize()
    out := make([]string, 0, len(toks))
    for _, t := range toks {
        out = append(out, string(t.Type)+":"+t.Lit)
    }
    return out
}
