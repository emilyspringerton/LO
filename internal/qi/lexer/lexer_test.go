package lexer

import "testing"

// TestLexParensAndBrackets -- the two structural delimiter pairs, with and without whitespace
// between them (standard Lisp-reader behavior: a bracket/paren is always its own token).
func TestLexParensAndBrackets(t *testing.T) {
	toks, err := Lex("([ ])[(")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []Kind{KindLParen, KindLBracket, KindRBracket, KindRParen, KindLBracket, KindLParen}
	if len(toks) != len(want) {
		t.Fatalf("got %d tokens, want %d: %+v", len(toks), len(want), toks)
	}
	for i, k := range want {
		if toks[i].Kind != k {
			t.Errorf("token %d: got %s, want %s", i, toks[i].Kind, k)
		}
	}
}

// TestLexRealLetExample -- QI_NORTHSTAR.md's own worked example: `(let [x 1] (+4 x x))`.
func TestLexRealLetExample(t *testing.T) {
	toks, err := Lex("(let [x 1] (+4 x x))")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	type want struct {
		kind Kind
		text string
		i    int
	}
	wants := []want{
		{KindLParen, "", 0},
		{KindSymbol, "let", 0},
		{KindLBracket, "", 0},
		{KindSymbol, "x", 0},
		{KindInt, "", 1},
		{KindRBracket, "", 0},
		{KindLParen, "", 0},
		{KindSymbol, "+4", 0},
		{KindSymbol, "x", 0},
		{KindSymbol, "x", 0},
		{KindRParen, "", 0},
		{KindRParen, "", 0},
	}
	if len(toks) != len(wants) {
		t.Fatalf("got %d tokens, want %d: %+v", len(toks), len(wants), toks)
	}
	for i, w := range wants {
		if toks[i].Kind != w.kind {
			t.Errorf("token %d: got kind %s, want %s", i, toks[i].Kind, w.kind)
		}
		if w.kind == KindSymbol && toks[i].Text != w.text {
			t.Errorf("token %d: got text %q, want %q", i, toks[i].Text, w.text)
		}
		if w.kind == KindInt && toks[i].Int != w.i {
			t.Errorf("token %d: got int %d, want %d", i, toks[i].Int, w.i)
		}
	}
}

// TestLexSymbolAdjacentToDelimiters -- no whitespace needed between a symbol and a paren/bracket.
func TestLexSymbolAdjacentToDelimiters(t *testing.T) {
	toks, err := Lex("(let[x 1])")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []Kind{KindLParen, KindSymbol, KindLBracket, KindSymbol, KindInt, KindRBracket, KindRParen}
	if len(toks) != len(want) {
		t.Fatalf("got %d tokens, want %d: %+v", len(toks), len(want), toks)
	}
	for i, k := range want {
		if toks[i].Kind != k {
			t.Errorf("token %d: got %s, want %s", i, toks[i].Kind, k)
		}
	}
	if toks[1].Text != "let" {
		t.Errorf("expected symbol 'let', got %q", toks[1].Text)
	}
}

// TestLexIntLiteralsInRange -- LO's own four base4 states, 0-3.
func TestLexIntLiteralsInRange(t *testing.T) {
	toks, err := Lex("0 1 2 3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(toks) != 4 {
		t.Fatalf("got %d tokens, want 4: %+v", len(toks), toks)
	}
	for i, want := range []int{0, 1, 2, 3} {
		if toks[i].Kind != KindInt || toks[i].Int != want {
			t.Errorf("token %d: got %+v, want INT %d", i, toks[i], want)
		}
	}
}

// TestLexIntLiteralOutOfRangeIsFatal -- a real, honest error, not silently wrapped or truncated.
func TestLexIntLiteralOutOfRangeIsFatal(t *testing.T) {
	if _, err := Lex("4"); err == nil {
		t.Fatal("expected a lex error for an out-of-range integer literal (4)")
	}
	if _, err := Lex("42"); err == nil {
		t.Fatal("expected a lex error for an out-of-range integer literal (42)")
	}
}

// TestLexIntLiteralLeadingZero -- a real, deliberate simplification named in lexInt's own doc
// comment: leading zeros are accepted.
func TestLexIntLiteralLeadingZero(t *testing.T) {
	toks, err := Lex("03")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(toks) != 1 || toks[0].Kind != KindInt || toks[0].Int != 3 {
		t.Errorf("got %+v, want a single INT token with value 3", toks)
	}
}

// TestLexStringLiteral -- double-quoted strings, no escape processing (same rule as LO's own
// LITERAL).
func TestLexStringLiteral(t *testing.T) {
	toks, err := Lex(`"cat"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(toks) != 1 || toks[0].Kind != KindString || toks[0].Text != "cat" {
		t.Errorf("got %+v, want a single STRING token with text \"cat\"", toks)
	}
}

func TestLexStringLiteralUnterminatedIsFatal(t *testing.T) {
	if _, err := Lex(`"cat`); err == nil {
		t.Fatal("expected a lex error for an unterminated string literal")
	}
}

// TestLexComment -- a `;` line comment runs to the next newline (or end of input) and is not
// itself a token.
func TestLexComment(t *testing.T) {
	toks, err := Lex("; a real comment\n(let [x 1] x) ; trailing comment")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []Kind{KindLParen, KindSymbol, KindLBracket, KindSymbol, KindInt, KindRBracket, KindSymbol, KindRParen}
	if len(toks) != len(want) {
		t.Fatalf("got %d tokens, want %d: %+v", len(toks), len(want), toks)
	}
	for i, k := range want {
		if toks[i].Kind != k {
			t.Errorf("token %d: got %s, want %s", i, toks[i].Kind, k)
		}
	}
}

// TestLexArithOperatorSymbols -- QI_NORTHSTAR.md's own ArithOp spellings are just ordinary
// symbols at the lexer level; the parser (Phase 2b) decides what they mean.
func TestLexArithOperatorSymbols(t *testing.T) {
	toks, err := Lex("+4 -4 &4 |4 ^4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"+4", "-4", "&4", "|4", "^4"}
	if len(toks) != len(want) {
		t.Fatalf("got %d tokens, want %d: %+v", len(toks), len(want), toks)
	}
	for i, w := range want {
		if toks[i].Kind != KindSymbol || toks[i].Text != w {
			t.Errorf("token %d: got %+v, want SYMBOL %q", i, toks[i], w)
		}
	}
}

// TestLexSymbolIsPermissive -- real, deliberate design (see Lex's own doc comment on the symbol
// branch): any non-delimiter character is valid inside a bare symbol, with no restricted charset
// invented -- there is no "unrecognized character" case for a lone punctuation rune like this.
func TestLexSymbolIsPermissive(t *testing.T) {
	toks, err := Lex("` @#%")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"`", "@#%"}
	if len(toks) != len(want) {
		t.Fatalf("got %d tokens, want %d: %+v", len(toks), len(want), toks)
	}
	for i, w := range want {
		if toks[i].Kind != KindSymbol || toks[i].Text != w {
			t.Errorf("token %d: got %+v, want SYMBOL %q", i, toks[i], w)
		}
	}
}
