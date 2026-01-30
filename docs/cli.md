# CLI Specification

The CLI is built with Cobra + Viper.

Binary name: `vizid`

## Config

Default config path:

- `~/.config/vizid/config.yaml`

Override via:

- `--config, -c /path/to/config.yaml`

Config keys (v1):

```yaml
timezone: "UTC"
custom: false
warn: true

components:
  year: true
  month: true
  day: true
  hour: true
  minute: true
  second: true
  ms: true
  uuid: true
```

## Commands

### `vizid gen`

Generate a new VIZID (prints visual form).

Flags:

- `--config, -c` alternate config file
- `--timezone, -tz` timezone (default `UTC`)
- `--custom, -C` enable custom component selection
- `--warn` warn if sort order may break
- component toggles (bool):
  - `--year`
  - `--month`
  - `--day`
  - `--hour`
  - `--minute`
  - `--second`
  - `--ms`
  - `--uuid`

### `vizid decode <vizid>`

Decode a VIZID into its ASCII form:

```
YYYYMMDDhhmmssmmm-PTTCCR
```

### `vizid encode <ascii>`

Encode an ASCII ID into VIZ form.

---

## Sort order warnings

If custom component toggles disable any high-significance timestamp component while leaving lower-significance components enabled, chronological sort order may break.

Examples of risky customization:

- disabling `year` but keeping `month`
- disabling `day` but keeping `hour`

The CLI MUST warn when `--warn=true`.
