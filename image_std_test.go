package pixelart

import (
	"image"
	"image/color"
	"testing"
)

func TestFromImageAndToImageRoundTrip(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 4, 3))
	src.Set(1, 1, color.RGBA{R: 10, G: 20, B: 30, A: 255})
	buf := FromImage(src)
	if buf.Width != 4 || buf.Height != 3 {
		t.Fatalf("size %dx%d", buf.Width, buf.Height)
	}
	i := (1*4 + 1) * 4
	if buf.Data[i] != 10 || buf.Data[i+1] != 20 || buf.Data[i+2] != 30 || buf.Data[i+3] != 255 {
		t.Fatalf("pixel = %v", buf.Data[i:i+4])
	}
	out, err := ToImage(buf)
	if err != nil {
		t.Fatal(err)
	}
	r, g, b, a := out.At(1, 1).RGBA()
	if uint8(r>>8) != 10 || uint8(g>>8) != 20 || uint8(b>>8) != 30 || uint8(a>>8) != 255 {
		t.Fatalf("round-trip failed")
	}
}

func TestImageToPixelArtImage(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 16, 16))
	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			src.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
		}
	}
	out, err := ImageToPixelArtImage(src, Options{BlockSize: 4, Grayscale: true})
	if err != nil {
		t.Fatal(err)
	}
	r, g, b, _ := out.At(0, 0).RGBA()
	if r != g || g != b {
		t.Fatalf("expected grayscale, got %d %d %d", r>>8, g>>8, b>>8)
	}
}
