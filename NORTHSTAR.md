# NORTHSTAR — LO

## Where this came from

New repo (2026-08-30, upstream pre-created by the founder, same real "empty repo before the ask"
pattern `SAND`/`JEWEL`/`BURROW`/`DUNG`/`CarePyre`/`ECOWAR` already followed), containing only
`LICENSE` (Unlicense — matching `PARENA`/`SKULDMARK`'s own real public-domain convention) and
`LoLanguageSpec.pdf` — a captured Gemini chat transcript designing an esolang. Founder real-time:
"LO upstream and design doc added - check it over it needs to compile into java ts c and go via
parena cli and burrow cli." This doc is that real, critical review, plus a real, phased scoping
pass — matching this monorepo's own standing "Spec Before Implementation" discipline
(`THE_EMILY_WAY.md` Principle 2) rather than starting to write a lexer against an
informally-specified, unreviewed design.

## What LO actually is, read from the real source doc

A three-tier compilation pipeline, as the spec's own final design pass states directly:

```
qi (high-level Lisp-like frontend: defn/let/match/cond/if)
  ↓
LO (strict, hyper-minimalist intermediate: emojis + colons + one literal string, nested
    ternaries only, no variables/loops/traditional control flow)
  ↓
PARENA (the real, existing base4 symbol-algebra backend)
```

LO's own real alphabet is deliberately almost nothing: emoji tokens (states, operators, type
assertions), the colon (`:`, ternary false-branch separator), and exactly one literal string,
`vec PARENA CONSTRUCT 312`, as the sole vector constructor keyword. Everything else — matrices,
patterns/regex, unions, functions — is built by composing those three token classes. All data
bottoms out in a 4-symbol state space (🌑🌒🌓🌔, mapped 1:1 to 2-bit values), which is a real,
direct callback to `PARENA/stdlib/base4.prn`'s own already-real base4 symbol algebra — LO isn't
inventing a new runtime model, it's building an esolang surface syntax over an existing one.

## Real, critical review — checked, not rubber-stamped

The spec is a genuinely coherent, internally consistent piece of esolang design — the string
encoding math checks out (a byte is exactly 4 base4 symbols since a base4 symbol carries 2 bits;
`I32` as 16 symbols since `4^16 = 2^32` is correct arithmetic, not asserted), the qi→LO→PARENA
layering is a sound compiler architecture, and the doc is honest about its own biggest wart (the
`let`-block lowering "explodes into a massive, stateless chain" — flagged directly in the source
material, not glossed over). It's also, explicitly and by design, an *esolang* — optimized for
being a weird, fun constraint puzzle, not for being pleasant or efficient to actually compile or
maintain. Real gaps and risks worth naming before any implementation starts:

1. **No formal grammar exists yet.** Every example in the source doc is prose-plus-illustration,
   not a BNF/PEG production rule set. A real lexer/parser needs an exact, unambiguous grammar —
   real work not done in the source material, and the single largest real blocker to Phase 0
   below.
2. **Emoji-as-token-alphabet is a real, non-trivial lexing problem, not a cosmetic detail.**
   Unicode emoji are not single codepoints in general — skin-tone modifiers, ZWJ sequences,
   variation selectors (`U+FE0F`), and font-dependent rendering of "the same" emoji all exist.
   The spec never states whether LO's lexer matches by exact codepoint sequence, by Unicode
   grapheme cluster, or by some normalized form — a real, load-bearing decision for whether
   `🌒` typed on two different keyboards/editors is guaranteed to be the same LO token. Needs a
   real answer before writing a single line of lexer code, not assumed.
3. **Loops/recursion and error handling are explicitly unfinished.** The source doc's own final
   message lists both as real, named, NOT-yet-designed remaining pieces. A language with no real
   iteration story beyond "nested ternaries" and no real error-reporting design is not a complete
   spec yet — real, honest, not glossed over.
4. **The real target-compilation claim needs one concrete correction.** The founder's own ask —
   "compile into java ts c and go via parena cli and burrow cli" — is achievable, but LO/`qi`
   itself should NOT try to emit Java/TS/C/Go directly. The sane, minimal-new-work architecture
   (and the one the source doc's own "clean pure parena" framing already points at) is: `qi`/LO's
   own real compiler emits real, valid `.prn` **source text** — nothing more — and the
   already-existing, already-working `parena`/`burrow` CLIs do 100% of the actual C/TS/Java (via
   `parena build`) and Go (via `burrow build`) emission, completely unchanged. LO's own real
   scope is "a new frontend that targets `.prn` text," not a second, parallel set of backend
   emitters — reusing everything this monorepo already built for PARENA/BURROW is the whole
   point, not an afterthought.
5. **Real, current backend capability sets which parts of LO are honestly reachable today, and
   this differs sharply by target — checked directly against both real compilers' own current
   state, not assumed:**
   - **`parena` (the primary, most complete compiler)** already has real `defstruct`/`defenum`/
     `match`/`loop`/`recur`/`Result`/`Vec`/FFI support (`PARENA`'s own real, current status).
     Most of LO's own real feature surface (vectors, matrices via `Vec`-of-`Vec`, unions via
     tagged `defenum`, pattern matching) is a plausible real target for `.prn` emitted toward
     `parena build`'s own C/TS/Java paths, once a real grammar and lowering pass exist.
   - **`burrow`'s own two emitters are real, but genuinely narrower** (this same session's own
     work): scalar `I32`/`F64`/`Bool`/`String`/`Unit` params, `if`, the real binop set, `not`,
     `defstruct`+`get-field` (read-only — no struct *construction*, no `let`, no `defenum`/
     `match`/`loop`/`Vec`). **Concretely: LO/qi programs limited to scalars and flat structs can
     realistically compile to Go via `burrow build` today; anything using LO's own real Matrix/
     Pattern/Union types, or `qi`'s own `let`-block lowering, cannot — not until `burrow` grows
     `defenum`/`match`/`loop`/`Vec`/construction/`let` support, real, separate, unstarted work.**
     This is the same honest boundary `stdlib/k8s/k8s.prn` vs. `stdlib/k8s/scaling.prn` already
     drew this same session — LO does not get to skip it by being a different frontend.
6. **The `let`-as-AST-duplication lowering is a real, potentially serious code-size/compile-time
   risk**, not just an aesthetic one — the source doc's own honesty about this ("explodes into a
   massive... chain") should be taken as a real engineering flag: a real program with a handful
   of nested `let`s could generate an exponentially large `.prn` (and downstream C/Go) file. Real,
   open question for Phase 1+: is a real, bounded De Bruijn/matrix-indexing scheme (as sketched)
   actually implementable without that blowup, or does `qi` need a smarter lowering (e.g., real
   PARENA `let` support, once/if that lands, instead of manual inlining)?

## Real, phased plan

**Phase 0 — a real, formal grammar.** Write the actual BNF/PEG grammar LO's own source doc never
produced — token classes (state, operator, type-assertion, the one literal string), a real,
explicit emoji-matching rule (item 2 above), real operator precedence/associativity for nested
ternaries. No code yet. Real acceptance bar: a grammar precise enough that two people (or a human
and this doc's own worked examples) parse the same LO snippet identically, every time.

**Phase 1 — the real, smallest proof point: LO → `.prn` text, scalars and flat structs only.**
A real, minimal LO compiler (target language for the compiler itself: real, not yet decided —
Go is the pragmatic default, matching this monorepo's own tooling) that lexes/parses the real
Phase 0 grammar and emits real `.prn` text for exactly the subset `burrow`'s own two emitters
already support today: scalar params, `if`, binops, `not`, flat `defstruct`+`get-field`. Real,
concrete acceptance test: a real LO program (the spec's own `on-thing`-shaped example, scoped to
its scalar `I32` path) compiles to real `.prn`, and that `.prn` compiles cleanly through BOTH
`parena build` (C/TS/Java) and `burrow build` (Go) — the real, founder-named bar, proven on the
narrowest real slice first, matching every other real Phase 0/1 in this monorepo.

**Phase 2 — `qi`'s own real frontend, same scalar/struct scope.** The Lisp-like `defn`/`let`/
`cond`/`match`/`if` surface lowering to the same Phase 1 LO/`.prn` subset — real, deferred
question from item 6 above (the `let`-lowering blowup risk) gets a real answer here, not assumed
away.

**Phase 3+ (design only, not detailed here)**: vectors/matrices/patterns/unions, gated on
`parena build`'s own already-real support for them (real, achievable per item 5) and, separately,
on `burrow` growing `defenum`/`match`/`loop`/`Vec`/construction/`let` support before those same
LO features could ever compile to Go — real, named, sequenced dependency, not glossed over.

## Related

- `PARENA` — the real base4 symbol algebra + compiler LO's own backend targets; `stdlib/
  base4.prn` is the real runtime LO's own state space already maps onto directly.
- `BURROW` — the real, narrower Go-target compiler; item 5 above is the concrete, current
  boundary on which real LO features can reach Go today.
- `stdlib/k8s/scaling.prn` / `stdlib/k8s/k8s.prn` — the real, already-proven precedent for
  "scalar-only compiles to both targets, String/Vec-heavy stays C-only" that LO's own real scope
  should follow, not reinvent.
- `EMILY` — RSI loop / backlog coordination for cross-repo work.
