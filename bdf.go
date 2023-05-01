package tigrfont

import (
	"fmt"
	"image"

	"github.com/zachomedia/go-bdf"
)

func tigrFromBDF(bdfBytes []byte, lowChar int, highChar int) (*image.NRGBA, error) {
	font, err := bdf.Parse(bdfBytes)
	if err != nil || font.Size == 0 {
		return nil, fmt.Errorf("failed to parse BDF")
	}

	face := font.NewFace()
	image, err := renderFontSheet(lowChar, highChar, face)
	if err != nil {
		return nil, fmt.Errorf("failed to render BDF: %w", err)
	}
	return image, nil
}
