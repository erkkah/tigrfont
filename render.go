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

	bg := image.NewUniform(color.NRGBA{0x00, 0x00, 0x00, 0x00})

	startDot := fixed.P(1, metrics.Ascent.Ceil()+1)
	drawer := xfont.Drawer{
		Dst:  image.NewRGBA(image.Rect(0, 0, 1, 1)),
		Src:  image.White,
		Face: face,
		Dot:  startDot,
	}

	// Render once to measure.
	destWidthPixels, renderedRows, maxAscent, maxDescent, runesRendered := renderFontChars(runes, drawer, 1, mode, false)
	// Row height, including bottom border
	rowHeightPixels := maxAscent + maxDescent + 1
	// Plus one for top border
	destHeightPixels := rowHeightPixels*renderedRows + 1
	// Plus one for top border
	startDot.Y = fixed.I(maxAscent + 1)

	// Render once more to the actual sheet size
	dest := image.NewNRGBA(image.Rect(0, 0, destWidthPixels, destHeightPixels))
	draw.Draw(dest, dest.Bounds(), bg, image.Point{}, draw.Src)

	drawer.Dst = dest
	drawer.Dot = startDot
	renderFontChars(runes, drawer, rowHeightPixels, mode, watermark)

	// left
	drawVerticalDivider(drawer.Dst, 0, 1, destHeightPixels)
	// top
	drawHorizontalDivider(drawer.Dst, 0, destWidthPixels, 0)

	if watermark {
		stampWatermark(drawer.Dst, 0, 1, rune(runesRendered))
	}

	return dest, runesRendered, nil
}

func renderFontChars(
	allRunes []rune, drawer xfont.Drawer, rowHeightPixels int, mode missingGlyphMode, watermark bool,
) (totalWidth, rows, maxAscent, maxDescent, runesRendered int) {

	dstMin := drawer.Dst.Bounds().Min

	minGlyphY := fixed.I(10000)
	maxGlyphY := fixed.I(-10000)

	maxWidthPixels := fixed.I(1000)
	rowHeight := fixed.I(rowHeightPixels)
	rowIndex := 0

	for _, r := range allRunes {
		if drawer.Dot.X >= maxWidthPixels {
			drawer.Dot.X = fixed.I(1)
			drawer.Dot.Y += rowHeight
			rowIndex++
		}

		s := string(r)
		bounds, advance, exists := drawer.Face.GlyphBounds(r)
		if !exists && mode == removeMissing {
			continue
		}

		if bounds.Min.Y < minGlyphY {
			minGlyphY = bounds.Min.Y
		}
		if bounds.Max.Y > maxGlyphY {
			maxGlyphY = bounds.Max.Y
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
		} else if bounds.Min.X < 0 {
			width -= bounds.Min.X
		}
		if width > advance {
			drawer.Dot.X -= advance
			drawer.Dot.X += width
		} else {
			width = advance
		}

		xEnd := drawer.Dot.X.Ceil() + dstMin.X
		xStart := xEnd - width.Ceil()
		yStart := rowHeightPixels*rowIndex + dstMin.Y + 1
		yEnd := yStart + rowHeightPixels

		// top
		drawHorizontalDivider(drawer.Dst, xStart, xEnd+1, yStart-1)
		// right
		drawVerticalDivider(drawer.Dst, xEnd, yStart, yEnd)
		// bottom
		drawHorizontalDivider(drawer.Dst, xStart, xEnd+1, yEnd-1)
		if watermark {
			stampWatermark(drawer.Dst, xEnd, yStart, r)
		}
		drawer.Dot.X += fixed.I(1)
		runesRendered++

		currentWidth := drawer.Dot.X.Ceil()
		if totalWidth < currentWidth {
			totalWidth = currentWidth
		}
	}

	rows = rowIndex + 1
	maxAscent = (-minGlyphY).Ceil()
	maxDescent = maxGlyphY.Ceil()
	return
}

func drawVerticalDivider(dest draw.Image, x, y0, y1 int) {
	draw.Draw(dest, image.Rect(x, y0, x+1, y1), Border, image.Point{}, draw.Src)
}

func drawHorizontalDivider(dest draw.Image, x0, x1, y int) {
	draw.Draw(dest, image.Rect(x0, y, x1, y+1), Border, image.Point{}, draw.Src)
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
