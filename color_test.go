package pixelart

import (
	"math"
	"testing"
)

func TestColorSimIdentical(t *testing.T) {
	if got := ColorSim(RGB{10, 20, 30}, RGB{10, 20, 30}); got != 0 {
		t.Fatalf("ColorSim identical = %v, want 0", got)
	}
}

func TestColorSimKnownDistance(t *testing.T) {
	got := ColorSim(RGB{0, 0, 0}, RGB{3, 4, 0})
	if math.Abs(got-5) > 1e-9 {
		t.Fatalf("ColorSim = %v, want 5", got)
	}
}

func TestSimilarColorNearestAndTies(t *testing.T) {
	palette := []RGB{{0, 0, 0}, {10, 0, 0}, {20, 0, 0}}
	got, err := SimilarColor(RGB{12, 0, 0}, palette)
	if err != nil {
		t.Fatal(err)
	}
	if got != (RGB{10, 0, 0}) {
		t.Fatalf("nearest = %v, want [10,0,0]", got)
	}
	// Exact midpoint between [10] and [20]: distance 5 either way → later wins
	got, err = SimilarColor(RGB{15, 0, 0}, palette)
	if err != nil {
		t.Fatal(err)
	}
	if got != (RGB{20, 0, 0}) {
		t.Fatalf("tie = %v, want [20,0,0]", got)
	}
}

func TestSimilarColorEmpty(t *testing.T) {
	_, err := SimilarColor(RGB{1, 2, 3}, nil)
	if err != ErrEmptyPalette {
		t.Fatalf("err = %v, want ErrEmptyPalette", err)
	}
}
