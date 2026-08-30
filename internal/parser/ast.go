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

// Program is GRAMMAR.md §2's top-level `TypedExpr` — an optional Door plus a body Expr.
type Program struct {
	DoorType TypeAtom // TypeInvalid if no Door was present
	HasDoor  bool
	Body     Expr
}
