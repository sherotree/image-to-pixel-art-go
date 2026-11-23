package pixelart

import (
	"math"
	"testing"
)

func TestNormalizeScale(t *testing.T) {
	cases := []struct {
		in   int
		want float64
	}{
		{7, 0.07},
		{50, 0.5},
		{1, 0.01},
		{0, 0.08},
		{51, 0.08},
		{-1, 0.08},
	}
	for _, tc := range cases {
		if got := NormalizeScale(tc.in); got != tc.want {
			t.Fatalf("NormalizeScale(%d) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestComputePixelateGeometrySmall(t *testing.T) {
	g := ComputePixelateGeometry(100, 80, 8)
	if g.WorkScale != 0.08 || g.FinalWidth != 100 || g.FinalHeight != 80 {
		t.Fatalf("unexpected geometry: %+v", g)
	}
}

func TestComputePixelateGeometryBorder(t *testing.T) {
	g := ComputePixelateGeometry(400, 200, 10)
	if g.WorkScale != 0.1 {
		t.Fatalf("workScale = %v, want 0.1", g.WorkScale)
	}
	wantW := 400 + 10/1.5
	if math.Abs(g.FinalWidth-wantW) > 1e-9 {
		t.Fatalf("finalWidth = %v, want %v", g.FinalWidth, wantW)
	}
	if g.FinalHeight != 200 {
		t.Fatalf("finalHeight = %v, want 200", g.FinalHeight)
	}
}

func TestComputePixelateGeometryHalve(t *testing.T) {
	g := ComputePixelateGeometry(1000, 500, 8)
	if g.WorkScale != 0.04 {
		t.Fatalf("workScale = %v, want 0.04", g.WorkScale)
	}
}

func TestPixelateImageRejects(t *testing.T) {
	if _, err := PixelateImage(RGBA{Width: 0, Height: 1, Data: []uint8{}}, 8); err == nil {
		t.Fatal("expected error for empty size")
	}
	if _, err := PixelateImage(RGBA{Width: 2, Height: 2, Data: make([]uint8, 4)}, 8); err == nil {
		t.Fatal("expected error for mismatched data")
	}
}

func TestPixelateImage1x1(t *testing.T) {
	src := CreateRGBA(1, 1)
	src.Data[0], src.Data[1], src.Data[2], src.Data[3] = 12, 34, 56, 255
	out, err := PixelateImage(src, 8)
	if err != nil {
		t.Fatal(err)
	}
	if out.Width != 1 || out.Height != 1 {
		t.Fatalf("size = %dx%d", out.Width, out.Height)
	}
	if out.Data[0] != 12 || out.Data[3] != 255 {
		t.Fatalf("data = %v", out.Data)
	}
}

func TestPixelateImageLength(t *testing.T) {
	src := CreateRGBA(64, 48)
	out, err := PixelateImage(src, 7)
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Data) != out.Width*out.Height*4 {
		t.Fatalf("len=%d want %d", len(out.Data), out.Width*out.Height*4)
	}
}
