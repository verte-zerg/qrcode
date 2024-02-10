package encode

import (
	"errors"
	"fmt"
	"regexp"
	"unicode/utf8"
)

type EncodingMode int

const (
	EncodingModeNumeric      EncodingMode = 1
	EncodingModeAlphaNumeric EncodingMode = 2
	EncodingModeLatin1       EncodingMode = 4
	EncodingModeKanji        EncodingMode = 8
	EncodingModeECI          EncodingMode = 7
)

// Count of length bits for each version and encoding mode.
// Structure: [encoding mode][version range number]
// Version range number is:
// * 0-3 for versions M1-M4,
// * 4 for versions 1-9
// * 5 for versions 10-26
// * 6 for versions 27-40.
var encodingModeLengthMap = map[EncodingMode][7]int{
	EncodingModeNumeric:      {3, 4, 5, 6, 10, 12, 14},
	EncodingModeAlphaNumeric: {0, 3, 4, 5, 9, 11, 13},
	EncodingModeLatin1:       {0, 0, 4, 5, 8, 16, 16},
	EncodingModeKanji:        {0, 0, 3, 4, 8, 10, 12},
}

var encodingModeEncoderMap = map[EncodingMode]QREncoder{
	EncodingModeNumeric:      &numericEncoder{},
	EncodingModeAlphaNumeric: &alphaNumericEncoder{},
	EncodingModeLatin1:       &latin1Encoder{},
	EncodingModeKanji:        &kanjiEncoder{},
}

var modeVersionValueBlockMap = map[EncodingMode][4]ValueBlock{
	EncodingModeNumeric: {
		{Value: 0, Bits: 0},
		{Value: 0, Bits: 1},
		{Value: 0, Bits: 2},
		{Value: 0, Bits: 3},
	},
	EncodingModeAlphaNumeric: {
		{Value: 0, Bits: 0},
		{Value: 1, Bits: 1},
		{Value: 1, Bits: 2},
		{Value: 1, Bits: 3},
	},
	EncodingModeLatin1: {
		{Value: 0, Bits: 0},
		{Value: 0, Bits: 0},
		{Value: 2, Bits: 2},
		{Value: 2, Bits: 3},
	},
	EncodingModeKanji: {
		{Value: 0, Bits: 0},
		{Value: 0, Bits: 0},
		{Value: 3, Bits: 2},
		{Value: 3, Bits: 3},
	},
}

type ValueBlock struct {
	Value int
	Bits  int
}

type ErrVersionInvalid struct {
	Version int
}

func (e ErrVersionInvalid) Error() string {
	return fmt.Sprintf("version %v invalid, must be between 1 and 40 inclusive", e.Version)
}

var ErrCannotDeterminEncodingMode = errors.New("cannot determine encoding mode")
var ErrUnknownEncodingMode = errors.New("unknown encoding mode")
var ErrVersionDoesNotSupportEncodingMode = errors.New("version does not support encoding mode")

var regexpNumeric *regexp.Regexp = regexp.MustCompile(`^[0-9]+$`)
var regexpAlphaNumeric *regexp.Regexp = regexp.MustCompile(`^[0-9A-Z $%*+\-./:]+$`)
var regexpLatin1 *regexp.Regexp = regexp.MustCompile(`^[\x00-\xFF]+$`)
var regexpKanji *regexp.Regexp = regexp.MustCompile(`^[\p{Hiragana}\p{Katakana}\p{Han}]+$`)

type QREncoder interface {
	Mode() EncodingMode
	Encode(s string, queue chan ValueBlock) error
	CanEncode(content string) bool
	Size(content string) int
}

type EncodeBlock struct {
	Mode EncodingMode
	Data string

	// Only for ECI mode
	SubMode          EncodingMode
	AssignmentNumber uint
}

func (b *EncodeBlock) GetSymbolsCount() int {
	if b.Mode == EncodingModeECI {
		enc := eciEncoder{
			AssignmentNumber: b.AssignmentNumber,
			DataMode:         b.SubMode,
		}
		return enc.Size(b.Data) / 8
	}

	return utf8.RuneCountInString(b.Data)
}

func (b *EncodeBlock) CalculateDataBitsCount() (int, error) {
	var enc QREncoder

	if b.Mode == EncodingModeECI {
		enc = eciEncoder{
			AssignmentNumber: b.AssignmentNumber,
			DataMode:         b.SubMode,
		}
	} else {
		var ok bool
		enc, ok = encodingModeEncoderMap[b.Mode]
		if !ok {
			return 0, ErrUnknownEncodingMode
		}
	}

	size := enc.Size(b.Data)

	if size == 0 {
		return 0, fmt.Errorf("failed to calculate data bits count: %w", ErrCannotDeterminEncodingMode)
	}

	return size, nil
}

func (b *EncodeBlock) GetLengthBits(version int) (int, error) {
	if version < -4 || version > 40 || version == 0 {
		return 0, ErrVersionInvalid{version}
	}

	mode := b.Mode
	if mode == EncodingModeECI {
		mode = b.SubMode
	}

	dataLength, ok := encodingModeLengthMap[mode]
	if !ok {
		return 0, ErrUnknownEncodingMode
	}

	if version <= 0 {
		bits := dataLength[-version-1]
		if bits == 0 {
			return 0, ErrVersionDoesNotSupportEncodingMode
		}
		return bits, nil
	}

	if version <= 9 {
		return dataLength[4], nil
	}
	if version <= 26 {
		return dataLength[5], nil
	}
	return dataLength[6], nil
}

func (b *EncodeBlock) GetModeBits(version int) int {
	if version < 0 {
		return -version - 1
	}

	if b.Mode == EncodingModeECI {
		return 16
	}
	return 4
}

func (b *EncodeBlock) GetBytesPrefix(
	version,
	lengthBits,
	itemsCount int,
	queue chan ValueBlock,
) {
	if version < 0 {
		queue <- modeVersionValueBlockMap[b.Mode][-version-1]
	} else {
		queue <- ValueBlock{
			Value: int(b.Mode),
			Bits:  4,
		}
	}

	if b.Mode == EncodingModeECI {
		queue <- ValueBlock{
			Value: int(b.AssignmentNumber),
			Bits:  8,
		}
		queue <- ValueBlock{
			Value: int(b.SubMode),
			Bits:  4,
		}
	}

	queue <- ValueBlock{
		Value: itemsCount,
		Bits:  lengthBits,
	}
}

// EncodeData transforms the content to the byte array according to the encoding mode.
func (b *EncodeBlock) EncodeData(queue chan ValueBlock) error {
	var enc QREncoder
	if b.Mode == EncodingModeECI {
		enc = eciEncoder{
			AssignmentNumber: b.AssignmentNumber,
			DataMode:         b.SubMode,
		}
	} else {
		var ok bool
		enc, ok = encodingModeEncoderMap[b.Mode]
		if !ok {
			return ErrUnknownEncodingMode
		}
	}

	err := enc.Encode(b.Data, queue)
	if err != nil {
		return fmt.Errorf("failed to encode data: %w", err)
	}

	return nil
}

func (b *EncodeBlock) Encode(version int, queue chan ValueBlock) (int, error) {
	symbolsCount := b.GetSymbolsCount()
	lengthBits, err := b.GetLengthBits(version)

	if err != nil {
		return 0, fmt.Errorf("failed to get length bits: %w", err)
	}
	dataBits, err := b.CalculateDataBitsCount()
	if err != nil {
		return 0, fmt.Errorf("failed to calculate data bits count: %w", err)
	}

	b.GetBytesPrefix(version, lengthBits, symbolsCount, queue)
	err = b.EncodeData(queue)
	if err != nil {
		return 0, fmt.Errorf("failed to encode data: %w", err)
	}

	allBits := dataBits + lengthBits + b.GetModeBits(version)

	return allBits, nil
}

// GetEncodingMode returns the encoding mode for the given string.
// EncodingMode is one of the EncodingMode constants (EncodingModeNumeric, EncodingModeAlphaNumeric, EncodingModeLatin1, EncodingModeKanji)
func GetEncodingMode(s string) (EncodingMode, error) {
	if regexpNumeric.MatchString(s) {
		return EncodingModeNumeric, nil
	}
	if regexpAlphaNumeric.MatchString(s) {
		return EncodingModeAlphaNumeric, nil
	}
	if regexpLatin1.MatchString(s) {
		return EncodingModeLatin1, nil
	}
	if regexpKanji.MatchString(s) {
		return EncodingModeKanji, nil
	}
	return 0, ErrCannotDeterminEncodingMode
}

// GetLengthBits returns the number of length bits for the given version and encoding mode.
// Version is an integer from 1 to 40 inclusive.
// Mode is one of the EncodingMode constants.

func GenerateData(queue chan ValueBlock, result chan []byte) {
	var data []byte
	freeBits := 0

	for v := range queue {
		if v.Bits == 0 {
			continue
		}

		var b byte

		blockSize := v.Bits
		value := v.Value

		if freeBits == 0 {
			data = append(data, 0)
			freeBits = 8
		}

		if blockSize > freeBits {
			b = byte(value >> (blockSize - freeBits))
		} else {
			b = byte(value << (freeBits - blockSize) & 0xff)
		}

		data[len(data)-1] |= b

		if blockSize > freeBits {
			if blockSize > freeBits+8 {
				b = byte(value >> ((blockSize - freeBits) - 8))
			} else {
				b = byte(value << (8 - (blockSize - freeBits)) & 0xff)
			}

			data = append(data, b)
		}

		if blockSize > freeBits+8 {
			data = append(data, byte(value<<(16-(blockSize-freeBits))&0xff))
		}

		freeBits = (freeBits + (16 - blockSize)) % 8
	}

	result <- data
}
