package pixelart

import (
	"errors"
	"math"
)

// ErrEmptyPalette is returned when a palette snap is requested with no colors.
var ErrEmptyPalette = errors.New("palette must contain at least one color")

// ColorSim returns the Euclidean distance between two RGB colors (lower is closer).
// Range is roughly [0, ~441.67].
func ColorSim(a, b RGB) float64 {
	var s float64
	for i := 0; i < 3; i++ {
		d := float64(a[i]) - float64(b[i])
		s += d * d
	}
	return math.Sqrt(s)
}

// SimilarColor picks the palette color closest to color.
// Ties prefer later entries, matching Pixel It.
func SimilarColor(color RGB, palette []RGB) (RGB, error) {
	if len(palette) == 0 {
		return RGB{}, ErrEmptyPalette
	}
	best := palette[0]
	bestD := ColorSim(color, best)
	for i := 1; i < len(palette); i++ {
		cur := palette[i]
		d := ColorSim(color, cur)
		if d <= bestD {
			best = cur
			bestD = d
		}
	}
	return best, nil
}
