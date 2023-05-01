package tigrfont

import (
	"image"
	"image/draw"
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

func shrinkToFit(img *image.NRGBA) *image.NRGBA {
	bounds := contentBounds(img)

	if bounds.Min.Y > 0 {
		bounds.Min.Y--
	}
	if bounds.Min.Y > 0 {
		bounds.Min.Y--
	}

	if bounds.Max.Y < img.Bounds().Dy() {
		bounds.Max.Y++
	}
	if bounds.Max.Y < img.Bounds().Dy() {
		bounds.Max.Y++
	}
	return (img.SubImage(bounds)).(*image.NRGBA)
}

func frame(dest draw.Image, border image.Image) {
	minX := dest.Bounds().Min.X
	minY := dest.Bounds().Min.Y
	maxX := minX + dest.Bounds().Dx()
	maxY := minY + dest.Bounds().Dy()

	draw.Draw(dest, image.Rect(minX, minY, maxX, minY+1), border, image.Point{}, draw.Src)
	draw.Draw(dest, image.Rect(maxX-1, minY, maxX, maxY), border, image.Point{}, draw.Src)
	draw.Draw(dest, image.Rect(minX, maxY-1, maxX, maxY), border, image.Point{}, draw.Src)
	draw.Draw(dest, image.Rect(minX, minY, 1, maxY), border, image.Point{}, draw.Src)
}
