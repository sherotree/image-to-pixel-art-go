package pixelart

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"math"
	"os"
	"testing"
)

type parityFile struct {
	NormalizeScale []struct {
		In  int     `json:"in"`
		Out float64 `json:"out"`
	} `json:"normalizeScale"`
	Geometry []struct {
		WorkScale    float64 `json:"workScale"`
		ScaledWidth  float64 `json:"scaledWidth"`
		ScaledHeight float64 `json:"scaledHeight"`
		FinalWidth   float64 `json:"finalWidth"`
		FinalHeight  float64 `json:"finalHeight"`
	} `json:"geometry"`
	ColorSim []struct {
		A   []float64 `json:"a"`
		B   []float64 `json:"b"`
		Out float64   `json:"out"`
	} `json:"colorSim"`
	SimilarColor []struct {
		Color   []float64   `json:"color"`
		Palette [][]float64 `json:"palette"`
		Out     []float64   `json:"out"`
	} `json:"similarColor"`
	Pixelate1x1 struct {
		Width  int     `json:"width"`
		Height int     `json:"height"`
		Data   []uint8 `json:"data"`
	} `json:"pixelate_1x1"`
	PipelineGrayscale dumpImg `json:"pipeline_grayscale"`
	PipelinePalette   dumpImg `json:"pipeline_palette"`
	PipelineNES       dumpImg `json:"pipeline_nes"`
	PipelineAlpha     dumpImg `json:"pipeline_alpha"`
	PipelineCaps      struct {
		Width  int     `json:"width"`
		Height int     `json:"height"`
		Sample []uint8 `json:"sample"`
	} `json:"pipeline_caps"`
	Resize []struct {
		W    int `json:"w"`
		H    int `json:"h"`
		MaxW any `json:"maxW"`
		MaxH any `json:"maxH"`
		Out  struct {
			W int `json:"w"`
			H int `json:"h"`
		} `json:"out"`
	} `json:"resize"`
	GrayscaleSolid    dumpImg   `json:"grayscale_solid"`
	PaletteSnap       dumpImg   `json:"palette_snap"`
	DefaultPalette    [][]uint8 `json:"defaultPalette"`
	PaletteNames      []string  `json:"paletteNames"`
	PixelateGradient7 struct {
		Width   int    `json:"width"`
		Height  int    `json:"height"`
		SHA256  string `json:"sha256"`
		DataB64 string `json:"data_b64"`
		Samples []struct {
			X    int     `json:"x"`
			Y    int     `json:"y"`
			RGBA []uint8 `json:"rgba"`
		} `json:"samples"`
	} `json:"pixelate_gradient_7"`
}

type dumpImg struct {
	Width  int     `json:"width"`
	Height int     `json:"height"`
	Data   []uint8 `json:"data"`
}

func loadParity(t *testing.T) parityFile {
	t.Helper()
	raw, err := os.ReadFile("testdata/parity.json")
	if err != nil {
		t.Fatal(err)
	}
	var p parityFile
	if err := json.Unmarshal(raw, &p); err != nil {
		t.Fatal(err)
	}
	return p
}

func rgbFromFloats(v []float64) RGB {
	return RGB{uint8(v[0]), uint8(v[1]), uint8(v[2])}
}

func TestParityNormalizeAndGeometry(t *testing.T) {
	p := loadParity(t)
	for _, tc := range p.NormalizeScale {
		if got := NormalizeScale(tc.In); got != tc.Out {
			t.Fatalf("NormalizeScale(%d)=%v want %v", tc.In, got, tc.Out)
		}
	}
	inputs := [][3]int{{100, 80, 8}, {400, 200, 10}, {1000, 500, 8}, {320, 320, 8}}
	for i, in := range inputs {
		g := ComputePixelateGeometry(in[0], in[1], in[2])
		want := p.Geometry[i]
		if math.Abs(g.WorkScale-want.WorkScale) > 1e-12 ||
			math.Abs(g.ScaledWidth-want.ScaledWidth) > 1e-9 ||
			math.Abs(g.ScaledHeight-want.ScaledHeight) > 1e-9 ||
			math.Abs(g.FinalWidth-want.FinalWidth) > 1e-9 ||
			math.Abs(g.FinalHeight-want.FinalHeight) > 1e-9 {
			t.Fatalf("geometry[%d]=%+v want %+v", i, g, want)
		}
	}
}

func TestParityColor(t *testing.T) {
	p := loadParity(t)
	for _, tc := range p.ColorSim {
		got := ColorSim(rgbFromFloats(tc.A), rgbFromFloats(tc.B))
		if math.Abs(got-tc.Out) > 1e-9 {
			t.Fatalf("ColorSim=%v want %v", got, tc.Out)
		}
	}
	for _, tc := range p.SimilarColor {
		pal := make([]RGB, len(tc.Palette))
		for i, c := range tc.Palette {
			pal[i] = rgbFromFloats(c)
		}
		got, err := SimilarColor(rgbFromFloats(tc.Color), pal)
		if err != nil {
			t.Fatal(err)
		}
		want := rgbFromFloats(tc.Out)
		if got != want {
			t.Fatalf("SimilarColor=%v want %v", got, want)
		}
	}
}

func assertDumpEqual(t *testing.T, label string, got RGBA, want dumpImg) {
	t.Helper()
	if got.Width != want.Width || got.Height != want.Height {
		t.Fatalf("%s size %dx%d want %dx%d", label, got.Width, got.Height, want.Width, want.Height)
	}
	if len(got.Data) != len(want.Data) {
		t.Fatalf("%s data len %d want %d", label, len(got.Data), len(want.Data))
	}
	for i := range want.Data {
		if got.Data[i] != want.Data[i] {
			t.Fatalf("%s byte[%d]=%d want %d", label, i, got.Data[i], want.Data[i])
		}
	}
}

func TestParityPipelines(t *testing.T) {
	p := loadParity(t)

	src1 := CreateRGBA(1, 1)
	src1.Data[0], src1.Data[1], src1.Data[2], src1.Data[3] = 12, 34, 56, 255
	out1, err := PixelateImage(src1, 8)
	if err != nil {
		t.Fatal(err)
	}
	assertDumpEqual(t, "pixelate_1x1", out1, dumpImg(p.Pixelate1x1))

	gray, err := ImageToPixelArt(solid(16, 16, 255, 0, 0, 255), Options{BlockSize: 4, Grayscale: true})
	if err != nil {
		t.Fatal(err)
	}
	assertDumpEqual(t, "pipeline_grayscale", gray, p.PipelineGrayscale)

	pal, err := ImageToPixelArt(solid(16, 16, 140, 143, 174, 255), Options{
		BlockSize: 4, UsePalette: true, Palette: DefaultPalette,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertDumpEqual(t, "pipeline_palette", pal, p.PipelinePalette)

	nes, ok := PaletteByName("NES")
	if !ok {
		t.Fatal("NES missing")
	}
	nesOut, err := ImageToPixelArt(solid(32, 32, 200, 40, 60, 255), Options{
		BlockSize: 10, UsePalette: true, Palette: nes.Colors,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertDumpEqual(t, "pipeline_nes", nesOut, p.PipelineNES)

	alpha, err := ImageToPixelArt(solid(8, 8, 200, 50, 50, 128), Options{
		BlockSize: 8, UsePalette: true, Palette: []RGB{{255, 0, 0}, {0, 0, 0}},
	})
	if err != nil {
		t.Fatal(err)
	}
	assertDumpEqual(t, "pipeline_alpha", alpha, p.PipelineAlpha)

	caps, err := ImageToPixelArt(solid(200, 100, 10, 20, 30, 255), Options{
		BlockSize: 8, MaxWidth: 50, MaxHeight: 40,
	})
	if err != nil {
		t.Fatal(err)
	}
	if caps.Width != p.PipelineCaps.Width || caps.Height != p.PipelineCaps.Height {
		t.Fatalf("caps size %dx%d want %dx%d", caps.Width, caps.Height, p.PipelineCaps.Width, p.PipelineCaps.Height)
	}
	for i, v := range p.PipelineCaps.Sample {
		if caps.Data[i] != v {
			t.Fatalf("caps sample[%d]=%d want %d", i, caps.Data[i], v)
		}
	}

	gSolid, err := ConvertGrayscale(solid(4, 4, 255, 0, 0, 255))
	if err != nil {
		t.Fatal(err)
	}
	assertDumpEqual(t, "grayscale_solid", gSolid, p.GrayscaleSolid)

	snap, err := ConvertPalette(solid(2, 2, 140, 143, 174, 255), DefaultPalette)
	if err != nil {
		t.Fatal(err)
	}
	assertDumpEqual(t, "palette_snap", snap, p.PaletteSnap)
}

func TestParityGradientPixelate(t *testing.T) {
	p := loadParity(t)
	w, h := 120, 80
	src := CreateRGBA(w, h)
	// Match JS Math.floor((x / W) * 255)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			i := (y*w + x) * 4
			src.Data[i] = uint8(math.Floor(float64(x) / float64(w) * 255))
			src.Data[i+1] = uint8(math.Floor(float64(y) / float64(h) * 255))
			src.Data[i+2] = 128
			src.Data[i+3] = 255
		}
	}

	out, err := PixelateImage(src, 7)
	if err != nil {
		t.Fatal(err)
	}
	if out.Width != p.PixelateGradient7.Width || out.Height != p.PixelateGradient7.Height {
		t.Fatalf("size %dx%d want %dx%d", out.Width, out.Height, p.PixelateGradient7.Width, p.PixelateGradient7.Height)
	}

	sum := sha256.Sum256(out.Data)
	gotHash := hex.EncodeToString(sum[:])
	if gotHash != p.PixelateGradient7.SHA256 {
		wantBytes, err := base64.StdEncoding.DecodeString(p.PixelateGradient7.DataB64)
		if err != nil {
			t.Fatal(err)
		}
		// Find first mismatch for a useful error.
		for i := 0; i < len(out.Data) && i < len(wantBytes); i++ {
			if out.Data[i] != wantBytes[i] {
				t.Fatalf("gradient mismatch at %d: got %d want %d (sha %s vs %s)",
					i, out.Data[i], wantBytes[i], gotHash, p.PixelateGradient7.SHA256)
			}
		}
		t.Fatalf("sha mismatch got %s want %s", gotHash, p.PixelateGradient7.SHA256)
	}

	for _, s := range p.PixelateGradient7.Samples {
		i := (s.Y*out.Width + s.X) * 4
		for c := 0; c < 4; c++ {
			if out.Data[i+c] != s.RGBA[c] {
				t.Fatalf("sample (%d,%d)[%d]=%d want %d", s.X, s.Y, c, out.Data[i+c], s.RGBA[c])
			}
		}
	}
}

func TestParityDefaultPaletteAndNames(t *testing.T) {
	p := loadParity(t)
	if len(DefaultPalette) != len(p.DefaultPalette) {
		t.Fatalf("default palette len %d want %d", len(DefaultPalette), len(p.DefaultPalette))
	}
	for i, c := range p.DefaultPalette {
		if DefaultPalette[i] != (RGB{c[0], c[1], c[2]}) {
			t.Fatalf("DefaultPalette[%d]=%v want %v", i, DefaultPalette[i], c)
		}
	}
	if len(Palettes) != len(p.PaletteNames) {
		t.Fatalf("palette count %d want %d", len(Palettes), len(p.PaletteNames))
	}
	for i, name := range p.PaletteNames {
		if Palettes[i].Name != name {
			t.Fatalf("Palettes[%d].Name=%q want %q", i, Palettes[i].Name, name)
		}
	}
}
