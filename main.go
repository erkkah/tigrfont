package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"

	"github.com/nfnt/resize"
)

var options struct {
	fontSize     int
	measure      bool
	dpi          int
	overSampling float64
}

func main() {
	flag.IntVar(&options.fontSize, "size", 12, "TTF font size in points (equals pixels at 72 DPI)")
	flag.BoolVar(&options.measure, "mx", false, "Measure an 'X' to adjust TTF point size")
	flag.IntVar(&options.dpi, "dpi", 72, "Render TTF at DPI")
	flag.Float64Var(&options.overSampling, "over", 1, "TTF oversampling factor")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: tigrfont [options] <source BDF/TTF> <target PNG>\n\nOptions:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		flag.Usage()
		os.Exit(1)
	}
	font := args[0]
	target := args[1]

	fontBytes, err := ioutil.ReadFile(font)
	if err != nil {
		fmt.Printf("Failed to load font file: %v\n", err)
		os.Exit(1)
	}

	const lowChar = 32
	const highChar = 255

	if bytes.Compare(fontBytes[0:4], []byte("OTTO")) == 0 {
		fmt.Printf("Open type font not supported.\n")
		os.Exit(1)
	}

	var image *image.NRGBA

	if bytes.Compare(fontBytes[0:4], []byte{0, 1, 0, 0}) == 0 {
		image, err = ttfToTigr(fontBytes, lowChar, highChar)
		if err != nil {
			fmt.Printf("Failed to render TTF: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Assume BDF file
		image, err = bdfToTigr(fontBytes, lowChar, highChar)

		if err != nil {
			fmt.Printf("Failed to render BDF: %v\n", err)
			os.Exit(1)
		}
	}

	image = shrinkToFit(image)
	frame(image, border)

	_ = resize.Bicubic

	pngFile, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, 0664)
	if err != nil {
		fmt.Printf("Failed to open target %q: %v\n", target, err)
		os.Exit(1)
	}

	encoder := png.Encoder{
		CompressionLevel: png.BestCompression,
	}
	err = encoder.Encode(pngFile, image)
	if err != nil {
		fmt.Printf("Failed to encode PNG: %v\n", err)
		os.Exit(1)
	}
	pngFile.Close()
}

var border = image.NewUniform(color.NRGBA{0xff, 0xff, 0x00, 0xff})

func contentBounds(img *image.NRGBA) image.Rectangle {
	minNonTransparentRow := img.Bounds().Max.Y
	maxNonTransparentRow := img.Bounds().Min.Y

rows:
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
	cols:
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			color := img.At(x, y)
			if color == border.C {
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

	draw.Draw(dest, image.Rect(minX, minY, maxX, minY+1), border, image.ZP, draw.Src)
	draw.Draw(dest, image.Rect(maxX-1, minY, maxX, maxY), border, image.ZP, draw.Src)
	draw.Draw(dest, image.Rect(minX, maxY-1, maxX, maxY), border, image.ZP, draw.Src)
	draw.Draw(dest, image.Rect(minX, minY, 1, maxY), border, image.ZP, draw.Src)
}
