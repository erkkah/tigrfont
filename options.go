package tigrfont

import (
	"image"
	"image/color"
)

type Options struct {
	FontSize int
	Measure  bool
	DPI      int
	Codepage int
}

var Border = image.NewUniform(color.NRGBA{0x00, 0xAA, 0xCC, 0xff})
