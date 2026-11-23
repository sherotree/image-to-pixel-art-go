package pixelart

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var hexRGB = regexp.MustCompile(`(?i)^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$`)

// HexToRGB parses #rrggbb / rrggbb into RGB, or ok=false if invalid.
func HexToRGB(hex string) (RGB, bool) {
	m := hexRGB.FindStringSubmatch(hex)
	if m == nil {
		return RGB{}, false
	}
	r, _ := strconv.ParseUint(m[1], 16, 8)
	g, _ := strconv.ParseUint(m[2], 16, 8)
	b, _ := strconv.ParseUint(m[3], 16, 8)
	return RGB{uint8(r), uint8(g), uint8(b)}, true
}

// RGBToHex formats RGB as lowercase #rrggbb.
func RGBToHex(rgb RGB) string {
	return fmt.Sprintf("#%02x%02x%02x", rgb[0], rgb[1], rgb[2])
}

// ParsePaletteText parses Paint.NET / .hex / .txt palette text.
// Splits on commas or newlines, skips ";" comments, takes the last 6 hex digits
// of each token.
func ParsePaletteText(text string) []RGB {
	parts := regexp.MustCompile(`[\r\n,]+`).Split(text, -1)
	out := make([]RGB, 0, len(parts))
	for _, s := range parts {
		s = strings.TrimSpace(s)
		if s == "" || s[0] == ';' {
			continue
		}
		last6 := s
		if len(s) > 6 {
			last6 = s[len(s)-6:]
		}
		if rgb, ok := HexToRGB(last6); ok {
			out = append(out, rgb)
		}
	}
	return out
}
