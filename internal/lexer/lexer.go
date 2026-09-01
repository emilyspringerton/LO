package lexer

import (
	"fmt"
)

// Variation selectors GRAMMAR.md §1.2 strips before token lookup.
const (
	vs15 = '︎'
	vs16 = '️'
)

// stateTable maps a base state emoji to its 0-3 value (GRAMMAR.md §1.1).
var stateTable = map[rune]int{
	0x1F311: 0, // 🌑
	0x1F312: 1, // 🌒
	0x1F313: 2, // 🌓
	0x1F314: 3, // 🌔
}

// tokenTable maps every non-state, non-VECLIT, non-colon base emoji to its Kind, per
// GRAMMAR.md §1.1's exact codepoint table.
var tokenTable = map[rune]Kind{
	0x2753:  KindQuery,    // ❓
	0x2795:  KindPlus4,    // ➕
	0x2796:  KindMinus4,   // ➖
	0x1F517: KindAnd4,     // 🔗
	0x1F52E: KindOr4,      // 🔮
	0x1F500: KindXor4,     // 🔀
	0x1F9F1: KindStack,    // 🧱
	0x1F3AF: KindDot,      // 🎯
	0x1F9EE: KindMatmul,   // 🧮
	0x2696:  KindDimlen,   // ⚖
	0x2693:  KindEq,       // ⚓ (founder-confirmed 2026-08-30, GRAMMAR.md §1.3)
	0x2728:  KindLet,      // ✨
	0x1F4A0: KindLambda,   // 💠
	0x1F4DE: KindCall,     // 📞
	0x1F518: KindSwitch,   // 🔘
	0x1F539: KindCase,     // 🔹
	0x1F538: KindDefault,  // 🔸
	0x1F50D: KindMatch,    // 🔍
	0x1F4A7: KindScalar,   // 💧
	0x1F30A: KindVector,   // 🌊
	0x1F9CA: KindMatrix,   // 🧊
	0x1F578: KindPattern,  // 🕸
	0x1F522: KindI32,      // 🔢
	0x26AB:  KindFloatType,  // ⚫
	0x26AA:  KindDoubleType, // ⚪
	0x1F4DC: KindString,   // 📜
	0x2699:  KindFunc,     // ⚙
	0x1F573: KindVoid,     // 🕳
	0x1F9EC: KindUnion,    // 🧬
	0x1F3F7: KindLabel,    // 🏷
	0x1F6AA: KindDoor,     // 🚪
	0x1F9F2: KindMagnet,   // 🧲
	0x1F0CF: KindWildcard, // 🃏
	0x1F30C: KindStar,     // 🌌
	0x2604:  KindOnePlus,  // ☄
	0x1F47B: KindOpt,      // 👻
	0x1F3C1: KindStart,    // 🏁
	0x1F6D1: KindEnd,      // 🛑
	0x1F6E4: KindAlt,      // 🛤
	0x1F5DC: KindGroup,    // 🗜
}

const vecLit = "vec PARENA CONSTRUCT 312"

// Error is a lex-time failure, per GRAMMAR.md §2's "any other character results in a fatal
// syntax error" and §1.2's ZWJ/combining-mark rejection.
type Error struct {
	Pos int
	Msg string
}

func (e *Error) Error() string { return fmt.Sprintf("lex error at rune %d: %s", e.Pos, e.Msg) }

// Lex tokenizes LO source per GRAMMAR.md §1. Whitespace between tokens is skipped (§1's own
// "not significant except inside the vector keyword" rule); the vector keyword itself is
// matched as one exact literal run, never split into word tokens.
func Lex(src string) ([]Token, error) {
	runes := []rune(src)
	var toks []Token
	i := 0
	for i < len(runes) {
		r := runes[i]

		if isSpace(r) {
			i++
			continue
		}

		if r == ':' {
			toks = append(toks, Token{Kind: KindColon, Pos: i})
			i++
			continue
		}

		if r == ';' {
			toks = append(toks, Token{Kind: KindSemi, Pos: i})
			i++
			continue
		}

		if matchVecLit(runes, i) {
			toks = append(toks, Token{Kind: KindVecLit, Pos: i})
			i += len([]rune(vecLit))
			continue
		}

		base, consumed, err := lexEmojiToken(runes, i)
		if err != nil {
			return nil, err
		}

		if state, ok := stateTable[base]; ok {
			toks = append(toks, Token{Kind: KindState, State: state, Pos: i})
		} else if kind, ok := tokenTable[base]; ok {
			toks = append(toks, Token{Kind: kind, Pos: i})
		} else {
			return nil, &Error{Pos: i, Msg: fmt.Sprintf("unrecognized codepoint U+%04X", base)}
		}
		i += consumed
	}
	return toks, nil
}

// lexEmojiToken reads one token's base codepoint starting at runes[i], applying GRAMMAR.md
// §1.2's matching rule: strip exactly one trailing VS15/VS16, reject a ZWJ or any other
// combining mark immediately after. Returns the base codepoint and how many runes were consumed.
func lexEmojiToken(runes []rune, i int) (rune, int, error) {
	base := runes[i]
	consumed := 1

	if i+1 < len(runes) {
		next := runes[i+1]
		if next == '‍' { // ZWJ -- GRAMMAR.md §1.2: "a fatal lex error, not silently absorbed"
			return 0, 0, &Error{Pos: i, Msg: "ZWJ immediately following a token codepoint is not permitted (no compound-emoji tokens are defined)"}
		}
		if next == vs15 || next == vs16 {
			consumed = 2
			// A second, unlisted combining mark right after is still a fatal error -- checked
			// by simply not consuming it here, so the next loop iteration hits the fallthrough
			// "unrecognized codepoint" case naturally if it isn't whitespace/a real token start.
		}
	}
	return base, consumed, nil
}

func matchVecLit(runes []rune, i int) bool {
	lit := []rune(vecLit)
	if i+len(lit) > len(runes) {
		return false
	}
	for j, r := range lit {
		if runes[i+j] != r {
			return false
		}
	}
	return true
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}
