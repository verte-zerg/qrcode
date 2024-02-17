package encode

import (
	"fmt"
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"
)

// byteEncoder is a struct that uses for converting string to byte data.
type byteEncoder struct{}

func (byteEncoder) Encode(content string, queue chan ValueBlock) error {
	enc := charmap.ISO8859_1.NewEncoder()
	buf, err := enc.Bytes([]byte(content))
	if err != nil {
		return fmt.Errorf("failed to encode string to byte: %w", err)
	}

	for _, b := range buf {
		queue <- ValueBlock{
			Bits:  8,
			Value: int(b),
		}
	}

	return nil
}

func (*byteEncoder) Size(content string) int {
	return utf8.RuneCountInString(content) * 8
}

func (*byteEncoder) CanEncode(content string) bool {
	return regexpByte.MatchString(content)
}

func (*byteEncoder) Mode() EncodingMode {
	return EncodingModeByte
}
