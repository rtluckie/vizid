# Worked Examples

## Example 1 — standard generation (UTC)

Assume timestamp:

- 2026-01-29 08:15:44.123 UTC

### Step A — convert each component to stored values

- Year = 2026
- Month stored = 1-1 = 0
- Day stored = 29-1 = 28
- Hour = 8
- Minute = 15
- Second = 44
- Millisecond = 123

### Step B — encode each component in base-36 with fixed widths

Widths (v1): Year=3, Month=1, Day=1, Hour=1, Minute=2, Second=2, Ms=2

- Year 2026 -> base36 = 1KQ (example; see implementation for exact conversion)
- Month 0 -> 0
- Day 28 -> S
- Hour 8 -> 8
- Minute 15 -> 0F
- Second 44 -> 18
- Ms 123 -> 3F

Timestamp ASCII digits (base36 form) becomes:

```
1KQ0S80F183F
```

### Step C — map base-36 digits to core-36 glyphs

Each base-36 digit maps to one glyph from `docs/alphabet.md`.

This yields `<VIZTIMESTAMP>` (12 glyphs).

### Step D — UUID

UUID structure: `P T T C C R`

- `P`: prefix glyph (from prefix alphabet)
- `TT`: time-mix derived from ms since minute
- `CC`: monotonic counter (starts at 00 for first event in ms)
- `R`: salt digit (process-local)

### Result

```
<VIZTIMESTAMP>-<VIZUUID>
```

In ASCII wire format:

```
20260129081544123-PTTCCR
```

---

## Example 2 — sorting demo

Given three files created in the same second:

- ... .100 (counter 00)
- ... .100 (counter 01)
- ... .101 (counter 00)

Their names will sort in the same order under `ls`:

1) earlier millisecond
2) then by counter within the same millisecond
3) then next millisecond

This is guaranteed by fixed-width encoding and the monotonic counter algorithm.
