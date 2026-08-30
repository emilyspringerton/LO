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

	if p.pos != len(p.toks) {
		return nil, p.errf("unexpected trailing tokens after a complete expression")
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

// expr implements GRAMMAR.md §2's `Expr ::= Ternary`.
func (p *parser) expr() (Expr, error) {
	return p.ternary()
}

// ternary implements GRAMMAR.md §2's `Ternary ::= Cond QUERY Expr COLON Expr | Value` and §3.2's
// rule that Cond is always an EQ application, never a bare Ternary — parsed here by first
// parsing a Value, then checking for a following EQ to decide whether this is a Cond at all.
func (p *parser) ternary() (Expr, error) {
	left, err := p.value()
	if err != nil {
		return nil, err
	}

	if p.peek().Kind != lexer.KindEq {
		return left, nil
	}
	p.next() // consume EQ
	right, err := p.value()
	if err != nil {
		return nil, err
	}
	cond := Eq{Left: left, Right: right}

	if p.peek().Kind != lexer.KindQuery {
		return nil, p.errf("expected QUERY after an EQ condition")
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

// value implements GRAMMAR.md §2's `Value` for this v0's covered cases: a bare State, VOID, or
// an Arith application (`Value ArithOp Value`, left-associative per GRAMMAR.md §3.3 — checked
// here by looking one token ahead for an arith operator after the first operand).
func (p *parser) value() (Expr, error) {
	tok := p.peek()

	if tok.Kind == lexer.KindVoid {
		p.next()
		return VoidExpr{}, nil
	}

	if tok.Kind != lexer.KindState {
		return nil, p.errf("expected a base4 state or VOID, got %s", tok.Kind)
	}
	p.next()
	left := Expr(State{Value: tok.State})

	if op := p.peek().Kind; isArithOp(op) {
		p.next()
		rightTok := p.peek()
		if rightTok.Kind != lexer.KindState {
			return nil, p.errf("expected a base4 state as the right operand of %s", op)
		}
		p.next()
		left = Arith{Op: op, Left: left, Right: State{Value: rightTok.State}}
	}
	return left, nil
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
