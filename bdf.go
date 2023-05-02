package tigrfont

import (
	"fmt"
	"image"

	"github.com/zachomedia/go-bdf"
)

func tigrFromBDF(bdfBytes []byte, runeSet []rune, mode missingGlyphMode) (*image.NRGBA, int, error) {
	font, err := bdf.Parse(bdfBytes)
	if err != nil || font.Size == 0 {
		return nil, 0, fmt.Errorf("failed to parse BDF")
	}

	face := font.NewFace()
	image, rendered, err := renderFontSheet(runeSet, face, mode)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to render BDF: %w", err)
	}
	return image, rendered, nil
}
