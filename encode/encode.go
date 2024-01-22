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
	EncodingModeNumeric:      NumericEncoder{},
	EncodingModeAlphaNumeric: AlphaNumericEncoder{},
	EncodingModeLatin1:       Latin1Encoder{},
	EncodingModeKanji:        KanjiEncoder{},
}

type QREncoder interface {
	Mode() EncodingMode
	Encode(s string) ([]byte, error)
	CanEncode(content string) bool
	Size(length int) int
}

type EncodeBlock struct {
	Mode EncodingMode
	Data string
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
func EncodeData(content string, mode EncodingMode) ([]byte, error) {
	enc, ok := EncodingModeEncoderMap[mode]
	if !ok {
		return nil, ErrUnknownEncodingMode
	}

	buf, err := enc.Encode(content)
	if err != nil {
		return nil, fmt.Errorf("failed to encode string: %w", err)
	}

	return buf, nil
}

// CalculateDataBitsCount returns the number of data bits for the given content and encoding mode.
// Content is the string to encode.
// Mode is one of the EncodingMode constants.
func CalculateDataBitsCount(content string, mode EncodingMode) (int, error) {
	enc, ok := EncodingModeEncoderMap[mode]
	if !ok {
		return 0, ErrUnknownEncodingMode
	}

	return enc.Size(utf8.RuneCountInString(content)), nil
}

// GetLengthBits returns the number of length bits for the given version and encoding mode.
// Version is an integer from 1 to 40 inclusive.
// Mode is one of the EncodingMode constants.
func GetLengthBits(version int, mode EncodingMode) (int, error) {
	if version < 1 || version > 40 {
		return 0, ErrVersionInvalid{version}
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
