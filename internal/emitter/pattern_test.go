package emitter

import (
	"testing"

	"github.com/emilyspringerton/LO/internal/lexer"
	"github.com/emilyspringerton/LO/internal/parser"
)

// parseMatchPattern lexes+parses a full LO program of the shape `S1 MATCH <pattern> ? S1 : S1;`
// and returns the Pattern out of its own Match Cond -- real, shared test plumbing so every case
// below only has to write the pattern's own emoji source.
func parseMatchPattern(t *testing.T, patternSrc string) parser.Pattern {
	t.Helper()
	toks, err := lexer.Lex("🌒 🔍 " + patternSrc + " ❓ 🌒 : 🌒;")
	if err != nil {
		t.Fatalf("unexpected lex error: %v", err)
	}
	prog, err := parser.Parse(toks)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	m, ok := prog.Body.(parser.Ternary).Cond.(parser.Match)
	if !ok {
		t.Fatalf("expected the Cond to be a Match, got %T", prog.Body.(parser.Ternary).Cond)
	}
	return m.Pattern
}

// Every case here is one of LO_Formal_Grammar_Phase_0_Complete.md's own §20-29 worked examples,
// checked against the EXACT PCRE text the doc itself shows as that pattern's "equivalent PCRE".
func TestPatternToPCREMatchesDocExamples(t *testing.T) {
	cases := []struct {
		name     string
		src      string
		wantPCRE string
	}{
		{"literal", `🔤"cat"`, "cat"},
		{"wildcard", `🔤"c" 🃏 🔤"t"`, "c.t"},
		{"star", `🔤"a" 🌌`, "a*"},
		{"plus", `🔤"a" ☄️`, "a+"},
		{"opt", `🔤"a" 👻`, "a?"},
		{"start-anchor", `🏁 🔤"cat"`, "^cat"},
		{"end-anchor", `🔤"cat" 🛑`, "cat$"},
		{"both-anchors", `🏁 🔤"cat" 🛑`, "^cat$"},
		{"alternation-2", `🔤"cat" 🛤️ 🔤"dog"`, "cat|dog"},
		{"alternation-3", `🔤"cat" 🛤️ 🔤"dog" 🛤️ 🔤"fox"`, "cat|dog|fox"},
		{"alternation-precedence", `🔤"a" 🔤"b" 🛤️ 🔤"c" 🔤"d"`, "ab|cd"},
		{"class-positive", `🅰️ 🔤"a" 🔤"b" 🔤"c"`, "[abc]"},
		{"class-negated", `🚫 🔤"a" 🔤"b" 🔤"c"`, "[^abc]"},
		{"class-one-range", `🅰️ 🔤"a" ↔️ 🔤"z"`, "[a-z]"},
		{"class-multi-range", `🅰️ 🔤"a" ↔️ 🔤"z" 🔤"A" ↔️ 🔤"Z" 🔤"0" ↔️ 🔤"9"`, "[a-zA-Z0-9]"},
		{"escaped-star", `🛡️ 🌌`, `\*`},
		{"escaped-wildcard", `🛡️ 🃏`, `\.`},
		{"group", `🗜️ 🔤"ab" 🗜️`, "(ab)"},
		{"quantified-group", `🗜️ 🔤"ab" 🗜️ ☄️`, "(ab)+"},
		{"exact-match", `🏁 🔤"cat" 🛑`, "^cat$"},
		{"one-or-more-digits", `🅰️ 🔤"0" ↔️ 🔤"9" ☄️`, "[0-9]+"},
		{"cat-or-dog", `🔤"cat" 🛤️ 🔤"dog"`, "cat|dog"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			pat := parseMatchPattern(t, c.src)
			got, err := patternToPCRE(pat)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != c.wantPCRE {
				t.Errorf("got %q, want %q", got, c.wantPCRE)
			}
		})
	}
}

// A literal containing its own PCRE metacharacters must come out escaped -- §20's own "provides
// Unicode literal content" (i.e. NOT PCRE syntax) is not exercised by any of the doc's own worked
// examples (none of them use a metacharacter inside a LITERAL), so this is this repo's own real,
// separate correctness check, not a doc-derived one.
func TestPatternToPCREEscapesLiteralMetacharacters(t *testing.T) {
	pat := parseMatchPattern(t, `🔤"a.b*c"`)
	got, err := patternToPCRE(pat)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := `a\.b\*c`; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// ESCAPE GROUP has no unambiguous literal-character mapping in this v0 (see escapedToPCRE's own
// doc comment) -- a real, honest error, not a silent guess.
func TestPatternToPCREEscapedGroupIsAnHonestError(t *testing.T) {
	pat := parseMatchPattern(t, `🛡️ 🗜️`)
	if _, err := patternToPCRE(pat); err == nil {
		t.Fatal("expected an error for ESCAPE GROUP, got none")
	}
}
