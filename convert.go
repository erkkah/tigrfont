package tigrfont

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

type Codepage int

const (
	ASCII   Codepage = 0
	CP1252  Codepage = 1252
	UNICODE Codepage = 12001
)

func (cp Codepage) String() string {
	switch cp {
	case ASCII:
		return "ascii"
	case CP1252:
		return "windows 1252"
	case UNICODE:
		return "unicode"
	default:
		return "invalid"
	}
}

func Convert(options Options, font, target string) (int, error) {
	var runeSet []rune
	var err error

	if len(options.Encoding) > 0 && len(options.SampleFile) > 0 {
		return 0, fmt.Errorf("cannot specify both encoding and sample file")
	}

	replaceMode := replaceMissing
	watermark := false

	if len(options.Encoding) > 0 {
		runeSet, err = runesFromEncoding(options.Encoding)
		if err != nil {
			return 0, fmt.Errorf("failed to extract characters from encoding %q: %w", options.Encoding, err)
		}
		replaceMode = removeMissing
		watermark = true
		options.Codepage = UNICODE
	}

	if len(options.SampleFile) > 0 {
		runeSet, err = runesFromFile(options.SampleFile)
		if err != nil {
			return 0, fmt.Errorf("failed to extract characters from sample %q: %w", options.SampleFile, err)
		}
		replaceMode = removeMissing
		watermark = true
		options.Codepage = UNICODE
	}

	if len(runeSet) == 0 {
		const lowChar = 32
		var highChar = 127

		switch options.Codepage {
		case ASCII:
			highChar = 127
		case CP1252:
			highChar = 255
		case UNICODE:
			return 0, fmt.Errorf("use encoding or sample file to create unicode sheet")
		default:
			return 0, fmt.Errorf("invalid TIGR codepage: %v", options.Codepage)
		}

		runeSet = runesFromRange(lowChar, highChar)
	}

	fontBytes, err := os.ReadFile(font)
	if err != nil {
		return 0, fmt.Errorf("failed to load font file: %w", err)
	}

	var img image.Image
	var rendered int

	// Try TTF first
	img, rendered, err = tigrFromTTF(options, fontBytes, runeSet, replaceMode, watermark)

	if err != nil {
		// Assume BDF file
		img, rendered, err = tigrFromBDF(fontBytes, runeSet, replaceMode, watermark)

		if err != nil {
			return 0, fmt.Errorf("failed to render font")
		}
	}

	if watermark {
		img, err = palettize(img)
		if err != nil {
			return 0, err
		}
	}

	pngFile, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0664)
	if err != nil {
		return 0, fmt.Errorf("failed to open target %q: %v", target, err)
	}
	defer pngFile.Close()

	encoder := png.Encoder{
		CompressionLevel: png.BestCompression,
	}

	err = encoder.Encode(pngFile, img)
	if err != nil {
		return 0, fmt.Errorf("failed to encode PNG: %v", err)
	}

	return rendered, nil
}
