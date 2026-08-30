// Package lexer implements LO/GRAMMAR.md §1's lexical grammar: LO Phase 1 (LO/NORTHSTAR.md),
// the real, minimal compiler frontend that (eventually) lexes/parses/emits `.prn` text for the
// scalar-only subset burrow's own two emitters already support.
//
// Real design decision: this lexer works directly on []rune (Unicode code points), not bytes or
// grapheme clusters — GRAMMAR.md §1.2's own matching rule ("strip exactly one trailing variation
// selector, U+FE0E or U+FE0F") operates at the code-point level, and every LO token is a single
// base emoji code point (never a multi-rune grapheme cluster, per §1.2's own reasoning: none of
// LO's tokens are skin-tone-modifiable or ZWJ-joined).
package lexer

// Kind identifies a lexical token class, per GRAMMAR.md §1.1's token table.
type Kind int

const (
	KindInvalid Kind = iota
	KindVecLit       // the literal ASCII "vec PARENA CONSTRUCT 312"
	KindColon        // :
	KindSemi         // ; -- required Program terminator (GRAMMAR.md §1 item 3, §2's Program rule)
	KindState        // 🌑🌒🌓🌔 -- see Token.State for which of 0-3
	KindQuery        // ❓
	KindPlus4        // ➕
	KindMinus4       // ➖
	KindAnd4         // 🔗
	KindOr4          // 🔮
	KindXor4         // 🔀
	KindStack        // 🧱
	KindDot          // 🎯
	KindMatmul       // 🧮
	KindDimlen       // ⚖️
	KindEq           // 🟰 (provisional glyph -- see GRAMMAR.md §1.3)
	KindMatch        // 🔍
	KindScalar       // 💧
	KindVector       // 🌊
	KindMatrix       // 🧊
	KindPattern      // 🕸️
	KindI32          // 🔢
	KindString       // 📜
	KindFunc         // ⚙️
	KindVoid         // 🕳️
	KindUnion        // 🧬
	KindLabel        // 🏷️
	KindDoor         // 🚪
	KindMagnet       // 🧲
	KindWildcard     // 🃏
	KindStar         // 🌌
	KindOnePlus      // ☄️
	KindOpt          // 👻
	KindStart        // 🏁
	KindEnd          // 🛑
	KindAlt          // 🛤️
	KindGroup        // 🗜️ (both open and close -- see GRAMMAR.md §5.4)
)

// Token is one lexed unit. Pos is the rune offset in the source where the token starts, for
// real error reporting once a parser exists (Phase 1's own next real step).
type Token struct {
	Kind  Kind
	State int // valid only when Kind == KindState: 0-3
	Pos   int
}

func (k Kind) String() string {
	switch k {
	case KindVecLit:
		return "VECLIT"
	case KindColon:
		return "COLON"
	case KindSemi:
		return "SEMI"
	case KindState:
		return "STATE"
	case KindQuery:
		return "QUERY"
	case KindPlus4:
		return "PLUS4"
	case KindMinus4:
		return "MINUS4"
	case KindAnd4:
		return "AND4"
	case KindOr4:
		return "OR4"
	case KindXor4:
		return "XOR4"
	case KindStack:
		return "STACK"
	case KindDot:
		return "DOT"
	case KindMatmul:
		return "MATMUL"
	case KindDimlen:
		return "DIMLEN"
	case KindEq:
		return "EQ"
	case KindMatch:
		return "MATCH"
	case KindScalar:
		return "SCALAR"
	case KindVector:
		return "VECTOR"
	case KindMatrix:
		return "MATRIX"
	case KindPattern:
		return "PATTERN"
	case KindI32:
		return "I32"
	case KindString:
		return "STRING"
	case KindFunc:
		return "FUNC"
	case KindVoid:
		return "VOID"
	case KindUnion:
		return "UNION"
	case KindLabel:
		return "LABEL"
	case KindDoor:
		return "DOOR"
	case KindMagnet:
		return "MAGNET"
	case KindWildcard:
		return "WILDCARD"
	case KindStar:
		return "STAR"
	case KindOnePlus:
		return "ONEPLUS"
	case KindOpt:
		return "OPT"
	case KindStart:
		return "START"
	case KindEnd:
		return "END"
	case KindAlt:
		return "ALT"
	case KindGroup:
		return "GROUP"
	default:
		return "INVALID"
	}
}
