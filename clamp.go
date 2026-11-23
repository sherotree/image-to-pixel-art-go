package pixelart

import "math"

// clampUint8 matches ECMAScript ToUint8Clamp / Uint8ClampedArray assignment:
// NaN → 0, ≤0 → 0, ≥255 → 255, otherwise round half to even.
func clampUint8(n float64) uint8 {
	if math.IsNaN(n) {
		return 0
	}
	if n <= 0 {
		return 0
	}
	if n >= 255 {
		return 255
	}
	f := math.Floor(n)
	if f+0.5 < n {
		return uint8(f + 1)
	}
	if n < f+0.5 {
		return uint8(f)
	}
	// Exactly halfway: round to even.
	if int64(f)%2 != 0 {
		return uint8(f + 1)
	}
	return uint8(f)
}

// roundHalfUp matches JavaScript Math.round for non-negative values
// (half rounds away from zero / toward +∞).
func roundHalfUp(n float64) int {
	if n < 0 {
		// Not used by the pixel pipeline for dimensions, but keep defined.
		return int(math.Ceil(n - 0.5))
	}
	return int(math.Floor(n + 0.5))
}
