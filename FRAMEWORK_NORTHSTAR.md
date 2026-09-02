# NORTHSTAR — LO "Batteries Included" (a Rails-like framework, as a dogfooding vehicle)

## 2026-09-02 update: real capability re-check + the concrete API definition

Founder real-time: "ok lets build SHITHUB using LO we need to build in that rails like framework
into LO... define the api first of the web framework... make it like rails build in the stdlib as
needed using single point code emojis when necessary to add to the language via the stdlib."

**Real capability re-check — a lot has changed since this doc's own original "no variables, no
functions, no strings" blockers (below) were written.** LO itself, not just `qi`, now has: real
`Let`/`LetRef` bindings (depth-indexed, not named yet — that's still `qi`'s own real job, see
`QI_NORTHSTAR.md`), real `Switch`/`Case`/`Default`, real `Lambda`/`Call` (single-param), a real
`StringLit` value and `DOOR STRING`, and real `Match`/`Pattern` wired into PARENA's own
`regex/pcre.prn`. What LO still genuinely lacks, unchanged: **multi-function programs** — each
`.llll` file still compiles to exactly one `defn`, so a real router/controller/model triad
*written in LO's own emoji syntax* still needs `qi`'s own Phase 2b (parser + lowering), which has
only reached Phase 2a (lexer) so far. This is the real reason the API below is being built as
PARENA stdlib `.prn` modules first (matching this session's own established, working pattern for
every other piece of LO's real backend — `base4/*.prn`, `regex/pcre.prn`, `http/router.prn`), not
LO source — `qi` Phase 2b is the real, separate bridge that will let application code (routes,
controllers, models) eventually be *authored* in a friendlier surface syntax that lowers into
calls against this same stdlib, once it lands.

**The "single point code emoji" question, answered rather than deferred**: no new LO grammar
tokens are needed for this pass. Every piece of the API below (`Route`, `Router`, `Model`,
migrations) is a real PARENA `defstruct`/`defn` in `.prn` source — LO's own emoji alphabet isn't
involved at all yet, since nothing here requires a new LO-level (not just stdlib-level) construct.
If a real, later piece of this framework genuinely needs a new LO-level construct that can't be
expressed as an ordinary stdlib function call (the kind of gap `SWITCH`/`LAMBDA`/`MATCH` each
closed for LO itself), it gets a real, single-codepoint emoji token added the same way those were
— named and justified in `GRAMMAR.md`, not invented silently. Not needed yet.

## Real, concrete Rails-like API definition

Four pillars, matching Rails' own MVC + routing shape. **Routing is real and shipped**
(`PARENA/stdlib/http/routes.prn`, this same pass); Models/Controllers/Migrations below are the
real, concrete API this pass DEFINES but has not yet built — named explicitly as not-yet-built,
not conflated with the routing work that is.

### 1. Routing — **real, shipped** (`stdlib/http/routes.prn`, on top of the already-real
`stdlib/http/router.prn`)

```
(defstruct Route (method : String) (pattern : String) (handler-name : String))

(add-route! router method pattern handler-name dest)   ;; Rails' `get "/repos", to: "repos#index"`
(match-route router method path dest)                   ;; -> route index, or -1 (not found)
(resource-routes router "repos" dest)                    ;; Rails' `resources :repos` --
                                                           ;; generates all 5 RESTful entries:
                                                           ;;   GET    /repos          repos#index
                                                           ;;   POST   /repos          repos#create
                                                           ;;   GET    /repos/:id      repos#show
                                                           ;;   PUT    /repos/:id      repos#update
                                                           ;;   DELETE /repos/:id      repos#destroy
```

Real, honest, named restriction (see `routes.prn`'s own header comment): `handler-name` is a
plain `String`, not a callable function value — PARENA's own `fn` literals are non-capturing and
file-scope-only, so a route table can't literally dispatch BY CALLING its own matched handler yet.
`match-route` returns WHICH route matched; the app's own dispatch code (a real, hand-written
`cond`/`match` chain today) decides what to call. Real, separate follow-up once `qi` lands: real
named function dispatch, closing this gap for real rather than working around it forever.

### 2. Models — real API DEFINED, not yet built

Per this doc's own already-made, founder-corrected decision (below): persisted via **IDUNA**, not
a new PARENA storage engine. A real, minimal ActiveRecord-lite shape:

```
(defstruct Model (kind : String) (id : String) (fields : (Vec (String String)) @ Region))

(model/save model dest)               ;; -> (Result Model IduaError) @ Region, POST/PUT to IDUNA
(model/find kind id dest)             ;; -> (Result Model IduaError) @ Region, GET from IDUNA
(model/all kind dest)                 ;; -> (Result (Vec Model) IduaError) @ Region
(model/destroy kind id dest)          ;; -> (Result Unit IduaError) @ Region, DELETE
```

Real, named, not-yet-decided question: IDUNA's own real HTTP client shape from PARENA — does
IDUNA need a new, real generic-record-storage endpoint (`/api/v1/records/:kind/:id`), or does each
SHITHUB model (`Repo`/`User`/`Issue`/`PullRequest`) get its own real, typed IDUNA endpoint (more
work per model, more type safety, matching how Apples/agent-secrets already work)? Real, honest,
deferred to whoever builds this — not assumed either way here.

### 3. Controllers — real API DEFINED, not yet built

```
(defstruct Request (method : String) (path : String) (params : (Vec (String String)) @ Region) (body : String))
(defstruct Response (status : I32) (body : String) (content-type : String))

;; one real PARENA defn per action, Rails' own controller-action shape:
(defn repos-index [(req : &Request) (dest : Arena @ Region)] : Response ...)
(defn repos-show   [(req : &Request) (dest : Arena @ Region)] : Response ...)
```

The real "framework glue" (Rails' own `ActionDispatch`): given a `Request`, call `match-route`,
then dispatch on the returned route's own `handler-name` string via a real, hand-written
`cond`/`match` chain to the right action `defn` — real, working, and named as a real interim
shape above, not a final design.

### 4. Migrations — real, direct reuse, unchanged from this doc's own original decision

`PARENA/stdlib/papercraft/note_version_mod.prn`'s own coalesce/eviction/conflict-detection
primitives (S215-02) are the same real shape a migration/versioning system needs generically,
already built. No new code needed for this pass — a real, separate follow-up is wiring a model's
own schema version through them concretely, once Models (above) are real.

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

## Real, phased plan (2026-09-02 update: A and B are real and done — PARENA itself already has
real, direct `defn`/`let`/multi-arg functions, the actual Phase A blocker was outdated the moment
this framework work targeted PARENA `.prn` stdlib directly instead of waiting on `qi`)

**Phase A — DONE.** `defn`/multi-arg functions/`let` real enough to hold a request value across a
few steps and call sibling handlers: this was never actually blocked on LO's own `qi` frontend —
PARENA itself has had real `defn`/`let`/multi-param functions the whole time. The original
blocker text above was written assuming the framework would be authored in LO's own emoji syntax
first; building it as PARENA stdlib instead (matching every other real LO backend piece this
session) sidesteps that entirely.

**Phase B — DONE.** The router (`stdlib/http/router.prn`, `stdlib/http/routes.prn`): real path
pattern matching, HTTP-method dispatch, and Rails' own `resources` RESTful-entry generation, all
real and tested (`make test-http-router`/`test-http-routes`).

**Phase B2 (not started) — Models.** `model/save`/`find`/`all`/`destroy` per the real API defined
above, once the real IDUNA HTTP client shape question (generic-record vs. per-model typed
endpoints) is decided.

**Phase B3 (not started) — Controllers.** The real `Request`/`Response` structs and a real,
hand-written `cond`/`match` dispatch chain from `match-route`'s own returned route to the right
action `defn`, per the real API defined above.

**Phase C (design only, unchanged)**: migrations wiring (real primitives already exist, per the
Migrations section above) and the real SHITHUB git-hosting domain logic itself (repos/users/
issues/pull requests) — not detailed here, a real, separate NORTHSTAR-worthy scope of its own once
B2/B3 are real.

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
