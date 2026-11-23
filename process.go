package pixelart

// ImageToPixelArt runs the full pipeline matching the hosted tool render order:
// pixelate → optional grayscale → optional palette snap → optional max size.
//
// It does not require DOM / Canvas; operate on [RGBA] buffers.
func ImageToPixelArt(source RGBA, opts Options) (RGBA, error) {
	if err := AssertRGBA(source, "source"); err != nil {
		return RGBA{}, err
	}

	out, err := PixelateImage(source, opts.BlockSize)
	if err != nil {
		return RGBA{}, err
	}

	if opts.Grayscale {
		out, err = ConvertGrayscaleInPlace(out)
		if err != nil {
			return RGBA{}, err
		}
	}

	if opts.UsePalette {
		palette := opts.Palette
		if palette == nil {
			palette = DefaultPalette
		}
		out, err = ConvertPaletteInPlace(out, palette)
		if err != nil {
			return RGBA{}, err
		}
	}

	if opts.MaxHeight > 0 || opts.MaxWidth > 0 {
		// Hosted UI applies height then width as separate passes when both are set.
		if opts.MaxHeight > 0 {
			out, err = ResizeImage(out, 0, opts.MaxHeight)
			if err != nil {
				return RGBA{}, err
			}
		}
		if opts.MaxWidth > 0 {
			out, err = ResizeImage(out, opts.MaxWidth, 0)
			if err != nil {
				return RGBA{}, err
			}
		}
	}

	return out, nil
}

// ImageToPixelArtCopy is the same as [ImageToPixelArt], but starts from a deep
// copy so callers can keep the source intact even if future steps mutate
// intermediate buffers.
func ImageToPixelArtCopy(source RGBA, opts Options) (RGBA, error) {
	cloned, err := CloneRGBA(source)
	if err != nil {
		return RGBA{}, err
	}
	return ImageToPixelArt(cloned, opts)
}
