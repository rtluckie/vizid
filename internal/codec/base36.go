package codec

import "fmt"

const digits = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func ToBase36(n int64, width int) (string, error) {
	if n < 0 {
		return "", fmt.Errorf("negative values not supported in v1: %d", n)
	}
	if width <= 0 {
		return "", fmt.Errorf("invalid width: %d", width)
	}
	buf := make([]byte, width)
	v := n
	for i := width - 1; i >= 0; i-- {
		d := v % 36
		buf[i] = digits[d]
		v = v / 36
	}
	if v != 0 {
		return "", fmt.Errorf("value %d overflows width %d base36 digits", n, width)
	}
	return string(buf), nil
}

func FromBase36(s string) (int64, error) {
	var v int64
	for _, r := range s {
		idx := int64(indexOf(byte(r)))
		if idx < 0 {
			return 0, fmt.Errorf("invalid base36 digit: %q", r)
		}
		v = v*36 + idx
	}
	return v, nil
}

func indexOf(b byte) int {
	if b >= '0' && b <= '9' {
		return int(b - '0')
	}
	if b >= 'A' && b <= 'Z' {
		return int(b-'A') + 10
	}
	return -1
}
