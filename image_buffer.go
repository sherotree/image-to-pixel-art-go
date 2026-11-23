package pixelart

import (
	"fmt"
	"math"
)

// CreateRGBA creates an empty transparent-black buffer of the given size.
// Width and height are floored and clamped to at least 1.
func CreateRGBA(width, height int) RGBA {
	w := width
	h := height
	if w < 1 {
		w = 1
	}
	if h < 1 {
		h = 1
	}
	return RGBA{
		Width:  w,
		Height: h,
		Data:   make([]uint8, w*h*4),
	}
}

// AssertRGBA validates buffer dimensions against data length.
func AssertRGBA(img RGBA, label string) error {
	if label == "" {
		label = "image"
	}
	if math.IsNaN(float64(img.Width)) || math.IsNaN(float64(img.Height)) {
		return fmt.Errorf("%s: width/height must be finite numbers", label)
	}
	if img.Width < 1 || img.Height < 1 {
		return fmt.Errorf("%s: width and height must be >= 1", label)
	}
	expected := img.Width * img.Height * 4
	if len(img.Data) != expected {
		return fmt.Errorf("%s: data length %d != %d", label, len(img.Data), expected)
	}
	return nil
}

// CloneRGBA deep-copies an RGBA buffer.
func CloneRGBA(img RGBA) (RGBA, error) {
	if err := AssertRGBA(img, "image"); err != nil {
		return RGBA{}, err
	}
	data := make([]uint8, len(img.Data))
	copy(data, img.Data)
	return RGBA{Width: img.Width, Height: img.Height, Data: data}, nil
}

// SampleNearest samples a pixel with nearest-neighbor and clamped coordinates.
func SampleNearest(img RGBA, x, y float64) [4]float64 {
	sx := int(math.Floor(x))
	sy := int(math.Floor(y))
	if sx < 0 {
		sx = 0
	}
	if sy < 0 {
		sy = 0
	}
	if sx > img.Width-1 {
		sx = img.Width - 1
	}
	if sy > img.Height-1 {
		sy = img.Height - 1
	}
	i := (sy*img.Width + sx) * 4
	return [4]float64{
		float64(img.Data[i]),
		float64(img.Data[i+1]),
		float64(img.Data[i+2]),
		float64(img.Data[i+3]),
	}
}

// SampleBilinear bilinear-samples (used for the downscale pass).
func SampleBilinear(img RGBA, x, y float64) [4]float64 {
	maxX := img.Width - 1
	maxY := img.Height - 1
	x0 := int(math.Floor(x))
	y0 := int(math.Floor(y))
	if x0 < 0 {
		x0 = 0
	}
	if y0 < 0 {
		y0 = 0
	}
	if x0 > maxX {
		x0 = maxX
	}
	if y0 > maxY {
		y0 = maxY
	}
	x1 := x0 + 1
	y1 := y0 + 1
	if x1 > maxX {
		x1 = maxX
	}
	if y1 > maxY {
		y1 = maxY
	}
	fx := x - float64(x0)
	fy := y - float64(y0)

	i00 := (y0*img.Width + x0) * 4
	i10 := (y0*img.Width + x1) * 4
	i01 := (y1*img.Width + x0) * 4
	i11 := (y1*img.Width + x1) * 4

	var out [4]float64
	for c := 0; c < 4; c++ {
		v00 := float64(img.Data[i00+c])
		v10 := float64(img.Data[i10+c])
		v01 := float64(img.Data[i01+c])
		v11 := float64(img.Data[i11+c])
		top := v00*(1-fx) + v10*fx
		bottom := v01*(1-fx) + v11*fx
		out[c] = top*(1-fy) + bottom*fy
	}
	return out
}

// ResizeBilinear resizes with bilinear filtering (downscale-friendly).
func ResizeBilinear(source RGBA, width, height int) (RGBA, error) {
	if err := AssertRGBA(source, "source"); err != nil {
		return RGBA{}, err
	}
	w := roundHalfUp(float64(width))
	h := roundHalfUp(float64(height))
	if w < 1 {
		w = 1
	}
	if h < 1 {
		h = 1
	}
	out := CreateRGBA(w, h)
	xRatio := float64(source.Width) / float64(w)
	yRatio := float64(source.Height) / float64(h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			samp := SampleBilinear(source, (float64(x)+0.5)*xRatio-0.5, (float64(y)+0.5)*yRatio-0.5)
			i := (y*w + x) * 4
			out.Data[i] = clampUint8(samp[0])
			out.Data[i+1] = clampUint8(samp[1])
			out.Data[i+2] = clampUint8(samp[2])
			out.Data[i+3] = clampUint8(samp[3])
		}
	}
	return out, nil
}

// ResizeNearest resizes with nearest-neighbor (pixel-art upscale).
func ResizeNearest(source RGBA, width, height int) (RGBA, error) {
	if err := AssertRGBA(source, "source"); err != nil {
		return RGBA{}, err
	}
	w := roundHalfUp(float64(width))
	h := roundHalfUp(float64(height))
	if w < 1 {
		w = 1
	}
	if h < 1 {
		h = 1
	}
	out := CreateRGBA(w, h)
	xRatio := float64(source.Width) / float64(w)
	yRatio := float64(source.Height) / float64(h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			samp := SampleNearest(source, float64(x)*xRatio, float64(y)*yRatio)
			i := (y*w + x) * 4
			out.Data[i] = clampUint8(samp[0])
			out.Data[i+1] = clampUint8(samp[1])
			out.Data[i+2] = clampUint8(samp[2])
			out.Data[i+3] = clampUint8(samp[3])
		}
	}
	return out, nil
}
