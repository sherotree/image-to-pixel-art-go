package pixelart

import (
	"image"
	"image/draw"
)

// FromImage converts a standard library image into an [RGBA] buffer.
// The result is always opaque-capable NRGBA/RGBA layout (4 bytes/pixel).
func FromImage(src image.Image) RGBA {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()
	if w < 1 {
		w = 1
	}
	if h < 1 {
		h = 1
	}
	out := CreateRGBA(w, h)

	// Fast path for *image.RGBA.
	if rgba, ok := src.(*image.RGBA); ok && b.Min.Eq(image.Pt(0, 0)) && rgba.Stride == w*4 {
		copy(out.Data, rgba.Pix[:w*h*4])
		return out
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b16, a := src.At(b.Min.X+x, b.Min.Y+y).RGBA()
			i := (y*w + x) * 4
			out.Data[i] = uint8(r >> 8)
			out.Data[i+1] = uint8(g >> 8)
			out.Data[i+2] = uint8(b16 >> 8)
			out.Data[i+3] = uint8(a >> 8)
		}
	}
	return out
}

// ToImage converts an [RGBA] buffer into *image.RGBA.
func ToImage(src RGBA) (*image.RGBA, error) {
	if err := AssertRGBA(src, "source"); err != nil {
		return nil, err
	}
	dst := image.NewRGBA(image.Rect(0, 0, src.Width, src.Height))
	copy(dst.Pix, src.Data)
	return dst, nil
}

// ImageToPixelArtImage runs [ImageToPixelArt] on a standard library image and
// returns *image.RGBA.
func ImageToPixelArtImage(src image.Image, opts Options) (*image.RGBA, error) {
	buf := FromImage(src)
	out, err := ImageToPixelArt(buf, opts)
	if err != nil {
		return nil, err
	}
	return ToImage(out)
}

// DrawOnto draws an [RGBA] buffer onto dst at the given point using Src.
func DrawOnto(dst draw.Image, src RGBA, at image.Point) error {
	img, err := ToImage(src)
	if err != nil {
		return err
	}
	draw.Draw(dst, img.Bounds().Add(at), img, image.Point{}, draw.Src)
	return nil
}
