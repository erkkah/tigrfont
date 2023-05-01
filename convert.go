package tigrfont

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

const (
	ASCII  = 0
	CP1252 = 1252
)

func Convert(options Options, font, target string) error {

	fontBytes, err := os.ReadFile(font)
	if err != nil {
		return fmt.Errorf("failed to load font file: %w", err)
	}

	const lowChar = 32
	var highChar = 127

	switch options.Codepage {
	case ASCII:
		highChar = 127
	case CP1252:
		highChar = 255
	default:
		return fmt.Errorf("invalid TIGR codepage: %v", options.Codepage)
	}

	var image *image.NRGBA

	// Try TTF first
	image, err = tigrFromTTF(options, fontBytes, lowChar, highChar)

	if err != nil {
		// Assume BDF file
		image, err = tigrFromBDF(fontBytes, lowChar, highChar)

		if err != nil {
			return fmt.Errorf("failed to render font")
		}
	}

	image = shrinkToFit(image)
	frame(image, Border)

	pngFile, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, 0664)
	if err != nil {
		return fmt.Errorf("failed to open target %q: %v", target, err)
	}
	defer pngFile.Close()

	encoder := png.Encoder{
		CompressionLevel: png.BestCompression,
	}
	err = encoder.Encode(pngFile, image)
	if err != nil {
		return fmt.Errorf("failed to encode PNG: %v", err)
	}
	return nil
}
