package encode

import (
	"fmt"
	"unicode/utf8"

	"golang.org/x/text/encoding/japanese"
)

// kanjiEncoder is a struct that uses for converting string to kanji data.
type kanjiEncoder struct{}

func (*kanjiEncoder) Encode(content string, queue chan ValueBlock) error {
	enc := japanese.ShiftJIS.NewEncoder()
	buf, err := enc.Bytes([]byte(content))

	if err != nil {
		return fmt.Errorf("failed to encode string to kanji: %w", err)
	}

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
			return fmt.Errorf("invalid byte: %v", high)
		}

		value := uint(high)*0xC0 + uint(low)
		queue <- ValueBlock{
			Bits:  13,
			Value: int(value),
		}
	}

	return nil
}

func (*kanjiEncoder) CanEncode(content string) bool {
	return regexpKanji.MatchString(content)
}

func (*kanjiEncoder) Size(content string) int {
	return utf8.RuneCountInString(content) * 13
}

func (*kanjiEncoder) Mode() EncodingMode {
	return EncodingModeKanji
}
