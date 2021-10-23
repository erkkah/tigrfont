package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/nfnt/resize"
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
	bg := image.NewUniform(color.NRGBA{0x00, 0x00, 0x00, 0x00})

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
	overSampling := options.overSampling
	ctx.SetDPI(float64(options.dpi) * overSampling)

	src := image.White
	ctx.SetSrc(src)

	cp := charmap.Windows1252
	scale := ctx.PointToFixed(float64(options.fontSize))
	baseline := scale.Ceil()

	// Assume no glyph is landscape
	bufferWidth := int(float64(dest.Bounds().Inset(1).Dy()) * overSampling)
	buffer := image.NewNRGBA(image.Rect(0, 0, bufferWidth, bufferWidth))
	ctx.SetDst(buffer)
	ctx.SetClip(buffer.Bounds())

	xOffset := int(math.Ceil(overSampling))
	buf2dest := func(x int) int {
		return int(float64(x) / overSampling)
	}

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

		if overSampling != 1 {
			drawnWidth := advance.X.Ceil()
			drawnImage := buffer.SubImage(image.Rect(0, 0, drawnWidth, buffer.Bounds().Dy()))
			destWidth := uint(buf2dest(drawnWidth))
			destHeight := uint(dest.Bounds().Inset(1).Dy())
			resized := resize.Resize(destWidth, destHeight, drawnImage, resize.Bicubic)
			buffer = resized.(*image.NRGBA)
		}

		destXOffset := buf2dest(xOffset)
		draw.Draw(dest, image.Rect(destXOffset, 1, destXOffset+buffer.Bounds().Dx(), dest.Bounds().Dy()-1), buffer, image.ZP, draw.Src)

		xOffset += advance.X.Ceil()

		destXOffset = buf2dest(xOffset)
		draw.Draw(dest, image.Rect(destXOffset, 0, destXOffset+1, dest.Bounds().Dy()), border, image.ZP, draw.Src)

		xOffset += int(math.Ceil(overSampling))
	}

	return buf2dest(xOffset), nil
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
