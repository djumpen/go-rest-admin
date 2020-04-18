package img

import (
	"bytes"
	"encoding/base64"
	"image"
	"io"
	"strings"
)

// ImageEncoder describes type which can encode & save to file image from image.Image
type ImageEncoder interface {
	Encode(w io.Writer, m image.Image) error
	Decode(r io.Reader) (image.Image, error)
}

// Converter is used to detect image ext & convert base64-encoded string to bytes buffer
type Converter struct {
	encoder  ImageEncoder
	imageExt string
}

// GetConverter is used to return correct converter based on image ext
func GetConverter(src string) *Converter {
	c := new(Converter)

	if strings.Contains(src, "data:image/jpeg") {
		c.encoder = new(JpgEncoder)
		c.imageExt = "jpeg"
	} else {
		c.encoder = new(PngEncoder)
		c.imageExt = "png"
	}

	return c
}

// BytesBufferFromString is used to convert base64-encoded image to bytes buffer
func (c *Converter) BytesBufferFromString(src string) (buff *bytes.Buffer, err error) {
	buff = new(bytes.Buffer)

	b64data := src[strings.IndexByte(src, ',')+1:]

	var unbased []byte
	if unbased, err = base64.StdEncoding.DecodeString(b64data); err != nil {
		return
	}

	im, err := c.encoder.Decode(bytes.NewReader(unbased))
	if err != nil {
		return
	}

	err = c.encoder.Encode(buff, im)
	return
}

// GetExt will return image ext
func (c *Converter) GetExt() string {
	return c.imageExt
}
