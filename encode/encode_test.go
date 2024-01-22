package encode

import (
	"fmt"
	"testing"
)

func EncodeDataWrapper(data string, mode EncodingMode) ([]byte, error) {
	queue := make(chan ValueBlock, 100)
	result := make(chan []byte)

	go GenerateData(queue, result)

	err := EncodeData(data, mode, queue)
	if err != nil {
		return nil, fmt.Errorf("failed to encode data: %w", err)
	}

	close(queue)
	return <-result, nil
}

func TestEncodingMode(t *testing.T) {
	tests := []struct {
		s    string
		mode EncodingMode
	}{
		{"1234567890", EncodingModeNumeric},
		{"1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ $%*+-./:", EncodingModeAlphaNumeric},
		{"1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ $%*+-./:\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0a\x0b\x0c\x0d\x0e\x0f", EncodingModeLatin1},
		{"あア亜", EncodingModeKanji},
	}

	for _, test := range tests {
		name := fmt.Sprintf("%v", test.s)
		t.Run(name, func(t *testing.T) {
			if mode, err := GetEncodingMode(test.s); err != nil {
				t.Error(err)
			} else if mode != test.mode {
				t.Errorf("Expected %v, got %v", test.mode, mode)
			}
		})
	}
}

func TestUnknownEncodingMode(t *testing.T) {
	if _, err := GetEncodingMode("abcABC123!@#あア亜"); err == nil {
		t.Error("Expected error")
	}
}

func TestLengthBits(t *testing.T) {
	tests := []struct {
		version int
		mode    EncodingMode
		bits    int
	}{
		{1, EncodingModeNumeric, 10},
		{1, EncodingModeAlphaNumeric, 9},
		{1, EncodingModeLatin1, 8},
		{1, EncodingModeKanji, 8},
		{10, EncodingModeNumeric, 12},
		{10, EncodingModeAlphaNumeric, 11},
		{10, EncodingModeLatin1, 16},
		{10, EncodingModeKanji, 10},
		{27, EncodingModeNumeric, 14},
		{27, EncodingModeAlphaNumeric, 13},
		{27, EncodingModeLatin1, 16},
		{27, EncodingModeKanji, 12},
	}

	for _, test := range tests {
		name := fmt.Sprintf("version %v, mode %v", test.version, test.mode)
		t.Run(name, func(t *testing.T) {
			if bits, err := GetLengthBits(test.version, test.mode); err != nil {
				t.Error(err)
			} else if bits != test.bits {
				t.Errorf("Expected %v, got %v", test.bits, bits)
			}
		})
	}
}
