package tigrfont

import (
	"image"
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

	return image.Rect(img.Bounds().Min.X, minNonTransparentRow, img.Bounds().Max.X, maxNonTransparentRow)
}
