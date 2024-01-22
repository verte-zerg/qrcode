package encode

import (
	"fmt"
)

type AlphaNumericEncoder struct{}

func (AlphaNumericEncoder) Encode(content string, queue chan ValueBlock) error {
	enc := &AlphaNumericConverter{}
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

func (AlphaNumericEncoder) CanEncode(content string) bool {
	enc := &AlphaNumericConverter{}
	_, err := enc.Convert(content)
	return err == nil
}

func (AlphaNumericEncoder) Size(length int) int {
	duplets := length / 2
	tail := length % 2
	extra := 0
	if tail != 0 {
		extra = 6
	}

	return duplets*11 + extra
}

func (AlphaNumericEncoder) Mode() EncodingMode {
	return EncodingModeAlphaNumeric
}
