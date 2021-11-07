package tigrfont

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
)

const (
	ASCII  = 0
	CP1252 = 1252
)

func Convert(options Options, font, target string) error {

	fontBytes, err := ioutil.ReadFile(font)
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

	if bytes.Compare(fontBytes[0:4], []byte("OTTO")) == 0 {
		return fmt.Errorf("open type font not supported")
	}

	var image *image.NRGBA

	if bytes.Compare(fontBytes[0:4], []byte{0, 1, 0, 0}) == 0 {
		image, err = tigrFromTTF(options, fontBytes, lowChar, highChar)
		if err != nil {
			return fmt.Errorf("failed to render TTF: %v", err)
		}
	} else {
		// Assume BDF file
		image, err = tigrFromBDF(fontBytes, lowChar, highChar)

		if err != nil {
			return fmt.Errorf("failed to render BDF: %v", err)
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
