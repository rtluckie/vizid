# VIZID Requirements (v1)

## Primary goals

1. **Sort correctness**  
   A directory of files named with VIZID must sort chronologically with default lexicographic ordering (e.g. `ls`).

2. **Fixed length**  
   All VIZIDs are the same length (given the same component selection).

3. **Human-decodable timestamp**  
   After learning the system, a human can decode the timestamp without a computer.

4. **Monochrome robustness**  
   Glyphs must be legible regardless of stroke thickness, fill style, and background (given reasonable contrast).

5. **Hand-drawable**  
   The glyphs must be easy to sketch accurately enough for decoding.

6. **Practical uniqueness**  
   UUID need not be cryptographically strong, but must be unique enough for real use and monotonic within the same millisecond.

---

## Sorting requirements (filesystem)

1. Timestamp components MUST appear most-significant → least-significant.
2. Each component MUST be fixed-width (zero padded).
3. Exactly one ASCII hyphen `-` MUST separate timestamp and UUID.
4. The UUID MUST include a fixed-width **monotonic counter** segment that determines ordering within a single millisecond.
5. If the user disables components in a way that could break sorting, the CLI MUST support `--warn` to surface it.

---

## Visual requirements

1. Glyphs MUST be monochrome and inversion-safe (black-on-white / white-on-black).
2. Glyphs MUST NOT rely on:
   - stroke weight
   - dotted outlines
   - subtle shading
3. Glyphs MUST be clearly distinguishable at common UI and terminal sizes.
4. Glyphs SHOULD be visually centered and roughly alphanumeric-sized.
5. Glyphs MUST be safe for filenames on common filesystems.

---

## Character set requirements

1. Core alphabet MUST support base-36 (`0–9`, `A–Z`).
2. Core glyphs must be non-confusable with ASCII alphanumerics.
3. UUID prefix glyphs MUST be distinct from core-36 and form a separate visual language.
4. No glyph may “contain” its decoded meaning (e.g. never use ⑤ to represent 5).

---

## Timestamp requirements

1. Calendar: **proleptic Gregorian** (no custom epoch).
2. Resolution: **milliseconds** (human-friendly; counter handles collisions).
3. Field order: `Year Month Day Hour Minute Second Millisecond`.
4. Timezone MUST be configurable; default is **UTC**.

---

## UUID requirements

1. Fixed length: **6 glyphs** (v1).
2. Structure: `P T T C C R`
   - `P`: prefix from separate visual language
   - `TT`: time-mix (timestamp-derived, opaque)
   - `CC`: monotonic counter (base36^2)
   - `R`: salt (base36)
3. Counter range: 0..1295; on overflow, generator MUST **block/spin** until next millisecond.

---

## Dual representation (VIZ ↔ ASCII)

1. Every VIZID MUST round-trip to ASCII with a 1:1 reversible mapping.
2. ASCII wire format: `YYYYMMDDhhmmssmmm-PTTCCR`
3. The ASCII form is canonical for:
   - tests
   - logging
   - copy/paste
   - interoperability
4. The VIZ form is canonical for filenames.

