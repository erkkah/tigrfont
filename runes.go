package tigrfont

import (
	"os"
	"sort"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/htmlindex"
)

func runesFromEncoding(encodingName string) ([]rune, error) {
	// https://encoding.spec.whatwg.org/#names-and-labels

	encoding, err := htmlindex.Get(encodingName)
	if err != nil {
		return nil, err
	}

	decoder := encoding.NewDecoder()

	allRunes := map[rune]bool{}

	decode := func(encoded []byte) (rune, bool) {
		decoded, err := decoder.Bytes(encoded)
		if err != nil {
			return utf8.RuneError, false
		}
		decodedRune, length := utf8.DecodeRune(decoded)
		if length > 0 && decodedRune != utf8.RuneError {
			return decodedRune, true
		}
		return utf8.RuneError, false
	}

	buffer := [2]byte{}

	// Try single and double char encodings only
	for first := 0; first <= 0xff; first++ {
		buffer[0] = byte(first)
		decoded, ok := decode(buffer[:1])
		if ok {
			allRunes[decoded] = true
		}

		for second := 0; second <= 0xff; second++ {
			buffer[1] = byte(second)
			decoded, ok := decode(buffer[:2])
			if ok {
				allRunes[decoded] = true
			}
		}
	}

	runeList := []rune{}

	for r := range allRunes {
		runeList = append(runeList, r)
	}

	return usableRunes(runeList), nil
}

func runesFromFile(file string) ([]rune, error) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	allRunes := map[rune]bool{}
	sample := string(bytes)
	for _, rune := range sample {
		allRunes[rune] = true
	}
	runeList := []rune{}

	for r := range allRunes {
		runeList = append(runeList, r)
	}

	return usableRunes(runeList), nil
}

func runesFromRange(lowChar, highChar int) []rune {
	cp := charmap.Windows1252
	runeList := []rune{}
	for char := lowChar; char <= highChar; char++ {
		decoded := cp.DecodeByte(byte(char))
		runeList = append(runeList, decoded)
	}
	return runeList
}

func usableRunes(runeSet []rune) []rune {
	usable := []rune{}

	for _, r := range runeSet {
		if unicode.IsGraphic(r) {
			usable = append(usable, r)
		}
	}

	sort.Slice(usable, func(i, j int) bool {
		return usable[i] < usable[j]
	})

	return usable
}
