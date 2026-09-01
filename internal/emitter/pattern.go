// pattern.go lowers a parsed LO Pattern (internal/parser's own §19-29 AST, see parser.pattern's
// own doc comment) into a real PCRE syntax STRING -- the format LO_Formal_Grammar_Phase_0_
// Complete.md's own §20-29 worked examples show as the pattern's "equivalent PCRE" text, and the
// format PARENA's own real, mature `regex/pcre.prn` (`compile`/`is-match`, confirmed present in
// S222-05's own BACKLOG entry) actually consumes.
//
// Real, honest, deliberately bounded scope for this pass: this file only produces the PCRE TEXT.
// Wiring that text into a real emitted `.prn` call to `regex/pcre/compile`+`is-match` (which
// needs a caller-supplied `Arena @ Region` -- a real, separate design problem, since LO's
// existing "name the defn main" verification trick needs a ZERO-parameter function, and an Arena
// parameter breaks that C entry-point signature) is a real, separate, larger follow-up, not
// attempted here. So is giving LO a real String-typed VALUE for MATCH's own subject (LO has no
// string values in the language at all yet). This file's own real, honest acceptance bar: every
// one of the source doc's own §20-29 worked examples round-trips to the EXACT PCRE text the doc
// itself shows.
package emitter

import (
	"fmt"
	"strings"

	"github.com/emilyspringerton/LO/internal/lexer"
	"github.com/emilyspringerton/LO/internal/parser"
)

// patternToPCRE implements §19/§24's `Pattern ::= PatternSequence | PatternSequence ALT
// PatternSequence+` -- one alternative per member, joined with a literal `|` (PCRE's own
// alternation operator, exactly matching every one of the doc's own §24 worked examples).
func patternToPCRE(pat parser.Pattern) (string, error) {
	parts := make([]string, len(pat.Alternatives))
	for i, seq := range pat.Alternatives {
		s, err := patternSequenceToPCRE(seq)
		if err != nil {
			return "", err
		}
		parts[i] = s
	}
	return strings.Join(parts, "|"), nil
}

func patternSequenceToPCRE(seq parser.PatternSequence) (string, error) {
	var b strings.Builder
	for _, item := range seq {
		s, err := patternItemToPCRE(item)
		if err != nil {
			return "", err
		}
		b.WriteString(s)
	}
	return b.String(), nil
}

// patternItemToPCRE implements §22/§23's Anchor and Quantifier lowering: `^`/`$` for a bare
// Anchor, or an Atom's own PCRE text with `*`/`+`/`?` appended per its Quantifier.
func patternItemToPCRE(item parser.PatternItem) (string, error) {
	if item.IsAnchor {
		if item.AnchorEnd {
			return "$", nil
		}
		return "^", nil
	}
	atomStr, err := patternAtomToPCRE(item.Atom)
	if err != nil {
		return "", err
	}
	switch item.Quant {
	case parser.QuantStar:
		return atomStr + "*", nil
	case parser.QuantPlus:
		return atomStr + "+", nil
	case parser.QuantOpt:
		return atomStr + "?", nil
	default:
		return atomStr, nil
	}
}

func patternAtomToPCRE(atom parser.PatternAtom) (string, error) {
	switch v := atom.(type) {
	case parser.PatternLiteral:
		// §20: "Quoted text provides Unicode literal content" -- the text is meant LITERALLY,
		// not as PCRE syntax, so any of its own characters that happen to be PCRE metacharacters
		// are escaped rather than passed through raw.
		return escapePCRELiteral(v.Text), nil

	case parser.PatternWildcard:
		return ".", nil

	case parser.PatternClass:
		return classToPCRE(v)

	case parser.PatternEscaped:
		return escapedToPCRE(v)

	case parser.PatternGroup:
		inner, err := patternSequenceToPCRE(v.Seq)
		if err != nil {
			return "", err
		}
		return "(" + inner + ")", nil

	default:
		return "", &Error{Msg: fmt.Sprintf("unsupported pattern atom %T", atom)}
	}
}

// classToPCRE implements §25/§26's `CharacterClass ::= CLASS ClassItem+ | NCLASS ClassItem+` and
// `Range ::= LiteralChar RANGE LiteralChar`.
func classToPCRE(c parser.PatternClass) (string, error) {
	var b strings.Builder
	b.WriteByte('[')
	if c.Negated {
		b.WriteByte('^')
	}
	for _, item := range c.Items {
		if item.IsRange {
			b.WriteString(escapePCRELiteral(item.From))
			b.WriteByte('-')
			b.WriteString(escapePCRELiteral(item.To))
		} else {
			b.WriteString(escapePCRELiteral(item.Char))
		}
	}
	b.WriteByte(']')
	return b.String(), nil
}

// escapedToPCRE implements §27's `Escaped ::= ESCAPE Escapable`, mapping each escapable TOKEN
// (not text) to the one literal character it represents when escaped -- only the source doc's
// own two worked examples (escaped STAR -> literal `*`, escaped WILDCARD -> literal `.`) are
// unambiguous derivations; the rest follow the same real pattern (the token's own PCRE meaning,
// backslash-escaped to force it literal). Real, honest, named exception: `ESCAPE GROUP` is NOT
// implemented -- `PatternEscaped` only records that a GROUP token was escaped, not whether the
// source meant the opening or closing paren (the parser has no way to know either, since GROUP
// is the same token for both), so this is a real, genuine ambiguity, not silently guessed at.
func escapedToPCRE(e parser.PatternEscaped) (string, error) {
	switch e.Kind {
	case lexer.KindLiteral:
		return escapePCRELiteral(e.Text), nil
	case lexer.KindWildcard:
		return `\.`, nil
	case lexer.KindStar:
		return `\*`, nil
	case lexer.KindOnePlus:
		return `\+`, nil
	case lexer.KindOpt:
		return `\?`, nil
	case lexer.KindStart:
		return `\^`, nil
	case lexer.KindEnd:
		return `\$`, nil
	case lexer.KindAlt:
		return `\|`, nil
	case lexer.KindEscape:
		return `\\`, nil
	default:
		return "", &Error{Msg: fmt.Sprintf("ESCAPE %s has no unambiguous literal-character mapping in this v0", e.Kind)}
	}
}

// pcreMetachars is every PCRE syntax character that must be backslash-escaped to force it
// literal, per real, standard PCRE metacharacter rules -- not LO's own invention.
const pcreMetachars = `.^$|()[]{}*+?\`

func escapePCRELiteral(s string) string {
	var b strings.Builder
	for _, r := range s {
		if strings.ContainsRune(pcreMetachars, r) {
			b.WriteByte('\\')
		}
		b.WriteRune(r)
	}
	return b.String()
}
