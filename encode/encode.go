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
// Version range number is 0 for versions 1-9, 1 for versions 10-26, 2 for versions 27-40.
var EncodingModeLengthMap = map[EncodingMode][3]int{
	EncodingModeNumeric:      {10, 12, 14},
	EncodingModeAlphaNumeric: {9, 11, 13},
	EncodingModeLatin1:       {8, 16, 16},
	EncodingModeKanji:        {8, 10, 12},
}

var EncodingModeEncoderMap = map[EncodingMode]QREncoder{
	EncodingModeNumeric:      &NumericEncoder{},
	EncodingModeAlphaNumeric: &AlphaNumericEncoder{},
	EncodingModeLatin1:       &Latin1Encoder{},
	EncodingModeKanji:        &KanjiEncoder{},
}

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
		enc := EciEncoder{
			AssignmentNumber: b.AssignmentNumber,
			DataMode:         b.SubMode,
		}
		return enc.Size(b.Data) / 8
	}

	return utf8.RuneCountInString(b.Data)
}

func (b *EncodeBlock) CalculateDataBitsCount() (int, error) {
	if b.Mode == EncodingModeECI {
		enc := EciEncoder{
			AssignmentNumber: b.AssignmentNumber,
			DataMode:         b.SubMode,
		}

		return enc.Size(b.Data), nil
	}

	enc, ok := EncodingModeEncoderMap[b.Mode]
	if !ok {
		return 0, ErrUnknownEncodingMode
	}

	return enc.Size(b.Data), nil
}

func (b *EncodeBlock) GetLengthBits(version int) (int, error) {
	if version < 1 || version > 40 {
		return 0, ErrVersionInvalid{version}
	}

	mode := b.Mode
	if mode == EncodingModeECI {
		mode = b.SubMode
	}

	dataLength, ok := EncodingModeLengthMap[mode]
	if !ok {
		return 0, ErrUnknownEncodingMode
	}

	if version <= 9 {
		return dataLength[0], nil
	}
	if version <= 26 {
		return dataLength[1], nil
	}
	return dataLength[2], nil
}

func (b *EncodeBlock) GetModeBits() int {
	if b.Mode == EncodingModeECI {
		return 16
	}
	return 4
}

func (b *EncodeBlock) GetBytesPrefix(
	lengthBits,
	itemsCount int,
	queue chan ValueBlock,
) {
	queue <- ValueBlock{
		Value: int(b.Mode),
		Bits:  4,
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

func (b *EncodeBlock) Encode(version int, queue chan ValueBlock) (int, error) {
	symbolsCount := b.GetSymbolsCount()
	lengthBits, err := b.GetLengthBits(version)

	if err != nil {
		return 0, fmt.Errorf("failed to get length bits: %w", err)
	}
	b.GetBytesPrefix(lengthBits, symbolsCount, queue)
	dataBits, err := b.CalculateDataBitsCount()
	if err != nil {
		return 0, fmt.Errorf("failed to calculate data bits count: %w", err)
	}

	err = EncodeData(b, queue)
	if err != nil {
		return 0, fmt.Errorf("failed to encode data: %w", err)
	}

	allBits := dataBits + lengthBits + b.GetModeBits()

	return allBits, nil
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

var regexpNumeric *regexp.Regexp = regexp.MustCompile(`^[0-9]+$`)
var regexpAlphaNumeric *regexp.Regexp = regexp.MustCompile(`^[0-9A-Z $%*+\-./:]+$`)
var regexpLatin1 *regexp.Regexp = regexp.MustCompile(`^[\x00-\xFF]+$`)
var regexpKanji *regexp.Regexp = regexp.MustCompile(`^[\p{Hiragana}\p{Katakana}\p{Han}]+$`)

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

// EncodeData transforms the content to the byte array according to the encoding mode.
func EncodeData(block *EncodeBlock, queue chan ValueBlock) error {
	if block.Mode == EncodingModeECI {
		enc := EciEncoder{
			AssignmentNumber: block.AssignmentNumber,
			DataMode:         block.SubMode,
		}

		err := enc.Encode(block.Data, queue)
		if err != nil {
			return fmt.Errorf("failed to encode data: %w", err)
		}

		return nil
	}

	enc, ok := EncodingModeEncoderMap[block.Mode]
	if !ok {
		return ErrUnknownEncodingMode
	}

	err := enc.Encode(block.Data, queue)
	if err != nil {
		return fmt.Errorf("failed to encode string: %w", err)
	}

	return nil
}

// GetLengthBits returns the number of length bits for the given version and encoding mode.
// Version is an integer from 1 to 40 inclusive.
// Mode is one of the EncodingMode constants.

func GenerateData(queue chan ValueBlock, result chan []byte) {
	var data []byte
	freeBits := 0

	for v := range queue {
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
