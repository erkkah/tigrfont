package tigrfont

import (
	"fmt"
	"image"

	xfont "golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

func tigrFromTTF(options Options, ttfBytes []byte, runeSet []rune, mode missingGlyphMode, watermark bool) (*image.NRGBA, int, error) {
	font, err := opentype.Parse(ttfBytes)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to parse TTF: %w", err)
	}

	if options.MeasureX {
		options.Measure = "X"
	}

	if len(options.Measure) > 0 {
		measure := []rune(options.Measure)[0]
		options.FontSize, err = getPointSizeFrom(font, options.FontSize, measure)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to measure char %q: %w", measure, err)
		}
	}

	face, err := opentype.NewFace(font, &opentype.FaceOptions{
		DPI:     float64(options.DPI),
		Size:    float64(options.FontSize),
		Hinting: xfont.HintingFull,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create font face: %w", err)
	}

	image, rendered, err := renderFontSheet(runeSet, face, mode, watermark)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to render TTF: %w", err)
	}

	return image, rendered, nil
}

func getPointSizeFrom(font *opentype.Font, fontSize int, char rune) (int, error) {
	face, err := opentype.NewFace(
		font,
		&opentype.FaceOptions{
			DPI: 72.0, Size: float64(fontSize), Hinting: xfont.HintingFull,
		})
	if err != nil {
		return 0, err
	}

	img, rendered, err := renderFontSheet([]rune{char}, face, removeMissing, false)
	if rendered == 0 {
		return 0, fmt.Errorf("cannot measure non-existant char %q", string(char))
	}

	if err != nil {
		return 0, err
	}

	bounds := contentBounds(img)
	actual := float64(bounds.Dy())
	expected := float64(fontSize)
	factor := expected / actual
	return int(expected * factor), nil
}
