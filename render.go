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

	startDot := fixed.P(0, metrics.Ascent.Ceil()+1)
	drawer := xfont.Drawer{
		Dst:  image.NewNRGBA(image.Rect(0, 0, 1, 1)),
		Src:  image.White,
		Face: face,
		Dot:  startDot,
	}

	// Render once to measure.
	destWidthPixels, renderedRows, maxAscent, maxDescent, runesRendered := renderFontChars(runes, drawer, 1, mode, watermark)

	rowHeightPixels := maxAscent + maxDescent
	destHeightPixels := rowHeightPixels * renderedRows
	startDot.Y = fixed.I(maxAscent)

	if watermark {
		// Skip initial watermark
		startDot.X = fixed.I(1)
	} else {
		// bottom border
		rowHeightPixels++

		// top border
		destHeightPixels += renderedRows + 1
		startDot.Y += fixed.I(1)
	}

	// Render once more to the actual sheet size
	dest := image.NewNRGBA(image.Rect(0, 0, destWidthPixels, destHeightPixels))
	clear(dest)

	drawer.Dst = dest
	drawer.Dot = startDot
	renderFontChars(runes, drawer, rowHeightPixels, mode, watermark)

	if watermark {
		stampWatermark(drawer.Dst, 0, 0, uint32(runesRendered), uint8(rowHeightPixels))
	} else {
		// left
		drawVerticalDivider(drawer.Dst, 0, 1, destHeightPixels)
		// top
		drawHorizontalDivider(drawer.Dst, 0, destWidthPixels, 0)
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
	yOffset := 1
	if watermark {
		yOffset = 0
	}

	rowIndex := 0

	// Always start on whole pixels
	drawer.Dot.X = fixed.I(drawer.Dot.X.Ceil())

	for _, r := range allRunes {
		if drawer.Dot.X >= maxWidthPixels {
			drawer.Dot.X = 0
			drawer.Dot.Y += rowHeight
			rowIndex++
		}

		xStart := drawer.Dot.X.Ceil() + dstMin.X

		// Skip left border / watermark
		drawer.Dot.X += fixed.I(1)

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

		// Glyph extends to left of box, must be shifted to the right.
		if bounds.Min.X < 0 {
			drawer.Dot.X += -bounds.Min.X
		}

		// Draw the glyph. This adds "advance" to Dot.X.
		drawer.DrawString(s)

		width := bounds.Max.X - bounds.Min.X
		if width == 0 {
			width = fixed.I(1)
		}

		if bounds.Min.X > 0 {
			width += bounds.Min.X
		} else if bounds.Min.X < 0 {
			width -= bounds.Min.X
		}

		if advance > width {
			width = advance
		}

		xEnd := xStart + width.Ceil() + 1
		currentWidth := xEnd + 1
		if totalWidth < currentWidth {
			totalWidth = currentWidth
		}

		drawer.Dot.X = fixed.I(xEnd)
		yStart := rowHeightPixels*rowIndex + dstMin.Y + yOffset
		yEnd := yStart + rowHeightPixels

		if watermark {
			stampWatermark(drawer.Dst, xStart, yStart, uint32(r), uint8(width.Ceil()))
		} else {
			// top
			drawHorizontalDivider(drawer.Dst, xStart, xEnd, yStart-1)
			// right
			drawVerticalDivider(drawer.Dst, xEnd, yStart, yEnd)
			// bottom
			drawHorizontalDivider(drawer.Dst, xStart, xEnd, yEnd-1)
		}

		runesRendered++
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

func stampWatermark(dest draw.Image, x0, y0 int, big uint32, small uint8) {
	mark := [7]uint8{
		uint8(0b10101010),
		uint8(big & 0xff),
		uint8((big >> 8) & 0xff),
		uint8((big >> 16) & 0xff),
		uint8((big >> 24) & 0xff),
		small,
		uint8(0b01010101),
	}

	for i, m := range mark {
		pixel := dest.At(x0, y0+i).(color.NRGBA)
		pixel.A = m
		dest.Set(x0, y0+i, pixel)
	}
}
