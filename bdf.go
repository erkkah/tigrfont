package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"github.com/zachomedia/go-bdf"
	xfont "golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"golang.org/x/text/encoding/charmap"
)

func bdfToTigr(bdfBytes []byte, lowChar int, highChar int) (*image.NRGBA, error) {
	font, err := bdf.Parse(bdfBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse BDF: %w", err)
	}

	image, err := renderBDFSheet(lowChar, highChar, font)
	if err != nil {
		return nil, fmt.Errorf("failed to render BDF: %w", err)
	}
	return image, nil
}

func renderBDFSheet(lowChar, highChar int, font *bdf.Font) (*image.NRGBA, error) {
	face := font.NewFace()
	metrics := face.Metrics()

	allChars := make([]byte, highChar-lowChar+1)
	for i := range allChars {
		allChars[i] = byte(lowChar + i)
	}
	destHeightPixels := metrics.Height.Ceil() + 2
	dest := image.NewNRGBA(image.Rect(0, 0, destHeightPixels, destHeightPixels))

	bg := image.NewUniform(color.NRGBA{0x00, 0x00, 0x00, 0x00})
	draw.Draw(dest, dest.Bounds(), bg, image.ZP, draw.Src)

	startDot := fixed.P(1, metrics.Ascent.Ceil()+1)
	drawer := xfont.Drawer{
		Dst:  dest,
		Src:  image.White,
		Face: face,
		Dot:  startDot,
	}

	// Render once to measure width
	destWidthPixels := renderBDFChars(allChars, drawer)

	// Render once more to get actual width image
	dest = image.NewNRGBA(image.Rect(0, 0, destWidthPixels, destHeightPixels))
	draw.Draw(dest, dest.Bounds(), bg, image.ZP, draw.Src)

	drawer.Dst = dest
	drawer.Dot = startDot
	renderBDFChars(allChars, drawer)

	return dest, nil
}

func renderBDFChars(allChars []byte, drawer xfont.Drawer) int {
	cp := charmap.Windows1252

	min := drawer.Dst.Bounds().Min
	max := drawer.Dst.Bounds().Max

	for _, c := range allChars {
		r := cp.DecodeByte(byte(c))
		s := string(r)
		bounds, advance, _ := drawer.Face.GlyphBounds(r)
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

		xPos := drawer.Dot.X.Ceil() + min.X

		draw.Draw(drawer.Dst, image.Rect(xPos, min.Y, xPos+1, max.Y+1), border, image.ZP, draw.Src)
		drawer.Dot.X += fixed.I(1)
	}

	return drawer.Dot.X.Ceil()
}
