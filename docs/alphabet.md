# Visual Alphabets

This document defines the **reference glyph tables** for VIZID.

## Important note on fonts

Unicode glyph rendering varies slightly across fonts.
These tables are the **reference set** for this repo and are chosen to be:

- monochrome / inversion-safe
- non-combining
- visually distinct
- commonly supported by Nerd Fonts and modern OS fonts

If you encounter a font where two glyphs are confusable, treat it as a font issue and switch fonts.

---

## Core-36 Glyph Alphabet (base-36 digits)

Core-36 renders base-36 values `0–9` and `A–Z`.

Values are grouped into 4 base-shape families:

1. ■ Square
2. ◆ Diamond
3. ▲ Triangle
4. ● Circle

Each family contains 9 glyphs (ordered left→right).

### Indexing

- `0..8`   => squares
- `9..17`  => diamonds
- `18..26` => triangles
- `27..35` => circles

### Base-36 mapping

| Value | ASCII | Glyph |
|------:|:-----:|:-----:|
| 0 | 0 | □ |
| 1 | 1 | ⊡ |
| 2 | 2 | ⊠ |
| 3 | 3 | ⊞ |
| 4 | 4 | ⊟ |
| 5 | 5 | ◫ |
| 6 | 6 | ◩ |
| 7 | 7 | ◪ |
| 8 | 8 | ■ |
| 9 | 9 | ◇ |
| 10 | A | ◈ |
| 11 | B | ◊ |
| 12 | C | ⟐ |
| 13 | D | ⟡ |
| 14 | E | ❖ |
| 15 | F | ⧫ |
| 16 | G | ◆ |
| 17 | H | ⬥ |
| 18 | I | △ |
| 19 | J | ◬ |
| 20 | K | ◭ |
| 21 | L | ◮ |
| 22 | M | ⟁ |
| 23 | N | ▲ |
| 24 | O | ◢ |
| 25 | P | ◣ |
| 26 | Q | ◤ |
| 27 | R | ○ |
| 28 | S | ◌ |
| 29 | T | ◍ |
| 30 | U | ◐ |
| 31 | V | ◑ |
| 32 | W | ◒ |
| 33 | X | ◓ |
| 34 | Y | ◔ |
| 35 | Z | ● |

> Note: triangles have fewer perfect “interior topology” variants in Unicode; the chosen set prioritizes distinctness and monotone readability.

---

## UUID Prefix Alphabet (separate visual language)

UUID prefix is NOT part of core-36. It signals the UUID boundary and adds entropy.

### ASCII prefix set (8)

These are easy to type and speak, and avoid `-`:

```
~ ! @ $ % ^ & *
```

### Glyph mapping

| ASCII | Glyph | Spoken |
|:----:|:-----:|---|
| ~ | ✦ | star |
| ! | ✧ | hollow-star |
| @ | ✱ | starburst |
| $ | ✲ | pinwheel |
| % | ✳ | asterisk-star |
| ^ | ✴ | eight-point |
| & | ✵ | sparkle |
| * | ✶ | six-point |

Implementations must round-trip prefix ASCII↔glyph with this exact table.
