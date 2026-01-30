# Font & Rendering Compatibility

## Goal

VIZID uses Unicode geometric glyphs directly in filenames.
This doc captures practical compatibility notes.

## What matters

- glyphs must not collapse into lookalikes
- glyphs must render at small sizes
- file managers must display them without truncation or replacement boxes

## Recommended fonts

Nerd Fonts generally render the glyph set well.

Commonly good options:

- JetBrainsMono Nerd Font
- Iosevka Nerd Font
- FiraCode Nerd Font

## What to avoid

- fonts that substitute geometric shapes with emoji-style glyphs
- fonts that render some shapes as empty tofu (□ replacement boxes)

## Troubleshooting checklist

1. Confirm your terminal/file manager is not using a fallback emoji font
2. Switch to a Nerd Font
3. Validate the core-36 alphabet table in `docs/alphabet.md` renders distinctly
4. If two glyphs appear identical in your environment, file an issue with:
   - OS + version
   - terminal/file manager
   - font name
   - screenshot
