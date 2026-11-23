// Pixelate then snap to the NES palette.
//
//	go run ./examples/palette
package main

import (
	"fmt"

	pixelart "github.com/sherotree/image-to-pixel-art-go"
)

func main() {
	nes, ok := pixelart.PaletteByName("NES")
	if !ok {
		panic("NES palette missing")
	}

	source := pixelart.CreateRGBA(32, 32)
	for i := 0; i < len(source.Data); i += 4 {
		source.Data[i] = 200
		source.Data[i+1] = 40
		source.Data[i+2] = 60
		source.Data[i+3] = 255
	}

	out, err := pixelart.ImageToPixelArt(source, pixelart.Options{
		BlockSize:  10,
		UsePalette: true,
		Palette:    nes.Colors,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("snapped color", pixelart.RGBToHex(pixelart.RGB{out.Data[0], out.Data[1], out.Data[2]}))
	fmt.Println("palette size", len(nes.Colors))
}
