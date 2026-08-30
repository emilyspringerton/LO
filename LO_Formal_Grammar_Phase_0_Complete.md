# LO — Formal Grammar
## Phase 0 — Complete Updated Grammar

**Status:** Proposed consolidated Phase 0 grammar.

This specification consolidates the LO grammar with the currently resolved additions for:

- `I32`
- `FLOAT`
- `DOUBLE`
- lambda functions
- function calls
- `SWITCH`
- `CASE`
- `DEFAULT`
- the expanded PCRE-oriented pattern language

---

# 1. Lexical Grammar

LO source is a sequence of Unicode scalar values.

The lexer recognizes:

1. the vector keyword `VECLIT`,
2. punctuation,
3. quoted literal payloads,
4. emoji tokens defined below,
5. insignificant whitespace between tokens.

Anything else is a fatal lexical error.

Whitespace consisting of spaces, tabs, or newlines is insignificant between tokens.

Whitespace inside the vector keyword is significant.

The exact vector keyword is:

`vec PARENA CONSTRUCT 312`

It is matched greedily as one `VECLIT` token.

---

# 2. Token Table

## 2.1 Base4 Tokens

| Token | Glyph | Unicode | Meaning |
|---|---|---|---|
| `S0` | 🌑 | U+1F311 | Base4 state 0 |
| `S1` | 🌒 | U+1F312 | Base4 state 1 |
| `S2` | 🌓 | U+1F313 | Base4 state 2 |
| `S3` | 🌔 | U+1F314 | Base4 state 3 |
| `QUERY` | ❓ | U+2753 | Ternary then |
| `PLUS4` | ➕ | U+2795 | Mod-4 addition |
| `MINUS4` | ➖ | U+2796 | Mod-4 subtraction |
| `AND4` | 🔗 | U+1F517 | Base4 AND |
| `OR4` | 🔮 | U+1F52E | Base4 OR |
| `XOR4` | 🔀 | U+1F500 | Base4 XOR |

---

## 2.2 Structural and Comparison Tokens

| Token | Glyph | Unicode | Meaning |
|---|---|---|---|
| `STACK` | 🧱 | U+1F9F1 | Bind vectors into a matrix |
| `DOT` | 🎯 | U+1F3AF | Dot product |
| `MATMUL` | 🧮 | U+1F9EE | Matrix multiplication |
| `DIMLEN` | ⚖️ | U+2696 + VS16 | Dimension reader |
| `EQ` | ⚓ | U+2693 | Equality |
| `LET` | ✨ | U+2728 | Variable binding |
| `MATCH` | 🔍 | U+1F50D | Pattern match |
| `MAGNET` | 🧲 | U+1F9F2 | Let/environment reference |

---

## 2.3 Type Tokens

| Token | Glyph | Unicode | Meaning |
|---|---|---|---|
| `SCALAR` | 💧 | U+1F4A7 | Generic scalar type |
| `VECTOR` | 🌊 | U+1F30A | Vector type |
| `MATRIX` | 🧊 | U+1F9CA | Matrix type |
| `PATTERN` | 🕸️ | U+1F578 + VS16 | Pattern type |
| `I32` | 🔢 | U+1F522 | 32-bit integer |
| `FLOAT` | ⚫ | U+26AB | Floating-point value |
| `DOUBLE` | ⚪ | U+26AA | Double-precision floating-point value |
| `STRING` | 📜 | U+1F4DC | String type |
| `FUNC` | ⚙️ | U+2699 + VS16 | Function type |
| `VOID` | 🕳️ | U+1F573 + VS16 | Void/nil |
| `UNION` | 🧬 | U+1F9EC | Union type |
| `LABEL` | 🏷️ | U+1F3F7 + VS16 | Type assertion |

### Numeric type hierarchy

```text
🔢 I32
⚫ FLOAT
⚪ DOUBLE
```

`FLOAT` and `DOUBLE` are distinct floating-point types.

`I32` is the explicit 32-bit integer type.

---

## 2.4 Function and Control-Flow Tokens

| Token | Glyph | Unicode | Meaning |
|---|---|---|---|
| `LAMBDA` | 💠 | U+1F4A0 | Anonymous/lambda function |
| `CALL` | 📞 | U+1F4DE | Function invocation |
| `SWITCH` | 🔘 | U+1F518 | Multi-way selector |
| `CASE` | 🔹 | U+1F539 | Switch alternative |
| `DEFAULT` | 🔸 | U+1F538 | Default switch alternative |

The functional distinction is:

```text
💠  LAMBDA
⚙️  FUNC
📞  CALL
```

The control-flow distinction is:

```text
🔘  SWITCH
🔹  CASE
🔸  DEFAULT
```

---

## 2.5 Pattern Tokens

| Token | Glyph | Unicode | PCRE |
|---|---|---|---|
| `WILDCARD` | 🃏 | U+1F0CF | `.` |
| `STAR` | 🌌 | U+1F30C | `*` |
| `ONEPLUS` | ☄️ | U+2604 + VS16 | `+` |
| `OPT` | 👻 | U+1F47B | `?` |
| `START` | 🏁 | U+1F3C1 | `^` |
| `END` | 🛑 | U+1F6D1 | `$` |
| `ALT` | 🛤️ | U+1F6E4 + VS16 | `\|` |
| `GROUP` | 🗜️ | U+1F5DC + VS16 | `(` / `)` |
| `LITERAL` | 🔤 | U+1F524 | Literal text |
| `CLASS` | 🅰️ | U+1F170 + VS16 | `[ ... ]` |
| `NCLASS` | 🚫 | U+1F6AB | `[^ ... ]` |
| `RANGE` | ↔️ | U+2194 + VS16 | `a-z` |
| `ESCAPE` | 🛡️ | U+1F6E1 + VS16 | Escape |

---

# 3. Punctuation

```ebnf
COLON ::= ":" ;
SEMI  ::= ";" ;
```

Every `Program` terminates with `SEMI`.

---

# 4. Complete EBNF

```ebnf
Program          ::= TypedExpr SEMI ;

TypedExpr        ::= Door Expr
                   | Expr ;

Door             ::= DOOR TypeSig ;

TypeSig          ::= UNION TypeAtom TypeAtom+
                   | TypeAtom ;

TypeAtom         ::= SCALAR
                   | VECTOR
                   | MATRIX
                   | PATTERN
                   | I32
                   | FLOAT
                   | DOUBLE
                   | STRING
                   | FUNC
                   | VOID ;

Expr             ::= Switch
                   | Let
                   | Ternary ;

Switch           ::= SWITCH Value Case+ Default? ;

Case             ::= CASE Value Expr ;

Default           ::= DEFAULT Expr ;

Let              ::= LET Value Expr ;

Ternary          ::= Cond QUERY Expr COLON Expr
                   | Value ;

Cond             ::= Value EQ Value
                   | Value MATCH Pattern ;

Value            ::= Lambda
                   | Call
                   | Labeled
                   | Arith
                   | LinAlg
                   | State
                   | VectorLit
                   | MagnetExpr
                   | LetRef
                   | PatternValue
                   | Switch
                   | Ternary ;

Lambda           ::= LAMBDA LambdaParams Expr ;

LambdaParams     ::= Param
                   | Param LambdaParams ;

Param            ::= LITERAL ;

Call             ::= CALL Value ArgList ;

ArgList          ::= Value
                   | Value ArgList ;

Labeled          ::= TypeAtom LABEL Value ;

Arith            ::= Value ArithOp Value ;

ArithOp          ::= PLUS4
                   | MINUS4
                   | AND4
                   | OR4
                   | XOR4 ;

LinAlg           ::= Value STACK Value STACK
                   | Value DOT Value
                   | Value MATMUL Value
                   | DIMLEN Value ;

State            ::= S0
                   | S1
                   | S2
                   | S3 ;

VectorLit        ::= VECLIT State+ ;

MagnetExpr       ::= VectorLit MAGNET Value ;

LetRef           ::= MAGNET ;

PatternValue     ::= Pattern ;

Pattern          ::= PatternSequence
                   | PatternSequence ALT PatternSequence+ ;

PatternSequence  ::= PatternItem+ ;

PatternItem      ::= Anchor
                   | Atom Quantifier?
                   | PatternGroup ;

Anchor           ::= START
                   | END ;

Atom             ::= State
                   | Literal
                   | WILDCARD
                   | CharacterClass
                   | Escaped ;

Quantifier       ::= STAR
                   | ONEPLUS
                   | OPT ;

PatternGroup     ::= GROUP PatternSequence GROUP ;

Literal          ::= LITERAL QuotedText ;

CharacterClass   ::= CLASS ClassItem+
                   | NCLASS ClassItem+ ;

ClassItem        ::= LiteralChar
                   | Range ;

Range            ::= LiteralChar RANGE LiteralChar ;

Escaped          ::= ESCAPE Escapable ;

Escapable        ::= Literal
                   | WILDCARD
                   | STAR
                   | ONEPLUS
                   | OPT
                   | START
                   | END
                   | ALT
                   | GROUP
                   | ESCAPE ;
```

---

# 5. Program and Return Types

A complete program is:

```ebnf
Program ::= TypedExpr SEMI ;
```

A top-level return type may be specified with `DOOR`.

```ebnf
Door ::= DOOR TypeSig ;
```

Examples:

```text
🚪 🔢 🌒 ;
🚪 ⚫ 🌒 ;
🚪 ⚪ 🌒 ;
🚪 📜 value ;
🚪 🧬 📜 🔢 expression ;
```

`DOOR` occurs only at the beginning of a program.

---

# 6. Base4 State Model

The four primitive states are:

```text
🌑 = 0
🌒 = 1
🌓 = 2
🌔 = 3
```

```ebnf
State ::= S0 | S1 | S2 | S3 ;
```

These states form the fundamental Base4 state space.

---

# 7. Numeric Types

LO defines three explicit numeric scalar types:

```text
🔢 I32
⚫ FLOAT
⚪ DOUBLE
```

## 7.1 I32

`I32` represents a 32-bit signed integer.

Token:

```text
🔢
```

---

## 7.2 FLOAT

`FLOAT` represents a floating-point numeric value.

Token:

```text
⚫
```

---

## 7.3 DOUBLE

`DOUBLE` represents a double-precision floating-point numeric value.

Token:

```text
⚪
```

The black/white pair is intentional:

```text
⚫ FLOAT
⚪ DOUBLE
```

Both are single Unicode code points.

---

# 8. Type Signatures

```ebnf
TypeSig ::= UNION TypeAtom TypeAtom+
          | TypeAtom ;
```

Examples:

```text
🚪 🔢
🚪 ⚫
🚪 ⚪
🚪 📜
🚪 🧬 📜 🔢
🚪 🧬 ⚫ ⚪
```

A union requires at least two member types.

---

# 9. Type Assertions

```ebnf
Labeled ::= TypeAtom LABEL Value ;
```

`LABEL` is the type assertion operator.

Examples:

```text
🔢 🏷️ 🌒
⚫ 🏷️ 🌒
⚪ 🏷️ 🌒
📜 🏷️ value
```

`LABEL` binds more tightly than binary operators.

---

# 10. Arithmetic and Base4 Logic

```ebnf
Arith   ::= Value ArithOp Value ;

ArithOp ::= PLUS4
          | MINUS4
          | AND4
          | OR4
          | XOR4 ;
```

Examples:

```text
🌒 ➕ 🌔
🌒 ➖ 🌔
🌒 🔗 🌓
🌒 🔮 🌓
🌒 🔀 🌓
```

Operations on vectors are elementwise when vector lengths agree.

Vector length mismatch is a compile-time error.

All binary arithmetic/logic operators are left-associative.

Thus:

```text
🌒 🔀 🌓 🔀 🌔
```

parses as:

```text
(🌒 🔀 🌓) 🔀 🌔
```

---

# 11. Linear Algebra

```ebnf
LinAlg ::= Value STACK Value STACK
         | Value DOT Value
         | Value MATMUL Value
         | DIMLEN Value ;
```

A matrix with N rows contains N+1 `STACK` tokens.

`DOT` requires vector operands of equal length.

`MATMUL` operates on matrices.

`DIMLEN` reads a dimension.

---

# 12. Equality and Matching

```ebnf
Cond ::= Value EQ Value
       | Value MATCH Pattern ;
```

Equality:

```text
🌒 ⚓ 🌒
```

Pattern matching:

```text
value 🔍 pattern
```

`EQ` requires compatible value shapes.

`MATCH` requires a `Pattern` on the right.

---

# 13. Ternary Expressions

```ebnf
Ternary ::= Cond QUERY Expr COLON Expr
          | Value ;
```

Example:

```text
🌒 ⚓ 🌒 ❓ 🌔 : 🌑
```

Conceptually:

```text
if 1 == 1 then 3 else 0
```

Ternary expressions are right-associative:

```text
A ❓ B : C ❓ D : E
```

means:

```text
A ❓ B : (C ❓ D : E)
```

A raw Base4 state is not automatically a condition.

---

# 14. Let Bindings

```ebnf
Let    ::= LET Value Expr ;
LetRef ::= MAGNET ;
```

`✨` creates a binding.

A bare `🧲` refers to the nearest active binding.

Example:

```text
✨ 🌒 🧲
```

Nested bindings shadow outer bindings.

The v1 grammar does not define named binding syntax beyond the existing environment mechanism.

---

# 15. Lambda Functions

```ebnf
Lambda       ::= LAMBDA LambdaParams Expr ;

LambdaParams ::= Param
               | Param LambdaParams ;

Param        ::= LITERAL ;
```

`💠` constructs an anonymous function.

Example:

```text
💠 🔤"x" x
```

Conceptually:

```text
lambda x -> x
```

Two parameters:

```text
💠 🔤"x" 🔤"y" body
```

Conceptually:

```text
lambda x, y -> body
```

`💠` is the lambda constructor.

`⚙️` remains the semantic function type.

---

# 16. Function Calls

```ebnf
Call    ::= CALL Value ArgList ;

ArgList ::= Value
          | Value ArgList ;
```

`📞` invokes a function.

Example:

```text
📞 function argument
```

Multiple arguments:

```text
📞 function arg1 arg2
```

Function arity checking is a semantic/compiler responsibility.

---

# 17. Switch / Case

```ebnf
Switch  ::= SWITCH Value Case+ Default? ;

Case    ::= CASE Value Expr ;

Default ::= DEFAULT Expr ;
```

Tokens:

```text
🔘 SWITCH
🔹 CASE
🔸 DEFAULT
```

A switch evaluates its selector once.

Cases are evaluated in source order.

The first matching case wins.

If no case matches and a default exists, the default expression is selected.

If no case matches and there is no default, the result is `VOID`.

---

# 18. Switch Example

```text
🔘 🌒
🔹 🌑 🌑
🔹 🌒 🌔
🔹 🌓 🌒
🔸 🌔 ;
```

Conceptually:

```c
switch (1) {
    case 0:
        result = 0;
        break;

    case 1:
        result = 3;
        break;

    case 2:
        result = 1;
        break;

    default:
        result = 3;
        break;
}
```

There is no fall-through between `CASE` clauses.

---

# 19. Pattern Grammar

LO supports a deliberately bounded PCRE-oriented pattern language.

```ebnf
Pattern         ::= PatternSequence
                  | PatternSequence ALT PatternSequence+ ;

PatternSequence ::= PatternItem+ ;

PatternItem     ::= Anchor
                  | Atom Quantifier?
                  | PatternGroup ;

Anchor          ::= START
                  | END ;

Atom            ::= State
                  | Literal
                  | WILDCARD
                  | CharacterClass
                  | Escaped ;

Quantifier      ::= STAR
                  | ONEPLUS
                  | OPT ;

PatternGroup    ::= GROUP PatternSequence GROUP ;
```

Pattern precedence:

1. atom
2. quantifier
3. concatenation
4. alternation

---

# 20. Pattern Literals

```ebnf
Literal ::= LITERAL QuotedText ;
```

Example:

```text
🔤"cat"
```

represents:

```regex
cat
```

Quoted text provides Unicode literal content.

---

# 21. Wildcard

```text
🃏
```

represents:

```regex
.
```

Example:

```text
🔤"c" 🃏 🔤"t"
```

represents:

```regex
c.t
```

---

# 22. Quantifiers

The three v1 quantifiers are:

```text
🌌 = *
☄️ = +
👻 = ?
```

Examples:

```text
🔤"a" 🌌
```

means:

```regex
a*
```

```text
🔤"a" ☄️
```

means:

```regex
a+
```

```text
🔤"a" 👻
```

means:

```regex
a?
```

A quantifier applies to exactly the preceding atom or capture group.

---

# 23. Anchors

```text
🏁 = ^
🛑 = $
```

Examples:

```text
🏁 🔤"cat"
```

means:

```regex
^cat
```

```text
🔤"cat" 🛑
```

means:

```regex
cat$
```

Together:

```text
🏁 🔤"cat" 🛑
```

means:

```regex
^cat$
```

---

# 24. Alternation

```text
🛤️ = |
```

Example:

```text
🔤"cat" 🛤️ 🔤"dog"
```

means:

```regex
cat|dog
```

Multiple alternatives:

```text
🔤"cat" 🛤️ 🔤"dog" 🛤️ 🔤"fox"
```

means:

```regex
cat|dog|fox
```

Alternation has lower precedence than concatenation.

Therefore:

```text
🔤"a" 🔤"b" 🛤️ 🔤"c" 🔤"d"
```

means:

```regex
ab|cd
```

---

# 25. Character Classes

```ebnf
CharacterClass ::= CLASS ClassItem+
                  | NCLASS ClassItem+ ;

ClassItem      ::= LiteralChar
                 | Range ;

Range          ::= LiteralChar RANGE LiteralChar ;
```

Positive class:

```text
🅰️ 🔤"a" 🔤"b" 🔤"c"
```

means:

```regex
[abc]
```

Negated class:

```text
🚫 🔤"a" 🔤"b" 🔤"c"
```

means:

```regex
[^abc]
```

---

# 26. Character Ranges

```text
↔️
```

represents a character range.

Example:

```text
🅰️ 🔤"a" ↔️ 🔤"z"
```

means:

```regex
[a-z]
```

Multiple ranges:

```text
🅰️
🔤"a" ↔️ 🔤"z"
🔤"A" ↔️ 🔤"Z"
🔤"0" ↔️ 🔤"9"
```

means:

```regex
[a-zA-Z0-9]
```

---

# 27. Escaping

```ebnf
Escaped   ::= ESCAPE Escapable ;

Escapable ::= Literal
            | WILDCARD
            | STAR
            | ONEPLUS
            | OPT
            | START
            | END
            | ALT
            | GROUP
            | ESCAPE ;
```

`🛡️` causes the following regex metacharacter to be interpreted literally.

Examples:

```text
🛡️ 🌌
```

represents a literal `*`.

```text
🛡️ 🃏
```

represents a literal `.`.

---

# 28. Capture Groups

```ebnf
PatternGroup ::= GROUP PatternSequence GROUP ;
```

Example:

```text
🗜️ 🔤"ab" 🗜️
```

means:

```regex
(ab)
```

A group may be quantified:

```text
🗜️ 🔤"ab" 🗜️ ☄️
```

means:

```regex
(ab)+
```

Nested capture groups are not permitted in v1.

---

# 29. MATCH Examples

Exact string match:

```text
value 🔍 🏁 🔤"cat" 🛑
```

Equivalent PCRE:

```regex
^cat$
```

One or more digits:

```text
🅰️ 🔤"0" ↔️ 🔤"9" ☄️
```

Equivalent PCRE:

```regex
[0-9]+
```

Cat or dog:

```text
🔤"cat" 🛤️ 🔤"dog"
```

Equivalent PCRE:

```regex
cat|dog
```

---

# 30. Precedence and Associativity

LO expression precedence, highest to lowest:

```text
1. State / literal / vector / pattern
2. Function call
3. Lambda
4. LABEL
5. Arithmetic / Base4 logic
6. DOT / MATMUL / DIMLEN
7. EQ / MATCH
8. SWITCH / CASE
9. QUERY / COLON
```

Binary arithmetic and linear-algebra operations are left-associative.

Ternary expressions are right-associative.

Regex has its own precedence:

```text
1. Atom
2. Quantifier
3. Concatenation
4. Alternation
```

---

# 31. Complete Token Summary

```text
BASE4
🌑 S0
🌒 S1
🌓 S2
🌔 S3

CONTROL
❓ QUERY
🔘 SWITCH
🔹 CASE
🔸 DEFAULT

ARITHMETIC
➕ PLUS4
➖ MINUS4
🔗 AND4
🔮 OR4
🔀 XOR4

LINEAR ALGEBRA
🧱 STACK
🎯 DOT
🧮 MATMUL
⚖️ DIMLEN

COMPARISON / BINDING
⚓ EQ
🔍 MATCH
✨ LET
🧲 MAGNET

TYPES
💧 SCALAR
🌊 VECTOR
🧊 MATRIX
🕸️ PATTERN
🔢 I32
⚫ FLOAT
⚪ DOUBLE
📜 STRING
⚙️ FUNC
🕳️ VOID
🧬 UNION
🏷️ LABEL

FUNCTIONS
💠 LAMBDA
📞 CALL

REGEX
🔤 LITERAL
🃏 WILDCARD
🌌 STAR
☄️ ONEPLUS
👻 OPT
🏁 START
🛑 END
🛤️ ALT
🗜️ GROUP
🅰️ CLASS
🚫 NCLASS
↔️ RANGE
🛡️ ESCAPE

STRUCTURAL
🚪 DOOR
: COLON
; SEMI
```

---

# 32. Resolved Design Decisions

The current Phase 0 token decisions are:

| Concept | Token |
|---|---|
| Base4 0 | 🌑 |
| Base4 1 | 🌒 |
| Base4 2 | 🌓 |
| Base4 3 | 🌔 |
| Equality | ⚓ |
| Let | ✨ |
| Let reference | 🧲 |
| Switch | 🔘 |
| Case | 🔹 |
| Default | 🔸 |
| Lambda | 💠 |
| Function | ⚙️ |
| Call | 📞 |
| I32 | 🔢 |
| Float | ⚫ |
| Double | ⚪ |

The numeric distinction is intentionally:

```text
🔢 I32
⚫ FLOAT
⚪ DOUBLE
```

The function distinction is intentionally:

```text
💠 LAMBDA
⚙️ FUNC
📞 CALL
```

---

# 33. Explicit v1 Boundaries

The following are intentionally outside this grammar:

- raw-state truthiness as a ternary condition;
- nested capture groups;
- regex lazy quantifiers (`*?`, `+?`, `??`);
- regex lookahead/lookbehind;
- regex backreferences;
- named capture groups;
- arbitrary PCRE conditionals;
- loops;
- recursion;
- multiple simultaneous `LET` bindings;
- explicit access to shadowed outer `LET` bindings;
- lambda closure semantics beyond the basic expression form;
- switch fall-through;
- switch ranges;
- switch pattern cases;
- implicit numeric coercion between `I32`, `FLOAT`, and `DOUBLE`;
- unchecked matrix-dimension enforcement at grammar level;
- implementation-specific floating-point representation details.

These should be specified in later semantic/compiler phases rather than silently inferred by the parser.
