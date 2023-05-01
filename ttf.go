package tigrfont

import (
	"fmt"
	"image"

	xfont "golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

func tigrFromTTF(options Options, ttfBytes []byte, lowChar int, highChar int) (*image.NRGBA, error) {
	font, err := opentype.Parse(ttfBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse TTF: %w", err)
	}

	if options.Measure {
		options.FontSize, err = getPointSizeFromX(font, options.FontSize)
		if err != nil {
			return nil, fmt.Errorf("failed to measure TTF font: %w", err)
		}
	}

	face, err := opentype.NewFace(font, &opentype.FaceOptions{
		DPI:     float64(options.DPI),
		Size:    float64(options.FontSize),
		Hinting: xfont.HintingFull,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create font face: %w", err)
	}

	image, err := renderFontSheet(lowChar, highChar, face)
	if err != nil {
		return nil, fmt.Errorf("failed to render TTF: %w", err)
	}

	return image, nil
}

func getPointSizeFromX(font *opentype.Font, fontSize int) (int, error) {
	face, err := opentype.NewFace(
		font,
		&opentype.FaceOptions{
			DPI: 72.0, Size: float64(fontSize), Hinting: xfont.HintingFull,
		})
	if err != nil {
		return 0, err
	}

	img, err := renderFontSheet('X', 'X', face)
	if err != nil {
		return 0, err
	}

	bounds := contentBounds(img)
	actual := float64(bounds.Dy())
	expected := float64(fontSize)
	factor := expected / actual
	return int(expected * factor), nil
}
