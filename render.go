package tigrfont

import (
	"image"
	"image/color"
	"image/draw"

	xfont "golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type missingGlyphMode int

const (
	replaceMissing missingGlyphMode = 1
	removeMissing  missingGlyphMode = 2
)

func renderFontSheet(runes []rune, face xfont.Face, mode missingGlyphMode, watermark bool) (*image.NRGBA, int, error) {
	metrics := face.Metrics()

	destHeightPixels := metrics.Height.Ceil() + 2
	dest := image.NewNRGBA(image.Rect(0, 0, destHeightPixels, destHeightPixels))

	bg := image.NewUniform(color.NRGBA{0x00, 0x00, 0x00, 0x00})
	draw.Draw(dest, dest.Bounds(), bg, image.Point{}, draw.Src)

	startDot := fixed.P(1, metrics.Ascent.Ceil()+1)
	drawer := xfont.Drawer{
		Dst:  dest,
		Src:  image.White,
		Face: face,
		Dot:  startDot,
	}

	// Render once to measure. We cannot trust metrics from above.
	destWidthPixels, maxAscent, maxDescent, runesRendered := renderFontChars(runes, drawer, mode, false)
	destHeightPixels = maxAscent + maxDescent + 2
	startDot.Y = fixed.I(maxAscent)

	// Render once more to get actual size image
	dest = image.NewNRGBA(image.Rect(0, 0, destWidthPixels, destHeightPixels))
	draw.Draw(dest, dest.Bounds(), bg, image.Point{}, draw.Src)

	drawer.Dst = dest
	drawer.Dot = startDot
	renderFontChars(runes, drawer, mode, watermark)

	drawHorizontalDivider(drawer.Dst, 0)
	if watermark {
		stampWatermark(drawer.Dst, 0, 1, rune(runesRendered))
	}

	return dest, runesRendered, nil
}

func renderFontChars(
	allRunes []rune, drawer xfont.Drawer, mode missingGlyphMode, watermark bool,
) (totalWidth, maxAscent, maxDescent, runesRendered int) {

	dstMin := drawer.Dst.Bounds().Min

	minY := fixed.I(10000)
	maxY := fixed.I(-10000)

	for _, r := range allRunes {
		s := string(r)
		bounds, advance, exists := drawer.Face.GlyphBounds(r)
		if !exists && mode == removeMissing {
			continue
		}

		if bounds.Min.Y < minY {
			minY = bounds.Min.Y
		}
		if bounds.Max.Y > maxY {
			maxY = bounds.Max.Y
		}
		if bounds.Min.X < 0 {
			drawer.Dot.X += -bounds.Min.X
		}
		drawer.DrawString(s)
		width := bounds.Max.X - bounds.Min.X
		if width == 0 {
			width = 1
		}
		if bounds.Min.X > 0 {
			width += bounds.Min.X
		}
		if width > advance {
			drawer.Dot.X -= advance
			drawer.Dot.X += width
		}

		xPos := drawer.Dot.X.Ceil() + dstMin.X

		drawHorizontalDivider(drawer.Dst, xPos)
		if watermark {
			stampWatermark(drawer.Dst, xPos, dstMin.Y+1, r)
		}
		drawer.Dot.X += fixed.I(1)
		runesRendered++
	}

	totalWidth = drawer.Dot.X.Ceil()
	maxAscent = (-minY).Ceil()
	maxDescent = maxY.Ceil()
	return
}

func drawHorizontalDivider(dest draw.Image, x int) {
	draw.Draw(dest, image.Rect(x, dest.Bounds().Min.Y, x+1, dest.Bounds().Max.Y+1), Border, image.Point{}, draw.Src)
}

func stampWatermark(dest draw.Image, x0, y0 int, char rune) {
	mark := [6]uint32{
		uint32(0b10101010),
		uint32(char & 0xff),
		uint32((char >> 8) & 0xff),
		uint32((char >> 16) & 0xff),
		uint32((char >> 24) & 0xff),
		uint32(0b01010101),
	}

	for i, m := range mark {
		pixel := dest.At(x0, y0+i)
		r, g, b, _ := pixel.RGBA()

		dest.Set(x0, y0+i, color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(m)})
	}
}
