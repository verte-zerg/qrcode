package encode

import (
	"fmt"
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"
)

type Latin1Encoder struct{}

func (Latin1Encoder) Encode(content string, queue chan ValueBlock) error {
	enc := charmap.ISO8859_1.NewEncoder()
	buf, err := enc.Bytes([]byte(content))
	if err != nil {
		return fmt.Errorf("failed to encode string to latin1: %w", err)
	}

	for _, b := range buf {
		queue <- ValueBlock{
			Bits:  8,
			Value: int(b),
		}
	}

	return nil
}

func (*Latin1Encoder) Size(content string) int {
	return utf8.RuneCountInString(content) * 8
}

func (*Latin1Encoder) CanEncode(content string) bool {
	return regexpLatin1.MatchString(content)
}

func (*Latin1Encoder) Mode() EncodingMode {
	return EncodingModeLatin1
}
