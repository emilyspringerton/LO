# NORTHSTAR — LO "Batteries Included" (a Rails-like framework, as a dogfooding vehicle)

## Where this came from

Founder real-time, same session as Phase 1 landing: "continue - now lets make LO batteries
included - write the stdlib using LO and PARENA to build a rails like framework into LO as a
dogfooding Northstar" → "using emojis for the stdlib mono emojis whatevs they are called" → "and
then as a dogfooding northstar for that work we need to standup SHITHUB our own version of GITHUB
or GITLAB" → "imean iduna" (SHITHUB's IAM should be IDUNA-powered) → "parena too i guess" →
"shitlab will be for our ci work copying the gitlab apis". Explicitly sequenced AFTER the real
LO compiler emitter landed (S214-02/S214-03, this same session) — "once we get the emitter
working" was the founder's own stated gate.

This is spec-only, matching `THE_EMILY_WAY.md` Principle 2 ("Spec Before Implementation") and
`GRAMMAR.md`/`NORTHSTAR.md`'s own already-established discipline for this repo: no framework
code is written here, because — checked directly, not assumed — LO cannot support one yet.

## Real, honest capability check before anything else

A "Rails-like framework" needs, at minimum: named routes dispatching to handlers, request/
response data flowing through multi-step logic, and persistent records with fields. Checked
against LO's own real, current state (Phase 1, just landed):

- **No variables.** LO's grammar has no `let`. Every real expression is one ternary tree; there
  is nowhere to hold "the current request" while computing a response.
- **No multi-argument, multi-statement functions in LO itself.** The Phase 1 emitter emits
  exactly one `(defn main [] : I32 <expr>)`; LO has no surface syntax for a user-defined function
  with parameters at all yet — that's `qi`'s own job (Phase 2), not landed.
- **No records/structs in LO.** `GRAMMAR.md`'s type system has `Vector`/`Matrix`/`Pattern`
  primitives over base4 states, not named fields — a "User" or "Repo" model has no real LO shape
  yet.
- **No strings, in the ordinary sense.** LO's `String` type assertion exists in `GRAMMAR.md`
  §3.2, but Phase 1's compiler doesn't implement string literals, string routing paths, or HTTP
  bodies — everything real so far is base4 state arithmetic.

None of this is a surprise — `NORTHSTAR.md`'s own phased plan already named Phase 2 (`qi`,
`defn`/`let`/`cond`/`match`/`if`) as the real gate for anything beyond scalar ternary expressions,
and Phase 3+ (vectors/matrices/patterns/unions) as design-only. A framework needs Phase 2 at an
absolute minimum, and realistically parts of Phase 3+ (real `String`/`Pattern` support) too.

## What "batteries included" concretely means here, once unblocked

Named now, built later, matching a real MVC/Rails shape rather than an abstract wish:

1. **Router** — path pattern matching. Real, natural reuse of already-built work: `PARENA/
   stdlib/base4/pattern.prn`'s own backtracking matcher (S208-07/S208-09, wildcard/quantifier/
   anchor/ALT/GROUP over base4 vectors) is architecturally the same problem as HTTP path routing
   (`/repos/:name/commits` matched against a request path) — but it operates on base4 STATE
   vectors, not characters. A real router needs either (a) a parallel PARENA `stdlib/http/
   router.prn` built the same way but over `String`/byte sequences, or (b) a real byte-to-base4
   encoding layer reusing `pattern.prn` directly (the string-encoding scheme `LoLanguageSpec.pdf`
   itself already designed — see `GRAMMAR.md`'s own string-literal section). Real, undecided;
   flagged for a founder call once Phase 2 lands, not resolved here.
2. **Models** — `qi`'s own `defstruct`-equivalent (Phase 2) for named-field records, persisted
   via IDUNA (see below), not a new PARENA-side storage engine.
3. **Controllers/handlers** — `qi` `defn` dispatched by the router; the real "framework" glue is
   this dispatch table, matching Rails' own routes-to-controllers convention.
4. **Migrations/versioning** — real, direct reuse, not reinvented: `PARENA/stdlib/papercraft/
   note_version_mod.prn`'s own coalesce/eviction/conflict-detection primitives (S215-02) are the
   same real shape a migration/versioning system needs generically, already built this session.

## Real, named dogfooding target: SHITHUB + SHITLAB

The founder's own concrete reason to build this framework at all, not an abstract exercise:

- **SHITHUB** — a self-hosted GitHub-shaped app (real git hosting + a real web GUI). Real,
  founder-corrected stack decision: IAM is **IDUNA**-powered (not PARENA — founder real-time:
  "imean iduna"), with PARENA "too i guess" — read as PARENA providing the framework/application
  layer (router/handlers/models per this doc) while IDUNA remains this monorepo's own single
  real trust authority for auth, matching root `CLAUDE.md`'s own standing rule ("IDUNA is the
  central trust authority... never trust tokens from other sources") rather than a new, separate
  auth system.
- **SHITLAB** — real, separate scope, founder-clarified: "shitlab will be for our ci work copying
  the gitlab apis" — a GitLab-CI-API-shaped service (pipelines/jobs/runners), not a second git
  host. Real, undecided whether SHITLAB is its own PARENA app on this same framework or a real
  CI layer bolted onto SHITHUB's own git storage — a founder call once there's a real framework
  to build either on.

Both are real, separate NORTHSTAR-worthy repos in their own right once this framework exists —
not scoped further here; this doc's own job is the framework underneath them, not the apps.

## The "mono emoji" naming question — real, open, not resolved

Founder real-time: "using emojis for the stdlib mono emojis whatevs they are called." Read as:
name new LO/framework-level concepts (a route, a model, a migration) with their own emoji
tokens, consistent with `GRAMMAR.md`'s own existing token table, rather than introducing
ASCII keywords into LO's alphabet. Real, honest, unresolved questions, named rather than
guessed at:
- "Mono emoji" most plausibly means **monochrome/outline-style emoji variants** (many emoji have
  a black-and-white glyph alongside the default full-color one, selected via the VS15 text-
  presentation selector `GRAMMAR.md` §1.2 already defined stripping rules for) — i.e. a real,
  visual style choice for which glyph variant represents each new framework concept, not a
  different Unicode block. Needs a founder confirmation before locking in specific glyphs.
- Any new framework-level emoji tokens (a hypothetical ROUTE/MODEL/MIGRATE token) would need
  their own `GRAMMAR.md` amendment (a real, versioned change to the Phase 0 spec, not a silent
  addition) — named here as the real process this doc expects, not skipped.

## Real, phased plan

**Phase A (blocked on LO Phase 2 — `qi`, not started)**: `defn`/`let` real enough to hold a
request value across a few steps and call sibling handlers. Nothing in this framework doc can
start before this.

**Phase B (blocked on Phase A + a real router decision above)**: the router alone, as the
smallest real "framework" proof point — one real, hand-written route pattern matched against one
real, hand-written request path, dispatching to one hard-coded handler. Matches this repo's own
established "narrowest real slice first" discipline (Phase 1's own `on-thing`/xor-check
precedent).

**Phase C (design only, not detailed here)**: models + IDUNA persistence, migrations reusing
`note_version_mod.prn`'s own primitives, and the real SHITHUB git-hosting logic itself.

## Related

- `NORTHSTAR.md` — LO's own core phased compiler plan (Phase 0-3+); this document is downstream
  of it, not a replacement.
- `PARENA/stdlib/base4/pattern.prn` — the real backtracking-matcher precedent a router would
  reuse the architecture of.
- `PARENA/stdlib/papercraft/note_version_mod.prn` — the real version-management decision logic a
  migration system would reuse directly.
- `IDUNA` — the real, intended IAM backend for SHITHUB, per root `CLAUDE.md`'s own "central trust
  authority" rule.
- `EMILY` — RSI loop / backlog coordination for cross-repo work.
