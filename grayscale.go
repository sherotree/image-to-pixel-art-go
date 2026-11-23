package pixelart

// ConvertGrayscaleInPlace averages RGB → grayscale (Pixel It convertGrayscale).
// Mutates image in place and returns it.
func ConvertGrayscaleInPlace(img RGBA) (RGBA, error) {
	if err := AssertRGBA(img, "image"); err != nil {
		return RGBA{}, err
	}
	d := img.Data
	for i := 0; i < len(d); i += 4 {
		r := float64(d[i])
		g := float64(d[i+1])
		b := float64(d[i+2])
		avg := (r + g + b) / 3
		v := clampUint8(avg)
		d[i] = v
		d[i+1] = v
		d[i+2] = v
	}
	return img, nil
}

// ConvertGrayscale returns a grayscale copy of image.
func ConvertGrayscale(img RGBA) (RGBA, error) {
	cloned, err := CloneRGBA(img)
	if err != nil {
		return RGBA{}, err
	}
	return ConvertGrayscaleInPlace(cloned)
}
