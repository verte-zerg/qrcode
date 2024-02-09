package encode

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"
)

func EncodeDataWrapper(block *EncodeBlock) ([]byte, error) {
	queue := make(chan ValueBlock, 100)
	result := make(chan []byte)

	go GenerateData(queue, result)

	err := block.EncodeData(queue)
	if err != nil {
		return nil, err
	}

	close(queue)
	return <-result, nil
}

func EncodeWrapper(block *EncodeBlock, version int) ([]byte, error) {
	queue := make(chan ValueBlock, 100)
	result := make(chan []byte)

	go GenerateData(queue, result)

	_, err := block.Encode(version, queue)
	if err != nil {
		return nil, err
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
		subMode EncodingMode
		bits    int
	}{
		{1, EncodingModeNumeric, 0, 10},
		{1, EncodingModeAlphaNumeric, 0, 9},
		{1, EncodingModeLatin1, 0, 8},
		{1, EncodingModeKanji, 0, 8},
		{1, EncodingModeECI, EncodingModeNumeric, 10},
		{10, EncodingModeNumeric, 0, 12},
		{10, EncodingModeAlphaNumeric, 0, 11},
		{10, EncodingModeLatin1, 0, 16},
		{10, EncodingModeECI, EncodingModeNumeric, 12},
		{10, EncodingModeKanji, 0, 10},
		{27, EncodingModeNumeric, 0, 14},
		{27, EncodingModeAlphaNumeric, 0, 13},
		{27, EncodingModeLatin1, 0, 16},
		{27, EncodingModeKanji, 0, 12},
		{27, EncodingModeECI, EncodingModeNumeric, 14},

		// Micro QR
		{-1, EncodingModeNumeric, 0, 3},
		{-2, EncodingModeNumeric, 0, 4},
		{-2, EncodingModeAlphaNumeric, 0, 3},
		{-3, EncodingModeNumeric, 0, 5},
		{-3, EncodingModeAlphaNumeric, 0, 4},
		{-3, EncodingModeLatin1, 0, 4},
		{-3, EncodingModeKanji, 0, 3},
		{-4, EncodingModeNumeric, 0, 6},
		{-4, EncodingModeAlphaNumeric, 0, 5},
		{-4, EncodingModeLatin1, 0, 5},
		{-4, EncodingModeKanji, 0, 4},
	}

	for _, test := range tests {
		name := fmt.Sprintf("version %v, mode %v", test.version, test.mode)
		t.Run(name, func(t *testing.T) {
			block := &EncodeBlock{
				Mode:    test.mode,
				SubMode: test.subMode,
			}
			if bits, err := block.GetLengthBits(test.version); err != nil {
				t.Error(err)
			} else if bits != test.bits {
				t.Errorf("Expected %v, got %v", test.bits, bits)
			}
		})
	}

	// Invalid version
	t.Run("invalid version", func(t *testing.T) {
		invalid_version := 0
		expected := ErrVersionInvalid{invalid_version}
		block := &EncodeBlock{
			Mode: EncodingModeNumeric,
		}
		_, err := block.GetLengthBits(invalid_version)
		if err != expected {
			t.Errorf("Expected %v, got %v", expected, err)
		}
	})

	// Invalid mode
	t.Run("invalid mode", func(t *testing.T) {
		block := &EncodeBlock{Mode: 0}
		_, err := block.GetLengthBits(1)
		if err != ErrUnknownEncodingMode {
			t.Errorf("Expected %v, got %v", ErrUnknownEncodingMode, err)
		}
	})

	// Invalid version for mode
	t.Run("invalid version for mode", func(t *testing.T) {
		block := &EncodeBlock{Mode: EncodingModeKanji}
		_, err := block.GetLengthBits(-1)
		if err != ErrVersionDoesNotSupportEncodingMode {
			t.Errorf("Expected %v, got %v", ErrVersionDoesNotSupportEncodingMode, err)
		}
	})
}

func TestPrefixBytes(t *testing.T) {
	tests := []struct {
		mode             EncodingMode
		subMode          EncodingMode
		assignmentNumber uint
		version          int
		lengthBits       int
		codewords        int
		prefix           []byte
	}{
		{EncodingModeNumeric, 0, 0, 1, 10, 8, []byte{0b00010000, 0b00100000}},
		{EncodingModeNumeric, 0, 0, 1, 10, 9, []byte{0b00010000, 0b00100100}},
		{EncodingModeNumeric, 0, 0, 1, 12, 10, []byte{0b00010000, 0b00001010}},
		{EncodingModeNumeric, 0, 0, 1, 14, 11, []byte{0b00010000, 0b00000010, 0b11000000}},
		{EncodingModeAlphaNumeric, 0, 0, 1, 16, 20, []byte{0b00100000, 0b00000001, 0b01000000}},
		{EncodingModeKanji, 0, 0, 1, 8, 10, []byte{0b10000000, 0b10100000}},
		{EncodingModeLatin1, 0, 0, 1, 8, 23, []byte{0b01000001, 0b01110000}},
		{EncodingModeECI, EncodingModeLatin1, 26, 1, 10, 8, []byte{0b01110001, 0b10100100, 0b00000010, 0b00000000}},

		// Micro QR
		{EncodingModeNumeric, 0, 0, -1, 3, 2, []byte{0b01000000}},
		{EncodingModeNumeric, 0, 0, -2, 4, 3, []byte{0b00011000}},
		{EncodingModeNumeric, 0, 0, -3, 5, 4, []byte{0b00001000}},
		{EncodingModeNumeric, 0, 0, -4, 6, 5, []byte{0b00000010, 0b10000000}},
		{EncodingModeAlphaNumeric, 0, 0, -2, 3, 2, []byte{0b10100000}},
		{EncodingModeAlphaNumeric, 0, 0, -3, 4, 3, []byte{0b01001100}},
		{EncodingModeAlphaNumeric, 0, 0, -4, 5, 4, []byte{0b00100100}},
		{EncodingModeLatin1, 0, 0, -3, 4, 3, []byte{0b10001100}},
		{EncodingModeLatin1, 0, 0, -4, 5, 4, []byte{0b01000100}},
		{EncodingModeKanji, 0, 0, -3, 3, 2, []byte{0b11010000}},
		{EncodingModeKanji, 0, 0, -4, 4, 3, []byte{0b01100110}},
	}

	for _, test := range tests {
		name := fmt.Sprintf("mode %v, lengthBits %v, codewords %v", test.mode, test.lengthBits, test.codewords)
		t.Run(name, func(t *testing.T) {
			queue := make(chan ValueBlock, 10)
			result := make(chan []byte)
			go GenerateData(queue, result)
			block := &EncodeBlock{
				Mode:             test.mode,
				SubMode:          test.subMode,
				AssignmentNumber: test.assignmentNumber,
			}
			block.GetBytesPrefix(test.version, test.lengthBits, test.codewords, queue)
			close(queue)
			data := <-result

			if !bytes.Equal(data, test.prefix) {
				t.Errorf("Expected %b, got %b", test.prefix, data)
			}
		})
	}
}

func TestGetSymbolsCount(t *testing.T) {
	tests := []struct {
		data             string
		mode             EncodingMode
		assignmentNumber uint
		subMode          EncodingMode
		symbol           int
	}{
		{"1234567890", EncodingModeNumeric, 0, EncodingModeNumeric, 10},
		{"1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ $%*+-./:", EncodingModeAlphaNumeric, 0, 0, 45},
		{"1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ $%*+-./:abcdefghijklmnopqrstuvwxyz", EncodingModeLatin1, 0, 0, 71},
		{"あア亜", EncodingModeKanji, 0, 0, 3},
		{"123ASDÄÖÏあア亜", EncodingModeECI, 26, EncodingModeLatin1, 21},
	}

	for _, test := range tests {
		name := fmt.Sprintf("data %v, mode %v, assignmentNumber %v, subMode %v", test.data, test.mode, test.assignmentNumber, test.subMode)
		t.Run(name, func(t *testing.T) {
			block := &EncodeBlock{
				Data:             test.data,
				Mode:             test.mode,
				AssignmentNumber: test.assignmentNumber,
				SubMode:          test.subMode,
			}

			if symbols := block.GetSymbolsCount(); symbols != test.symbol {
				t.Errorf("Expected %v, got %v", test.symbol, symbols)
			}
		})
	}
}

func TestCalculateDataBitsCount(t *testing.T) {
	tests := []struct {
		data             string
		mode             EncodingMode
		assignmentNumber uint
		subMode          EncodingMode
		bits             int
	}{
		{"1234567890", EncodingModeNumeric, 0, 0, 10*3 + 4},
		{"123ASD:", EncodingModeAlphaNumeric, 0, 0, 11*3 + 6},
		{"123ASDasd", EncodingModeLatin1, 0, 0, 9 * 8},
		{"あア亜", EncodingModeKanji, 0, 0, 13 * 3},
		{"123ASDÄÖÏあア亜", EncodingModeECI, 26, EncodingModeLatin1, 21 * 8},
	}

	for _, test := range tests {
		name := fmt.Sprintf("data %v, mode %v, assignmentNumber %v, subMode %v", test.data, test.mode, test.assignmentNumber, test.subMode)
		t.Run(name, func(t *testing.T) {
			block := &EncodeBlock{
				Data:             test.data,
				Mode:             test.mode,
				AssignmentNumber: test.assignmentNumber,
				SubMode:          test.subMode,
			}

			bits, err := block.CalculateDataBitsCount()
			if err != nil {
				t.Error(err)
			} else if bits != test.bits {
				t.Errorf("Expected %v, got %v", test.bits, bits)
			}
		})
	}

	// Invalid mode
	t.Run("invalid mode", func(t *testing.T) {
		block := &EncodeBlock{
			Data: "1234567890",
			Mode: 0,
		}

		_, err := block.CalculateDataBitsCount()

		if err != ErrUnknownEncodingMode {
			t.Errorf("Expected %v, got %v", ErrUnknownEncodingMode, err)
		}
	})
}

func TestGetModeBits(t *testing.T) {
	tests := []struct {
		mode    EncodingMode
		version int
		bits    int
	}{
		{EncodingModeNumeric, 1, 4},
		{EncodingModeAlphaNumeric, 1, 4},
		{EncodingModeLatin1, 1, 4},
		{EncodingModeKanji, 1, 4},
		{EncodingModeECI, 1, 16},

		// Micro QR
		{EncodingModeNumeric, -1, 0},
		{EncodingModeKanji, -4, 3},
	}

	for _, test := range tests {
		name := fmt.Sprintf("mode %v", test.mode)
		t.Run(name, func(t *testing.T) {
			b := &EncodeBlock{Mode: test.mode}
			if bits := b.GetModeBits(test.version); bits != test.bits {
				t.Errorf("Expected %b, got %b", test.bits, bits)
			}
		})
	}
}

func TestEncodeData(t *testing.T) {
	tests := []struct {
		content          string
		data             []byte
		mode             EncodingMode
		subMode          EncodingMode
		assignmentNumber uint
	}{
		// Alphanumeric
		{"A", []byte{0b00101000}, EncodingModeAlphaNumeric, 0, 0},
		{"AB", []byte{0b00111001, 0b10100000}, EncodingModeAlphaNumeric, 0, 0},
		{"ABCDEFGHIJKLM", []byte{0b00111001, 0b10101000, 0b10100101, 0b01000010, 0b10101110, 0b00010110, 0b01111010, 0b11100110, 0b01010110}, EncodingModeAlphaNumeric, 0, 0},
		{"AC-42", []byte{0b00111001, 0b11011100, 0b11100100, 0b00100000}, EncodingModeAlphaNumeric, 0, 0},

		// Numeric
		{"0", []byte{0x00}, EncodingModeNumeric, 0, 0},
		{"1", []byte{0b00010000}, EncodingModeNumeric, 0, 0},
		{"12", []byte{0b00011000}, EncodingModeNumeric, 0, 0},
		{"123", []byte{0b00011110, 0b11000000}, EncodingModeNumeric, 0, 0},
		{"012345", []byte{0b00000011, 0b00010101, 0b10010000}, EncodingModeNumeric, 0, 0},
		{"01234567", []byte{0b00000011, 0b00010101, 0b10011000, 0b01100000}, EncodingModeNumeric, 0, 0},
		{"0123456789012345", []byte{0b00000011, 0b00010101, 0b10011010, 0b10011011, 0b10000101, 0b00111010, 0b10010100}, EncodingModeNumeric, 0, 0},

		// Latin1
		{"a÷åäö", []byte{0x61, 0xf7, 0xe5, 0xe4, 0xf6}, EncodingModeLatin1, 0, 0},

		// Kanji
		{"点", []byte{0b01101100, 0b11111000}, EncodingModeKanji, 0, 0},
		{"茗", []byte{0b11010101, 0b01010000}, EncodingModeKanji, 0, 0},
		{"茗点", []byte{0b11010101, 0b01010011, 0b01100111, 0b11000000}, EncodingModeKanji, 0, 0},

		// ECI
		{"Ä点", []byte{0xc3, 0x84, 0xe7, 0x82, 0xb9}, EncodingModeECI, EncodingModeLatin1, 26},
	}

	for _, test := range tests {
		name := fmt.Sprintf("%v", test.content)
		t.Run(name, func(t *testing.T) {
			b := &EncodeBlock{
				Data:             test.content,
				Mode:             test.mode,
				SubMode:          test.subMode,
				AssignmentNumber: test.assignmentNumber,
			}

			data, err := EncodeDataWrapper(b)
			if err != nil {
				t.Error(err)
				return
			}
			if !bytes.Equal(data, test.data) {
				t.Errorf("Expected %b, got %b", test.data, data)
			}
		})
	}

	// Invalid mode
	t.Run("invalid mode", func(t *testing.T) {
		b := &EncodeBlock{
			Data: "1234567890",
			Mode: 0,
		}

		_, err := EncodeDataWrapper(b)
		if err != ErrUnknownEncodingMode {
			t.Errorf("Expected %v, got %v", ErrUnknownEncodingMode, err)
		}
	})

	// Invalid ECI assignment number
	t.Run("invalid ECI assignment number", func(t *testing.T) {
		b := &EncodeBlock{
			Data:             "1234567890",
			Mode:             EncodingModeECI,
			AssignmentNumber: 1000,
		}

		_, err := EncodeDataWrapper(b)
		if errors.Is(err, ErrUnknownAssignmentNumber) {
			t.Errorf("Expected %v, got %v", ErrUnknownAssignmentNumber, err)
		}
	})

	t.Run("failed to encode data", func(t *testing.T) {
		b := &EncodeBlock{
			Data: "ABC",
			Mode: EncodingModeNumeric,
		}

		_, err := EncodeDataWrapper(b)
		if err == nil {
			t.Error("Expected error")
		}
	})

	t.Run("failed to encode ECI data", func(t *testing.T) {
		b := &EncodeBlock{
			Data:             "Å",
			Mode:             EncodingModeECI,
			SubMode:          EncodingModeLatin1,
			AssignmentNumber: 7,
		}

		_, err := EncodeDataWrapper(b)
		if err == nil {
			t.Error("Expected error")
		}
	})

	// Failed to get length bits
	t.Run("failed to get length bits", func(t *testing.T) {
		b := &EncodeBlock{
			Data: "ABC",
			Mode: EncodingModeNumeric,
		}

		_, err := b.Encode(0, nil)
		if err == nil {
			t.Fatal("Expected error")
		}

		if !strings.Contains(err.Error(), "failed to get length bits") {
			t.Fatalf("Expected error to contain 'failed to get length bits', got %v", err)
		}
	})

	// Failed to calculate data bits count
	t.Run("failed to calculate data bits count", func(t *testing.T) {
		b := &EncodeBlock{
			Data:             "АБВГД",
			Mode:             EncodingModeECI,
			SubMode:          EncodingModeLatin1,
			AssignmentNumber: 5,
		}

		_, err := b.Encode(1, nil)
		if err == nil {
			t.Fatal("Expected error")
		}

		if !strings.Contains(err.Error(), "failed to calculate data bits count") {
			t.Fatalf("Expected error to contain 'failed to calculate data bits count', got %v", err)
		}
	})

	// Failed to encode data
	t.Run("failed to encode data", func(t *testing.T) {
		b := &EncodeBlock{
			Data: "ABC",
			Mode: EncodingModeNumeric,
		}

		queue := make(chan ValueBlock, 100)
		_, err := b.Encode(1, queue)
		close(queue)
		if err == nil {
			t.Fatal("Expected error")
		}

		if !strings.Contains(err.Error(), "failed to encode data") {
			t.Fatalf("Expected error to contain 'failed to encode data', got %v", err)
		}
	})
}

func TestEncode(t *testing.T) {
	tests := []struct {
		content          string
		mode             EncodingMode
		subMode          EncodingMode
		assignmentNumber uint
		version          int
		data             []byte
	}{
		// Alphanumeric
		{"A", EncodingModeAlphaNumeric, 0, 0, 1, []byte{0b00100000, 0b00001001, 0b01000000}},
		{"AB", EncodingModeAlphaNumeric, 0, 0, 1, []byte{0b00100000, 0b00010001, 0b11001101}},
		{"ABC", EncodingModeAlphaNumeric, 0, 0, 1, []byte{0b00100000, 0b00011001, 0b11001101, 0b00110000}},

		// Numeric
		{"1", EncodingModeNumeric, 0, 0, 1, []byte{0b00010000, 0b00000100, 0b01000000}},
		{"12", EncodingModeNumeric, 0, 0, 1, []byte{0b00010000, 0b00001000, 0b01100000}},
		{"123", EncodingModeNumeric, 0, 0, 1, []byte{0b00010000, 0b00001100, 0b01111011}},
		{"012345", EncodingModeNumeric, 0, 0, 1, []byte{0b00010000, 0b00011000, 0b00001100, 0b01010110, 0b01000000}},

		// Latin1
		{"abc", EncodingModeLatin1, 0, 0, 1, []byte{0b01000000, 0b00110110, 0b00010110, 0b00100110, 0b00110000}},
		{"äöå", EncodingModeLatin1, 0, 0, 1, []byte{0b01000000, 0b00111110, 0b01001111, 0b01101110, 0b01010000}},

		// Kanji
		{"点", EncodingModeKanji, 0, 0, 1, []byte{0b10000000, 0b00010110, 0b11001111, 0b10000000}},
		{"茗", EncodingModeKanji, 0, 0, 1, []byte{0b10000000, 0b00011101, 0b01010101, 0b00000000}},
		{"茗点", EncodingModeKanji, 0, 0, 1, []byte{0b10000000, 0b00101101, 0b01010101, 0b00110110, 0b01111100}},

		// ECI
		{"Ä点", EncodingModeECI, EncodingModeLatin1, 26, 1, []byte{0b01110001, 0b10100100, 0b00000101, 0b11000011, 0b10000100, 0b11100111, 0b10000010, 0b10111001}},
		{"abc", EncodingModeECI, EncodingModeLatin1, 26, 1, []byte{0b01110001, 0b10100100, 0b00000011, 0b01100001, 0b01100010, 0b01100011}},
		{"äbcöå", EncodingModeECI, EncodingModeLatin1, 26, 1, []byte{0b01110001, 0b10100100, 0b00001000, 0b11000011, 0b10100100, 0b01100010, 0b01100011, 0b11000011, 0b10110110, 0b11000011, 0b10100101}},
	}

	for _, test := range tests {
		name := fmt.Sprintf("%v", test.content)
		t.Run(name, func(t *testing.T) {
			b := &EncodeBlock{
				Data:             test.content,
				Mode:             test.mode,
				SubMode:          test.subMode,
				AssignmentNumber: test.assignmentNumber,
			}

			data, err := EncodeWrapper(b, test.version)
			if err != nil {
				t.Error(err)
				return
			}
			if !bytes.Equal(data, test.data) {
				t.Errorf("Expected %08b, got %08b", test.data, data)
			}
		})
	}
}
