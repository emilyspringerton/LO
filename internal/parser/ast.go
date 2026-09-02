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
// over the source doc's own `LambdaParams ::= Param+` with `Param ::= LITERAL` (a quoted-string
// name): this v0 instead reuses the already-real `Let`/`LetRef` depth-index binding mechanism —
// a `Lambda` introduces exactly ONE parameter at the next nesting depth (same as `Let`),
// referenced inside `Body` via `LetRef` exactly like a `Let`-bound value — the only real
// difference from `Let` is that a `Lambda`'s "bound value" isn't supplied until `Call` applies
// it, so `Lambda` carries no `Bound` field. Real quoted-string lexing (LITERAL/StringLit) DOES
// exist now (added 2026-09-01, S222-08) — this is no longer blocked on that. Not revisited since:
// the doc's own `LambdaParams Expr` grammar (one-or-more `Param` LITERALs immediately followed by
// the body `Expr`) is genuinely ambiguous to parse as written — nothing in the grammar says how
// many trailing LITERAL tokens are params versus the start of a LITERAL-valued body, and no
// worked example in the doc disambiguates it either. A real, deliberately DEFERRED decision, not
// silently guessed at either way — multi-param `Lambda`/multi-arg `Call` stay this v0's own
// single-parameter/single-argument depth-index scheme until that ambiguity is resolved.
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

// StringLit is a bare `LITERAL` (🔤"...") used directly as a `Value` — LO's first real
// String-typed value, added 2026-09-01 (founder real-time: "continue"). See the parser's own
// `primary` doc comment for why this is a real, flagged, minimal extension of the source doc's
// own `Value` production (which never lists a bare `Literal`, only `PatternValue`), not a guess
// at something the doc specifies differently.
type StringLit struct{ Text string }

func (StringLit) isExpr() {}

// Match is LO_Formal_Grammar_Phase_0_Complete.md §12/§29's `Cond ::= Value MATCH Pattern`, the
// pattern-matching sibling of `Eq` — used only as a `Ternary`'s own `Cond` in this v0, same as
// `Eq`. Real, honest, narrow scope: `Subject` is parsed as a `Value` exactly like `Eq`'s own
// operands, even though the doc's own real intent is a String subject (LO has no String VALUES
// at all yet — a real, separate, larger follow-up, same one `GRAMMAR.md`'s own 2026-09-01 update
// names). No emitter support exists yet either (`internal/emitter` doesn't handle this node) —
// this v0 is parser-only, real progress toward the doc's own §19-29 Pattern grammar without
// pretending the whole feature is done in one pass.
type Match struct {
	Subject Expr
	Pattern Pattern
}

func (Match) isExpr() {}

// Pattern is LO_Formal_Grammar_Phase_0_Complete.md §19's `Pattern ::= PatternSequence |
// PatternSequence ALT PatternSequence+` — a deliberately bounded PCRE-oriented pattern language,
// used as `Match`'s own right-hand operand. `Alternatives` has exactly one member when the
// source has no top-level `ALT` (§24: alternation binds looser than concatenation).
//
// Real, honest, flagged-not-silently-resolved narrowing versus the source doc's own EBNF:
// `Atom ::= State | Literal | WILDCARD | CharacterClass | Escaped` includes a bare base4 `State`
// as a pattern atom, but every one of the doc's own worked examples (§20-29) uses only
// text-oriented atoms (Literal/Wildcard/CharacterClass/Escaped) — this looks like a real leftover
// from LO's OLDER base4-vector pattern grammar (GRAMMAR.md's own pre-2026-08-31 Pattern/
// PatternAtom, matched against base4 Vectors, a completely different concept from this doc's own
// PCRE-over-text design). `patternAtom` below does NOT parse a bare `State` for this reason —
// named here for the founder's own awareness, not decided unilaterally either way.
type Pattern struct {
	Alternatives []PatternSequence
}

func (Pattern) isExpr() {}

// PatternSequence is §19's `PatternSequence ::= PatternItem+` — one concatenated run of pattern
// items (anchors, atoms, quantified atoms, or groups).
type PatternSequence []PatternItem

// PatternItem is one element of a PatternSequence: either a bare Anchor, or an Atom (optionally
// followed by a Quantifier) — extended, per §28's own worked example (`🗜️ 🔤"ab" 🗜️ ☄️`, a
// quantified capture group), so `Atom` may itself be a `PatternGroup`, which the doc's own EBNF
// `PatternItem ::= Anchor | Atom Quantifier? | PatternGroup` technically lists as a THIRD, its
// own unquantifiable alternative — a real, minor grammar refinement made here because the doc's
// own worked example directly contradicts its own EBNF otherwise (flagged, not silently patched).
type PatternItem struct {
	IsAnchor  bool // when true, Atom/Quant are unused; see AnchorEnd
	AnchorEnd bool // valid when IsAnchor: false = START (^), true = END ($)
	Atom      PatternAtom
	Quant     QuantKind
}

// QuantKind is one of §22's three v1 quantifiers, or none.
type QuantKind int

const (
	QuantNone QuantKind = iota
	QuantStar          // 🌌 *
	QuantPlus          // ☄️ +
	QuantOpt           // 👻 ?
)

// PatternAtom is one atomic pattern element: a Literal, WILDCARD, a CharacterClass, an Escaped
// token, or a parenthesized PatternGroup (grouped so it can carry its own Quantifier, per
// PatternItem's own doc comment above).
type PatternAtom interface{ isPatternAtom() }

// PatternLiteral is §20's `Literal ::= LITERAL QuotedText` (e.g. `🔤"cat"`).
type PatternLiteral struct{ Text string }

func (PatternLiteral) isPatternAtom() {}

// PatternWildcard is §21's bare WILDCARD (🃏), matching any one character.
type PatternWildcard struct{}

func (PatternWildcard) isPatternAtom() {}

// PatternClass is §25/§26's `CharacterClass ::= CLASS ClassItem+ | NCLASS ClassItem+`.
type PatternClass struct {
	Negated bool // true for NCLASS (🚫), false for CLASS (🅰️)
	Items   []ClassItem
}

func (PatternClass) isPatternAtom() {}

// ClassItem is one CLASS/NCLASS member per §25/§26: a single literal character, or an inclusive
// a-to-z Range. Real, narrow v0: every `LiteralChar` here must be a single-rune LITERAL — the
// source doc's own worked examples never show a multi-character LITERAL inside a class, and
// nothing else in the doc gives multi-char class members a defined meaning.
type ClassItem struct {
	IsRange bool
	Char    string // valid when !IsRange: the single literal character
	From    string // valid when IsRange
	To      string // valid when IsRange
}

// PatternEscaped is §27's `Escaped ::= ESCAPE Escapable`. Real, deliberate design: rather than
// deciding HERE which literal character each escaped structural/quantifier/anchor token maps to
// (STAR -> "*", WILDCARD -> ".", etc.), this just records which token was escaped (plus its own
// Text, when it was itself a LITERAL) and defers the actual character mapping to the emitter —
// the parser's job is shape, not PCRE-string lowering policy.
type PatternEscaped struct {
	Kind lexer.Kind // one of KindLiteral/KindWildcard/KindStar/KindOnePlus/KindOpt/KindStart/
	// KindEnd/KindAlt/KindGroup/KindEscape, per §27's own Escapable production
	Text string // valid only when Kind == KindLiteral
}

func (PatternEscaped) isPatternAtom() {}

// PatternGroup is §28's `PatternGroup ::= GROUP PatternSequence GROUP` ("(ab)"). Held as a real
// PatternAtom (see PatternItem's own doc comment) so it can carry its own Quantifier, matching
// the doc's own worked quantified-group example. Real, honest, NOT-yet-enforced boundary: §33
// lists nested capture groups as out of v1 scope, but nothing in `patternSequence`/`patternAtom`
// below actually rejects a GROUP nested inside another GROUP — it parses structurally fine today.
// Flagged as a real follow-up (add the depth check) rather than silently enforced or ignored.
type PatternGroup struct{ Seq PatternSequence }

func (PatternGroup) isPatternAtom() {}

// Program is GRAMMAR.md §2's top-level `TypedExpr` — an optional Door plus a body Expr.
type Program struct {
	DoorType TypeAtom // TypeInvalid if no Door was present
	HasDoor  bool
	Body     Expr
}
