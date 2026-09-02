# NORTHSTAR — `qi` (LO's own Phase 2 frontend)

## Status

**Phase 2a (lexer) is real and shipped**: `internal/qi/lexer` (`Lex`, `Token`, `Kind`) tokenizes
the real surface grammar below — parens/brackets, symbols (permissive, no restricted charset:
`+4`/`foo-bar`/`x?`/a bare `` ` `` are all valid symbol text), decimal int literals range-checked
to 0-3, double-quoted strings (no escapes, matching LO's own `LITERAL` rule exactly), and `;` line
comments. 11 real tests, `GOWORK=off go build/vet/test ./...` clean. Phase 2b (parser + the real
name-resolution lowering to `internal/parser`'s existing AST) has NOT started — that's the next
real step, not attempted in the same pass as the lexer.

Design-only below this point, matching this repo's own standing "Spec Before Implementation"
discipline (`THE_EMILY_WAY.md` Principle 2) — the same real precedent `NORTHSTAR.md` and
`GRAMMAR.md` already set for LO itself. Written 2026-09-02 (founder real-time: "continue", against
`NORTHSTAR.md`'s own already-named Phase 2: *"`qi`'s own real frontend, same scalar/struct
scope"*), once LO's own Phase 1 work (S214 through S223) reached a real, natural pause point —
every well-scoped, unambiguous LO grammar item was either shipped or is genuinely blocked (a
source-doc ambiguity best left for founder confirmation, or a larger, separately-scoped PARENA
compiler project — see `EMILY/BACKLOG.md` SECTION 223).

## Why `qi` at all, given PARENA's own `.prn` is *already* Lisp-like

A real question worth answering directly, not glossed over: PARENA's own `.prn` syntax already
has real `defn`/`let`/`match`/`if`/`cond` forms — so what does `qi` actually add over just writing
`.prn` by hand? Two real, distinct things:

1. **`qi` stays LO's own conceptual world (nested ternaries over a 4-symbol base4 state space),
   just with a friendlier ASCII surface than the emoji alphabet** — not a general-purpose PARENA
   authoring frontend. A `qi` program still only ever produces the same, deliberately narrow
   value space LO itself does (`State`, `Arith`/base4 logic, `Let`/`LetRef`, `Switch`, `Lambda`/
   `Call`, `StringLit`, `Match`/`Pattern`) — it's a second FRONTEND onto LO's own existing AST and
   emitter, not a rewrite of either.

2. **`qi` can give LO real, NAMED variables** — the one concrete capability gap named repeatedly
   throughout LO's own Phase 1 work (`Let`/`LetRef`/`Lambda`'s own doc comments each flag this).
   LO's emoji grammar has no lexical room for a name (no quoted-identifier token distinct from
   `LITERAL`'s own STRING-VALUE role), so `Let`/`Lambda` bindings are addressed purely by
   nesting-depth position (`x0`, `x1`, ...). `qi`'s own frontend can trivially keep a real
   compile-time `name -> depth` environment during its own lowering pass — turning `(let [x 5]
   (+4 x x))` into exactly the same `parser.Let{Bound, Body}` / `parser.LetRef{Depth}` tree LO's
   emitter already knows how to emit, with ZERO changes to `internal/parser`'s AST or
   `internal/emitter`. This is `qi`'s single most concrete, real value-add — not a vague
   "friendlier syntax" claim.

## Design principle: `qi` targets LO's existing AST, not new AST or a new emitter

`qi`'s own lowering pass produces `internal/parser.Expr`/`Program` values directly — the exact
same types `internal/parser`'s own emoji-token parser already produces, and the exact same types
`internal/emitter.Emit` already knows how to turn into `.prn` text. `qi` is purely a second
FRONTEND (lexer + parser + a real name-resolution/lowering pass) feeding the SAME back end. This
is a real, deliberate architectural choice, not an oversight: it means every real, hard-won
finding from Phase 1 (the FLOAT/DOUBLE F64 coincidence, the `#define main` verification
technique, the Arena-threading `lo-program` rename, the `MatchBudget`-needs-a-`let` compiler
constraint, the if-condition-can't-hold-a-binding-form compiler bug) is inherited for free, not
re-discovered or re-solved for a second time.

## Real, concrete `qi` surface syntax (first draft)

Standard s-expression lexing: `(`/`)`/`[`/`]`, whitespace-insignificant, `;` line comments (a
real, useful addition LO's own emoji grammar has no room for — `;` is already LO's own SEMI
token), double-quoted string literals (reusing the exact same "no escapes, `"` cannot appear
inside" rule `LO_Formal_Grammar_Phase_0_Complete.md`'s own `LITERAL` already resolved), bare
symbols (identifiers), and decimal integers restricted to `0`-`3` (mapping directly to LO's own
four base4 states — a `qi` integer literal outside that range is a real, honest parse error, not
silently wrapped).

```
Program   ::= "(" "defn" Symbol "[" "]" ":" TypeAtom Expr ")"
            | Expr                              (* no DOOR -- matches LO's own optional Door *)

TypeAtom  ::= "I32" | "FLOAT" | "DOUBLE" | "STRING"   (* the Door types LO's emitter supports today *)

Expr      ::= IntLit                             (* 0-3 -> State *)
            | StringLit                          (* "..." -> StringLit *)
            | Symbol                             (* a name bound by an enclosing let/fn -> LetRef *)
            | "(" "let" "[" Symbol Expr "]" Expr ")"           (* -> Let/LetRef *)
            | "(" ArithOp Expr Expr ")"                        (* -> Arith *)
            | "(" "if" "(" "=" Expr Expr ")" Expr Expr ")"     (* -> Ternary/Eq *)
            | "(" "if" "(" "match" Expr Pattern ")" Expr Expr ")"  (* -> Ternary/Match *)
            | "(" "switch" Expr Case+ "(" "default" Expr ")" ")"   (* -> Switch *)
            | "(" "fn" "[" Symbol "]" Expr ")"                 (* -> Lambda/LetRef *)
            | "(" "call" Expr Expr ")"                         (* -> Call *)

ArithOp   ::= "+4" | "-4" | "&4" | "|4" | "^4"
Case      ::= "(" "case" IntLit Expr ")"

Pattern   ::= a real, separate ASCII mirror of GRAMMAR.md's own Pattern/PatternItem/PatternAtom/
              ClassItem AST (Literal/Wildcard/quantifiers/anchors/alternation/classes/escapes/
              groups) -- real, deliberately NOT drafted in full here; a real, separate follow-up
              once the scalar/Let/Switch/Lambda slice above is real and tested, matching Phase 1's
              own "prove the narrowest slice first" discipline.
```

Real, deliberate narrowing versus what a full Lisp frontend might support, named rather than
silently assumed: exactly one `defn`, no top-level `def`/globals, one `let` binding per form
(matching LO's own `Let`, itself real per `GRAMMAR.md` §2's own doc comment on why multiple
simultaneous bindings are out of v1 scope), one `fn` parameter (matching LO's own `Lambda`, for
the exact same reason `LO/internal/parser/ast.go`'s own `Lambda` doc comment names: multi-param
is real, deferred, ambiguous grammar territory upstream, not `qi`'s call to resolve unilaterally
either).

## Real, phased plan

**Phase 2a — lexer.** A real, standalone `internal/qi/lexer` (s-expression tokens: parens/
brackets, symbols, decimal int literals, double-quoted strings, `;` comments). No parser yet.
Real acceptance bar, matching every other Phase 1 slice in this repo: real tests against a
handful of concrete example snippets, not just "it doesn't panic."

**Phase 2b — parser + lowering, scalars/Let/Switch/Lambda only.** Parse the surface grammar above
(minus `Pattern`/`Match`, deferred to 2c) directly into `internal/parser.Expr`/`Program` values —
the real name-resolution pass (a `[]string` scope stack mapping each in-scope name to its own
depth, mirroring `internal/emitter`'s own `[]exprType` convention) is the one genuinely new piece
of logic; everything downstream (emission, verification technique) is already real and unchanged.
Real, concrete acceptance test, matching Phase 1's own: a real `qi` program compiles to the exact
same `.prn` text (or a semantically-equivalent one) its own hand-written LO/emoji equivalent
already does, verified through the SAME real `parena build` + `cc` + execution pipeline
`internal/emitter`'s own tests already use — not a second, parallel verification story.

**Phase 2c — `Pattern`/`Match` support.** The ASCII mirror of `GRAMMAR.md`'s own Pattern grammar,
named as deliberately deferred above. Real, separate acceptance test once 2b is real and tested.

**Phase 3 (design only, unchanged from `NORTHSTAR.md`'s own existing text)**: vectors/matrices/
patterns/unions, gated on the same real, already-named blockers (no real Vector/Matrix VALUES in
LO yet at all; `base4/matrix.prn` has no row-extraction primitive yet either — see
`EMILY/BACKLOG.md`'s own `MagnetExpr` follow-up note).

## Related docs

- `NORTHSTAR.md` — LO's own Phase 0/1 critical review and phased plan; this document is its own
  named Phase 2, scoped out separately once Phase 2 actually became the real next step.
- `GRAMMAR.md` — LO's own emoji grammar and real implementation-status tracker; every `qi`
  surface form above lowers to exactly the AST node that document's own EBNF already names.
- `LO_Formal_Grammar_Phase_0_Complete.md` — the founder-uploaded authoritative Phase 0 grammar
  `qi`'s own `Pattern` mirror (Phase 2c) will need to match glyph-for-concept, not just LO's own
  narrower current implementation.
