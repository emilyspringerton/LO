# LO

## What this is

A hyper-minimalist esolang: emojis, colons, and one literal string (`vec PARENA CONSTRUCT 312`)
as the entire alphabet, everything expressed as nested ternaries over a 4-symbol base4 state
space, sitting under a real, higher-level Lisp-like frontend (`qi`) and compiling down to real
`.prn` source that the existing `parena`/`burrow` CLIs turn into C/TS/Java/Go — no new backend
emitters, LO/`qi` is purely a new frontend. **Read `NORTHSTAR.md` before writing any code** — it
has the full real critical review of the source design doc (`LoLanguageSpec.pdf`, a captured
Gemini chat transcript), the real gaps found (no formal grammar yet, emoji-tokenization
ambiguity, loops/error-handling explicitly unfinished in the source material, a real code-size
risk in the spec's own `let`-lowering scheme), and the real, phased plan.

## Status

Real, live compiler (this doc was stale — corrected 2026-09-03). Phase 0 (`GRAMMAR.md`, plus the
founder-uploaded, more authoritative `LO_Formal_Grammar_Phase_0_Complete.md`) and Phase 1
(`internal/lexer`/`internal/parser`/`internal/emitter`, `cmd/lo`) are both real and shipped:
`State`/`Arith`/`Eq`/`Ternary`/`Switch`/`Case`/`Default`/`Let`/`Lambda`/`Call`/`Float`/`Double`
Doors, plus a PCRE-lite pattern lexer+parser (emitter not wired yet). Every real language feature
above has an actual `lo build` → `.prn` → `parena build` → `cc` → run end-to-end test, not just a
shape check. A real, higher-level `qi` frontend (Phase 2) has its own lexer landed. A real,
separate Rails-like framework effort (`FRAMEWORK_NORTHSTAR.md`, Controllers + a live-verified
SQLite projector shipped) is being built PARENA-native ahead of `qi`/LO themselves reaching that
far. **Real, current ceiling, found 2026-09-03** (see `NORTHSTAR.md`'s own "DUNG integration"
section): LO's own arithmetic is mod-4 by design, not general integer arithmetic, and `lo
build`'s compiled output has no way to be invoked with a runtime argument at all yet — every
compiled program is a single, self-contained, zero-parameter computation. Both matter for any
future host-integration use (e.g. `DUNG`).

## Real, current backend capability (checked directly, not assumed — see NORTHSTAR.md item 5)

- `parena` (the primary compiler): real `defstruct`/`defenum`/`match`/`loop`/`recur`/`Result`/
  `Vec`/FFI support already exists — most of LO's own real feature surface is a plausible target
  here, once a real grammar/lowering pass exists.
- `burrow` (this monorepo's own Go+PARENA reimplementation): real, but narrower — scalar params,
  `if`, binops, `not`, flat `defstruct`+`get-field` (read-only). No `let`, no construction, no
  `defenum`/`match`/`loop`/`Vec` yet. Only LO programs limited to scalars/flat structs can
  realistically reach Go via `burrow build` today.

## Related Repos

- `PARENA` — the real language/compiler LO's own backend targets; `stdlib/base4.prn` is the real
  runtime state-space model LO's own design already maps onto.
- `BURROW` — the real, currently narrower Go-target compiler; sets the real, current ceiling on
  which LO features can reach Go.
- `EMILY` — RSI loop / backlog coordination for cross-repo work.

## Founder Real-Time Direction

Whenever the founder gives real-time direction — a new ask, a correction, a "can we also..." —
route it through `emily observe -s info "Founder real-time: <summary>"` first, even if it isn't
this repo's usual domain, then sprint-plan it into `EMILY/BACKLOG.md` (`emily backlog curate`,
scoped into a real SECTION/sub-item, not just a one-line log), and only then implement. See
`EMILY/docs/THE_EMILY_WAY.md` Principle 18 ("Pave the Cow Paths").

## Apple Filing Protocol

After any meaningful change, file an Apple:
```bash
emily apples post -t completion -repo LO "<title>" "<body with commit hash>"
```
Then mark the item done in `EMILY/BACKLOG.md` and commit.

## CHANGELOG Protocol

After any meaningful change, update CHANGELOG.md:
```bash
emily changelog add LO "<what changed>"
# or manually: append a dated bullet under ## YYYY-MM-DD in LO/CHANGELOG.md
```

## Frame-Break Reframing

Founder-sourced prompting technique (REDGARDEN/NORTHSTAR.md §28, full origin in
REDGARDEN/docs2/MULTI_AGENT_RD_RESEARCH_NOTES.md §5): given a request, name the underlying
structural/systemic pattern it's one instance of — one level of abstraction up — as an added
lens during planning/triage/judgment calls. Use it to spot the general case behind a specific
ask. It augments judgment, it does not replace doing the work: direct, concrete execution of
the literal task asked for still happens every time.

## Commit Protocol (standing instruction)

Always commit and push completed work immediately — don't wait to be asked. This is the default for every repo in this monorepo.

Every commit — human-written or produced by automated code paths (git-commit helpers in emily-agent, emily.cli, IDUNA handlers, etc.) — must carry the active `emily session` fingerprint as a `session: <tag>` trailer (blank line, then the trailer). This was silently missing from several independently-implemented automated commit helpers across the monorepo until an audit on 2026-08-10 (founder, real-time: "where in the fuck is my llm session id anywhere"). If you add a new automated git-commit code path anywhere, wire in the session tag the same way — don't assume an existing helper already does it.
