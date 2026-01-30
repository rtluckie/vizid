package codec

import (
	"fmt"
	"strconv"
	"strings"
)

// ASCII wire format: YYYYMMDDhhmmssmmm-PTTCCR
// VIZ format:       <12 glyphs>-<6 glyphs>
//
// Timestamp glyph layout (12):
//   Year(3) Month(1) Day(1) Hour(1) Minute(2) Second(2) Millis(2)
//
// UUID glyph layout (6):
//   Prefix(1) + Core36(5)  -> PTTCCR in ASCII
func DecodeVIZToASCII(viz string) (string, error) {
	parts := strings.Split(viz, "-")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid VIZID: expected single '-' delimiter")
	}
	tsGlyphs := []rune(parts[0])
	uidGlyphs := []rune(parts[1])
	if len(tsGlyphs) != 12 {
		return "", fmt.Errorf("invalid timestamp glyph length: got %d, want 12", len(tsGlyphs))
	}
	if len(uidGlyphs) != 6 {
		return "", fmt.Errorf("invalid uuid glyph length: got %d, want 6", len(uidGlyphs))
	}

	// Decode timestamp glyphs -> base36 ASCII digits
	b36 := make([]byte, 12)
	for i, g := range tsGlyphs {
		val, err := CoreGlyphToVal(g)
		if err != nil {
			return "", fmt.Errorf("timestamp: %w", err)
		}
		b36[i] = valToASCII(val)
	}

	year, err := FromBase36(string(b36[0:3])); if err != nil { return "", err }
	month0, err := FromBase36(string(b36[3:4])); if err != nil { return "", err }
	day0, err := FromBase36(string(b36[4:5])); if err != nil { return "", err }
	hour, err := FromBase36(string(b36[5:6])); if err != nil { return "", err }
	minute, err := FromBase36(string(b36[6:8])); if err != nil { return "", err }
	second, err := FromBase36(string(b36[8:10])); if err != nil { return "", err }
	ms, err := FromBase36(string(b36[10:12])); if err != nil { return "", err }

	// Convert stored values back to calendar values
	month := month0 + 1
	day := day0 + 1

	asciiTS := fmt.Sprintf("%04d%02d%02d%02d%02d%02d%03d",
		year, month, day, hour, minute, second, ms)

	// Decode UUID
	pAscii, ok := PrefixGlyph[uidGlyphs[0]]
	if !ok {
		return "", fmt.Errorf("uuid: unknown prefix glyph %q", string(uidGlyphs[0]))
	}
	b := []byte{pAscii}
	for i := 1; i < 6; i++ {
		val, err := CoreGlyphToVal(uidGlyphs[i])
		if err != nil {
			return "", fmt.Errorf("uuid: %w", err)
		}
		b = append(b, valToASCII(val))
	}
	asciiUUID := string(b)

	return asciiTS + "-" + asciiUUID, nil
}

func EncodeASCIIToVIZ(ascii string) (string, error) {
	parts := strings.Split(ascii, "-")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid ASCII ID: expected single '-' delimiter")
	}
	ts := parts[0]
	uuid := parts[1]
	if len(ts) != 17 {
		return "", fmt.Errorf("timestamp must be 17 chars YYYYMMDDhhmmssmmm")
	}
	if len(uuid) != 6 {
		return "", fmt.Errorf("uuid must be 6 chars PTTCCR")
	}

	// Parse decimal timestamp components
	year, err := strconv.Atoi(ts[0:4]); if err != nil { return "", err }
	month, err := strconv.Atoi(ts[4:6]); if err != nil { return "", err }
	day, err := strconv.Atoi(ts[6:8]); if err != nil { return "", err }
	hour, err := strconv.Atoi(ts[8:10]); if err != nil { return "", err }
	minute, err := strconv.Atoi(ts[10:12]); if err != nil { return "", err }
	second, err := strconv.Atoi(ts[12:14]); if err != nil { return "", err }
	ms, err := strconv.Atoi(ts[14:17]); if err != nil { return "", err }

	if year < 0 || year > 9999 { return "", fmt.Errorf("year out of v1 range: %d", year) }
	if month < 1 || month > 12 { return "", fmt.Errorf("invalid month: %d", month) }
	if day < 1 || day > 31 { return "", fmt.Errorf("invalid day: %d", day) }
	if hour < 0 || hour > 23 { return "", fmt.Errorf("invalid hour: %d", hour) }
	if minute < 0 || minute > 59 { return "", fmt.Errorf("invalid minute: %d", minute) }
	if second < 0 || second > 59 { return "", fmt.Errorf("invalid second: %d", second) }
	if ms < 0 || ms > 999 { return "", fmt.Errorf("invalid ms: %d", ms) }

	// Stored values
	month0 := int64(month - 1)
	day0 := int64(day - 1)

	y3, err := ToBase36(int64(year), 3); if err != nil { return "", err }
	m1, err := ToBase36(month0, 1); if err != nil { return "", err }
	d1, err := ToBase36(day0, 1); if err != nil { return "", err }
	h1, err := ToBase36(int64(hour), 1); if err != nil { return "", err }
	mn2, err := ToBase36(int64(minute), 2); if err != nil { return "", err }
	s2, err := ToBase36(int64(second), 2); if err != nil { return "", err }
	ms2, err := ToBase36(int64(ms), 2); if err != nil { return "", err }

	b36 := y3 + m1 + d1 + h1 + mn2 + s2 + ms2
	vizTS, err := mapBase36ToGlyphs(b36)
	if err != nil { return "", err }

	// UUID
	p := uuid[0]
	pg, ok := PrefixASCII[p]
	if !ok {
		return "", fmt.Errorf("unknown UUID prefix: %q", string(p))
	}
	vizUUID := []rune{pg}
	glyphs, err := mapBase36ToGlyphs(uuid[1:])
	if err != nil { return "", err }
	vizUUID = append(vizUUID, []rune(glyphs)...)

	return vizTS + "-" + string(vizUUID), nil
}

func mapBase36ToGlyphs(b36 string) (string, error) {
	r := make([]rune, 0, len(b36))
	for _, ch := range b36 {
		val := -1
		if ch >= '0' && ch <= '9' { val = int(ch - '0') }
		if ch >= 'A' && ch <= 'Z' { val = int(ch-'A') + 10 }
		if val < 0 || val >= 36 {
			return "", fmt.Errorf("invalid base36 char: %q", string(ch))
		}
		g, err := CoreValToGlyph(val)
		if err != nil { return "", err }
		r = append(r, g)
	}
	return string(r), nil
}

func valToASCII(val int) byte {
	if val < 10 {
		return byte('0' + val)
	}
	return byte('A' + (val - 10))
}
