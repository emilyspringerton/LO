package parser

import (
	"fmt"

	"github.com/emilyspringerton/LO/internal/lexer"
)

// Error is a parse-time failure.
type Error struct {
	Pos int
	Msg string
}

func (e *Error) Error() string { return fmt.Sprintf("parse error at token position %d: %s", e.Pos, e.Msg) }

type parser struct {
	toks []lexer.Token
	pos  int
}

// Parse implements GRAMMAR.md §2's `Program ::= TypedExpr` over this v0's covered subset (see
// this package's own doc comment in ast.go for exactly what's covered).
func Parse(toks []lexer.Token) (*Program, error) {
	p := &parser{toks: toks}
	prog := &Program{}

	if p.peek().Kind == lexer.KindDoor {
		p.next()
		t, ok := p.typeAtom()
		if !ok {
			return nil, p.errf("expected a TypeAtom after DOOR")
		}
		prog.HasDoor = true
		prog.DoorType = t
	}

	expr, err := p.expr()
	if err != nil {
		return nil, err
	}
	prog.Body = expr

	// GRAMMAR.md §2's `Program ::= TypedExpr SEMI` -- founder real-time, 2026-08-30: "also
	// require semicolons in LO." A required trailing terminator, not optional.
	if p.peek().Kind != lexer.KindSemi {
		return nil, p.errf("expected a trailing SEMI (;) to terminate the program")
	}
	p.next()

	if p.pos != len(p.toks) {
		return nil, p.errf("unexpected trailing tokens after the terminating SEMI")
	}
	return prog, nil
}

func (p *parser) typeAtom() (TypeAtom, bool) {
	tok := p.peek()
	switch tok.Kind {
	case lexer.KindScalar:
		p.next()
		return TypeScalar, true
	case lexer.KindVector:
		p.next()
		return TypeVector, true
	case lexer.KindMatrix:
		p.next()
		return TypeMatrix, true
	case lexer.KindPattern:
		p.next()
		return TypePattern, true
	case lexer.KindI32:
		p.next()
		return TypeI32, true
	case lexer.KindFloatType:
		p.next()
		return TypeFloat, true
	case lexer.KindDoubleType:
		p.next()
		return TypeDouble, true
	case lexer.KindString:
		p.next()
		return TypeString, true
	case lexer.KindFunc:
		p.next()
		return TypeFunc, true
	case lexer.KindVoid:
		p.next()
		return TypeVoid, true
	default:
		return TypeInvalid, false
	}
}

// expr implements GRAMMAR.md §2's `Expr ::= Ternary | Let`.
func (p *parser) expr() (Expr, error) {
	if p.peek().Kind == lexer.KindLet {
		return p.let_()
	}
	if p.peek().Kind == lexer.KindSwitch {
		return p.switch_()
	}
	return p.ternary()
}

// switch_ implements LO_Formal_Grammar_Phase_0_Complete.md §17's
// `Switch ::= SWITCH Value Case+ Default?` (founder real-time: "add switch and case"). Real,
// narrow v0 differences from the source doc, named rather than silently assumed: each `Case`'s
// own match value is a bare `State` (0-3), not the doc's own full `Value` generality -- every
// worked example in the source doc itself uses a bare state; and `Default` is REQUIRED here
// rather than optional, sidestepping the doc's own separate "no case matches, no default ->
// VOID" fallback-typing question for this v0 (this compiler's scalar-I32-only emitter has no
// real way to produce a VOID result anyway -- see the emitter's own existing VoidExpr handling).
func (p *parser) switch_() (Expr, error) {
	p.next() // consume SWITCH
	selector, err := p.primary()
	if err != nil {
		return nil, err
	}

	var cases []SwitchCase
	for p.peek().Kind == lexer.KindCase {
		p.next()
		matchTok := p.peek()
		if matchTok.Kind != lexer.KindState {
			return nil, p.errf("expected a base4 state as a CASE's own match value, got %s", matchTok.Kind)
		}
		p.next()
		body, err := p.expr()
		if err != nil {
			return nil, err
		}
		cases = append(cases, SwitchCase{Match: matchTok.State, Body: body})
	}
	if len(cases) == 0 {
		return nil, p.errf("expected at least one CASE after SWITCH's own selector")
	}

	if p.peek().Kind != lexer.KindDefault {
		return nil, p.errf("expected a required DEFAULT after SWITCH's own CASE clauses (this v0 doesn't support an implicit VOID fallback)")
	}
	p.next()
	defaultExpr, err := p.expr()
	if err != nil {
		return nil, err
	}

	return Switch{Selector: selector, Cases: cases, Default: defaultExpr}, nil
}

// let_ implements GRAMMAR.md §2's `Let ::= LET Value Expr` (added 2026-08-30, founder real-time:
// "use ✨ for LET"). Named `let_` (trailing underscore) since `let` collides with no Go keyword
// but reads awkwardly as a method name otherwise.
func (p *parser) let_() (Expr, error) {
	p.next() // consume LET
	bound, err := p.value()
	if err != nil {
		return nil, err
	}
	body, err := p.expr()
	if err != nil {
		return nil, err
	}
	return Let{Bound: bound, Body: body}, nil
}

// ternary implements GRAMMAR.md §2's `Ternary ::= Cond QUERY Expr COLON Expr | Value` and §3.2's
// rule that Cond is always an EQ/MATCH application, never a bare Ternary — parsed here by first
// parsing a Value, then checking for a following EQ or MATCH to decide whether this is a Cond at
// all. MATCH added per LO_Formal_Grammar_Phase_0_Complete.md §12's `Cond ::= Value EQ Value |
// Value MATCH Pattern` (real, honest, parser-only for now — see Match's own doc comment in
// ast.go for why no emitter support exists yet).
func (p *parser) ternary() (Expr, error) {
	left, err := p.value()
	if err != nil {
		return nil, err
	}

	var cond Expr
	switch p.peek().Kind {
	case lexer.KindEq:
		p.next()
		right, err := p.value()
		if err != nil {
			return nil, err
		}
		cond = Eq{Left: left, Right: right}
	case lexer.KindMatch:
		p.next()
		pat, err := p.pattern()
		if err != nil {
			return nil, err
		}
		cond = Match{Subject: left, Pattern: pat}
	default:
		return left, nil
	}

	if p.peek().Kind != lexer.KindQuery {
		return nil, p.errf("expected QUERY after an EQ/MATCH condition")
	}
	p.next()
	trueExpr, err := p.expr()
	if err != nil {
		return nil, err
	}
	if p.peek().Kind != lexer.KindColon {
		return nil, p.errf("expected COLON after a ternary's true branch")
	}
	p.next()
	falseExpr, err := p.expr()
	if err != nil {
		return nil, err
	}
	return Ternary{Cond: cond, True: trueExpr, False: falseExpr}, nil
}

// value implements GRAMMAR.md §2's `Value` for this v0's covered cases: a real, left-associative
// chain of `ArithOp` applications over `primary()` operands (GRAMMAR.md §3.3: "Arith and LinAlg
// operators are all left-associative, single precedence tier"). Real, found-live fix over an
// earlier draft: that version only ever accepted ONE operator with a bare `State` on both sides,
// so neither a chain (`a XOR4 b XOR4 c`) nor combining two `Let`-bound values (`x XOR4 y`, both
// `LetRef`s) could parse — both are real, ordinary uses of GRAMMAR.md's own `Value ArithOp Value`
// production, not exotic cases.
func (p *parser) value() (Expr, error) {
	left, err := p.primary()
	if err != nil {
		return nil, err
	}
	for isArithOp(p.peek().Kind) {
		op := p.next().Kind
		right, err := p.primary()
		if err != nil {
			return nil, err
		}
		left = Arith{Op: op, Left: left, Right: right}
	}
	return left, nil
}

// primary is one ArithOp operand: a bare State, VOID, or a LetRef (bare MAGNET). GRAMMAR.md's
// own `VectorLit`/`MagnetExpr`/`LinAlg`/`Labeled` productions are real, separate, not-yet-covered
// operand shapes — a real, named follow-up, not silently assumed unreachable.
func (p *parser) primary() (Expr, error) {
	tok := p.peek()

	if tok.Kind == lexer.KindVoid {
		p.next()
		return VoidExpr{}, nil
	}

	// GRAMMAR.md §2's `LetRef ::= MAGNET` (bare, no operands): the nearest (Depth 0) enclosing
	// Let's own bound value.
	if tok.Kind == lexer.KindMagnet {
		p.next()
		return LetRef{Depth: 0}, nil
	}

	// GRAMMAR.md §2's `LetRef ::= VectorLit MAGNET`, extended 2026-08-30 to reach OUTER Let
	// bindings: a single-state VectorLit immediately followed by MAGNET sets Depth to that
	// state's value. Real, narrow v0: exactly one state -- a multi-state VectorLit here (or one
	// not immediately followed by MAGNET) is `MagnetExpr`'s own real row-extraction shape,
	// genuinely not yet implemented, so this returns a real, honest parse error rather than
	// guessing at a different production.
	if tok.Kind == lexer.KindVecLit {
		p.next()
		stateTok := p.peek()
		if stateTok.Kind != lexer.KindState {
			return nil, p.errf("expected exactly one base4 state after the vector keyword (a depth-index LetRef), got %s", stateTok.Kind)
		}
		p.next()
		if p.peek().Kind != lexer.KindMagnet {
			return nil, p.errf("expected MAGNET immediately after a single-state vector (MagnetExpr's own multi-state row-extraction shape is not yet implemented), got %s", p.peek().Kind)
		}
		p.next()
		return LetRef{Depth: stateTok.State}, nil
	}

	// LO_Formal_Grammar_Phase_0_Complete.md §15's `Lambda ::= LAMBDA LambdaParams Expr`, real,
	// narrow v0: exactly one implicit parameter (see the Lambda AST node's own doc comment for
	// why this doesn't parse the source doc's own quoted-string `LambdaParams`).
	if tok.Kind == lexer.KindLambda {
		p.next()
		body, err := p.expr()
		if err != nil {
			return nil, err
		}
		return Lambda{Body: body}, nil
	}

	// LO_Formal_Grammar_Phase_0_Complete.md §16's `Call ::= CALL Value ArgList`, real, narrow v0:
	// exactly one argument (see the Call AST node's own doc comment).
	if tok.Kind == lexer.KindCall {
		p.next()
		fn, err := p.primary()
		if err != nil {
			return nil, err
		}
		arg, err := p.primary()
		if err != nil {
			return nil, err
		}
		return Call{Fn: fn, Arg: arg}, nil
	}

	// StringLit -- a bare LITERAL used directly as a Value, not just inside a Pattern. Real,
	// flagged-not-silently-invented extension: LO_Formal_Grammar_Phase_0_Complete.md §4's own
	// `Value` production list never actually includes a bare `Literal` (only `PatternValue`,
	// itself just `Pattern`) — every one of the doc's own real STRING-Door examples (§5, e.g.
	// `🚪 📜 value ;`) writes the placeholder ENGLISH WORD "value", never real LO source for one.
	// Adding this is the minimal, obvious reading of what a STRING Door's own body must be able
	// to produce (LO otherwise has NO way to write a string value at all), not a guess at
	// something the doc actively specifies differently.
	if tok.Kind == lexer.KindLiteral {
		p.next()
		return StringLit{Text: tok.Text}, nil
	}

	if tok.Kind != lexer.KindState {
		return nil, p.errf("expected a base4 state, VOID, MAGNET, a depth-index vector, LAMBDA, CALL, or LITERAL, got %s", tok.Kind)
	}
	p.next()
	return State{Value: tok.State}, nil
}

// pattern implements LO_Formal_Grammar_Phase_0_Complete.md §19/§24's
// `Pattern ::= PatternSequence | PatternSequence ALT PatternSequence+`.
func (p *parser) pattern() (Pattern, error) {
	first, err := p.patternSequence()
	if err != nil {
		return Pattern{}, err
	}
	alts := []PatternSequence{first}
	for p.peek().Kind == lexer.KindAlt {
		p.next()
		seq, err := p.patternSequenceIn(false)
		if err != nil {
			return Pattern{}, err
		}
		alts = append(alts, seq)
	}
	return Pattern{Alternatives: alts}, nil
}

func (p *parser) patternSequence() (PatternSequence, error) { return p.patternSequenceIn(false) }

// patternSequenceIn implements §19's `PatternSequence ::= PatternItem+`. Real, honest stopping
// rule the source doc never states explicitly: greedy, ending at the first token that can't
// start a PatternItem (ALT, a GROUP's own closer, or the end of input). `insideGroup` is true
// only when parsing a PatternGroup's own body -- see isPatternItemStart's own doc comment for
// why GROUP itself needs this to disambiguate "open a nested group" from "this is my closer".
func (p *parser) patternSequenceIn(insideGroup bool) (PatternSequence, error) {
	var items PatternSequence
	for isPatternItemStart(p.peek().Kind, insideGroup) {
		item, err := p.patternItem(insideGroup)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if len(items) == 0 {
		return nil, p.errf("expected at least one pattern item")
	}
	return items, nil
}

// isPatternItemStart reports whether k can start a new PatternItem. GROUP is context-sensitive:
// nested capture groups are explicitly out of v1 scope (doc §33), so while already inside a
// group's own body (insideGroup == true) a GROUP token is always read as THIS group's own
// closing delimiter, never as opening a new nested one -- outside a group, it opens one.
func isPatternItemStart(k lexer.Kind, insideGroup bool) bool {
	switch k {
	case lexer.KindStart, lexer.KindEnd, lexer.KindLiteral, lexer.KindWildcard,
		lexer.KindClass, lexer.KindNClass, lexer.KindEscape:
		return true
	case lexer.KindGroup:
		return !insideGroup
	default:
		return false
	}
}

// patternItem implements §19's `PatternItem ::= Anchor | Atom Quantifier? | PatternGroup` — see
// PatternItem's own ast.go doc comment for the real, minor refinement (a PatternGroup may itself
// carry a trailing Quantifier, per §28's own worked example).
func (p *parser) patternItem(insideGroup bool) (PatternItem, error) {
	switch p.peek().Kind {
	case lexer.KindStart:
		p.next()
		return PatternItem{IsAnchor: true, AnchorEnd: false}, nil
	case lexer.KindEnd:
		p.next()
		return PatternItem{IsAnchor: true, AnchorEnd: true}, nil
	case lexer.KindGroup:
		p.next() // consume the opening GROUP
		seq, err := p.patternSequenceIn(true)
		if err != nil {
			return PatternItem{}, err
		}
		if p.peek().Kind != lexer.KindGroup {
			return PatternItem{}, p.errf("expected a closing GROUP (🗜️) to end this pattern group")
		}
		p.next() // consume the closing GROUP
		return PatternItem{Atom: PatternGroup{Seq: seq}, Quant: p.patternQuantifier()}, nil
	default:
		atom, err := p.patternAtom()
		if err != nil {
			return PatternItem{}, err
		}
		return PatternItem{Atom: atom, Quant: p.patternQuantifier()}, nil
	}
}

func (p *parser) patternQuantifier() QuantKind {
	switch p.peek().Kind {
	case lexer.KindStar:
		p.next()
		return QuantStar
	case lexer.KindOnePlus:
		p.next()
		return QuantPlus
	case lexer.KindOpt:
		p.next()
		return QuantOpt
	default:
		return QuantNone
	}
}

// patternAtom implements §19's `Atom ::= State | Literal | WILDCARD | CharacterClass | Escaped`
// — real, deliberate narrowing: the `State` alternative is not parsed here. See Pattern's own
// ast.go doc comment for why (a likely leftover from LO's older base4-vector pattern grammar,
// flagged rather than silently either implemented or dropped).
func (p *parser) patternAtom() (PatternAtom, error) {
	tok := p.peek()
	switch tok.Kind {
	case lexer.KindLiteral:
		p.next()
		return PatternLiteral{Text: tok.Text}, nil
	case lexer.KindWildcard:
		p.next()
		return PatternWildcard{}, nil
	case lexer.KindClass, lexer.KindNClass:
		return p.patternClass()
	case lexer.KindEscape:
		return p.patternEscaped()
	default:
		return nil, p.errf("expected a pattern atom (LITERAL, WILDCARD, CLASS/NCLASS, or ESCAPE), got %s", tok.Kind)
	}
}

// patternClass implements §25's `CharacterClass ::= CLASS ClassItem+ | NCLASS ClassItem+` and
// §26's `Range ::= LiteralChar RANGE LiteralChar`. Real, narrow v0: a `LiteralChar` here must be
// a single-rune LITERAL — the source doc's own worked examples never show a multi-char LITERAL
// inside a class, so a longer one is a real, honest parse error rather than silently truncated.
func (p *parser) patternClass() (PatternAtom, error) {
	negated := p.peek().Kind == lexer.KindNClass
	p.next() // consume CLASS or NCLASS

	var items []ClassItem
	for p.peek().Kind == lexer.KindLiteral {
		lit := p.next()
		ch, err := p.singleRune(lit.Text)
		if err != nil {
			return nil, err
		}
		if p.peek().Kind == lexer.KindRange {
			p.next()
			toTok := p.peek()
			if toTok.Kind != lexer.KindLiteral {
				return nil, p.errf("expected a LITERAL after RANGE (↔️), got %s", toTok.Kind)
			}
			p.next()
			toCh, err := p.singleRune(toTok.Text)
			if err != nil {
				return nil, err
			}
			items = append(items, ClassItem{IsRange: true, From: ch, To: toCh})
		} else {
			items = append(items, ClassItem{Char: ch})
		}
	}
	if len(items) == 0 {
		return nil, p.errf("expected at least one class member after CLASS/NCLASS")
	}
	return PatternClass{Negated: negated, Items: items}, nil
}

func (p *parser) singleRune(s string) (string, error) {
	if len([]rune(s)) != 1 {
		return "", p.errf("expected a single-character LITERAL, got %q", s)
	}
	return s, nil
}

// patternEscaped implements §27's `Escaped ::= ESCAPE Escapable`.
func (p *parser) patternEscaped() (PatternAtom, error) {
	p.next() // consume ESCAPE
	tok := p.peek()
	switch tok.Kind {
	case lexer.KindLiteral, lexer.KindWildcard, lexer.KindStar, lexer.KindOnePlus, lexer.KindOpt,
		lexer.KindStart, lexer.KindEnd, lexer.KindAlt, lexer.KindGroup, lexer.KindEscape:
		p.next()
		return PatternEscaped{Kind: tok.Kind, Text: tok.Text}, nil
	default:
		return nil, p.errf("expected an Escapable token after ESCAPE (🛡️), got %s", tok.Kind)
	}
}

func isArithOp(k lexer.Kind) bool {
	switch k {
	case lexer.KindPlus4, lexer.KindMinus4, lexer.KindAnd4, lexer.KindOr4, lexer.KindXor4:
		return true
	default:
		return false
	}
}

func (p *parser) peek() lexer.Token {
	if p.pos >= len(p.toks) {
		return lexer.Token{Kind: lexer.KindInvalid, Pos: -1}
	}
	return p.toks[p.pos]
}

func (p *parser) next() lexer.Token {
	tok := p.peek()
	p.pos++
	return tok
}

func (p *parser) errf(format string, args ...any) error {
	return &Error{Pos: p.pos, Msg: fmt.Sprintf(format, args...)}
}
