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

const (
	CP437             = 0 // also have a value of 2
	ISO8859_1         = 1 // also have a value of 3
	ISO8859_2         = 4
	ISO8859_3         = 5
	ISO8859_4         = 6
	ISO8859_5         = 7
	ISO8859_6         = 8
	ISO8859_7         = 9
	ISO8859_8         = 10
	ISO8859_9         = 11
	ISO8859_10        = 12
	ISO8859_11        = 13 // ISO8859-11, but Windows-874 is a superset
	REVERSED_14       = 14 // does not support
	ISO8859_13        = 15
	ISO8859_14        = 16
	ISO8859_15        = 17
	ISO8859_16        = 18
	REVERSED_19       = 19 // does not support
	ShiftJIS          = 20
	Windows1250       = 21
	Windows1251       = 22
	Windows1252       = 23
	Windows1256       = 24
	UTF16BigEndian    = 25
	UTF8              = 26
	ASCII             = 27
	Big5              = 28
	GBT2312           = 29
	KSX1001           = 30
	GBK               = 31
	GB18030           = 32
	UTF16LittleEndian = 33
	UTF32BigEndian    = 34
	UTF32LittleEndian = 35
)

// assigmentNumbersEncodings is a map of ECI assignment numbers to their respective encodings.
var assigmentNumbersEncodings = map[uint]encoding.Encoding{
	CP437:             charmap.CodePage437,
	ISO8859_1:         charmap.ISO8859_1,
	2:                 charmap.CodePage437,
	3:                 charmap.ISO8859_1,
	ISO8859_2:         charmap.ISO8859_2,
	ISO8859_3:         charmap.ISO8859_3,
	ISO8859_4:         charmap.ISO8859_4,
	ISO8859_5:         charmap.ISO8859_5,
	ISO8859_6:         charmap.ISO8859_6,
	ISO8859_7:         charmap.ISO8859_7,
	ISO8859_8:         charmap.ISO8859_8,
	ISO8859_9:         charmap.ISO8859_9,
	ISO8859_10:        charmap.ISO8859_10,
	ISO8859_11:        charmap.Windows874, // ISO8859-11, but Windows-874 is a superset
	REVERSED_14:       nil,                // does not support
	ISO8859_13:        charmap.ISO8859_13,
	ISO8859_14:        charmap.ISO8859_14,
	ISO8859_15:        charmap.ISO8859_15,
	ISO8859_16:        charmap.Windows1258,
	REVERSED_19:       nil, // does not support
	ShiftJIS:          japanese.ShiftJIS,
	Windows1250:       charmap.Windows1250,
	Windows1251:       charmap.Windows1251,
	Windows1252:       charmap.Windows1252,
	Windows1256:       charmap.Windows1256,
	UTF16BigEndian:    unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM),
	UTF8:              unicode.UTF8,
	ASCII:             charmap.Windows1252, // ASCII, but Windows-1252 is a superset of ASCII
	Big5:              traditionalchinese.Big5,
	GBT2312:           simplifiedchinese.GBK, // should be GB/T 2312, but GBK is a superset of GB/T 2312
	KSX1001:           korean.EUCKR,          // KS X 1001, not fully supported
	GBK:               simplifiedchinese.GBK,
	GB18030:           simplifiedchinese.GB18030,
	UTF16LittleEndian: unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM),
	UTF32BigEndian:    utf32.UTF32(utf32.BigEndian, utf32.IgnoreBOM),
	UTF32LittleEndian: utf32.UTF32(utf32.LittleEndian, utf32.IgnoreBOM),
}

var ErrUnknownAssignmentNumber = fmt.Errorf("unknown assignment number")

// eciEncoder is an encoder for ECI (Extended Channel Interpretation) mode.
type eciEncoder struct {
	AssignmentNumber uint
	DataMode         EncodingMode
}

func (e eciEncoder) Encode(content string, queue chan ValueBlock) error {
	enc, ok := assigmentNumbersEncodings[e.AssignmentNumber]
	if !ok {
		return fmt.Errorf("unknown assignment number: %d", e.AssignmentNumber)
	}
	if enc == nil {
		return fmt.Errorf("the assignment number %d has not supported", e.AssignmentNumber)
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

func (e eciEncoder) CanEncode(content string) bool {
	enc, ok := assigmentNumbersEncodings[e.AssignmentNumber]
	if !ok || enc == nil {
		return false
	}
	encoder := enc.NewEncoder()

	_, err := encoder.Bytes([]byte(content))
	return err == nil
}

func (e eciEncoder) Size(content string) int {
	enc, ok := assigmentNumbersEncodings[e.AssignmentNumber]
	if !ok || enc == nil {
		return 0
	}
	encoder := enc.NewEncoder()

	data, err := encoder.Bytes([]byte(content))
	if err != nil {
		return 0
	}

	return len(data) * 8
}

func (e eciEncoder) Mode() EncodingMode {
	return EncodingModeECI
}
