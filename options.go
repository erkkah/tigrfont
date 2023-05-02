package tigrfont

import (
	"image"
	"image/color"
)

type Options struct {
	FontSize   int
	MeasureX   bool
	Measure    string
	DPI        int
	Codepage   Codepage
	Encoding   string
	SampleFile string
}

var Border = image.NewUniform(color.NRGBA{0x00, 0xAA, 0xCC, 0xff})
