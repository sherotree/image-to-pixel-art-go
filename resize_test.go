package pixelart

import "testing"

func TestResizeImageNoCaps(t *testing.T) {
	img := CreateRGBA(10, 10)
	out, err := ResizeImage(img, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if out.Width != 10 || out.Height != 10 {
		t.Fatalf("size = %dx%d", out.Width, out.Height)
	}
}

func TestResizeImageMaxWidth(t *testing.T) {
	img := CreateRGBA(100, 50)
	out, err := ResizeImage(img, 50, 0)
	if err != nil {
		t.Fatal(err)
	}
	if out.Width != 50 || out.Height != 25 {
		t.Fatalf("size = %dx%d, want 50x25", out.Width, out.Height)
	}
}

func TestResizeImageHeightWins(t *testing.T) {
	img := CreateRGBA(100, 100)
	out, err := ResizeImage(img, 50, 25)
	if err != nil {
		t.Fatal(err)
	}
	if out.Width != 25 || out.Height != 25 {
		t.Fatalf("size = %dx%d, want 25x25", out.Width, out.Height)
	}
}

func TestResizeImageWithinCaps(t *testing.T) {
	img := CreateRGBA(20, 20)
	out, err := ResizeImage(img, 100, 100)
	if err != nil {
		t.Fatal(err)
	}
	if out.Width != 20 || out.Height != 20 {
		t.Fatalf("size = %dx%d", out.Width, out.Height)
	}
}
