package img

import (
	"testing"
)

func TestGetConverter(t *testing.T) {
	img := "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=="
	converter := GetConverter(img)

	if converter.GetExt() != "png" {
		t.Error("wrong converter returned")
	}
}

func TestConverter_BytesBufferFromString(t *testing.T) {
	img := "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=="
	converter := GetConverter(img)

	_, err := converter.BytesBufferFromString(img)
	if err != nil {
		t.Errorf("[error] converter.BytesBufferFromString - %s", err.Error())
	}

	img = "data:image/pAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=="
	converter = GetConverter(img)

	_, err = converter.BytesBufferFromString(img)
	if err == nil {
		t.Errorf("converter.BytesBufferFromString - error expected, got nil")
	}
}

func TestConverter_GetExt(t *testing.T) {
	converter := &Converter{}
	converter.imageExt = "png"
	if converter.GetExt() != "png" {
		t.Error("converter.GetExt() - wrong ext returned")
	}
}
