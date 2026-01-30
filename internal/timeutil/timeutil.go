package timeutil

import (
	"fmt"
	"regexp"
	"time"
)

var offsetRe = regexp.MustCompile(`^[+-]\d{2}:?\d{2}$`)

// LoadLocation accepts:
// - IANA zone names: America/Chicago
// - UTC offsets: +02:00, -0500
// - UTC
func LoadLocation(spec string) (*time.Location, error) {
	if spec == "" || spec == "UTC" || spec == "Z" {
		return time.UTC, nil
	}
	if offsetRe.MatchString(spec) {
		// normalize +HHMM to +HH:MM
		if len(spec) == 5 && (spec[0] == '+' || spec[0] == '-') {
			spec = spec[:3] + ":" + spec[3:]
		}
		sign := 1
		if spec[0] == '-' {
			sign = -1
		}
		hh := int((spec[1]-'0')*10 + (spec[2] - '0'))
		mm := int((spec[4]-'0')*10 + (spec[5] - '0'))
		off := sign * ((hh * 3600) + (mm * 60))
		return time.FixedZone(spec, off), nil
	}
	loc, err := time.LoadLocation(spec)
	if err != nil {
		return nil, fmt.Errorf("invalid timezone %q: %w", spec, err)
	}
	return loc, nil
}
