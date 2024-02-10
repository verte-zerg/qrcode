package encode

import (
	"fmt"
	"unicode/utf8"
)

// alphaNumericConverter is a struct that uses for converting string to alphanumeric data.
type alphaNumericEncoder struct{}

func (*alphaNumericEncoder) Encode(content string, queue chan ValueBlock) error {
	enc := &alphaNumericConverter{}
	encoded, err := enc.Convert(content)
	if err != nil {
		return fmt.Errorf("failed to convert string to alphanumeric: %w", err)
	}

	duplets := len(encoded) / 2
	if len(encoded)%2 != 0 {
		duplets++
	}

	for i := 0; i < duplets; i++ {
		right := i*2 + 2
		dupletBytesSize := 6
		number := uint(encoded[i*2])
		if right <= len(encoded) {
			dupletBytesSize = 11
			number = number*45 + uint(encoded[i*2+1])
		}

		queue <- ValueBlock{
			Bits:  dupletBytesSize,
			Value: int(number),
		}
	}

	return nil
}

func (*alphaNumericEncoder) CanEncode(content string) bool {
	enc := &alphaNumericConverter{}
	_, err := enc.Convert(content)
	return err == nil
}

func (*alphaNumericEncoder) Size(content string) int {
	length := utf8.RuneCountInString(content)
	duplets := length / 2
	tail := length % 2
	extra := 0
	if tail != 0 {
		extra = 6
	}

	return duplets*11 + extra
}

func (*alphaNumericEncoder) Mode() EncodingMode {
	return EncodingModeAlphaNumeric
}
