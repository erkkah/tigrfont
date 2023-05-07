package tigrfont

import (
	"image"
	"image/color"
)

func contentBounds(img *image.NRGBA) image.Rectangle {
	minNonTransparentRow := img.Bounds().Max.Y
	maxNonTransparentRow := img.Bounds().Min.Y

rows:
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
	cols:
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			color := img.At(x, y)
			if color == Border.C {
				continue cols
			}
			_, _, _, a := color.RGBA()
			if a != 0 {
				// non-transparent
				if y < minNonTransparentRow {
					minNonTransparentRow = y
				}
				if y > maxNonTransparentRow {
					maxNonTransparentRow = y
				}
				continue rows
			}
		}
	}

	return image.Rect(img.Bounds().Min.X, minNonTransparentRow, img.Bounds().Max.X, maxNonTransparentRow+1)
}

func whitePalette() color.Palette {
	p := make(color.Palette, 256)

	for a := 0; a < 256; a++ {
		p[a] = color.NRGBA{255, 255, 255, uint8(a)}
	}

	return p
}

func palettize(img *image.NRGBA) (image.Image, error) {
	palette := whitePalette()
	pi := image.NewPaletted(img.Bounds(), palette)

	for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			alpha := img.Pix[img.PixOffset(x, y)+3]

			pi.Pix[pi.PixOffset(x, y)] = alpha
		}
	}

	return pi, nil
}

func clear(img *image.NRGBA) {
	rgba := [...]uint8{0xff, 0xff, 0xff, 0x00}

	for i := range img.Pix {
		img.Pix[i] = rgba[i%4]
	}
}
