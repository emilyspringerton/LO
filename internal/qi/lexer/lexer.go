package lexer

import (
	"fmt"
	"strconv"
	"strings"
)

// Error is a lex-time failure.
type Error struct {
	Pos int
	Msg string
}

func (e *Error) Error() string { return fmt.Sprintf("qi lex error at rune %d: %s", e.Pos, e.Msg) }

// delimiters are the characters that always end a bare symbol or integer literal, even with no
// intervening whitespace -- standard Lisp-reader behavior (`(let[x 1])` tokenizes the same as
// `(let [x 1])`).
const delimiters = "()[] \t\n\r;\""

// Lex tokenizes `qi` source per QI_NORTHSTAR.md's own first-draft surface grammar. Whitespace
// between tokens is insignificant; a `;` starts a line comment running to the next newline (or
// end of input); anything else that isn't a paren/bracket/string/int-literal start is read as one
// bare KindSymbol token, stopping at the next delimiter.
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

		if r == ';' {
			for i < len(runes) && runes[i] != '\n' {
				i++
			}
			continue
		}

		switch r {
		case '(':
			toks = append(toks, Token{Kind: KindLParen, Pos: i})
			i++
			continue
		case ')':
			toks = append(toks, Token{Kind: KindRParen, Pos: i})
			i++
			continue
		case '[':
			toks = append(toks, Token{Kind: KindLBracket, Pos: i})
			i++
			continue
		case ']':
			toks = append(toks, Token{Kind: KindRBracket, Pos: i})
			i++
			continue
		case '"':
			text, consumed, err := lexQuotedText(runes, i)
			if err != nil {
				return nil, err
			}
			toks = append(toks, Token{Kind: KindString, Text: text, Pos: i})
			i += consumed
			continue
		}

		if isDigit(r) {
			tok, consumed, err := lexInt(runes, i)
			if err != nil {
				return nil, err
			}
			toks = append(toks, tok)
			i += consumed
			continue
		}

		// Real, deliberate permissiveness, matching how real Lisp readers treat symbols: any
		// character that isn't a delimiter (paren/bracket/whitespace/`;`/`"`) is valid inside a
		// bare symbol, with no restricted charset invented here -- QI_NORTHSTAR.md's own grammar
		// doesn't specify one, and every real symbol spelling it names (`let`, `+4`, `foo-bar`,
		// `x?`) already needs punctuation beyond plain letters/digits. This can never fail: the
		// loop below always consumes at least one rune, since every delimiter character was
		// already handled by an earlier branch above.
		start := i
		for i < len(runes) && !strings.ContainsRune(delimiters, runes[i]) {
			i++
		}
		toks = append(toks, Token{Kind: KindSymbol, Text: string(runes[start:i]), Pos: start})
	}
	return toks, nil
}

// lexQuotedText implements the same real, resolved quoting rule
// LO_Formal_Grammar_Phase_0_Complete.md's own LITERAL already settled (see
// LO/internal/lexer/lexer.go's own lexQuotedText doc comment): no backslash-escape processing
// inside the quotes, a bare `"` can't appear inside one, and an unterminated quote is a fatal lex
// error. Returns the text and how many runes were consumed starting at `start` (the opening
// quote itself).
func lexQuotedText(runes []rune, start int) (string, int, error) {
	var sb strings.Builder
	i := start + 1
	for i < len(runes) {
		if runes[i] == '"' {
			return sb.String(), i + 1 - start, nil
		}
		sb.WriteRune(runes[i])
		i++
	}
	return "", 0, &Error{Pos: start, Msg: "unterminated quoted string literal (missing closing '\"')"}
}

// lexInt reads a maximal run of decimal digits starting at `start`, parses it as a base-10
// integer, and range-checks it against LO's own four base4 states (0-3) -- a real, honest parse
// error otherwise, not silently wrapped or truncated. Real, deliberate simplification: leading
// zeros are accepted (`"03"` reads as the in-range value 3) since QI_NORTHSTAR.md's own grammar
// never gives them separate meaning.
func lexInt(runes []rune, start int) (Token, int, error) {
	i := start
	for i < len(runes) && isDigit(runes[i]) {
		i++
	}
	text := string(runes[start:i])
	value, err := strconv.Atoi(text)
	if err != nil {
		// Real, honest defensive case -- isDigit already guarantees text is all decimal digits,
		// so this can only fire on an implausibly long literal overflowing a real Go int.
		return Token{}, 0, &Error{Pos: start, Msg: fmt.Sprintf("invalid integer literal %q: %v", text, err)}
	}
	if value < 0 || value > 3 {
		return Token{}, 0, &Error{Pos: start, Msg: fmt.Sprintf("integer literal %d is out of range (LO's own base4 states are only 0-3)", value)}
	}
	return Token{Kind: KindInt, Int: value, Pos: start}, i - start, nil
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}
