# LO — Formal Grammar (Phase 0)

Status: **Phase 0 of `NORTHSTAR.md`'s phased plan.** This is the real, formal grammar the source
design doc (`LoLanguageSpec.pdf`) never produced — see `NORTHSTAR.md` finding #1. No lexer/parser
code exists yet; that's Phase 1. Every production below is checked against a worked derivation of
one of the source doc's own example programs in §7, and every open ambiguity the source left
unresolved is either closed here (with the decision stated) or explicitly deferred (named as such,
not silently assumed).

Acceptance bar this doc is held to (`NORTHSTAR.md`'s own words): *a grammar precise enough that
two people — or a human and this doc's own worked examples — parse the same LO snippet
identically, every time.*

## 1. Lexical grammar

LO source is a sequence of Unicode scalar values. The lexer recognizes exactly four token
shapes; anything else is a fatal lex error (matching `LoLanguageSpec.pdf` §2's own "any other
character results in a fatal syntax error").

1. **The vector keyword** — the literal ASCII sequence `vec PARENA CONSTRUCT 312`: the four
   words `vec`, `PARENA`, `CONSTRUCT`, `312`, each separated by exactly one U+0020 SPACE, no
   leading/trailing space, case-sensitive. This is matched greedily as a single token
   (`VECLIT`) — LO has no standalone identifiers, so there is no ambiguity with matching `vec`
   or `PARENA` on their own.
2. **The colon** — U+003A `:` (`COLON`), the ternary false-branch separator.
3. **The semicolon** — U+003B `;` (`SEMI`), required to terminate every top-level `Program`
   (§2). Founder real-time, 2026-08-30, after Phase 1's compiler landed: "also require
   semicolons in LO." Real, deliberate design call, not the source material's own idea — every
   worked example in `LoLanguageSpec.pdf` and every prior version of this grammar omitted a
   terminator entirely, relying on end-of-input to mark the end of the one real expression a
   program was. A required trailing `SEMI` is added here as a real, explicit statement
   terminator, matching the common C-family convention, ahead of Phase 2 (`qi`) ever needing to
   sequence more than one top-level form.
4. **Emoji tokens** — every other token in the language. §1.1 gives the exact codepoint table;
   §1.2 gives the matching rule that resolves `NORTHSTAR.md` finding #2 (emoji tokenization
   ambiguity).

Whitespace (space, tab, newline) between tokens is not significant and may be repeated freely,
**except** inside the vector keyword, where exactly single spaces in the exact four-word sequence
above are required — nothing else is a valid `VECLIT` match.

### 1.1 Token table (exact codepoints)

Every LO emoji token is defined by its **base codepoint**. Where the base codepoint is a
"text-presentation-by-default" character, Unicode commonly pairs it with U+FE0F (VS16, emoji
presentation) when typed from an emoji picker — §1.2 states how those are handled.

| Token | Emoji | Codepoint | Meaning |
|---|---|---|---|
| `S0` | 🌑 | U+1F311 | Base4 state 0 |
| `S1` | 🌒 | U+1F312 | Base4 state 1 |
| `S2` | 🌓 | U+1F313 | Base4 state 2 |
| `S3` | 🌔 | U+1F314 | Base4 state 3 |
| `QUERY` | ❓ | U+2753 | Ternary "then" |
| `PLUS4` | ➕ | U+2795 | Mod-4 addition |
| `MINUS4` | ➖ | U+2796 | Mod-4 subtraction |
| `AND4` | 🔗 | U+1F517 | Base4 AND |
| `OR4` | 🔮 | U+1F52E | Base4 OR |
| `XOR4` | 🔀 | U+1F500 | Base4 XOR |
| `STACK` | 🧱 | U+1F9F1 | Bind vectors into a matrix |
| `DOT` | 🎯 | U+1F3AF | Dot product |
| `MATMUL` | 🧮 | U+1F9EE | Matrix multiplication |
| `DIMLEN` | ⚖️ | U+2696 (+VS16) | Dimension reader |
| `EQ` | 🟰 | U+1F7F0 | Equality — **provisional, see §1.3** |
| `MATCH` | 🔍 | U+1F50D | Pattern match |
| `SCALAR` | 💧 | U+1F4A7 | Scalar type |
| `VECTOR` | 🌊 | U+1F30A | Vector type |
| `MATRIX` | 🧊 | U+1F9CA | Matrix type |
| `PATTERN` | 🕸️ | U+1F578 (+VS16) | Pattern type |
| `I32` | 🔢 | U+1F522 | 32-bit integer semantic type |
| `STRING` | 📜 | U+1F4DC | String semantic type |
| `FUNC` | ⚙️ | U+2699 (+VS16) | Callable/function semantic type |
| `VOID` | 🕳️ | U+1F573 (+VS16) | Void/nil |
| `UNION` | 🧬 | U+1F9EC | Union type binder |
| `LABEL` | 🏷️ | U+1F3F7 (+VS16) | Type assertion (prefix) |
| `DOOR` | 🚪 | U+1F6AA | Return-type assertion (top-level) |
| `MAGNET` | 🧲 | U+1F9F2 | Environment-matrix row extraction |
| `WILDCARD` | 🃏 | U+1F0CF | PCRE `.` |
| `STAR` | 🌌 | U+1F30C | PCRE `*` |
| `ONEPLUS` | ☄️ | U+2604 (+VS16) | PCRE `+` |
| `OPT` | 👻 | U+1F47B | PCRE `?` |
| `START` | 🏁 | U+1F3C1 | PCRE `^` |
| `END` | 🛑 | U+1F6D1 | PCRE `$` |
| `ALT` | 🛤️ | U+1F6E4 (+VS16) | PCRE `\|` |
| `GROUP` | 🗜️ | U+1F5DC (+VS16) | PCRE capture-group delimiter (both ends — see §5.4) |

### 1.2 Emoji-matching rule (resolves NORTHSTAR finding #2)

The source spec never stated whether a lexed emoji must match by exact codepoint sequence, by
Unicode grapheme cluster, or by some normalized form. Phase 0 decides this concretely:

1. **Match on the base codepoint**, after stripping exactly one trailing variation selector if
   present — U+FE0E (VS15, text presentation) or U+FE0F (VS16, emoji presentation). Both
   `⚖️` (U+2696 U+FE0F) and a bare `⚖` (U+2696) lex to the same `DIMLEN` token. This is the
   correct call for LO specifically because every token in §1.1 is a symbol/object emoji, not a
   human gesture or a flag — none of them have skin-tone modifiers or regional-indicator
   sequences to worry about, so a single trailing-selector-strip rule fully closes the ambiguity
   for this token set.
2. **No ZWJ (U+200D) sequences are defined.** LO has no compound-emoji tokens. A ZWJ
   immediately following any token codepoint is a fatal lex error — it is not silently absorbed
   or ignored.
3. **No other combining marks are permitted** after a token codepoint. Any codepoint LO does not
   recognize (per §1.1, plus VECLIT and COLON) is a fatal lex error, full stop — this is the
   literal reading of the source spec's own "any other character results in a fatal syntax
   error," made precise instead of assumed.

This means `⚖`, `⚖️`, and any run containing `⚖` + a second, unlisted variation selector are
respectively: valid, valid (same token), and a lex error — deterministic in all three cases.

### 1.3 Open item carried forward, not glossed over

`EQ` (🟰) is a **provisional assignment**, not a verified fact from the source. The source PDF's
own text-extraction lost the actual equality glyph everywhere it appears (rendered as an empty
parenthetical in five separate places — e.g. "checks if a vector component is equal ( ) to State
2"), which is itself the same class of "silent character loss" the source's own final answer
worried about for emoji encoding in general. 🟰 (Heavy Equals Sign, U+1F7F0) is chosen here
because it is visually and semantically the closest unclaimed "equals" emoji and does not collide
with any other token in §1.1 — but this needs a founder look at the original Gemini chat
transcript (the live page, not the exported PDF) to confirm the glyph actually used, before Phase
1 treats it as final. Flagged, not blocking: every example in §7 parses identically regardless of
which glyph ultimately wins, since `EQ` behaves as an ordinary infix comparison token throughout.

## 2. Grammar overview (EBNF)

```
Program      ::= TypedExpr SEMI

TypedExpr    ::= Door Expr
               | Expr

Door         ::= DOOR TypeSig

TypeSig      ::= UNION TypeAtom TypeAtom+      (* 2+ members, e.g. DOOR UNION STRING I32 *)
               | TypeAtom

TypeAtom     ::= SCALAR | VECTOR | MATRIX | PATTERN | I32 | STRING | FUNC | VOID

Expr         ::= Ternary

Ternary      ::= Cond QUERY Expr COLON Expr
               | Value

Cond         ::= Value EQ Value
               | Value MATCH Value            (* target MATCH pattern *)

Value        ::= Labeled
               | Arith
               | LinAlg
               | State
               | VectorLit
               | MagnetExpr
               | Ternary                        (* a ternary nests as a value: parenthesization
                                                    is purely by token adjacency, see §3 *)

Labeled      ::= TypeAtom LABEL Value

Arith        ::= Value ArithOp Value
ArithOp      ::= PLUS4 | MINUS4 | AND4 | OR4 | XOR4

LinAlg       ::= Value STACK Value STACK          (* one matrix "row group"; see §5.2 *)
               | Value DOT Value                   (* dot product *)
               | Value MATMUL Value                 (* matrix multiplication *)
               | DIMLEN Value                        (* dimension read, unary *)

State        ::= S0 | S1 | S2 | S3

VectorLit    ::= VECLIT State+

MagnetExpr   ::= VectorLit MAGNET Value           (* index-vector MAGNET matrix -> row *)

Pattern      ::= VECLIT PatternAtom+
PatternAtom  ::= State | WILDCARD | STAR | ONEPLUS | OPT | START | END
               | ALT
               | GROUP PatternAtom+ GROUP          (* see §5.4 — non-nesting in v1 *)
```

`Value` is deliberately one unified nonterminal rather than a precedence-climbing chain (`Term`,
`Factor`, ...) because the source material never establishes relative precedence between, say,
`Arith` and `LinAlg` used as sub-expressions of each other — every real example in
`LoLanguageSpec.pdf` nests exactly one operator deep before either terminating or handing off to
another `Ternary`. §3 states the concrete rule Phase 0 adopts for the cases the source never
actually exercises.

## 3. Precedence, associativity, and grouping

LO has **no explicit grouping/parenthesis tokens** — the source spec never introduces any (its
own worked "nested ternary" examples group purely by *replacing* a branch with another full
ternary, never by parenthesizing a sub-expression inline with siblings). Phase 0 therefore adopts
the only rule that makes every source example parse unambiguously:

1. **`QUERY`/`COLON` (ternary) is right-associative and lowest precedence.** `A ❓ B : C ❓ D : E`
   parses as `A ❓ B : (C ❓ D : E)` — matching every nested example in the source doc, which
   always nests in the false branch, never the true branch or the condition.
2. **`Cond` (the `EQ`/`MATCH` production) binds tighter than `QUERY`/`COLON`** — a condition is
   always a single `EQ` or `MATCH` application, never itself a ternary. (The source never writes
   a ternary as a raw condition without wrapping it in `EQ` first — e.g. the matrix-multiply
   example wraps two `DIMLEN` reads in `EQ` rather than testing a ternary result directly. Using
   a bare `Ternary` result as a condition — is a `State1` result "true"? — is a real, load-bearing
   truthiness question the source material never answers. Phase 0 leaves it **unresolved and
   out of scope**, not silently assumed: v1's grammar simply doesn't permit it. A later phase
   that wants "truthy state" conditions needs to define which of the 4 states count as true
   first.)
3. **`Arith` and `LinAlg` operators are all left-associative, single precedence tier, no source
   example chains more than one operator.** Where Phase 1 needs to compile a real qi-lowered
   expression with a longer chain (e.g. `a + b + c`), it associates left: `(a PLUS4 b) PLUS4 c`.
   This is the conventional default and doesn't contradict anything in the source; it is a Phase
   0 decision made because the source is silent, not a discovered fact.
4. **`LABEL` binds tighter than everything except a bare `State`/`VectorLit`** — `TypeAtom LABEL
   Value` always takes the smallest well-formed `Value` immediately to its right, per every
   worked example (`📜 🏷 🌊 vec ...` never has `🏷` reach across an operator).
5. **`DOOR` only ever appears once, at the very start of `Program`.** The source never nests or
   repeats it mid-expression; Phase 0 encodes that as a hard grammar constraint (`TypedExpr` is
   the sole production point for `Door`), not a convention left to discipline.

## 4. Vector, matrix, and row-extraction shapes

- A `VectorLit` is `VECLIT` followed by **one or more** `State` tokens. Phase 0 deliberately
  disallows a zero-length `VectorLit` — an empty collection is represented by the dedicated
  `VOID` token instead, so there is exactly one way to say "nothing," not two competing ones
  (a real ambiguity the source doc never raised but would have hit the first time someone wrote
  `vec PARENA CONSTRUCT 312` with nothing after it).
- A matrix (`LinAlg`'s first alternative) is `STACK Value STACK` repeated per row, matching the
  source's own `🧱 vec ... 🧱 vec ... 🧱` shape literally: `STACK` appears once *before* the first
  row and once *after every* row, including the last — i.e. an N-row matrix has exactly N+1
  `STACK` tokens, not N-1. This matches every worked matrix example in the source byte-for-byte
  (§7.3 shows the derivation).
- `MagnetExpr`'s left operand is a `VectorLit` holding exactly one `State` — the row index —
  never a bare `State`. This matches the source's own `let`-lowering example, where the index is
  always wrapped (`vec PARENA CONSTRUCT 312 🌑 🧲 ...`), never written as a bare `🌑 🧲 ...`.

## 5. Operators, in full

### 5.1 Base4 arithmetic/logic — `PLUS4 MINUS4 AND4 OR4 XOR4`

Binary infix, `Value ArithOp Value`. Operate elementwise when both operands are `VectorLit`s of
equal length (source's own vector-addition example); operate directly on two `State`s otherwise.
Mismatched vector lengths are a **compile-time error** in Phase 0's grammar — there is no runtime
in an intermediate language, so "what happens" for a length mismatch is a static rejection, not a
value.

### 5.2 Linear algebra — `STACK DOT MATMUL DIMLEN`

- `DOT`: `Value DOT Value`, both operands `VectorLit`, equal length, produces a `SCALAR`.
- `MATMUL`: `Value MATMUL Value`, both operands the matrix shape from §4, produces a matrix.
  Per the source's own dimension-checking convention, a bare `MATMUL` application does **not**
  itself check dimensions — the source always gates it behind an explicit `EQ` of two `DIMLEN`
  reads first (§7.3). Phase 0 keeps that as a *convention*, not a grammar-enforced invariant:
  the grammar permits an unchecked `MATMUL`; whether Phase 1's compiler should reject one that
  isn't dimension-gated is a real open question for that phase, named here rather than decided
  by fiat.
- `DIMLEN`: unary prefix, `DIMLEN Value`, produces a `SCALAR`-or-`I32`-shaped length reading
  (the source calls it a "dimension reader" without pinning down the exact output type; left as
  Phase 1's concern).

### 5.3 Equality and matching — `EQ MATCH`

`Value EQ Value` and `Value MATCH Value` are the only two `Cond` productions (§3.2). `MATCH`'s
right operand must be a `Pattern` (§5.4); `EQ`'s two operands must be the same `Value` shape
(both `State`, both `VectorLit` of equal length, etc.) — a shape mismatch is a compile-time error
for the same reason as §5.1.

### 5.4 PCRE pattern grammar

`Pattern` (§2) is a `VectorLit`-like sequence but drawn from `PatternAtom`, not `State`, directly
matching the source's own pattern-vector examples. One real gap found and closed here, not left
implicit: **`GROUP` (🗜️) is specified in the source as the *same* glyph for both the open and
close of a capture group** ("clamps placed before and after the sequence you want to capture").
A single non-distinguishing delimiter cannot support nested capture groups — there is no way to
tell an inner close from an outer close by glyph alone. Phase 0's decision: **v1 patterns support
at most one capture-group region, and it may not contain another `GROUP` pair.** The grammar
above encodes this directly (`GROUP PatternAtom+ GROUP` with no recursive `GROUP` inside the
inner `PatternAtom+`). Nested capture groups are deferred — they need either a second,
distinct "close" glyph or a counting scheme, and the source material never proposed either.

## 6. What Phase 0 resolves vs. defers

**Resolved (this document):**
- Exact lexical token set and codepoints (§1.1).
- The emoji-matching rule closing NORTHSTAR finding #2 (§1.2).
- Ternary associativity/precedence and the "no bare-state-as-condition" boundary (§3).
- Vector/matrix/row-extraction shape rules, including the empty-vector-vs-VOID resolution (§4).
- The `GROUP` non-nesting restriction for capture groups (§5.4).

**Explicitly deferred, not silently assumed away:**
- Truthiness of a raw `State` used directly as a `Cond` (§3.2) — out of scope for v1's grammar.
- Whether `MATMUL` without a preceding `DIMLEN`/`EQ` gate should be a compiler error (§5.2).
- The `let`-lowering blowup risk (`NORTHSTAR.md` item 6) — this grammar defines `MagnetExpr` as a
  legal production so Phase 2's `qi` lowering has somewhere to land, but does not attempt to
  solve the exponential-duplication problem itself; that is Phase 2's named job.
- Loops/recursion and error handling — still not designed (`NORTHSTAR.md` items 3, unchanged);
  no grammar production exists for either, matching the source material's own honesty that these
  were unfinished.
- The `EQ` glyph identity (§1.3) — provisionally 🟰, pending a founder check against the original
  chat transcript.

## 7. Worked derivations

Three of the source doc's own examples, re-derived against the grammar above, to satisfy the
Phase 0 acceptance bar directly (not just assert it). All three examples below predate the
required trailing `SEMI` (§1's item 3, added later, after Phase 1's compiler landed) — each
`Source:` line is quoted faithfully as the original source material wrote it; under the current
grammar every one of them needs a trailing `;` appended before `Program`'s own `TypedExpr SEMI`
production accepts it.

### 7.1 The first nested-ternary example

Source: `🌒 EQ 🌓 ❓ 🌒 🔀 🌔 : 🌔 EQ 🌔 ❓ 🔗 🌒 : 🌑` — "checks if a vector component is equal to
State 2, and if so, performs a Base4 XOR; if not, falls back through a second check."

```
Program
└─ TypedExpr → Expr → Ternary
   ├─ Cond:  S1 EQ S2                      (🌒 EQ 🌓)
   ├─ true:  Expr → Arith: S1 XOR4 S3      (🌒 🔀 🌔)
   └─ false: Expr → Ternary                (right-associated per §3.1)
             ├─ Cond:  S3 EQ S3            (🌔 EQ 🌔)
             ├─ true:  Expr → Value: AND4 applied — (source elides the left operand of 🔗;
             │         treated in Phase 1 as a unary/partial form — flagged, not silently
             │         resolved: the source's own text drops an operand here too)
             └─ false: S0                  (🌑)
```
One real thing this derivation surfaces that a prose read misses: the true-branch of the inner
ternary (`🔗 🌒`) is a `PLUS4`/AND4-shaped token followed by a single `State`, which is **not** a
valid `Arith` production (`Arith` requires two `Value` operands). This is a genuine gap in the
source's own worked example, not an artifact of this grammar — flagged here for Phase 1 rather
than quietly patched over by inventing an operand that was never stated.

### 7.2 Hello World (the fully typed example)

Source: `🚪 📜 🌑 EQ 🌑 ❓ 📜 🏷 🌊 vec PARENA CONSTRUCT 312 <16 states> : 🕳️`

```
Program
└─ TypedExpr
   ├─ Door: DOOR STRING                              (🚪 📜)
   └─ Expr → Ternary
      ├─ Cond:  S0 EQ S0                              (🌑 EQ 🌑)
      ├─ true:  Value → Labeled
      │         STRING LABEL (VECTOR LABEL VectorLit) — reading 🏷 right-associatively per §3.4:
      │         📜 🏷 (🌊 🏷 (vec ... 16 states))
      └─ false: Value → VOID                          (🕳️)
```
This is a clean, unambiguous parse under the grammar — no gaps, unlike §7.1 — which is itself
useful signal: the source's more carefully-worked final example is fully well-formed against this
grammar, while its earlier, more casual example (§7.1) is not. That is exactly the kind of
discrepancy a formal grammar is supposed to surface.

### 7.3 Ternary-gated matrix multiplication

Source: `DIMLEN [MatrixA] EQ DIMLEN [MatrixB] ❓ [MatrixA] MATMUL [MatrixB] : 🌑`

```
Program
└─ TypedExpr → Expr → Ternary
   ├─ Cond:  DimEq
   │         (DIMLEN MatrixA) EQ (DIMLEN MatrixB)
   ├─ true:  Value → LinAlg
   │         MatrixA MATMUL MatrixB
   └─ false: S0                                       (🌑)
```
Where `MatrixA` = `STACK VectorLit STACK VectorLit STACK` (a 2-row matrix, matching §4's "N+1
STACK tokens for N rows" rule exactly against the source's own `🧱 vec ... 🧱 vec ... 🧱`
literal token sequence) and `MatrixB` the same shape. This derivation is unambiguous and requires
no repair — it's the example §4's matrix-shape rule was reverse-engineered from.

## Related

- `NORTHSTAR.md` — the critical review and phased plan this grammar implements Phase 0 of.
- `LoLanguageSpec.pdf` — the source design doc every production above is checked against.
- `PARENA/stdlib/base4.prn` — the real base4 runtime LO's state space maps onto; unaffected by
  this document (Phase 0 is grammar-only, no compiler/runtime code).
