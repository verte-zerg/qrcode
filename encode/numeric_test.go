package encode

import (
	"bytes"
	"fmt"
	"testing"
)

func TestGetNumericData(t *testing.T) {
	tests := []struct {
		content string
		data    []byte
	}{
		{"0", []byte{0x00}},
		{"1", []byte{0b00010000}},
		{"12", []byte{0b00011000}},
		{"123", []byte{0b00011110, 0b11000000}},
		{"012345", []byte{0b00000011, 0b00010101, 0b10010000}},
		{"01234567", []byte{0b00000011, 0b00010101, 0b10011000, 0b01100000}},
		{"0123456789012345", []byte{0b00000011, 0b00010101, 0b10011010, 0b10011011, 0b10000101, 0b00111010, 0b10010100}},
	}

	for _, test := range tests {
		name := fmt.Sprintf("%v", test.content)
		t.Run(name, func(t *testing.T) {
			data, err := EncodeDataWrapper(test.content, EncodingModeNumeric)
			if err != nil {
				t.Error(err)
				return
			}
			if !bytes.Equal(data, test.data) {
				t.Errorf("Expected %b, got %b", test.data, data)
			}
		})
	}
}
