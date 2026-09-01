// Package parser implements LO/GRAMMAR.md §2's EBNF over the lexer's token stream: Phase 1's
// real next step after the lexer (S214-01). Real, deliberately narrow first slice, matching
// Phase 1's own concrete acceptance test (compile a real LO program to .prn, verify through both
// parena build and burrow build): covers exactly the Door/Ternary/Cond/Arith/State/Void shape
// the repo's own examples/xor_check.llll example needs, not GRAMMAR.md's full §2 grammar yet.
// Real, named, not-yet-covered productions: VectorLit, MagnetExpr, Pattern, LinAlg (STACK/DOT/
// MATMUL/DIMLEN), Labeled — real, separate follow-ups once a real program needs them.
package parser

import "github.com/emilyspringerton/LO/internal/lexer"

// TypeAtom is one of GRAMMAR.md §2's TypeAtom productions.
type TypeAtom int

const (
	TypeInvalid TypeAtom = iota
	TypeScalar
	TypeVector
	TypeMatrix
	TypePattern
	TypeI32
	TypeFloat  // ⚫, LO_Formal_Grammar_Phase_0_Complete.md §7.2 -- see the emitter's own doc
	// comment on why this and TypeDouble both map to the same real PARENA F64 (checked directly
	// against src/emit.c: PARENA has no separate single-precision float type at all).
	TypeDouble // ⚪, same doc §7.3
	TypeString
	TypeFunc
	TypeVoid
)

// Expr is any GRAMMAR.md §2 Value/Ternary/Cond node this narrow v0 covers.
type Expr interface{ isExpr() }

// State is a bare base4 state literal (0-3).
type State struct{ Value int }

func (State) isExpr() {}

// Arith is a binary base4 operator application (GRAMMAR.md §5.1).
type Arith struct {
	Op    lexer.Kind // KindPlus4/KindMinus4/KindAnd4/KindOr4/KindXor4
	Left  Expr
	Right Expr
}

func (Arith) isExpr() {}

// Eq is GRAMMAR.md §5.3's EQ comparison, used only as a Cond in this v0.
type Eq struct {
	Left  Expr
	Right Expr
}

func (Eq) isExpr() {}

// Ternary is GRAMMAR.md §2's core Ternary production: `Cond QUERY Expr COLON Expr`.
type Ternary struct {
	Cond  Expr // an Eq in this v0 (GRAMMAR.md §3.2: Cond is always EQ/MATCH, never a bare Ternary)
	True  Expr
	False Expr
}

func (Ternary) isExpr() {}

// VoidExpr is the bare VOID (🕳️) token used as a value.
type VoidExpr struct{}

func (VoidExpr) isExpr() {}

// Let is GRAMMAR.md §2's `Let ::= LET Value Expr`, added 2026-08-30 (founder real-time: "use ✨
// for LET"). Each `Let` introduces one new binding at the next nesting depth; a `LetRef` inside
// `Body` refers to any enclosing binding by depth (0 = innermost), not just the nearest one --
// see `LetRef`'s own doc comment for the real reason this changed from the original "innermost
// only, no reaching outer" v0.
type Let struct {
	Bound Expr
	Body  Expr
}

func (Let) isExpr() {}

// LetRef is GRAMMAR.md §2's `LetRef ::= MAGNET | VectorLit MAGNET`, extended 2026-08-30 to reach
// OUTER Let bindings, not just the innermost. A bare MAGNET is `Depth: 0` (the nearest enclosing
// Let — unchanged, backward compatible with every existing bare-MAGNET program). `VectorLit
// MAGNET` (a single-state vector immediately followed by MAGNET, e.g. `vec PARENA CONSTRUCT 312
// 🌒 🧲`) sets `Depth` to that state's own value (1, 2, or 3) — real, deliberate reuse of
// `GRAMMAR.md`'s own already-specified `MagnetExpr ::= VectorLit MAGNET Value` token SHAPE
// (index-vector immediately before MAGNET), simplified for this v0's real, narrow scope: a
// single-state vector as a small depth index, not a full row-extraction target. `Depth` counts
// OUTWARD from the innermost active Let at the point this LetRef appears — referencing a depth
// that doesn't exist (fewer than `Depth+1` enclosing Lets) is a real, honest emit-time error, not
// silently clamped or wrapped.
type LetRef struct {
	Depth int
}

func (LetRef) isExpr() {}

// Switch is LO_Formal_Grammar_Phase_0_Complete.md §17's `Switch ::= SWITCH Value Case+ Default?`
// (2026-08-31, founder real-time: "add switch and case"). Real, narrow v0: `Selector` is a
// `primary()` (State/VOID/LetRef), each `Case`'s own match value is a bare `State` (the doc's own
// full `Value` generality for a case label is a real, separate follow-up — every worked example
// in the source doc itself uses a bare state), and `Default` is required in this v0 rather than
// optional (the doc's own §17 "if no case matches and there is no default, the result is VOID" is
// a real, separate fallback-typing question this v0 sidesteps by just requiring one). Lowers to a
// real nested PARENA `if`/`=` chain — the same shape `Ternary` already emits, just chained.
type Switch struct {
	Selector Expr
	Cases    []SwitchCase
	Default  Expr
}

func (Switch) isExpr() {}

// SwitchCase is one `Case ::= CASE Value Expr` arm (its own match `Value` narrowed to a bare
// `State` in this v0 — see Switch's own doc comment).
type SwitchCase struct {
	Match int // 0-3, the base4 state this case matches
	Body  Expr
}

// Lambda is LO_Formal_Grammar_Phase_0_Complete.md §15's `Lambda ::= LAMBDA LambdaParams Expr`
// (2026-08-31, founder real-time: "🐪 LAMBDA", later formalized as 💠 in the uploaded grammar
// doc — see this package's own doc comment on that discrepancy). Real, deliberate simplification
// over the source doc's own `LambdaParams ::= Param+` with `Param ::= LITERAL` (a quoted string
// name): real string-literal lexing (quoted text) doesn't exist in this compiler yet — a real,
// separate, genuinely bigger lexer feature, not attempted here. This v0 instead reuses the
// already-real `Let`/`LetRef` depth-index binding mechanism: a `Lambda` introduces exactly ONE
// parameter at the next nesting depth (same as `Let`), referenced inside `Body` via `LetRef`
// exactly like a `Let`-bound value — the only real difference from `Let` is that a `Lambda`'s
// "bound value" isn't supplied until `Call` applies it, so `Lambda` carries no `Bound` field.
// Real, honest scope: exactly one parameter (the source doc's own multi-param `LambdaParams` is
// a real, separate follow-up).
type Lambda struct {
	Body Expr
}

func (Lambda) isExpr() {}

// Call is LO_Formal_Grammar_Phase_0_Complete.md §16's `Call ::= CALL Value ArgList` (real,
// narrowed to exactly one argument in this v0, matching `Lambda`'s own single-parameter scope —
// the source doc's own multi-argument `ArgList` is a real, separate follow-up).
type Call struct {
	Fn  Expr // must be a Lambda in this v0 -- a real, separate follow-up once LO can hold a
	// function value in a Let/pass one through a chain of calls
	Arg Expr
}

func (Call) isExpr() {}

// Program is GRAMMAR.md §2's top-level `TypedExpr` — an optional Door plus a body Expr.
type Program struct {
	DoorType TypeAtom // TypeInvalid if no Door was present
	HasDoor  bool
	Body     Expr
}
