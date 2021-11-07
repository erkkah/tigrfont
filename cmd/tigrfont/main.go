package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/erkkah/tigrfont"
)

func main() {
	var options tigrfont.Options

	flag.IntVar(&options.FontSize, "size", 18, "TTF font size in points (equals pixels at 72 DPI)")
	flag.BoolVar(&options.Measure, "mx", false, "Measure an 'X' to get TTF point size")
	flag.IntVar(&options.DPI, "dpi", 72, "Render TTF at DPI")
	flag.IntVar(&options.Codepage, "cp", tigrfont.CP1252, "Font sheet codepage, 0 or 1252")
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

	err := tigrfont.Convert(options, font, target)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
