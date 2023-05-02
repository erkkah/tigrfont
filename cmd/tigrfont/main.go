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
	flag.BoolVar(&options.MeasureX, "mx", false, "Measure an 'X' to get TTF point size")
	flag.StringVar(&options.Measure, "m", "", "Measure specified character to get TTF point size")
	flag.IntVar(&options.DPI, "dpi", 72, "Render TTF at DPI")
	flag.IntVar((*int)(&options.Codepage), "cp", (int)(tigrfont.CP1252), "Font sheet codepage: 0, 1252")
	flag.StringVar(&options.Encoding, "encoding", "", "Create sheet using the characters from this encoding")
	flag.StringVar(&options.SampleFile, "sample", "", "Create sheet using the UTF-8 characters in a file")

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

	generated, err := tigrfont.Convert(options, font, target)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	var source string
	if len(options.Encoding) > 0 {
		options.Codepage = tigrfont.UNICODE
		source = fmt.Sprintf(" from encoding %q", options.Encoding)
	} else if len(options.SampleFile) > 0 {
		options.Codepage = tigrfont.UNICODE
		source = fmt.Sprintf(" from sample %q", options.SampleFile)
	}

	fmt.Printf("Generated a %s font sheet for %v characters%s\n", options.Codepage, generated, source)
}
