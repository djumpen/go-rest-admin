package img

import (
	"image"
	"image/jpeg"
	"image/png"
	"io"
)

// PngEncoder is just realization or ImageEncoder intarface based on png.* lib
type PngEncoder struct{}

// Encode encodes image via png.Encode call
func (e *PngEncoder) Encode(w io.Writer, m image.Image) error {
	return png.Encode(w, m)
}

// Decode decodes image via png.Decode call
func (e *PngEncoder) Decode(r io.Reader) (image.Image, error) {
	return png.Decode(r)
}

// JpgEncoder is just realization or ImageEncoder intarface based on jpeg.* lib
type JpgEncoder struct{}

// Encode encodes image via jpeg.Encode call
func (e *JpgEncoder) Encode(w io.Writer, m image.Image) error {
	return jpeg.Encode(w, m, nil)
}

// Decode decodes image via jpeg.Decode call
func (e *JpgEncoder) Decode(r io.Reader) (image.Image, error) {
	return jpeg.Decode(r)
}
