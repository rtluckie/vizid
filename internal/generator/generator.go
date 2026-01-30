package generator

import (
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	"github.com/ryanl/vizid/internal/codec"
	"github.com/ryanl/vizid/internal/model"
	"github.com/ryanl/vizid/internal/timeutil"
)

var (
	mu        sync.Mutex
	lastMs    int64
	counter   int
	saltDigit int // 0..35, chosen at init
)

func init() {
	b := make([]byte, 1)
	_, _ = rand.Read(b)
	saltDigit = int(b[0] % 36)
}

func Generate(opts model.Options) (vizID string, asciiWire string, warnMsg string, err error) {
	loc, err := timeutil.LoadLocation(opts.Timezone)
	if err != nil {
		return "", "", "", err
	}
	now := time.Now().In(loc)
	vizTS, asciiTS, err := encodeTimestamp(now, opts.Components)
	if err != nil {
		return "", "", "", err
	}

	vizUUID, asciiUUID, err := generateUUID(now, loc)
	if err != nil {
		return "", "", "", err
	}
	if !opts.Components.UUID {
		vizUUID = ""
		asciiUUID = ""
	}

	if opts.Warn {
		warnMsg = warnIfSortBroken(opts.Components)
	}

	vizID = vizTS
	asciiWire = asciiTS
	if opts.Components.UUID {
		vizID = vizTS + "-" + vizUUID
		asciiWire = asciiTS + "-" + asciiUUID
	}
	return vizID, asciiWire, warnMsg, nil
}

// encodeTimestamp renders the fixed-width visual timestamp.
// v1 field widths: Year=3, Month=1, Day=1, Hour=1, Minute=2, Second=2, Ms=2  (12 total glyphs).
func encodeTimestamp(t time.Time, c model.Components) (viz string, asciiWire string, err error) {
	year := int64(t.Year())
	if year < 0 || year > 9999 {
		return "", "", fmt.Errorf("year out of v1 range: %d", year)
	}
	month := int64(t.Month()) - 1 // 0..11
	day := int64(t.Day()) - 1     // 0..30
	hour := int64(t.Hour())       // 0..23
	min := int64(t.Minute())      // 0..59
	sec := int64(t.Second())      // 0..59
	ms := int64(t.Nanosecond() / 1e6)

	// ASCII wire timestamp remains decimal digits for human typing/logs.
	asciiWire = fmt.Sprintf("%04d%02d%02d%02d%02d%02d%03d", t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute(), t.Second(), ms)

	// Encode in base36 fixed widths
	y3, err := codec.ToBase36(year, 3); if err != nil { return "", "", err }
	m1, err := codec.ToBase36(month, 1); if err != nil { return "", "", err }
	d1, err := codec.ToBase36(day, 1); if err != nil { return "", "", err }
	h1, err := codec.ToBase36(hour, 1); if err != nil { return "", "", err }
	mn2, err := codec.ToBase36(min, 2); if err != nil { return "", "", err }
	s2, err := codec.ToBase36(sec, 2); if err != nil { return "", "", err }
	ms2, err := codec.ToBase36(ms, 2); if err != nil { return "", "", err }

	b36 := y3 + m1 + d1 + h1 + mn2 + s2 + ms2

	// Map base36 digits to core glyphs
	runes := make([]rune, 0, len(b36))
	for _, r := range b36 {
		val := -1
		if r >= '0' && r <= '9' { val = int(r-'0') }
		if r >= 'A' && r <= 'Z' { val = int(r-'A') + 10 }
		g, e := codec.CoreValToGlyph(val)
		if e != nil { return "", "", e }
		runes = append(runes, g)
	}
	viz = string(runes)

	// TODO: component toggles (custom mode) — v1 scaffold keeps full timestamp.
	_ = c
	return viz, asciiWire, nil
}

func generateUUID(now time.Time, loc *time.Location) (viz string, ascii string, err error) {
	// Prefix: choose from 8 options based on saltDigit (stable per process)
	prefixes := []byte{'~','!','@','$','%','^','&','*'}
	p := prefixes[saltDigit%len(prefixes)]
	pg := codec.PrefixASCII[p]

	// Time-mix uses ms since start of minute
	in := now.In(loc)
	t := int64(in.Second()*1000 + in.Nanosecond()/1e6) // 0..59999
	mixed := mixTime(t)
	// reduce to 36^2
	mod := int64(36 * 36)
	m2 := mixed % mod
	T, err := codec.ToBase36(m2, 2)
	if err != nil { return "", "", err }

	// Counter (monotonic within same millisecond)
	cc, err := nextCounter(in)
	if err != nil { return "", "", err }
	C, err := codec.ToBase36(int64(cc), 2)
	if err != nil { return "", "", err }

	// Salt digit as one base36 char
	R := string(codecDigit(saltDigit))

	ascii = string([]byte{p}) + T + C + R

	// Render ASCII UUID to VIZ:
	// P -> prefix glyph, T/C/R digits -> core glyphs
	vizRunes := []rune{pg}
	for _, ch := range (T + C + R) {
		val := -1
		if ch >= '0' && ch <= '9' { val = int(ch-'0') }
		if ch >= 'A' && ch <= 'Z' { val = int(ch-'A') + 10 }
		g, e := codec.CoreValToGlyph(val)
		if e != nil { return "", "", e }
		vizRunes = append(vizRunes, g)
	}
	viz = string(vizRunes)
	return viz, ascii, nil
}

func mixTime(t int64) int64 {
	// Simple deterministic mix (not cryptographic)
	// XOR constant + rotate-ish via shifts
	x := t ^ 0x5A5A
	x = x + (x >> 3) + (x << 2)
	return x & 0x7FFFFFFF
}

func nextCounter(now time.Time) (int, error) {
	mu.Lock()
	defer mu.Unlock()
	ms := now.UnixMilli()
	if ms > lastMs {
		lastMs = ms
		counter = 0
		return 0, nil
	}
	if ms < lastMs {
		// clock moved backward; in v1, treat as error
		return 0, fmt.Errorf("clock moved backwards: %d < %d", ms, lastMs)
	}
	counter++
	if counter <= 1295 {
		return counter, nil
	}
	// block/spin until next millisecond tick
	for {
		mu.Unlock()
		time.Sleep(time.Microsecond * 200)
		mu.Lock()
		ms2 := time.Now().UnixMilli()
		if ms2 > lastMs {
			lastMs = ms2
			counter = 0
			return 0, nil
		}
	}
}

func codecDigit(val int) byte {
	if val < 0 || val >= 36 { return '0' }
	if val < 10 { return byte('0' + val) }
	return byte('A' + (val - 10))
}

func warnIfSortBroken(c model.Components) string {
	// Very conservative: if you disable a more-significant field but keep any less-significant fields,
	// warn that sorting could break.
	sig := []struct{
		name string
		enabled bool
	}{
		{"year", c.Year},
		{"month", c.Month},
		{"day", c.Day},
		{"hour", c.Hour},
		{"minute", c.Minute},
		{"second", c.Second},
		{"ms", c.Ms},
	}
	seenDisabled := false
	for _, s := range sig {
		if !s.enabled {
			seenDisabled = true
			continue
		}
		if seenDisabled {
			return "custom components disable a more-significant field while keeping a less-significant field; lexicographic sorting may not match time ordering"
		}
	}
	return ""
}
