package pixelart

// ResizeImage proportionally resizes with max width / max height caps
// (Pixel It resizeImage). When both caps apply, height ratio overwrites width
// ratio (same precedence as Pixel It). Zero caps are treated as unset.
func ResizeImage(img RGBA, maxWidth, maxHeight int) (RGBA, error) {
	if err := AssertRGBA(img, "image"); err != nil {
		return RGBA{}, err
	}
	if maxWidth == 0 && maxHeight == 0 {
		return img, nil
	}

	ratio := 1.0
	if maxWidth > 0 && img.Width > maxWidth {
		ratio = float64(maxWidth) / float64(img.Width)
	}
	if maxHeight > 0 && img.Height > maxHeight {
		ratio = float64(maxHeight) / float64(img.Height)
	}
	if ratio == 1 {
		return img, nil
	}

	w := roundHalfUp(float64(img.Width) * ratio)
	h := roundHalfUp(float64(img.Height) * ratio)
	if w < 1 {
		w = 1
	}
	if h < 1 {
		h = 1
	}
	return ResizeNearest(img, w, h)
}
