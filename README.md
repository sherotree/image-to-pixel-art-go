# image-to-pixel-art-go

Idiomatic Go library that turns RGBA image buffers (and `image.Image`) into pixel art using the Pixel It algorithm (MIT).

**Reference tool:** [Image to Pixel Art](https://www.uwarp.design/image-to-pixel-art)

Module path: [`github.com/sherotree/image-to-pixel-art-go`](https://pkg.go.dev/github.com/sherotree/image-to-pixel-art-go) (indexed by the Go module proxy after the first tagged release — see [Publishing](#publishing-pkggo.dev)).

TypeScript sibling: [`sherotree/image-to-pixel-art`](https://github.com/sherotree/image-to-pixel-art).

## Install (after the first tagged release)

```bash
go get github.com/sherotree/image-to-pixel-art-go@latest
```

Until a semver tag is published, clone this repo and import the module path locally (Go workspace / `replace`, or just `go test` / `go run` inside the checkout).

## Quick start

```go
package main

import (
	"fmt"

	pixelart "github.com/sherotree/image-to-pixel-art-go"
)

func main() {
	source := pixelart.CreateRGBA(64, 64)
	// …fill source.Data (width × height × 4)…

	nes, _ := pixelart.PaletteByName("NES")
	out, err := pixelart.ImageToPixelArt(source, pixelart.Options{
		BlockSize:  7,
		Grayscale:  false,
		UsePalette: true,
		Palette:    nes.Colors,
		MaxWidth:   128,
		MaxHeight:  128,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(out.Width, out.Height, len(out.Data))
}
```

Pipeline order: **pixelate → optional grayscale → optional palette snap → optional max width/height**.

### `image.Image` support

```go
import "image"

out, err := pixelart.ImageToPixelArtImage(srcImage, pixelart.Options{BlockSize: 8})
// out is *image.RGBA
```

## API

| Export | Role |
| --- | --- |
| `ImageToPixelArt(source, options)` | Full pipeline on `RGBA` |
| `ImageToPixelArtImage(src, options)` | Full pipeline on `image.Image` → `*image.RGBA` |
| `PixelateImage(source, blockSize)` | Block-size pixelation only |
| `ConvertGrayscale` / `ConvertPalette` | Color steps |
| `ResizeImage(image, maxWidth, maxHeight)` | Proportional caps (`0` = unset) |
| `DefaultPalette`, `Palettes` | Built-in palettes |
| `HexToRGB`, `RGBToHex`, `ParsePaletteText` | Palette file / hex helpers |
| `ColorSim`, `SimilarColor` | Euclidean nearest-color match |
| `CreateRGBA`, `CloneRGBA`, `FromImage`, `ToImage` | Buffer helpers |

### Options

- `BlockSize` — 1–50 (Pixel It scale; invalid → 8)
- `Grayscale` — average RGB after pixelation
- `UsePalette` — snap to nearest palette color
- `Palette` — `[]RGB` (defaults to Pixel It default)
- `MaxWidth` / `MaxHeight` — optional output caps (height applied before width when both set; `0` = unset)

## Examples

```bash
go run ./examples/basic
go run ./examples/palette
go run ./examples/grayscale
```

## Checks

```bash
gofmt -l .
go vet ./...
go test ./...
go test -race ./...
go build ./...
```

Parity vectors under `testdata/parity.json` were generated from the TypeScript reference implementation.

## Publishing (pkg.go.dev)

Go does **not** use a separate package-registry upload. Public modules are discovered when:

1. The code lives at the module path URL (`github.com/sherotree/image-to-pixel-art-go`).
2. A **semantic version git tag** (for example `v0.1.0`) is pushed.
3. That version is fetched through **proxy.golang.org**, which caches it and records it on **index.golang.org**.
4. **pkg.go.dev** polls the index and documents the module (usually within a few minutes).

Do **not** retag or rewrite a published version. Coordinate the first tag with the maintainer before releasing.

Exact next steps once publication is approved are listed in [PUBLISHING.md](./PUBLISHING.md).

## License

MIT. Algorithm credit: [giventofly/pixelit](https://github.com/giventofly/pixelit). See [NOTICE](./NOTICE) and [LICENSE](./LICENSE).
