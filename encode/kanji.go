package encode

import (
	"fmt"

	"golang.org/x/text/encoding/japanese"
)

type KanjiEncoder struct{}

func (KanjiEncoder) Encode(content string) ([]byte, error) {
	enc := japanese.ShiftJIS.NewEncoder()
	buf, err := enc.Bytes([]byte(content))
	if err != nil {
		return nil, fmt.Errorf("failed to encode string to kanji: %w", err)
	}

	const itemSize = 13

	var data []byte
	freeBits := 0

	for i := 0; i < len(buf); i += 2 {
		high, low := buf[i], buf[i+1]

		// subtract 0x8140 from each byte between 0x8140 and 0x9FFC
		// subtract 0xC140 from each byte between 0xE040 and 0xEBBF
		if high >= 0x81 && high <= 0x9F {
			high -= 0x81
			low -= 0x40
		} else if high >= 0xE0 && high <= 0xEB {
			high -= 0xC1
			low -= 0x40
		} else {
			return nil, fmt.Errorf("invalid byte: %v", high)
		}

		value := uint(high)*0xC0 + uint(low)

		if freeBits == 0 {
			data = append(data, 0)
			freeBits = 8
		}

		data[len(data)-1] |= byte(value >> (itemSize - freeBits))

		var b byte

		if itemSize > freeBits+8 {
			b = byte(value >> ((itemSize - freeBits) - 8))
		} else {
			b = byte(value << (8 - (itemSize - freeBits)) & 0xff)
		}

		data = append(data, b)

		if itemSize > freeBits+8 {
			data = append(data, byte(value<<(16-(itemSize-freeBits))&0xff))
		}

		freeBits = (freeBits + 3) % 8
	}

	return data, nil
}

func (KanjiEncoder) CanEncode(content string) bool {
	return regexpKanji.MatchString(content)
}

func (KanjiEncoder) Size(length int) int {
	return length * 13
}

func (KanjiEncoder) Mode() EncodingMode {
	return EncodingModeKanji
}
