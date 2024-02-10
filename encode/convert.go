package encode

import "errors"

var ErrCannotEncodeToAlphaNumeric = errors.New("cannot encode to alphanumeric")

var utfToAlphaNumeric = map[rune]rune{
	'0': 0,
	'1': 1,
	'2': 2,
	'3': 3,
	'4': 4,
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
	'9': 9,
	'A': 10,
	'B': 11,
	'C': 12,
	'D': 13,
	'E': 14,
	'F': 15,
	'G': 16,
	'H': 17,
	'I': 18,
	'J': 19,
	'K': 20,
	'L': 21,
	'M': 22,
	'N': 23,
	'O': 24,
	'P': 25,
	'Q': 26,
	'R': 27,
	'S': 28,
	'T': 29,
	'U': 30,
	'V': 31,
	'W': 32,
	'X': 33,
	'Y': 34,
	'Z': 35,
	' ': 36,
	'$': 37,
	'%': 38,
	'*': 39,
	'+': 40,
	'-': 41,
	'.': 42,
	'/': 43,
	':': 44,
}

type alphaNumericConverter struct{}

func (e *alphaNumericConverter) Convert(s string) ([]byte, error) {
	var result []byte

	for _, r := range s {
		if v, ok := utfToAlphaNumeric[r]; ok {
			result = append(result, byte(v))
		} else {
			return nil, ErrCannotEncodeToAlphaNumeric
		}
	}

	return result, nil
}
