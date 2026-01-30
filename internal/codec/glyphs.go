package codec

import "fmt"

// Core36Glyphs maps base-36 value -> rune glyph.
var Core36Glyphs = []rune{
	'□','⊡','⊠','⊞','⊟','◫','◩','◪','■',
	'◇','◈','◊','⟐','⟡','❖','⧫','◆','⬥',
	'△','◬','◭','◮','⟁','▲','◢','◣','◤',
	'○','◌','◍','◐','◑','◒','◓','◔','●',
}

// PrefixASCII maps UUID prefix ASCII -> glyph.
var PrefixASCII = map[byte]rune{
	'~': '✦',
	'!': '✧',
	'@': '✱',
	'$': '✲',
	'%': '✳',
	'^': '✴',
	'&': '✵',
	'*': '✶',
}

// PrefixGlyph maps UUID prefix glyph -> ASCII.
var PrefixGlyph = func() map[rune]byte {
	m := map[rune]byte{}
	for k, v := range PrefixASCII {
		m[v] = k
	}
	return m
}()

var coreGlyphToVal = func() map[rune]int {
	m := map[rune]int{}
	for i, r := range Core36Glyphs {
		m[r] = i
	}
	return m
}()

func CoreValToGlyph(val int) (rune, error) {
	if val < 0 || val >= 36 {
		return 0, fmt.Errorf("invalid base36 value: %d", val)
	}
	return Core36Glyphs[val], nil
}

func CoreGlyphToVal(g rune) (int, error) {
	v, ok := coreGlyphToVal[g]
	if !ok {
		return 0, fmt.Errorf("unknown glyph: %q", string(g))
	}
	return v, nil
}
