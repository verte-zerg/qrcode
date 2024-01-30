package encode

import (
	"testing"
)

func TestAlphaNumericEncoder_Encode(t *testing.T) {
	tests := []struct {
		content  string
		expected []ValueBlock
	}{
		{
			content: "ABC",
			expected: []ValueBlock{
				{Bits: 11, Value: 461},
				{Bits: 6, Value: 12},
			},
		},
	}

	enc := &AlphaNumericEncoder{}

	for _, test := range tests {
		t.Run(test.content, func(t *testing.T) {
			queue := make(chan ValueBlock, 100)
			err := enc.Encode(test.content, queue)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			close(queue)

			i := 0
			for block := range queue {
				if i >= len(test.expected) {
					t.Fatalf("unexpected block: %v", block)
				}

				if block != test.expected[i] {
					t.Errorf("expected %v, got %v", test.expected[i], block)
				}
				i++
			}
		})
	}

	// Invalid content
	t.Run("invalid content", func(t *testing.T) {
		queue := make(chan ValueBlock, 100)
		err := enc.Encode("abc", queue)
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestAlphaNumericEncoder_CanEncode(t *testing.T) {
	tests := []struct {
		content  string
		expected bool
	}{
		{
			content:  "ABC",
			expected: true,
		},
		{
			content:  "123",
			expected: true,
		},
		{
			content:  "ABC:/",
			expected: true,
		},
		{
			content:  "abc",
			expected: false,
		},
		{
			content:  "Å",
			expected: false,
		},
	}

	enc := &AlphaNumericEncoder{}

	for _, test := range tests {
		t.Run(test.content, func(t *testing.T) {
			canEncode := enc.CanEncode(test.content)
			if canEncode != test.expected {
				t.Errorf("expected %v, got %v", test.expected, canEncode)
			}
		})
	}
}

func TestAlphaNumericEncoder_Size(t *testing.T) {
	tests := []struct {
		content  string
		expected int
	}{
		{
			content:  "ABC",
			expected: 17,
		},
		{
			content:  "AB",
			expected: 11,
		},
		{
			content:  "A",
			expected: 6,
		},
		{
			content:  "ABC:",
			expected: 22,
		},
	}

	enc := &AlphaNumericEncoder{}

	for _, test := range tests {
		t.Run(test.content, func(t *testing.T) {
			size := enc.Size(test.content)
			if size != test.expected {
				t.Errorf("expected %v, got %v", test.expected, size)
			}
		})
	}
}

func TestAlphaNumericEncoder_Mode(t *testing.T) {
	enc := &AlphaNumericEncoder{}
	if enc.Mode() != EncodingModeAlphaNumeric {
		t.Errorf("expected %v, got %v", EncodingModeAlphaNumeric, enc.Mode())
	}
}
