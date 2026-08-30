package parser

import (
	"testing"

	"github.com/emilyspringerton/LO/internal/lexer"
)

// Real, found-live discrepancy, not silently resolved either way: the repo's own
// examples/xor_check.llll (`🚪 🔢 🌒 ⚓ 🌒 ❓ 🔀 🌒 🌔 : 🌑`) uses PREFIX notation for the arith
// operator (XOR4 S1 S3), but GRAMMAR.md §5.1 (and the original LoLanguageSpec.pdf's own worked
// examples) specify INFIX (`Value ArithOp Value`) — confirmed directly by decoding the file's own
// codepoints, not assumed. Since GRAMMAR.md is the reviewed, canonical Phase 0 spec, this parser
// implements infix, matching the spec; the example file is left as-is (not authored by this
// pass) rather than silently rewritten to agree with one side. This test uses a fresh,
// GRAMMAR.md-consistent infix string with the same real shape instead of that file, and is a
// real, named follow-up: reconcile examples/xor_check.llll with GRAMMAR.md once decided which
// notation LO actually wants.
func TestParseGrammarConsistentXorExample(t *testing.T) {
	toks, err := lexer.Lex("🚪 🔢 🌒 ⚓ 🌒 ❓ 🌒 🔀 🌔 : 🌑 ;")
	if err != nil {
		t.Fatalf("unexpected lex error: %v", err)
	}
	prog, err := Parse(toks)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if !prog.HasDoor || prog.DoorType != TypeI32 {
		t.Fatalf("expected a DOOR I32, got HasDoor=%v DoorType=%v", prog.HasDoor, prog.DoorType)
	}
	tern, ok := prog.Body.(Ternary)
	if !ok {
		t.Fatalf("expected a top-level Ternary, got %T", prog.Body)
	}
	cond, ok := tern.Cond.(Eq)
	if !ok {
		t.Fatalf("expected an Eq condition, got %T", tern.Cond)
	}
	if cond.Left.(State).Value != 1 || cond.Right.(State).Value != 1 {
		t.Errorf("expected EQ(S1, S1), got %+v", cond)
	}
	trueArith, ok := tern.True.(Arith)
	if !ok {
		t.Fatalf("expected the true branch to be an Arith, got %T", tern.True)
	}
	if trueArith.Op != lexer.KindXor4 || trueArith.Left.(State).Value != 1 || trueArith.Right.(State).Value != 3 {
		t.Errorf("expected S1 XOR4 S3, got %+v", trueArith)
	}
	falseState, ok := tern.False.(State)
	if !ok || falseState.Value != 0 {
		t.Fatalf("expected the false branch to be S0, got %+v", tern.False)
	}
}

func TestParseTrailingTokensIsError(t *testing.T) {
	toks, err := lexer.Lex("🌑 🌑;")
	if err != nil {
		t.Fatalf("unexpected lex error: %v", err)
	}
	if _, err := Parse(toks); err == nil {
		t.Fatal("expected a parse error for trailing tokens after a complete expression")
	}
}

// TestParseMissingSemiIsError -- GRAMMAR.md §2's `Program ::= TypedExpr SEMI` (founder
// real-time: "also require semicolons in LO") -- a program with no trailing `;` is a real
// parse error, not silently accepted.
func TestParseMissingSemiIsError(t *testing.T) {
	toks, err := lexer.Lex("🌑")
	if err != nil {
		t.Fatalf("unexpected lex error: %v", err)
	}
	if _, err := Parse(toks); err == nil {
		t.Fatal("expected a parse error for a program with no trailing SEMI")
	}
}

// TestParseWithSemiSucceeds -- the same program as above, but with the now-required trailing
// SEMI, should parse cleanly.
func TestParseWithSemiSucceeds(t *testing.T) {
	toks, err := lexer.Lex("🌑;")
	if err != nil {
		t.Fatalf("unexpected lex error: %v", err)
	}
	prog, err := Parse(toks)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if s, ok := prog.Body.(State); !ok || s.Value != 0 {
		t.Errorf("expected the body to be S0, got %+v", prog.Body)
	}
}

// TestParseLet -- GRAMMAR.md §2's `Let ::= LET Value Expr` (founder real-time: "use ✨ for
// LET"). `✨ S1 🧲` binds S1, then the bare MAGNET (🧲) refers to it.
func TestParseLet(t *testing.T) {
	toks, err := lexer.Lex("✨ 🌒 🧲;")
	if err != nil {
		t.Fatalf("unexpected lex error: %v", err)
	}
	prog, err := Parse(toks)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	let, ok := prog.Body.(Let)
	if !ok {
		t.Fatalf("expected a Let, got %T", prog.Body)
	}
	if s, ok := let.Bound.(State); !ok || s.Value != 1 {
		t.Errorf("expected the bound value to be S1, got %+v", let.Bound)
	}
	if _, ok := let.Body.(LetRef); !ok {
		t.Errorf("expected the body to be a LetRef, got %T", let.Body)
	}
}

// TestParseNestedLetShadows -- real, documented GRAMMAR.md semantics: an inner Let's MAGNET
// refers to the innermost binding, shadowing the outer one entirely.
func TestParseNestedLetShadows(t *testing.T) {
	toks, err := lexer.Lex("✨ 🌒 ✨ 🌓 🧲;")
	if err != nil {
		t.Fatalf("unexpected lex error: %v", err)
	}
	prog, err := Parse(toks)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	outer, ok := prog.Body.(Let)
	if !ok {
		t.Fatalf("expected an outer Let, got %T", prog.Body)
	}
	inner, ok := outer.Body.(Let)
	if !ok {
		t.Fatalf("expected the outer Let's body to be an inner Let, got %T", outer.Body)
	}
	if _, ok := inner.Body.(LetRef); !ok {
		t.Errorf("expected the inner Let's body to be a LetRef, got %T", inner.Body)
	}
}
