package pixelart

import "testing"

func TestHexToRGB(t *testing.T) {
	got, ok := HexToRGB("#ff7755")
	if !ok || got != (RGB{255, 119, 85}) {
		t.Fatalf("got %v ok=%v", got, ok)
	}
	got, ok = HexToRGB("ff7755")
	if !ok || got != (RGB{255, 119, 85}) {
		t.Fatalf("got %v ok=%v", got, ok)
	}
	if _, ok := HexToRGB(""); ok {
		t.Fatal("expected invalid empty")
	}
	if _, ok := HexToRGB("#fff"); ok {
		t.Fatal("expected invalid short")
	}
}

func TestRGBToHex(t *testing.T) {
	if got := RGBToHex(RGB{0, 15, 255}); got != "#000fff" {
		t.Fatalf("got %q", got)
	}
}

func TestParsePaletteText(t *testing.T) {
	text := "; Paint.NET Palette\nFF112233\n; skip\naabbcc, #ddeeff\ngarbage\n"
	got := ParsePaletteText(text)
	want := []RGB{{0x11, 0x22, 0x33}, {0xaa, 0xbb, 0xcc}, {0xdd, 0xee, 0xff}}
	if len(got) != len(want) {
		t.Fatalf("len=%d want %d (%v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("[%d] = %v want %v", i, got[i], want[i])
		}
	}
}

func TestParsePaletteTextEmpty(t *testing.T) {
	if got := ParsePaletteText("; only\n; comments"); len(got) != 0 {
		t.Fatalf("got %v", got)
	}
	if got := ParsePaletteText(""); len(got) != 0 {
		t.Fatalf("got %v", got)
	}
}
