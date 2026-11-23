package pixelart

// PixelateGeometry holds working scale and destination size using the same
// rules as Pixel It (pixelate() in giventofly/pixelit and the hosted tool port).
type PixelateGeometry struct {
	WorkScale    float64
	ScaledWidth  float64
	ScaledHeight float64
	FinalWidth   float64
	FinalHeight  float64
}

// NormalizeScale maps Pixel It block-size control (1–50) to the internal scale
// factor. Invalid values fall back to 0.08 (block size 8), matching Pixel It.
func NormalizeScale(blockSize int) float64 {
	if blockSize > 0 && blockSize <= 50 {
		return float64(blockSize) * 0.01
	}
	return 0.08
}

// ComputePixelateGeometry computes working scale and destination size.
func ComputePixelateGeometry(width, height, blockSize int) PixelateGeometry {
	workScale := NormalizeScale(blockSize)
	scaledWidth := float64(width) * workScale
	scaledHeight := float64(height) * workScale

	if width > 900 || height > 900 {
		workScale = workScale * 0.5
		scaledWidth = float64(width) * workScale
		scaledHeight = float64(height) * workScale
	}

	finalWidth := float64(width)
	if width > 300 {
		if width > height {
			finalWidth += float64(int(float64(width)/(float64(width)*workScale))) / 1.5
		} else {
			finalWidth += float64(int(float64(width) / (float64(width) * workScale)))
		}
	}

	finalHeight := float64(height)
	if height > 300 {
		if height > width {
			finalHeight += float64(int(float64(height)/(float64(height)*workScale))) / 1.5
		} else {
			finalHeight += float64(int(float64(height) / (float64(height) * workScale)))
		}
	}

	return PixelateGeometry{
		WorkScale:    workScale,
		ScaledWidth:  scaledWidth,
		ScaledHeight: scaledHeight,
		FinalWidth:   finalWidth,
		FinalHeight:  finalHeight,
	}
}

// PixelateImage performs portable pixelation: bilinear downsample → nearest
// upsample. Geometry matches Pixel It; filtering approximates canvas drawImage
// (smooth down, crisp up).
func PixelateImage(source RGBA, blockSize int) (RGBA, error) {
	if err := AssertRGBA(source, "source"); err != nil {
		return RGBA{}, err
	}
	g := ComputePixelateGeometry(source.Width, source.Height, blockSize)

	sw := roundHalfUp(g.ScaledWidth)
	sh := roundHalfUp(g.ScaledHeight)
	dw := roundHalfUp(g.FinalWidth)
	dh := roundHalfUp(g.FinalHeight)
	if sw < 1 {
		sw = 1
	}
	if sh < 1 {
		sh = 1
	}
	if dw < 1 {
		dw = 1
	}
	if dh < 1 {
		dh = 1
	}

	small, err := ResizeBilinear(source, sw, sh)
	if err != nil {
		return RGBA{}, err
	}
	return ResizeNearest(small, dw, dh)
}
