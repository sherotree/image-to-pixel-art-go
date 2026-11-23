// Pixelate then convert to grayscale.
//
//	go run ./examples/grayscale
package main

import (
	"fmt"

	pixelart "github.com/sherotree/image-to-pixel-art-go"
)

func main() {
	source := pixelart.CreateRGBA(16, 16)
	for i := 0; i < len(source.Data); i += 4 {
		source.Data[i] = 255
		source.Data[i+1] = 0
		source.Data[i+2] = 0
		source.Data[i+3] = 255
	}

	out, err := pixelart.ImageToPixelArt(source, pixelart.Options{BlockSize: 4, Grayscale: true})
	if err != nil {
		panic(err)
	}
	r, g, b := out.Data[0], out.Data[1], out.Data[2]
	fmt.Printf("grayscale channels r=%d g=%d b=%d\n", r, g, b)
	fmt.Println("equal?", r == g && g == b)
}
