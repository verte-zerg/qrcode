package encode

import (
	"fmt"
	"strconv"
	"unicode/utf8"
)

// numericEncoder is a struct that uses for converting string to numeric data.
type numericEncoder struct{}

func (*numericEncoder) Encode(content string, queue chan ValueBlock) error {
	triplets := len(content) / 3
	if len(content)%3 != 0 {
		triplets++
	}

	for i := 0; i < triplets; i++ {
		right := i*3 + 3
		if right > len(content) {
			right = len(content)
		}

		triplet := content[i*3 : right]
		tripletBytesSize := 1 + len(triplet)*3 // 4, 7, 10
		number, err := strconv.Atoi(triplet)
		if err != nil {
			return fmt.Errorf("failed to convert string to int: %w", err)
		}

		queue <- ValueBlock{
			Bits:  tripletBytesSize,
			Value: number,
		}
	}

	return nil
}

func (*numericEncoder) Size(content string) int {
	length := utf8.RuneCountInString(content)
	triplets := length / 3
	tail := length % 3
	extra := 0
	if tail != 0 {
		extra = 1 + tail*3
	}

	return triplets*10 + extra
}

func (*numericEncoder) CanEncode(content string) bool {
	for _, r := range content {
		if r < '0' || r > '9' {
			return false
		}
	}

	return true
}

func (*numericEncoder) Mode() EncodingMode {
	return EncodingModeNumeric
}
