package encode

import (
	"fmt"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/encoding/unicode/utf32"
)

var AssigmentNumbersEncodings = map[uint]encoding.Encoding{
	0:  charmap.CodePage437,
	1:  charmap.ISO8859_1,
	2:  charmap.CodePage437,
	3:  charmap.ISO8859_1,
	4:  charmap.ISO8859_2,
	5:  charmap.ISO8859_3,
	6:  charmap.ISO8859_4,
	7:  charmap.ISO8859_5,
	8:  charmap.ISO8859_6,
	9:  charmap.ISO8859_7,
	10: charmap.ISO8859_8,
	11: charmap.ISO8859_9,
	12: charmap.ISO8859_10,
	13: charmap.Windows874, // ISO8859-11, but Windows-874 is a superset of ISO8859-11
	15: charmap.ISO8859_13,
	16: charmap.ISO8859_14,
	17: charmap.ISO8859_15,
	18: charmap.Windows1258,
	20: japanese.ShiftJIS,
	21: charmap.Windows1250,
	22: charmap.Windows1251,
	23: charmap.Windows1252,
	24: charmap.Windows1256,
	25: unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM),
	26: unicode.UTF8,
	27: charmap.Windows1252, // ASCII, but Windows-1252 is a superset of ASCII
	28: traditionalchinese.Big5,
	29: simplifiedchinese.GBK, // should be GB/T 2312, but GBK is a superset of GB/T 2312
	30: korean.EUCKR,          // KS X 1001, not fully supported
	31: simplifiedchinese.GBK,
	32: simplifiedchinese.GB18030,
	33: unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM),
	34: utf32.UTF32(utf32.BigEndian, utf32.IgnoreBOM),
	35: utf32.UTF32(utf32.LittleEndian, utf32.IgnoreBOM),
}

type EciEncoder struct {
	AssignmentNumber uint
	DataMode         EncodingMode
}

func (e EciEncoder) Encode(content string, queue chan ValueBlock) error {
	enc, ok := AssigmentNumbersEncodings[e.AssignmentNumber]
	if !ok {
		return fmt.Errorf("unknown assignment number: %d", e.AssignmentNumber)
	}
	encoder := enc.NewEncoder()

	buf, err := encoder.Bytes([]byte(content))
	if err != nil {
		return fmt.Errorf("failed to encode string: %w", err)
	}

	for _, b := range buf {
		queue <- ValueBlock{
			Bits:  8,
			Value: int(b),
		}
	}

	return nil
}

func (e EciEncoder) CanEncode(content string) bool {
	enc, ok := AssigmentNumbersEncodings[e.AssignmentNumber]
	if !ok {
		return false
	}
	encoder := enc.NewEncoder()

	_, err := encoder.Bytes([]byte(content))
	return err == nil
}

func (e EciEncoder) Size(content string) int {
	enc, ok := AssigmentNumbersEncodings[e.AssignmentNumber]
	if !ok {
		return 0
	}
	encoder := enc.NewEncoder()

	data, err := encoder.Bytes([]byte(content))
	if err != nil {
		return 0
	}

	return len(data) * 8
}

func (e EciEncoder) Mode() EncodingMode {
	return EncodingModeECI
}
