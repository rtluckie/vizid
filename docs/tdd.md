# VIZID Technical Design (v1)

## Overview

VIZID is a **fixed-length visual identifier** composed of:

- a sortable visual timestamp (milliseconds resolution)
- a time-entangled, monotonic UUID (not cryptographic)

Format:

```
<VIZTIMESTAMP>-<VIZUUID>
```

ASCII wire format (canonical):

```
YYYYMMDDhhmmssmmm-PTTCCR
```

---

## Calendar and resolution

- Calendar: **proleptic Gregorian**
- Resolution: **milliseconds**

### Why milliseconds (not nanoseconds)

Nanoseconds are machine-native but human-useless. Milliseconds:

- shorten the identifier
- reduce transcription errors
- remain “high-resolution” for practical use
- allow monotonic correctness via a counter segment

---

## Sorting guarantee (why `ls` works)

Lexicographic sorting works because:

1. Timestamp fields are fixed-width and ordered most→least significant.
2. `-` cleanly separates timestamp from UUID.
3. Within the same millisecond, UUID ordering is determined by a fixed-width monotonic counter.

---

## Base-36 encoding

### Digit set

| Value | ASCII |
|------:|------|
| 0–9 | `0`–`9` |
| 10–35 | `A`–`Z` |

### Rules

1. All numeric fields are encoded as base-36 integers.
2. Encoding is big-endian (most-significant digit first).
3. All fields are fixed width; pad with `0`.

---

## Timestamp field definitions (v1)

All fields are stored as *calendar values*, but encoded as **zero-based** integers for sorting.

| Component | Range | Stored value | Base-36 width |
|---|---:|---:|---:|
| Year | 0..9999 | year | 3 |
| Month | 1..12 | month-1 | 1 |
| Day | 1..31 | day-1 | 1 |
| Hour | 0..23 | hour | 1 |
| Minute | 0..59 | minute | 2 |
| Second | 0..59 | second | 2 |
| Millisecond | 0..999 | ms | 2 |

Total timestamp width: **12 glyphs**.

> Note: year width 3 supports up to 46655. v1 constrains to 0..9999 for sanity; can be expanded later without changing width.

---

## Visual alphabets

VIZID uses two separate visual alphabets:

1. **Core-36 glyph alphabet** — renders base-36 digits
2. **UUID prefix alphabet** — renders the leading UUID prefix

See `docs/alphabet.md`.

---

## UUID design (v1)

UUID width: **6 glyphs**.

Structure:

```
P T T C C R
```

- `P` (1): prefix glyph (separate visual language)
- `TT` (2): time-mix (opaque, deterministic, timestamp-derived)
- `CC` (2): monotonic counter (base36^2)
- `R` (1): salt digit (base36)

### Monotonic counter algorithm (required)

Definitions:

- `current_ts_ms`: current timestamp in milliseconds
- `last_ts_ms`: last timestamp observed
- `counter`: integer counter (initially 0)

Counter width: 2 base-36 digits. Range 0..1295.

Algorithm:

1. read `current_ts_ms`
2. if `current_ts_ms > last_ts_ms`:
   - `counter = 0`
   - `last_ts_ms = current_ts_ms`
3. else (`current_ts_ms == last_ts_ms`):
   - `counter++`
4. if `counter > 1295`:
   - **BLOCK/SPIN** until `current_ts_ms` advances
   - then reset `counter = 0`

Guarantees:

- strict ordering within a millisecond
- stable lexicographic sorting
- fixed UUID length

### Time-mix (recommended)

Goal: entangle UUID with time while remaining visually opaque.

Recommended input:

- milliseconds since start of minute: `t = (second * 1000) + ms` (0..59999)

Recommended mix:

- apply a cheap, deterministic mixing function (XOR constant + rotate)
- reduce modulo 36^2
- encode as two base-36 digits

Constants may vary, but must be deterministic.

### Salt digit

One base-36 digit chosen at process start (or derived from host fingerprint). Purpose: reduce cross-process collisions.

---

## Dual representation mapping

Implementations MUST provide:

- `EncodeASCIIToVIZ(asciiID) -> vizID`
- `DecodeVIZToASCII(vizID) -> asciiID`

Where ASCII wire format is:

```
YYYYMMDDhhmmssmmm-PTTCCR
```

and VIZ format is:

```
<VIZTIMESTAMP>-<VIZUUID>
```

---

## Component toggles

The CLI may disable components (custom mode). If a user chooses a component set that breaks chronological sorting, CLI must surface it via `--warn`.

