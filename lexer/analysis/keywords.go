package lexanalysis

import lex "codeberg.org/clockwise-lang/clockwise/lexer"

// ExtraKeywords can be extended to include future reserved words for the
// language. Keeping as a separate file makes it easy to add feature-based
// keywords without touching core token definitions.
var ExtraKeywords = map[string]lex.TokenType{
    "import": lex.TokenType("IMPORT"),
    "type":   lex.TokenType("TYPE"),
}
