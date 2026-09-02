// Package lexer implements QI_NORTHSTAR.md's Phase 2a: the real, standalone lexer for `qi`'s own
// ASCII s-expression surface syntax. Real, deliberate scope, matching QI_NORTHSTAR.md's own
// phased plan: lexer only -- no parser, no lowering to internal/parser's AST. That's Phase 2b, a
// real, separate next step once this package's own tokens are real and tested.
//
// `qi` is not a new value language -- it's a second, friendlier ASCII FRONTEND onto the exact
// same value space LO's own emoji grammar already covers (see QI_NORTHSTAR.md's own "why qi at
// all" section). Its lexer is accordingly a real, standard s-expression reader: parens/brackets,
// bare symbols (identifiers AND operator spellings like `+4`/`let`/`if` alike -- the PARSER, not
// the lexer, decides what a given symbol spelling means, the same real division every Lisp reader
// uses), decimal integer literals restricted to LO's own four base4 states (0-3), double-quoted
// strings, and `;` line comments (a real, useful addition LO's own emoji grammar has no lexical
// room for -- `;` is already LO's own SEMI token there).
package lexer

// Kind identifies a lexical token class, per QI_NORTHSTAR.md's own first-draft surface grammar.
type Kind int

const (
	KindInvalid Kind = iota
	KindLParen       // (
	KindRParen       // )
	KindLBracket     // [
	KindRBracket     // ]
	KindSymbol       // a bare identifier or operator spelling, e.g. `let`, `if`, `+4`, `x`, `foo-bar`
	KindInt          // a decimal integer literal restricted to 0-3 (see Token.Int)
	KindString       // a double-quoted string literal (see Token.Text)
)

func (k Kind) String() string {
	switch k {
	case KindLParen:
		return "LPAREN"
	case KindRParen:
		return "RPAREN"
	case KindLBracket:
		return "LBRACKET"
	case KindRBracket:
		return "RBRACKET"
	case KindSymbol:
		return "SYMBOL"
	case KindInt:
		return "INT"
	case KindString:
		return "STRING"
	default:
		return "INVALID"
	}
}

// Token is one lexed unit. Pos is the rune offset in the source where the token starts, for real
// error reporting once a parser exists (Phase 2b's own next real step).
type Token struct {
	Kind Kind
	Text string // valid when Kind == KindSymbol (the symbol's own spelling) or KindString (its
	// own unescaped content -- see Lex's own doc comment on why there's no escape processing)
	Int int // valid only when Kind == KindInt: 0-3
	Pos int
}
