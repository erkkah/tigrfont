# tigrfont - bitmap font sheet generator for TIGR

`tigrfont` is a commandline tool for creating bitmap font sheets
for [TIGR] from TTF/OTF or [BDF] font files.

TIGR font sheets are simply PNG files with rows of white characters on a transparent background, separated by single-colored borders:

![](tigrsheet.png)

> This is the default font included in TIGR, and has a black drop shadow. The font sheets created by `tigrfont` don't have shadows.

## Installation

Get [pre-built binaries](https://github.com/erkkah/tigrfont/releases) for Linux, Windows or OSX or install using your local golang setup:

```console
$ go install github.com/erkkah/tigrfont
```

## BDF to TIGR

Creating font sheets from BDF files is straightforward, since they are bitmap fonts already:
```console
$ tigrfont 5x7.bdf 5x7.png
```

## TTF to TIGR

Converting from TTF files often requires a bit more testing and tweaking, depending on the specifics of the font.

Since TTF fonts are vector fonts, they are rendered to a bitmap before being exported as the final font sheet.

The rendering uses anti-aliasing, which will cause visible semi-transparent smudges at the low resolutions typically used with TIGR. 

YMMV :car:

### Font resolution and size

The font is rendered at a given dpi, by default 72.

The font size is specified in points, by default 18.

Since apparent character height for a given point size varies a lot between fonts, `tigrfont` can measure the height of an 'X' and adjust the effective point size to make the 'X' render with a height of the given point size.

For example, running
```console
$ tigrfont -mx -size 20 myfont.ttf myfont.png
```
will render a font sheet at a size where a capital 'X' is 20 pixels high, since pixels equal points at 72 DPI.

> You can also use `-m <char>` to measure using any character.

## Unicode and codepages

TIGR, and `tigrfont` traditionally support two codepages, ASCII (code points 32 to 127) and CP-1252 (code points 32 to 255).
Font sheets created using CP-1252 (the default) are loaded like this:

```C
Tigr* fontImage = tigrLoadImage("font.png");
TigrFont* font = tigrLoadFont(fontImage, TCP_1252);
```

### Unicode sheets

Since version 1.0, `tigrfont` supports sparse unicode-encoded font sheets.
Note that [TIGR] version 3.1 is needed to use these font sheets.

Instead of simply enumerating code points, as in the ASCII and CP-1252 cases, the set of code points to include is specified using either the `-encoding` or the `-sample` option.

The `-encoding` argument accepts an [HTML5 encoding name] and tries to extract the set of code points covered by that encoding. This often generates a superset of code points. For example, specifying "gbk" or "gb3212" results in the same large set of code points.

Using the `-sample` option, you can specify a UTF-8 encoded text file containing the code points you want in the font sheet. Since duplicates are allowed, you can simply specify a sample text file with the code points needed.

> If you look at the generated unicode font sheets, you might notice that there are semi-transparent sections in the borders around the characters. This is since the alpha channel is used to store code point info. :brain:

[HTML5 encoding name]: https://encoding.spec.whatwg.org/#names-and-labels
[TIGR]: https://github.com/erkkah/tigr
[BDF]: https://en.wikipedia.org/wiki/Glyph_Bitmap_Distribution_Format
