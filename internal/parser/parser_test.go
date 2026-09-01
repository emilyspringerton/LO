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

// TestParseArithChain -- GRAMMAR.md §3.3's own left-associative chain rule, real, found-live
// fix over an earlier draft that only ever accepted one operator: `S1 XOR4 S2 XOR4 S3` must
// parse as `(S1 XOR4 S2) XOR4 S3`, not error out on the second operator.
func TestParseArithChain(t *testing.T) {
	toks, err := lexer.Lex("🌒 🔀 🌓 🔀 🌔;")
	if err != nil {
		t.Fatalf("unexpected lex error: %v", err)
	}
	prog, err := Parse(toks)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	outer, ok := prog.Body.(Arith)
	if !ok {
		t.Fatalf("expected the top-level node to be an Arith, got %T", prog.Body)
	}
	if outer.Right.(State).Value != 3 {
		t.Errorf("expected the outer Arith's right operand to be S3, got %+v", outer.Right)
	}
	inner, ok := outer.Left.(Arith)
	if !ok {
		t.Fatalf("expected left-associativity: outer.Left should itself be an Arith, got %T", outer.Left)
	}
	if inner.Left.(State).Value != 1 || inner.Right.(State).Value != 2 {
		t.Errorf("expected the inner Arith to be S1 XOR4 S2, got %+v", inner)
	}
}

// TestParseArithOverTwoLetRefs -- real, found-live fix: a bare MAGNET (LetRef) on BOTH sides of
// an Arith must parse -- an earlier draft only accepted a bare State as an Arith's right
// operand, so this failed. Real, honest naming caveat, UPDATED: both MAGNETs here are bare
// (Depth 0), so they both still resolve to the SAME (innermost) binding once emitted -- that's
// no longer LetRef's own real limit, though (see TestParseLetRefReachesOuterBinding below):
// LetRef now has a real Depth field, and a `VectorLit MAGNET` form reaches an outer binding.
// This test only proves a bare-MAGNET-on-both-sides Arith parses at all, nothing more.
func TestParseArithOverTwoLetRefs(t *testing.T) {
	toks, err := lexer.Lex("✨ 🌒 ✨ 🌓 🧲 🔀 🧲;")
	if err != nil {
		t.Fatalf("unexpected lex error: %v", err)
	}
	prog, err := Parse(toks)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	outer := prog.Body.(Let)
	inner := outer.Body.(Let)
	arith, ok := inner.Body.(Arith)
	if !ok {
		t.Fatalf("expected the innermost body to be an Arith, got %T", inner.Body)
	}
	if _, ok := arith.Left.(LetRef); !ok {
		t.Errorf("expected the Arith's left operand to be a LetRef, got %T", arith.Left)
	}
	if _, ok := arith.Right.(LetRef); !ok {
		t.Errorf("expected the Arith's right operand to be a LetRef, got %T", arith.Right)
	}
}

// TestParseLetRefReachesOuterBinding -- the real point of the depth-index LetRef extension: a
// `VectorLit MAGNET` (single-state vector immediately before MAGNET) reaches an OUTER Let's
// binding, closing the exact limitation `TestParseArithOverTwoLetRefs`'s own doc comment named.
// `vec PARENA CONSTRUCT 312 S1 🧲` inside the innermost body means Depth 1 (one level out);
// a bare `🧲` still means Depth 0 (innermost).
func TestParseLetRefReachesOuterBinding(t *testing.T) {
	toks, err := lexer.Lex("✨ 🌓 ✨ 🌒 vec PARENA CONSTRUCT 312 🌒 🧲 🔀 🧲;")
	if err != nil {
		t.Fatalf("unexpected lex error: %v", err)
	}
	prog, err := Parse(toks)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	outer := prog.Body.(Let)
	inner := outer.Body.(Let)
	arith, ok := inner.Body.(Arith)
	if !ok {
		t.Fatalf("expected the innermost body to be an Arith, got %T", inner.Body)
	}
	left, ok := arith.Left.(LetRef)
	if !ok || left.Depth != 1 {
		t.Errorf("expected the left operand to be LetRef{Depth: 1} (the outer binding), got %+v", arith.Left)
	}
	right, ok := arith.Right.(LetRef)
	if !ok || right.Depth != 0 {
		t.Errorf("expected the right operand to be LetRef{Depth: 0} (the inner binding), got %+v", arith.Right)
	}
}

// TestParseSwitch -- LO_Formal_Grammar_Phase_0_Complete.md §17/18's own worked example
// (founder real-time: "add switch and case"): `SWITCH S1 (CASE S0 -> S0) (CASE S1 -> S3)
// (CASE S2 -> S1) (DEFAULT -> S3)`.
func TestParseSwitch(t *testing.T) {
	toks, err := lexer.Lex("🔘 🌒 🔹 🌑 🌑 🔹 🌒 🌔 🔹 🌓 🌒 🔸 🌔;")
	if err != nil {
		t.Fatalf("unexpected lex error: %v", err)
	}
	prog, err := Parse(toks)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	sw, ok := prog.Body.(Switch)
	if !ok {
		t.Fatalf("expected a Switch, got %T", prog.Body)
	}
	if sw.Selector.(State).Value != 1 {
		t.Errorf("expected the selector to be S1, got %+v", sw.Selector)
	}
	if len(sw.Cases) != 3 {
		t.Fatalf("expected 3 cases, got %d", len(sw.Cases))
	}
	wantMatches := []int{0, 1, 2}
	wantBodies := []int{0, 3, 1}
	for i, c := range sw.Cases {
		if c.Match != wantMatches[i] {
			t.Errorf("case %d: expected match %d, got %d", i, wantMatches[i], c.Match)
		}
		if c.Body.(State).Value != wantBodies[i] {
			t.Errorf("case %d: expected body S%d, got %+v", i, wantBodies[i], c.Body)
		}
	}
	if sw.Default.(State).Value != 3 {
		t.Errorf("expected the default to be S3, got %+v", sw.Default)
	}
}

// TestParseSwitchRequiresDefault -- real, deliberate v0 restriction named in switch_'s own doc
// comment: a SWITCH with no DEFAULT is a real parse error here, not an implicit VOID fallback.
func TestParseSwitchRequiresDefault(t *testing.T) {
	toks, err := lexer.Lex("🔘 🌒 🔹 🌒 🌔;")
	if err != nil {
		t.Fatalf("unexpected lex error: %v", err)
	}
	if _, err := Parse(toks); err == nil {
		t.Fatal("expected a parse error for a SWITCH with no DEFAULT")
	}
}

// TestParseLambdaAndCall -- LO_Formal_Grammar_Phase_0_Complete.md §15/16 (founder real-time:
// "🐪 LAMBDA", formalized as 💠 in the uploaded grammar doc, plus "add ... LAMBDA"): `CALL
// (LAMBDA x -> x XOR4 S1) S2`.
func TestParseLambdaAndCall(t *testing.T) {
	toks, err := lexer.Lex("📞 💠 🧲 🔀 🌒 🌓;")
	if err != nil {
		t.Fatalf("unexpected lex error: %v", err)
	}
	prog, err := Parse(toks)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	call, ok := prog.Body.(Call)
	if !ok {
		t.Fatalf("expected a Call, got %T", prog.Body)
	}
	lambda, ok := call.Fn.(Lambda)
	if !ok {
		t.Fatalf("expected Call.Fn to be a Lambda, got %T", call.Fn)
	}
	arith, ok := lambda.Body.(Arith)
	if !ok {
		t.Fatalf("expected the Lambda's body to be an Arith, got %T", lambda.Body)
	}
	if _, ok := arith.Left.(LetRef); !ok {
		t.Errorf("expected the Arith's left operand to be a LetRef (the lambda's own parameter), got %T", arith.Left)
	}
	if call.Arg.(State).Value != 2 {
		t.Errorf("expected the Call's own argument to be S2, got %+v", call.Arg)
	}
}
