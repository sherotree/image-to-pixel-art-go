package pixelart

import "testing"

func solid(w, h int, r, g, b, a uint8) RGBA {
	img := CreateRGBA(w, h)
	for i := 0; i < len(img.Data); i += 4 {
		img.Data[i] = r
		img.Data[i+1] = g
		img.Data[i+2] = b
		img.Data[i+3] = a
	}
	return img
}

func TestImageToPixelArtGrayscale(t *testing.T) {
	out, err := ImageToPixelArt(solid(16, 16, 255, 0, 0, 255), Options{BlockSize: 4, Grayscale: true})
	if err != nil {
		t.Fatal(err)
	}
	if out.Data[0] != out.Data[1] || out.Data[1] != out.Data[2] {
		t.Fatalf("channels not equal: %v %v %v", out.Data[0], out.Data[1], out.Data[2])
	}
}

func TestImageToPixelArtPalette(t *testing.T) {
	out, err := ImageToPixelArt(solid(16, 16, 140, 143, 174, 255), Options{
		BlockSize:  4,
		UsePalette: true,
		Palette:    DefaultPalette,
	})
	if err != nil {
		t.Fatal(err)
	}
	first := DefaultPalette[0]
	if out.Data[0] != first[0] || out.Data[1] != first[1] || out.Data[2] != first[2] {
		t.Fatalf("got %v %v %v want %v", out.Data[0], out.Data[1], out.Data[2], first)
	}
}

func TestImageToPixelArtEmptyPalette(t *testing.T) {
	_, err := ImageToPixelArt(solid(8, 8, 1, 2, 3, 255), Options{UsePalette: true, Palette: []RGB{}})
	if err != ErrEmptyPalette {
		t.Fatalf("err = %v, want ErrEmptyPalette", err)
	}
}

func TestImageToPixelArtCaps(t *testing.T) {
	out, err := ImageToPixelArt(solid(200, 100, 10, 20, 30, 255), Options{
		BlockSize: 8,
		MaxWidth:  50,
		MaxHeight: 40,
	})
	if err != nil {
		t.Fatal(err)
	}
	if out.Width > 50 || out.Height > 40 {
		t.Fatalf("size %dx%d exceeds caps", out.Width, out.Height)
	}
}

func TestImageToPixelArtPreservesAlpha(t *testing.T) {
	out, err := ImageToPixelArt(solid(8, 8, 200, 50, 50, 128), Options{
		BlockSize:  8,
		UsePalette: true,
		Palette:    []RGB{{255, 0, 0}, {0, 0, 0}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if out.Data[3] != 128 {
		t.Fatalf("alpha = %d, want 128", out.Data[3])
	}
}

func TestPalettesIncludeExpected(t *testing.T) {
	names := map[string]bool{}
	for _, p := range Palettes {
		names[p.Name] = true
	}
	for _, want := range []string{"Pixelit Default", "NES", "Endesga 16"} {
		if !names[want] {
			t.Fatalf("missing palette %q", want)
		}
	}
}
