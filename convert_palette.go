package pixelart

// ConvertPaletteInPlace snaps each opaque RGB sample to the nearest palette
// color (Pixel It convertPalette). Alpha is preserved. Mutates image in place.
func ConvertPaletteInPlace(img RGBA, palette []RGB) (RGBA, error) {
	if err := AssertRGBA(img, "image"); err != nil {
		return RGBA{}, err
	}
	if len(palette) == 0 {
		return RGBA{}, ErrEmptyPalette
	}
	d := img.Data
	for i := 0; i < len(d); i += 4 {
		matched, err := SimilarColor(RGB{d[i], d[i+1], d[i+2]}, palette)
		if err != nil {
			return RGBA{}, err
		}
		d[i] = matched[0]
		d[i+1] = matched[1]
		d[i+2] = matched[2]
	}
	return img, nil
}

// ConvertPalette returns a palette-snapped copy of image.
func ConvertPalette(img RGBA, palette []RGB) (RGBA, error) {
	cloned, err := CloneRGBA(img)
	if err != nil {
		return RGBA{}, err
	}
	return ConvertPaletteInPlace(cloned, palette)
}
