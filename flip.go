package main

// from http://www.upsidedowntext.com/unicode
var flipped = map[rune]rune{
	'a':  'ɐ',
	'b':  'q',
	'c':  'ɔ',
	'd':  'p',
	'e':  'ǝ',
	'f':  'ɟ',
	'g':  'ƃ',
	'h':  'ɥ',
	'i':  'ᴉ',
	'j':  'ɾ',
	'k':  'ʞ',
	'l':  'l',
	'm':  'ɯ',
	'n':  'u',
	'o':  'o',
	'p':  'd',
	'q':  'b',
	'r':  'ɹ',
	's':  's',
	't':  'ʇ',
	'u':  'n',
	'v':  'ʌ',
	'w':  'ʍ',
	'x':  'x',
	'y':  'ʎ',
	'z':  'z',
	'A':  '∀',
	'B':  'B',
	'C':  'Ɔ',
	'D':  'D',
	'E':  'Ǝ',
	'F':  'Ⅎ',
	'G':  'פ',
	'H':  'H',
	'I':  'I',
	'J':  'ſ',
	'K':  'K',
	'L':  '˥',
	'M':  'W',
	'N':  'N',
	'O':  'O',
	'P':  'Ԁ',
	'Q':  'Q',
	'R':  'R',
	'S':  'S',
	'T':  '┴',
	'U':  '∩',
	'V':  'Λ',
	'W':  'M',
	'X':  'X',
	'Y':  '⅄',
	'Z':  'Z',
	'0':  '0',
	'1':  'Ɩ',
	'2':  'ᄅ',
	'3':  'Ɛ',
	'4':  'ㄣ',
	'5':  'ϛ',
	'6':  '9',
	'7':  'ㄥ',
	'8':  '8',
	'9':  '6',
	',':  '\'',
	'.':  '˙',
	'?':  '¿',
	'!':  '¡',
	'"':  '"',
	'\'': ',',
	'`':  ',',
	'(':  ')',
	')':  '(',
	'[':  ']',
	']':  '[',
	'{':  '}',
	'}':  '{',
	'<':  '>',
	'>':  '<',
	'&':  '⅋',
	'_':  '‾',
}

func flip(input string) string {
	output := make([]rune, len(input))
	// rewrite the input with flipped runes
	for i, r := range input {
		f, ok := flipped[r]
		if !ok {
			// if we don't have a flipped rune,
			// use the original
			f = r
		}
		output[i] = f
	}
	// reverse the string
	for i, j := 0, len(output)-1; i < j; i, j = i+1, j-1 {
		output[i], output[j] = output[j], output[i]
	}
	return "(ﾉಥ益ಥ）ﾉ ┻━┻ " + string(output)
}
