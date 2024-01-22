package qrcode

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/verte-zerg/qrcode/encode"
)

func TestPrefixBytes(t *testing.T) {
	tests := []struct {
		mode       encode.EncodingMode
		lengthBits int
		codewords  int
		prefix     []byte
	}{
		{encode.EncodingModeNumeric, 10, 8, []byte{0b00010000, 0b00100000}},
		{encode.EncodingModeNumeric, 10, 9, []byte{0b00010000, 0b00100100}},
		{encode.EncodingModeNumeric, 12, 10, []byte{0b00010000, 0b00001010}},
		{encode.EncodingModeNumeric, 14, 11, []byte{0b00010000, 0b00000010, 0b11000000}},
		{encode.EncodingModeAlphaNumeric, 16, 20, []byte{0b00100000, 0b00000001, 0b01000000}},
		{encode.EncodingModeLatin1, 8, 23, []byte{0b01000001, 0b01110000}},
	}

	for _, test := range tests {
		name := fmt.Sprintf("mode %v, lengthBits %v, codewords %v", test.mode, test.lengthBits, test.codewords)
		t.Run(name, func(t *testing.T) {
			if prefix := GetBytesPrefix(test.mode, test.lengthBits, test.codewords); !bytes.Equal(prefix, test.prefix) {
				t.Errorf("Expected %b, got %b", test.prefix, prefix)
			}
		})
	}
}

// func TestGetEDCData(t *testing.T) {  #TODO
// 	tests := []struct {
// 		codewords int
// 		data      []byte
// 		edc       []byte
// 	}{
// 		{17, []byte{31, 175, 212, 178, 236, 46, 49, 178}, []byte{211, 142, 30, 87, 151, 42, 34, 28, 0}},
// 		{20, []byte{186, 202, 37, 191, 189, 83}, []byte{115, 57, 254, 13, 198, 48, 178, 135, 120, 155, 54, 76, 144, 119}},
// 		{20, []byte{89, 176, 126, 141, 157, 255, 163, 212, 181, 5}, []byte{104, 111, 144, 5, 25, 157, 147, 163, 208, 73}},
// 		{18, []byte{129, 161, 199, 110, 185, 193, 230}, []byte{121, 36, 245, 220, 24, 105, 81, 239, 57, 68, 209}},
// 		{18, []byte{0, 227, 216, 165, 171, 89, 48, 55, 151}, []byte{170, 92, 76, 226, 243, 240, 151, 248, 200}},
// 		{20, []byte{250, 79, 10, 159, 91, 93, 6, 179, 69, 18}, []byte{196, 201, 159, 169, 240, 156, 169, 50, 211, 219}},
// 		{10, []byte{6, 181, 58, 199}, []byte{224, 115, 47, 239, 113, 108}},
// 		{18, []byte{68, 46, 213, 155}, []byte{176, 131, 232, 148, 243, 241, 84, 85, 176, 191, 217, 92, 175, 77}},
// 		{18, []byte{26, 136, 197, 136}, []byte{91, 165, 92, 61, 92, 26, 152, 26, 43, 154, 18, 59, 46, 50}},
// 		{20, []byte{200, 254, 35, 175, 175, 105, 226, 169}, []byte{100, 104, 126, 129, 171, 67, 22, 14, 98, 79, 69, 92}},
// 	}

// 	for _, test := range tests {
// 		name := fmt.Sprintf("%v", test.codewords)
// 		t.Run(name, func(t *testing.T) {
// 			edc := GetEDCData(test.data, test.codewords, ErrorCorrectionLevelMedium)
// 			if !bytes.Equal(edc, test.edc) {
// 				t.Errorf("Expected %b, got %b", test.edc, edc)
// 			}
// 		})
// 	}
// }
