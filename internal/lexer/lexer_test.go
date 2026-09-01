package lexer

import (
	"os"
	"testing"
)

// Real, hand-verified against GRAMMAR.md §7.1's own worked derivation.
func TestLexFirstXorExample(t *testing.T) {
	src := "🌒 ⚓ 🌓 ❓ 🌒 🔀 🌔 : 🌔 ⚓ 🌔 ❓ 🔗 🌒 : 🌑"
	toks, err := Lex(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []Kind{
		KindState, KindEq, KindState, KindQuery, KindState, KindXor4, KindState, KindColon,
		KindState, KindEq, KindState, KindQuery, KindAnd4, KindState, KindColon, KindState,
	}
	if len(toks) != len(want) {
		t.Fatalf("got %d tokens, want %d: %+v", len(toks), len(want), toks)
	}
	for i, k := range want {
		if toks[i].Kind != k {
			t.Errorf("token %d: got %s, want %s", i, toks[i].Kind, k)
		}
	}
}

// The repo's own real xor_check.llll sample (examples/xor_check.llll) -- lexing this is the
// real, concrete smoke test that today's LO source in this repo actually tokenizes.
func TestLexXorCheckExample(t *testing.T) {
	data, err := os.ReadFile("../../examples/xor_check.llll")
	if err != nil {
		t.Fatalf("could not read examples/xor_check.llll: %v", err)
	}
	toks, err := Lex(string(data))
	if err != nil {
		t.Fatalf("unexpected error lexing the repo's own example file: %v", err)
	}
	if len(toks) == 0 {
		t.Fatal("expected at least one token")
	}
	if toks[0].Kind != KindDoor {
		t.Errorf("first token: got %s, want DOOR", toks[0].Kind)
	}
}

// GRAMMAR.md §1.2's own matching rule: bare and VS16-suffixed forms of the same base emoji
// must lex to the identical token.
func TestVariationSelectorStripped(t *testing.T) {
	bare, err := Lex("⚖")
	if err != nil {
		t.Fatalf("bare form: %v", err)
	}
	withVS16, err := Lex("⚖️")
	if err != nil {
		t.Fatalf("VS16 form: %v", err)
	}
	if len(bare) != 1 || len(withVS16) != 1 || bare[0].Kind != withVS16[0].Kind {
		t.Errorf("bare=%+v withVS16=%+v, want identical single DIMLEN token", bare, withVS16)
	}
	if bare[0].Kind != KindDimlen {
		t.Errorf("got %s, want DIMLEN", bare[0].Kind)
	}
}

// GRAMMAR.md §1.2: a ZWJ immediately after a token codepoint is a fatal lex error, not
// silently absorbed.
func TestZWJIsFatal(t *testing.T) {
	_, err := Lex("🌑‍🌒") // new-moon ZWJ waxing-crescent
	if err == nil {
		t.Fatal("expected a lex error for a ZWJ-joined sequence, got none")
	}
}

// GRAMMAR.md §2: the vector keyword is matched as one exact literal, never split into word
// tokens, and an unrecognized character anywhere is a fatal error.
func TestVecLitAndUnknownChar(t *testing.T) {
	toks, err := Lex("vec PARENA CONSTRUCT 312 🌒🌓🌔")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []Kind{KindVecLit, KindState, KindState, KindState}
	if len(toks) != len(want) {
		t.Fatalf("got %d tokens, want %d: %+v", len(toks), len(want), toks)
	}
	for i, k := range want {
		if toks[i].Kind != k {
			t.Errorf("token %d: got %s, want %s", i, toks[i].Kind, k)
		}
	}

	if _, err := Lex("hello"); err == nil {
		t.Fatal("expected a lex error for plain ASCII text outside the exact vector keyword")
	}
}

// LO_Formal_Grammar_Phase_0_Complete.md §20's own worked example: 🔤"cat".
func TestLexLiteralQuotedText(t *testing.T) {
	toks, err := Lex(`🔤"cat"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(toks) != 1 {
		t.Fatalf("got %d tokens, want 1: %+v", len(toks), toks)
	}
	if toks[0].Kind != KindLiteral {
		t.Errorf("got %s, want LITERAL", toks[0].Kind)
	}
	if toks[0].Text != "cat" {
		t.Errorf("got Text=%q, want %q", toks[0].Text, "cat")
	}
}

// Real, resolved decision (see lexQuotedText's own doc comment): LITERAL's opening quote must
// immediately follow with no intervening whitespace, unlike every other token pair.
func TestLexLiteralRejectsSpaceBeforeQuote(t *testing.T) {
	if _, err := Lex(`🔤 "cat"`); err == nil {
		t.Fatal("expected a lex error for whitespace between LITERAL and its opening quote")
	}
}

func TestLexLiteralRejectsUnterminatedQuote(t *testing.T) {
	if _, err := Lex(`🔤"cat`); err == nil {
		t.Fatal("expected a lex error for an unterminated quoted literal")
	}
}

// LO_Formal_Grammar_Phase_0_Complete.md §2.5's remaining pattern tokens not already covered by
// the pre-existing base4-pattern glyph set (WILDCARD/STAR/ONEPLUS/OPT/START/END/ALT/GROUP).
func TestLexPatternTokens(t *testing.T) {
	toks, err := Lex("🅰️ 🚫 ↔️ 🛡️")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []Kind{KindClass, KindNClass, KindRange, KindEscape}
	if len(toks) != len(want) {
		t.Fatalf("got %d tokens, want %d: %+v", len(toks), len(want), toks)
	}
	for i, k := range want {
		if toks[i].Kind != k {
			t.Errorf("token %d: got %s, want %s", i, toks[i].Kind, k)
		}
	}
}
