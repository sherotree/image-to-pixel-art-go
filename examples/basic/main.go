// Pixelate a synthetic gradient without DOM / Canvas.
//
//	go run ./examples/basic
package main

import (
	"fmt"

	pixelart "github.com/sherotree/image-to-pixel-art-go"
)

func main() {
	const W, H = 120, 80
	source := pixelart.CreateRGBA(W, H)
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			i := (y*W + x) * 4
			source.Data[i] = uint8(float64(x) / float64(W) * 255)
			source.Data[i+1] = uint8(float64(y) / float64(H) * 255)
			source.Data[i+2] = 128
			source.Data[i+3] = 255
		}
	}

	blockSize := 7
	geom := pixelart.ComputePixelateGeometry(W, H, blockSize)
	out, err := pixelart.ImageToPixelArt(source, pixelart.Options{BlockSize: blockSize})
	if err != nil {
		panic(err)
	}

	fmt.Println("blockSize", blockSize)
	fmt.Printf("geometry %+v\n", geom)
	fmt.Printf("output width=%d height=%d bytes=%d\n", out.Width, out.Height, len(out.Data))
	fmt.Println("sample pixel [0]", out.Data[0], out.Data[1], out.Data[2], out.Data[3])
}
