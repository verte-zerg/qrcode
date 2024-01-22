package encode

import (
	"bytes"
	"fmt"
	"testing"
)

func TestGetKanjiData(t *testing.T) {
	tests := []struct {
		content string
		data    []byte
	}{
		{"点", []byte{0b01101100, 0b11111000}},
		{"茗", []byte{0b11010101, 0b01010000}},
		{"茗点", []byte{0b11010101, 0b01010011, 0b01100111, 0b11000000}},
	}

	for _, test := range tests {
		name := fmt.Sprintf("%v", test.content)
		t.Run(name, func(t *testing.T) {
			data, err := EncodeDataWrapper(test.content, EncodingModeKanji)
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
