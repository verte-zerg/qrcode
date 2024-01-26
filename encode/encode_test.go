package encode

import (
	"bytes"
	"fmt"
	"testing"
)

func EncodeDataWrapper(data string, mode EncodingMode) ([]byte, error) {
	queue := make(chan ValueBlock, 100)
	result := make(chan []byte)

	go GenerateData(queue, result)

	err := EncodeData(&EncodeBlock{Data: data, Mode: mode}, queue)
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
			block := &EncodeBlock{
				Mode: test.mode,
			}
			if bits, err := block.GetLengthBits(test.version); err != nil {
				t.Error(err)
			} else if bits != test.bits {
				t.Errorf("Expected %v, got %v", test.bits, bits)
			}
		})
	}
}

func TestPrefixBytes(t *testing.T) {
	tests := []struct {
		mode       EncodingMode
		lengthBits int
		codewords  int
		prefix     []byte
	}{
		{EncodingModeNumeric, 10, 8, []byte{0b00010000, 0b00100000}},
		{EncodingModeNumeric, 10, 9, []byte{0b00010000, 0b00100100}},
		{EncodingModeNumeric, 12, 10, []byte{0b00010000, 0b00001010}},
		{EncodingModeNumeric, 14, 11, []byte{0b00010000, 0b00000010, 0b11000000}},
		{EncodingModeAlphaNumeric, 16, 20, []byte{0b00100000, 0b00000001, 0b01000000}},
		{EncodingModeLatin1, 8, 23, []byte{0b01000001, 0b01110000}},
	}

	for _, test := range tests {
		name := fmt.Sprintf("mode %v, lengthBits %v, codewords %v", test.mode, test.lengthBits, test.codewords)
		t.Run(name, func(t *testing.T) {
			queue := make(chan ValueBlock, 10)
			result := make(chan []byte)
			go GenerateData(queue, result)
			block := &EncodeBlock{
				Mode: test.mode,
			}
			block.GetBytesPrefix(test.lengthBits, test.codewords, queue)
			close(queue)
			data := <-result

			if !bytes.Equal(data, test.prefix) {
				t.Errorf("Expected %b, got %b", test.prefix, data)
			}
		})
	}
}
