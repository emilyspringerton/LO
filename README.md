# LO

A hyper-minimalist esolang — emojis, colons, and one literal string as the entire alphabet,
everything a nested ternary over a 4-symbol base4 state space — sitting under a real, higher-level
Lisp-like frontend (`qi`) and compiling down to real `.prn` source for the existing
[`parena`](https://github.com/emilyspringerton/PARENA)/[`burrow`](https://github.com/emilyspringerton/BURROW)
CLIs to turn into C, TypeScript, Java, and Go.

New repo (2026-08-30). `LoLanguageSpec.pdf` is the real source design doc (a captured design
conversation). See `NORTHSTAR.md` for the full critical review of that spec and the real, phased
implementation plan.

## Status

Scoping/review only — no compiler code yet.

## License

Unlicense (public domain) — see `LICENSE`.
