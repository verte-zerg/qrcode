package encode

import (
	"bytes"
	"fmt"
	"testing"
)

func TestGetAlphaNumericData(t *testing.T) {
	tests := []struct {
		content string
		data    []byte
	}{
		{"A", []byte{0b00101000}},
		{"AB", []byte{0b00111001, 0b10100000}},
		{"ABCDEFGHIJKLM", []byte{0b00111001, 0b10101000, 0b10100101, 0b01000010, 0b10101110, 0b00010110, 0b01111010, 0b11100110, 0b01010110}},
		{"AC-42", []byte{0b00111001, 0b11011100, 0b11100100, 0b00100000}},
	}

	for _, test := range tests {
		name := fmt.Sprintf("%v", test.content)
		t.Run(name, func(t *testing.T) {
			data, err := AlphaNumericEncoder{}.Encode(test.content)
			if err != nil {
				t.Error(err)
			} else if !bytes.Equal(data, test.data) {
				t.Errorf("Expected %b, got %b", test.data, data)
			}
		})
	}
}
