# gVIZID — visual timestamps that *sort right* (golang)

**VIZID** is a filename-friendly ID format that combines:

- a **visually encoded timestamp** you can learn to read by eye
- a short **time-entangled UUID** that stays **monotonic within the same millisecond**
- a **Unicode glyph alphabet** that looks great in terminals and file managers

This is **not cryptography**. It’s “human-parseable structure + practical uniqueness”.

---

## What it looks like

A VIZID is:

```
<VIZTIMESTAMP>-<VIZUUID>
```

Where:

- the timestamp is fixed-width and ordered most→least significant
- the UUID is fixed-width and includes a monotonic counter

Result: **`ls` sorting == chronological sorting**.

---

## Filenames and font expectations

VIZID uses **visual glyphs directly in filenames**.

This project assumes a modern Unicode-capable environment:

- macOS Finder
- Windows Explorer
- Linux file managers
- terminal emulators with Nerd Fonts (or equivalent)

All VIZ glyphs are:

- simple geometric Unicode characters
- non-combining
- left-to-right
- safe for filenames on common filesystems

If it renders legibly in common file managers, it’s “in spec”.

---

## Dual representation (VIZ ↔ ASCII)

Every VIZID has a reversible ASCII form for typing, logging, tests, and copy/paste:

- Visual (primary): `<VIZTIMESTAMP>-<VIZUUID>`
- ASCII (canonical wire format): `YYYYMMDDhhmmssmmm-PTTCCR`

The ASCII form is canonical for **implementations/tests**.
The VIZ form is canonical for **filenames**.

---

## CLI (planned + scaffolded)

This repo includes a Go CLI scaffold using **Cobra + Viper**.

- config defaults at: `~/.config/vizid/config.yaml`
- flags:
  - `--config, -c` alternate config file
  - `--timezone, -t` timezone (default `UTC`)
  - `--user-defined, -u` start with all components disabled
  - per-component toggles (bool): `--year --month --day --hour --minute --second --ms --uuid`
  - `--warn` warn if your chosen components would break sort order

---

## Docs

- [`docs/requirements.md`](docs/requirements.md) — what the system must do
- [`docs/tdd.md`](docs/tdd.md) — design + algorithms (base36, counter, UUID)
- [`docs/alphabet.md`](docs/alphabet.md) — the glyph tables (core-36 + UUID prefix set)
- [`docs/cli.md`](docs/cli.md) — CLI behavior, config file, flags
- [`docs/examples.md`](docs/examples.md) — worked examples + sort demonstrations
- [`docs/font-compat.md`](docs/font-compat.md) — practical rendering notes + “known good” fonts

---

## Usage

### Building

```
make build
```

### 'Round Trip'

```
./bin/vizid gen
⊡◭◈□◍⟐□◣□◭◮◢-✱◮◢□□◔

./bin/vizid decode
20260130122520780-@LO00Y

./bin/vizid encode
⊡◭◈□◍⟐□◣□◭◮◢-✱◮◢□□◔
```
---

## Quick status

- ✅ specification + docs (v1)
- ✅ Go project layout + CLI scaffold
- ✅ reference mapping tables (v1)
- ⏳ full implementation (planned next)

---

## Non-goals

- security claims
- resisting brute force
- steganography
- perfect font consistency across every environment

If it’s easy to crack but easy to read and sorts right: **that’s success**.

