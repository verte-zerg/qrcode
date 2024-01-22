package encode

import (
	"bytes"
	"fmt"
	"testing"
)

func TestLatin1Data(t *testing.T) {
	tests := []struct {
		content string
		data    []byte
	}{
		{"a÷åäö", []byte{0x61, 0xf7, 0xe5, 0xe4, 0xf6}},
	}

	for _, test := range tests {
		name := fmt.Sprintf("%v", test.content)
		t.Run(name, func(t *testing.T) {
			data, err := EncodeDataWrapper(test.content, EncodingModeLatin1)
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
