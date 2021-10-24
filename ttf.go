package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/math/fixed"
	"golang.org/x/text/encoding/charmap"
)

func ttfToTigr(ttfBytes []byte, lowChar int, highChar int) (*image.NRGBA, error) {
	font, err := freetype.ParseFont(ttfBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse TTF: %w", err)
	}

	ctx := freetype.NewContext()
	ctx.SetFont(font)
	ctx.SetDPI(float64(options.dpi))

	if options.measure {
		options.fontSize, err = getPointSizeFromX(font)
		if err != nil {
			return nil, fmt.Errorf("failed to measure TTF font: %w", err)
		}
	}
	ctx.SetFontSize(float64(options.fontSize))

	image, err := renderTTFSheet(lowChar, highChar, ctx, font)
	if err != nil {
		return nil, fmt.Errorf("failed to render TTF: %w", err)
	}

	return image, nil
}

func renderTTFSheet(lowChar, highChar int, ctx *freetype.Context, font *truetype.Font) (*image.NRGBA, error) {
	bg := image.NewUniform(color.NRGBA{0x00, 0x00, 0xFF, 0x00})

	destHeightPixels := ctx.PointToFixed(float64(options.fontSize)*1.5).Ceil() + 2
	dest := image.NewNRGBA(image.Rect(0, 0, destHeightPixels, destHeightPixels))

	actualWidth, err := renderTTFChars(lowChar, highChar, ctx, font, dest, border, bg)
	if err != nil {
		return nil, err
	}

	dest = image.NewNRGBA(image.Rect(0, 0, actualWidth, destHeightPixels))

	draw.Draw(dest, dest.Bounds(), bg, image.ZP, draw.Src)

	_, err = renderTTFChars(lowChar, highChar, ctx, font, dest, border, bg)
	if err != nil {
		return nil, err
	}

	return dest, nil
}

func renderTTFChars(lowChar, highChar int, ctx *freetype.Context, font *truetype.Font, dest draw.Image, border image.Image, bg image.Image) (int, error) {
	ctx.SetDPI(float64(options.dpi))

	src := image.White
	ctx.SetSrc(src)

	cp := charmap.Windows1252
	scale := ctx.PointToFixed(float64(options.fontSize))
	baseline := scale.Ceil()

	// Assume no glyph is landscape
	bufferWidth := int(float64(dest.Bounds().Inset(1).Dy()))
	buffer := image.NewNRGBA(image.Rect(0, 0, bufferWidth, bufferWidth))

	ctx.SetDst(buffer)
	ctx.SetClip(buffer.Bounds())

	xOffset := 1

	for c := lowChar; c <= highChar; c++ {
		r := cp.DecodeByte(byte(c))
		index := font.Index(r)
		hMetric := font.HMetric(scale, index)
		leftSideAdjustment := fixed.I(0)
		if hMetric.LeftSideBearing < 0 {
			leftSideAdjustment = -hMetric.LeftSideBearing
		}

		// Fill with background
		draw.Draw(buffer, buffer.Bounds(), bg, image.ZP, draw.Src)

		// Draw glyph
		advance, err := ctx.DrawString(string(r), fixed.Point26_6{X: leftSideAdjustment, Y: fixed.I(baseline)})
		if err != nil {
			return 0, err
		}
		advance.X += leftSideAdjustment

		draw.Draw(dest, image.Rect(xOffset, 1, xOffset+buffer.Bounds().Dx(), dest.Bounds().Dy()-1), buffer, image.ZP, draw.Src)

		xOffset += advance.X.Ceil()

		draw.Draw(dest, image.Rect(xOffset, 0, xOffset+1, dest.Bounds().Dy()), border, image.ZP, draw.Src)

		xOffset += 1.0
	}

	return xOffset, nil
}

func getPointSizeFromX(font *truetype.Font) (int, error) {
	ctx := freetype.NewContext()
	ctx.SetFont(font)
	ctx.SetDPI(72.0)
	ctx.SetFontSize(float64(options.fontSize))

	img, err := renderTTFSheet('X', 'X', ctx, font)
	if err != nil {
		return 0, err
	}

	bounds := contentBounds(img)
	actual := float64(bounds.Dy())
	expected := float64(options.fontSize)
	factor := expected / actual
	return int(expected * factor), nil
}
