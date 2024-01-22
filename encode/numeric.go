package encode

import (
	"fmt"
	"strconv"
)

type NumericEncoder struct{}

func (NumericEncoder) Encode(content string) ([]byte, error) {
	triplets := len(content) / 3
	if len(content)%3 != 0 {
		triplets++
	}

	var data []byte
	freeBits := 0

	for i := 0; i < triplets; i++ {
		right := i*3 + 3
		if right > len(content) {
			right = len(content)
		}

		triplet := content[i*3 : right]
		tripletBytesSize := 1 + len(triplet)*3 // 4, 7, 10
		number, err := strconv.Atoi(triplet)
		if err != nil {
			return nil, fmt.Errorf("failed to convert string to int: %w", err)
		}

		var b byte

		if freeBits == 0 {
			data = append(data, 0)
			freeBits = 8
		}

		if tripletBytesSize > freeBits {
			b = byte(number >> (tripletBytesSize - freeBits))
		} else {
			b = byte(number << (freeBits - tripletBytesSize) & 0xff)
		}

		data[len(data)-1] |= b

		if tripletBytesSize > freeBits {
			data = append(data, byte(number<<(8-(tripletBytesSize-freeBits))&0xff))
		}

		freeBits = (freeBits + 6) % 8
	}

	return data, nil
}

func (NumericEncoder) Size(length int) int {
	triplets := length / 3
	tail := length % 3
	extra := 0
	if tail != 0 {
		extra = 1 + tail*3
	}

	return triplets*10 + extra
}

func (NumericEncoder) CanEncode(content string) bool {
	for _, r := range content {
		if r < '0' || r > '9' {
			return false
		}
	}

	return true
}

func (NumericEncoder) Mode() EncodingMode {
	return EncodingModeNumeric
}
