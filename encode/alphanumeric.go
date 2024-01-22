package encode

import (
	"fmt"
)

type AlphaNumericEncoder struct{}

func (AlphaNumericEncoder) Encode(content string) ([]byte, error) {
	enc := &AlphaNumericConverter{}
	encoded, err := enc.Convert(content)
	if err != nil {
		return nil, fmt.Errorf("failed to convert string to alphanumeric: %w", err)
	}

	duplets := len(encoded) / 2
	if len(encoded)%2 != 0 {
		duplets++
	}

	var data []byte
	freeBits := 0

	for i := 0; i < duplets; i++ {
		right := i*2 + 2
		dupletBytesSize := 6
		number := uint(encoded[i*2])
		if right <= len(encoded) {
			dupletBytesSize = 11
			number = number*45 + uint(encoded[i*2+1])
		}

		var b byte

		if freeBits == 0 {
			data = append(data, 0)
			freeBits = 8
		}

		if dupletBytesSize > freeBits {
			b = byte(number >> (dupletBytesSize - freeBits))
		} else {
			b = byte(number << (freeBits - dupletBytesSize) & 0xff)
		}

		data[len(data)-1] |= b

		if dupletBytesSize > freeBits {
			if dupletBytesSize > freeBits+8 {
				b = byte(number >> ((dupletBytesSize - freeBits) - 8))
			} else {
				b = byte(number << (8 - (dupletBytesSize - freeBits)) & 0xff)
			}

			data = append(data, b)
		}

		if dupletBytesSize > freeBits+8 {
			data = append(data, byte(number<<(16-(dupletBytesSize-freeBits))&0xff))
		}

		freeBits = (freeBits + 5) % 8
	}

	return data, nil
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
