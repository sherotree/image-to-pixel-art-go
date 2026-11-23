// Package pixelart turns RGBA image buffers into pixel art using the Pixel It
// algorithm (MIT), matching the portable TypeScript library semantics.
//
// Reference tool page: https://www.uwarp.design/image-to-pixel-art
//
// Pipeline order for [ImageToPixelArt]:
//
//	pixelate → optional grayscale → optional palette snap → optional max size
//
// The package operates on an [RGBA] buffer (row-major, 4 bytes per pixel) and
// also accepts standard library [image.Image] values via [FromImage] /
// [ImageToPixelArtImage].
//
// Upstream algorithm: https://github.com/giventofly/pixelit (MIT).
package pixelart
