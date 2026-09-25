# NORTHSTAR — LO-2D (bidirectional grid execution model)

## Where this came from

Founder real-time, 2026-09-25 (`EMILY/BACKLOG.md` SECTION 549, routed via `emily observe`, obs
`2026-09-25T10-46-42Z`, Apple #20808): "make LO work both left to right and bottom to top (a
programming language that goes 2 directions and can resolve at the same time) (imagine
programming a matrix of emojis or like a dominos like crossroads setup) its like lisp the actual
programming language is data." SECTION 549 scopes this as a real design pass, not implementation —
"design-only counts as done for this item... say clearly no compiler code was written." That's
exactly this document. The other SECTION 549 item (stdlib gaps + "cannon programming" into
`DEADWEIGHT_2`) is explicitly out of scope here — held for a separate pass.

## A real correction to this pass's own briefing, checked before designing anything

The task that spawned this doc described LO as "NORTHSTAR-only... no compiler code has been
written yet." **That's stale, checked directly against this repo, not assumed.** As of this pass,
LO has a real, formal Phase 0 grammar (`GRAMMAR.md`, superseded/extended for everything it covers
by the founder-uploaded `LO_Formal_Grammar_Phase_0_Complete.md`), a real, shipped Phase 1 compiler
(`internal/lexer`/`internal/parser`/`internal/emitter`, `cmd/lo`) verified end-to-end through
`parena build`/`cc`/execution for `State`/`Arith`/`Eq`/`Ternary`/`Switch`/`Let`/`Lambda`/`Call`/
`Float`/`Double`/`StringLit`/`Match`, plus a real Phase 2a `qi` lexer (`internal/qi/lexer`,
`QI_NORTHSTAR.md`). This matters directly for scope: LO-2D is not designed against a blank slate —
it has to sit ALONGSIDE a real, working 1D grammar/compiler without breaking it. That constraint
drives the single most load-bearing design decision below (§2): **LO-2D is a strict additive
superset of LO's existing 1D grammar, not a rewrite or a second language.**

## Unpacking the ask, concretely

- **Two-dimensional layout.** Source isn't a token stream; it's a grid. LO's existing emoji
  alphabet already leans this way visually — this makes it literal.
- **Bidirectional execution.** Evaluation proceeds along BOTH axes: rows left-to-right, columns
  bottom-to-top (not top-to-bottom — a real, deliberate reading of the founder's own exact words,
  named explicitly in §3).
- **Simultaneous resolution.** Not "row first, then column" (that just secretly privileges one
  axis) and not unspecified/nondeterministic order (a real correctness hazard the moment two
  axes touch the same cell). Needs a real model. §1 below states which one and why.
- **"Dominoes crossroads."** A cell shared by a horizontal chain and a vertical chain needs real,
  defined behavior at the join — not silently picking whichever axis "wins."
- **"Like Lisp, the program is data."** The grid itself should be a real, inspectable data
  structure LO/`qi` code can eventually construct and manipulate, not just a textual layout
  convention with no runtime representation. §6 names how far this design gets that, and how far
  it honestly doesn't (yet).

## 1. The evaluation model chosen, and why

**Chosen: a checkerboard-structured dependency graph, resolved by worklist fixed-point iteration,
with explicit equality-checked junctions at cells reachable from both axes.** Call it the
**Crossroads model** in the rest of this doc.

Why this, and not the obvious alternatives:

- **Not "row-major then column-major" (sequential, one axis privileged).** Rejected because it
  isn't actually bidirectional — it's linear execution wearing a 2D costume, and it can't express
  the one case that matters most: a blank cell in the MIDDLE of a row that can only be filled by
  consulting its column, which itself needs something from a row. A fixed traversal order cannot
  express that dependency at all; only a real graph can. This is the single fact that forces a
  dependency-graph answer rather than a "pick an order" one.
- **Not unspecified/nondeterministic order.** Rejected per the task's own framing — that's a
  correctness hazard, not a feature, unless the underlying computation is provably order-
  independent. The Crossroads model gets to claim genuine order-independence (see "confluence"
  below), which is a much stronger and more honest claim than "we didn't specify an order."
- **Not a general constraint solver / SMT.** Considered: let junction cells be resolved by
  arbitrary algebraic equation-solving over mod-4 arithmetic (e.g., `row-op` and `col-op` both
  nontrivially combine into the shared cell, and the compiler solves for the unique value that
  satisfies both). Rejected for v1 — real, unscoped extra work (equation solving over a small
  finite ring, not obviously always uniquely solvable, no existing precedent anywhere in
  PARENA/LO), and it isn't needed to satisfy the founder's actual ask, which is about two chains
  MEETING and AGREEING, not about deriving a value no direct chain-walk could produce. Named here
  as a real, deliberate, considered-and-rejected v2+ direction, not silently unavailable.
- **Not "grid cells always self-verify against both axes."** Considered: make literal cells
  checked against both axes' running folds too (not just blanks). Rejected — it turns every
  literal into a two-axis constraint by default, which is a much stronger, more failure-prone
  system for very little real gain (you'd have to painstakingly balance every literal in the grid
  against both its row's and column's accumulated arithmetic, for cells the programmer never
  intended to be junctions at all). The chosen model keeps literals as unconditional axioms —
  junction checking only ever applies to cells the programmer explicitly marked (🧩, §2) — which
  is also the more Lisp-honest reading of "code is data, not magic": junctions are asked for, not
  inferred.

**Why this is honestly "simultaneous," not just "we didn't say."** The equations in §3 below have
the property that any cell's resolved value depends only on OTHER cells' resolved values, never on
itself (once cycles are rejected, §3), and every cell is written by exactly one rule (ordinary
axis resolution, or the junction EQ-check). That makes the system **confluent**: any legal
resolution order — row-first, column-first, interleaved, or genuinely concurrent (one goroutine
per row, one per column, synchronizing only at 🧩 cells via a promise/future) — reaches the exact
same fixed point. "Resolves at the same time" is therefore a true claim about the *result* being
order-independent, which is what licenses an implementation to actually run both axes
concurrently later, not a claim this doc is making without backing it.

## 2. Grid grammar — lexical and structural layer

LO-2D is entered by a new top-level mode marker, reusing the **existing** `MATRIX` token (🧊,
U+1F9CA) as a prefix — real, deliberate reuse (🧊 already means "matrix" in LO's vocabulary; no
new glyph needed for this). Everything after `🧊` up to the closing `SEMI` is grid-mode source.

**A real, deliberate, flagged break from §1's existing whitespace-insignificance rule**: inside
grid mode, **newline becomes significant** — it is the row separator. This is scoped to grid mode
only; ordinary LO/`qi` programs (which never see a leading `🧊`) are completely unaffected. Within
a row, cells need no separator beyond ordinary token boundaries — LO's existing emoji tokens are
already unambiguous adjacent to each other (each is a fixed codepoint run, §1.2 of `GRAMMAR.md`),
so `🌒🔀🌔` already lexes cleanly with no whitespace needed, same as today.

**Column alignment is by token index, not visual/character column.** Row `r`'s `i`-th token
occupies the same grid column as row `r'`'s `i`-th token, for every row. This closes exactly the
kind of ambiguity `NORTHSTAR.md`'s own finding #2 (emoji tokenization) already had to close for the
1D grammar — named and resolved here rather than left to "however it happens to look."

**Well-formedness (compile-time, static rejection — matching `GRAMMAR.md` §5.1's own "no runtime
in an intermediate language" convention):**
- Every row has the same token count `C`; the grid has `R` rows. Jagged grids are a compile error.
- `R` and `C` must both be odd, and each ≥ 1 — every row and column must begin and end on an
  operand cell (see checkerboard below).

**Checkerboard cell classification** (row index `r`: 0 = top/first source line, `R-1` = bottom/
last source line; column index `c`: 0 = leftmost):

| `r` parity | `c` parity | Role | Holds |
|---|---|---|---|
| even | even | **Operand** | a `State` literal, `⬜` BLANK, `🧩` JUNCTION, or `📍` REF (below) |
| even | odd | **Row-operator** | one of the existing `ArithOp` tokens (➕➖🔗🔮🔀) — belongs to that row's own fold only |
| odd | even | **Column-operator** | one of the same `ArithOp` tokens — belongs to that column's own fold only |
| odd | odd | **Corner** | reserved, must be the existing `VOID` (🕳️) token — unused in v1, named as a real reserved slot for a possible future third axis (diagonal), not designed here |

**Three new glyphs, minimal and justified — added here for §1.1's token table:**

| Token | Emoji | Codepoint | Meaning |
|---|---|---|---|
| `BLANK` | ⬜ | U+2B1C | Operand cell resolved from whichever single axis reaches it |
| `JUNCTION` | 🧩 | U+1F9E9 | Operand cell that MUST be reachable from both axes, and both must agree |
| `PIN` | 📍 | U+1F4CD | 2D coordinate reference (below) |

Both new operand markers get the same trailing-VS16-strip matching rule §1.2 of `GRAMMAR.md`
already established for every other token — no new lexing concept, just two more table rows.

**Output marking.** Reusing the existing `DOOR` (🚪) convention (`GRAMMAR.md` §3 item 5: "Door
only ever appears once"), extended to 2D: exactly one operand cell in the grid is prefixed with
🚪. That cell's resolved value is the program's result. A grid with zero or more than one
🚪-marked cell is a compile error — considered inferring the output from graph sinks instead and
rejected for the same "explicit over inferred" reason junctions aren't auto-detected either.

**2D reference (`PIN`, minimal, secondary — not required for the core crossroads case, named for
completeness and for §6's homoiconic story):** `REF ::= PIN VectorLit VectorLit` — a cell whose
value is defined as *equal to* another named cell's value by absolute `(row, col)` coordinate,
reusing `MagnetExpr`'s own already-established "a small `VectorLit` as an index" shape rather than
inventing a new addressing scheme. Creates an explicit dependency edge in the same graph as
everything else in §3; a REF cycle is a compile error, same as any other cycle.

**Backward compatibility, checked, not asserted:** an `R=1` grid has no odd row index at all — zero
column-operator cells, zero corners exist. It degenerates byte-for-byte to LO's existing 1D `Arith`
chain grammar. Example 1 in §4 verifies this by hand, not just by argument.

## 3. Formal evaluation semantics

For each operand cell, its **resolved value** `V(r,c)` and each axis's own **running accumulator**
(`acc` for rows, `accCol` for columns) are two different quantities — conflating them was the
actual dead end worked through and discarded while designing this (an earlier draft tried to make
every operand cell simultaneously the fold's input AND its output at the same position, which is
circular; keeping them separate is what makes the model well-defined).

**Row fold** (left to right), for row `r`, operand slots at `c = 0, 2, ..., C-1`:
- `acc(r,0) = V(r,0)`.
- For `c ≥ 2`: `acc(r,c) = ArithOp(r,c-1) ( acc(r,c-2), V(r,c) )` — the ORDINARY existing
  left-associative `Arith` fold (`GRAMMAR.md` §3 item 3), applied per slot, unchanged from 1D LO.
- `rowResult(r) = acc(r, C-1)`.

**Column fold** (bottom to top), for column `c`, operand slots at `r = R-1, R-3, ..., 0`: same
shape, mirrored — seed at the bottom row, `accCol` walking upward, `colResult(c) = accCol(c, 0)`.

**Resolving `V(r,c)` itself:**
- **Literal**: `V(r,c) = ` its own written State. Always an axiom — never checked against either
  axis's fold (§1's rejected-alternative note explains why).
- **REF**: `V(r,c) = V(refR, refC)` of its target cell (which must itself be resolved first — a
  real dependency edge, participates in the same worklist as everything below).
- **BLANK (⬜) or JUNCTION (🧩)**: has no written value of its own. Each axis, if it can reach this
  cell, proposes a **candidate** by **pass-through**: `rowCandidate = acc(r, c-2)` (the row's own
  accumulator immediately before this slot, carried forward UNCHANGED — the row-operator token
  immediately preceding a blank/junction slot is real, syntactically present, and semantically
  **inert** for this one purpose: it is not applied to produce the blank's own value). Symmetric
  for `colCandidate = accCol(c, r+2)`.
  - A candidate is only available once its own predecessor position is itself resolved — this is
    exactly the recursive, may-need-the-other-axis-first structure that makes fixed-point
    iteration necessary rather than a single pass (§1's "not sequential" argument, made concrete).
  - **Exactly one candidate available** → `V(r,c) := ` that candidate. Ordinary propagation, no
    check. **A cell marked 🧩 with only one axis able to reach it is a compile error** ("declared a
    junction, but only one axis actually reaches it — use ⬜ instead").
  - **Both candidates available** → cell **must** be marked 🧩 (a cell reachable both ways but only
    marked ⬜ is a compile error too: "reachable from both axes — mark 🧩 to declare this is
    intentional"). If `rowCandidate == colCandidate`: `V(r,c) := ` that value — the crossroads
    resolves, the dominoes fit. If they differ: **compile-time JUNCTION MISMATCH** — the dominoes
    do not fit. Not a runtime value, not silently picking one side; a static rejection, matching
    every other shape-mismatch convention already in `GRAMMAR.md` §5.1/§5.3.
  - **Neither candidate available** → compile-time **UNDERDETERMINED CELL** error.
- Once `V(r,c)` is known (however it was resolved), each axis's own fold continues normally past it
  — `acc(r,c) := V(r,c)`, and the row's OWN operator applies fully and un-inertly to whatever comes
  next. Inertness is a one-time thing, only for producing the blank/junction's own value, never for
  what happens afterward.

**Algorithm**: standard worklist/fixed-point iteration (the same family as classic dataflow
analysis, spreadsheet recalculation, or a Kahn Process Network) — repeatedly attempt every not-yet-
resolved cell's row/column candidates; a pass that resolves at least one cell triggers another
pass; termination when no cell resolves in a full pass. Any cell still unresolved is
UNDERDETERMINED. **Cycle detection is required** (a REF chain, or — structurally impossible for
plain row/column pass-through since axis position strictly increases/decreases, but real for REF —
any dependency loop) as a separate, explicit compile pass before the worklist runs, matching every
other "static rejection, not a runtime surprise" convention this repo already holds itself to.

## 4. Worked examples, hand-derived

**Example 1 — degenerate `R=1`, backward compatibility.** Row: `2 ➕ 3 ➕ 1` (no blanks).
`acc(0,0)=2`, `acc(0,2)=PLUS4(2,3)=1`, `acc(0,4)=PLUS4(1,1)=2`. `rowResult=2`. This is exactly
`(2+3+1) mod 4 = 2` — the same left-associative fold `GRAMMAR.md` §3 item 3 already defines for 1D
LO, unchanged, confirming LO-2D adds nothing and breaks nothing for any program that doesn't use a
second row.

**Example 2 — a plain `⬜` blank, single axis, 3×3.** Grid (PLUS4 throughout):
```
r=0 (top):     ⬜    ➕    S0
r=1:           ➕          ➕
r=2 (bottom):  S1    ➕    S0
```
Column 0 (bottom→top): `acc(2,0)=1` (seed). `(0,0)` is blank: `rowCandidate` — none, `(0,0)` is
row 0's OWN seed, nothing to its left, the row cannot propose anything. `colCandidate =
acc(2,0) = 1`. Exactly one candidate → `V(0,0) := 1`, ordinary propagation (marking this cell 🧩
would itself be a compile error, since only one axis reaches it — this example is deliberately the
contrast case for Example 3/4 below). Row 0 continues: `acc(0,0):=1`, `acc(0,2)=PLUS4(1,0)=1`,
`rowResult(row0)=1`. `colResult(col0) = V(0,0) = 1` (row 0 is column 0's topmost slot).

**Example 3 — a real `🧩` junction, both axes agree, 5×5.** Grid (PLUS4 throughout, `S`=literal):
```
r=0 (top):     S2    ➕    🧩🚪  ➕    S1
r=1:           ➕          ➕          ➕
r=2:           S0    ➕    S3    ➕    S0
r=3:           ➕          ➕          ➕
r=4 (bottom):  S0    ➕    S3    ➕    S0
```
Column 2 (bottom→top): `acc(4,2)=3` (seed). `acc(2,2)=PLUS4(3,3)=6 mod 4=2`. Reaching `(0,2)`:
`colCandidate = acc(2,2) = 2`. Row 0: `acc(0,0)=2`. Reaching `(0,2)`: `rowCandidate = acc(0,0) =
2`. **Both candidates = 2** → junction resolves, `V(0,2) := 2`. Row 0 continues: `acc(0,2):=2`,
`acc(0,4)=PLUS4(2,1)=3`, `rowResult(row0)=3`. Column 2: `(0,2)` is its topmost slot,
`colResult(col2) = V(0,2) = 2`. `(0,2)` is also the 🚪 cell → **program result = 2.** (Rows 2/4 and
columns 0/4 are fully literal and uneventful — `rowResult(row2)=rowResult(row4)=3`,
`colResult(col0)=2`, `colResult(col4)=1` — included for completeness, not because they're
interesting; every literal cell is an unconditional axiom, per §3.)

**Example 4 — the same grid, one cell changed, junction mismatch (compile error).** Change `(2,2)`
from `S3` to `S1`. `acc(2,2) = PLUS4(3,1) = 4 mod 4 = 0`. `colCandidate` at `(0,2)` becomes `0`.
`rowCandidate` is unchanged (row 0 wasn't touched) = `2`. `0 ≠ 2` → **compile-time JUNCTION
MISMATCH at `(0,2)`: row-side candidate = S2 (2), column-side candidate = S0 (0).** A single-cell
change three rows away breaks the crossroads — the dominoes literally don't fit, and the compiler
says so statically, never at runtime.

## 5. Mapping to `.prn` — honest answer to "is this expressible directly"

**No — it needs a new compilation strategy, not just a new grammar.** `.prn` is linear
S-expression text; PARENA itself has no 2D/grid concept and doesn't need one. The real, concrete
answer: **the fixed-point resolution in §3 is a compile-time elaboration pass, entirely inside a
new LO-2D compiler stage (Go, not yet written), that runs to completion BEFORE any `.prn` is
emitted.** Its output is an ordinary, linear, topologically-ordered chain of `let` bindings — one
per resolved cell, in the exact order the worklist resolved them — terminating in the 🚪-marked
cell's own reference, reusing LO's EXISTING `Let`/`LetRef` depth-index AST nodes and
`internal/emitter` completely unchanged. The grid concept is fully erased by the time PARENA ever
sees the program; PARENA/`.prn` stay exactly as narrow and grid-ignorant as they are today.

**A real, named, open scoping question this doc does not resolve**: the worked examples above are
all fully static (every operand is a literal, blank, junction, or REF — nothing depends on a
runtime input). That's a deliberate v1 simplification, not an oversight, but it's worth stating as
open rather than silently permanent: does LO-2D v1 stay fully static (all resolution happens at
LO-2D-compile-time, §3's worklist never sees an unknown), or does it need to support a
runtime-parameterized cell (e.g., a `Lambda` parameter feeding a grid cell)? If the latter, a
junction's EQ-check stops being a compile-time-only rejection and needs a real runtime assert
(PARENA's `Result`/`Err`, matching how `MATCH` already threads errors) — real, separate,
significantly harder work, named here as genuinely open, not assumed either way.

**Reaches `parena build` (C/TS/Java) only, not `burrow build` (Go) — a direct, mechanical
consequence of an already-documented fact, not new investigation.** `NORTHSTAR.md`'s own item 5
already states `burrow` has no `let` support at all. Since the elaboration pass's entire output is
a `let`-chain, LO-2D inherits the exact same "scalar-only reaches both, anything using `Let`
reaches only `parena build`" boundary the `k8s.prn`/`k8s/scaling.prn` precedent already set for
plain LO — not a new limitation, the same one, extended.

## 6. Homoiconicity — "the program is data," honestly scoped

Two real, different levels, not conflated:

- **What this design gives for free (weak sense)**: the grid, as parsed, is an ordinary AST the Go
  compiler can inspect/transform — true of any LO/`.prn` program today, not a new property.
- **What "like Lisp" actually asks for (strong sense)**: the grid itself constructible/inspectable
  as a first-class **runtime** `qi` value — real macro-writing-macro territory, S-expression-as-
  list's actual precedent. **Checked directly against `PARENA/STDLIB.md` rather than assumed**:
  PARENA's own GENERAL `Vec` already supports real `Vec`-of-`Vec` today (`flatten [(v :
  &(Vec (Vec T)))...]`, `zip` returning `(Vec (T U))`, a real `NDArray` with shape/`matmul`) — this
  is NOT the same thing as LO's own base4-specific `stdlib/base4/vector.prn`/matrix construction,
  which `NORTHSTAR.md` already names as not-started (S208-03); PARENA's GENERAL Vec is a separate,
  already-real capability LO-2D could lean on directly. So the strong, homoiconic version is closer
  than `NORTHSTAR.md`'s existing Phase-3 gating implied — the real remaining gap is entirely on
  `qi`'s own side: a grid-literal surface syntax and a lowering pass that materializes a grid as a
  real `(Vec (Vec Cell))` PARENA value, not a PARENA stdlib gap. Named as real, phased, NOT
  attempted in this pass (§8, Phase G3).

## 7. Stdlib / compiler gap audit

- **Three new glyphs** (⬜, 🧩, 📍) — mechanical, `GRAMMAR.md` §1.1 table addition, no new lexing
  concept beyond what §1.2 already covers.
- **New AST layer** — `Grid`/`GridCell`/`RowOp`/`ColOp`/`Corner` types in `internal/parser/ast.go`,
  additive, alongside (never replacing) the existing 1D `Expr` tree.
- **The single genuinely new piece of engineering**: the §3 worklist/fixed-point elaboration pass.
  Nothing in LO, `qi`, PARENA, or `burrow` today does dependency-graph/fixed-point resolution — the
  entire existing toolchain is one-pass recursive-descent parse-then-emit. This is real, unstarted,
  and is the actual size of this ask, not the glyph additions.
- **No PARENA-level stdlib gap for v1** (fully static grids) — confirmed by §5: the elaboration
  pass's output is plain `let`-chains PARENA already runs today.
- **A real `qi`-level gap for the strong homoiconic version (§6, Phase G3, not v1)**: a grid-literal
  surface syntax in `qi` and a lowering pass targeting PARENA's real, already-existing `Vec`-of-
  `Vec` — genuinely new `qi` frontend work, not a PARENA stdlib gap (PARENA already has what's
  needed here, per the STDLIB.md check above).

## 8. Open questions, named rather than guessed at

1. **Static-only v1 vs. runtime-parameterized grids** (§5) — left open, not decided either way.
2. **Interior-blank pass-through inertness** (§3): the row/column operator immediately before a
   blank/junction slot is defined as inert for producing that slot's own value. The alternative
   (apply the operator to the incoming candidate too) was worked through and rejected because it
   turns a junction's EQ-check into a real algebraic equation-solve, not a simple comparison — but
   this trade-off is a real design choice, not a discovered fact, and is named as revisitable.
3. **Corner cells (odd, odd) are reserved, unused in v1** — a natural slot for a future third
   (diagonal) axis the founder didn't ask for. Explicitly out of scope, not precluded.
4. **Jagged (non-rectangular) grids** — disallowed in v1 (§2's well-formedness rule). No real use
   case was found that needs them for the founder's own ask; named as deliberately narrow, not
   an oversight.
5. **Multi-junction chains** (a junction's resolved value feeding into ANOTHER junction elsewhere
   in the grid) — the §3 worklist algorithm should handle this correctly by construction (nothing
   in its definition assumes exactly one junction), but this is **untested** — there is no
   implementation yet, and no worked example above exercises more than one junction at once. Named
   honestly as unverified, not silently assumed fine.

## 9. Real, phased plan — design only, none of this started

**Phase G0** — `GRAMMAR.md`/`LO_Formal_Grammar_Phase_0_Complete.md` amendment: the three new
glyphs, the grid EBNF, the checkerboard well-formedness rules, formalized as real productions (this
doc states the model; G0 is transcribing it into the grammar docs' own EBNF style once the founder
has had a chance to react to the open questions in §8).

**Phase G1** — the §3 worklist/fixed-point elaboration pass, as a new, standalone, unit-testable Go
package — real acceptance bar, matching every other Phase in this repo: the four worked examples in
§4 above, reproduced as real Go tests with the exact hand-derived values checked, before any `.prn`
emission is attempted.

**Phase G2** — wire G1's output (a `Let`-chain `Expr`) into the EXISTING `internal/emitter`,
unchanged — real acceptance bar: a grid program compiles through `parena build`/`cc`/execution and
produces the exact hand-derived result from one of §4's examples, live-verified, not just unit-
tested in isolation (matching Phase 1's own "shape check is not enough" discipline).

**Phase G3** (design only, not detailed here) — the strong homoiconic version from §6: a `qi`
grid-literal surface syntax lowering to a real PARENA `Vec`-of-`Vec` runtime value, gated on `qi`
itself reaching Phase 2b/2c first (`QI_NORTHSTAR.md`).

## Related

- `NORTHSTAR.md` — LO's own Phase 0/1 critical review; this doc's own "real correction" in the
  intro is checked directly against it.
- `GRAMMAR.md` / `LO_Formal_Grammar_Phase_0_Complete.md` — the 1D grammar LO-2D is an additive
  superset of; Phase G0 above amends these, doesn't replace them.
- `QI_NORTHSTAR.md` — Phase G3's real dependency; also the home of the strong homoiconic grid-
  literal idea named in §6.
- `PARENA/STDLIB.md` — checked directly (not assumed) for this doc's §6 finding that PARENA's
  general `Vec`-of-`Vec` already exists, separate from LO's own not-yet-started base4-vector work.
- `EMILY/BACKLOG.md` SECTION 549 — the founder ask and scoping this doc answers.
