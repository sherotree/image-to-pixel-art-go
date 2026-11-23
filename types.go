package pixelart

// RGB is a triplet in the 0–255 channel range.
type RGB [3]uint8

// RGBA is a row-major RGBA raster buffer (4 bytes per pixel).
// Data length must be Width * Height * 4.
type RGBA struct {
	Width  int
	Height int
	Data   []uint8
}

// NamedPalette is a built-in or custom named color set.
type NamedPalette struct {
	Name   string
	Colors []RGB
}

// Options configures the high-level [ImageToPixelArt] pipeline.
// Mirrors the hosted tool controls (block size, grayscale, palette snap, max size).
type Options struct {
	// BlockSize is the Pixel It scale control (1–50). Internally mapped to
	// blockSize * 0.01. Invalid values fall back to 8 (scale 0.08).
	BlockSize int
	// Grayscale converts to average RGB after pixelation when true.
	Grayscale bool
	// UsePalette snaps each pixel to the nearest palette color when true.
	UsePalette bool
	// Palette is used when UsePalette is true. Defaults to [DefaultPalette].
	Palette []RGB
	// MaxWidth is an optional output width cap (px), applied after color steps.
	// Zero means unset.
	MaxWidth int
	// MaxHeight is an optional output height cap (px). When both caps are set,
	// height is applied before width (same precedence as the TypeScript port).
	// Zero means unset.
	MaxHeight int
}
