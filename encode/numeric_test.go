package encode

import (
	"reflect"
	"testing"
)

func TestNumericEncoder_Encode(t *testing.T) {
	tests := []struct {
		content string
		blocks  []ValueBlock
	}{
		{
			content: "1",
			blocks: []ValueBlock{
				{Bits: 4, Value: 1},
			},
		},
		{
			content: "12",
			blocks: []ValueBlock{
				{Bits: 7, Value: 12},
			},
		},
		{
			content: "123",
			blocks: []ValueBlock{
				{Bits: 10, Value: 123},
			},
		},
		{
			content: "1234",
			blocks: []ValueBlock{
				{Bits: 10, Value: 123},
				{Bits: 4, Value: 4},
			},
		},
	}

	enc := &NumericEncoder{}

	for _, test := range tests {
		t.Run(test.content, func(t *testing.T) {
			queue := make(chan ValueBlock, 100)
			err := enc.Encode(test.content, queue)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			close(queue)

			var blocks []ValueBlock
			for b := range queue {
				blocks = append(blocks, b)
			}

			if !reflect.DeepEqual(blocks, test.blocks) {
				t.Errorf("expected %v, got %v", test.blocks, blocks)
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

func TestNumericEncoder_Size(t *testing.T) {
	tests := []struct {
		content string
		size    int
	}{
		{
			content: "1",
			size:    4,
		},
		{
			content: "12",
			size:    7,
		},
		{
			content: "123",
			size:    10,
		},
		{
			content: "1234",
			size:    14,
		},
	}

	enc := &NumericEncoder{}

	for _, test := range tests {
		t.Run(test.content, func(t *testing.T) {
			size := enc.Size(test.content)
			if size != test.size {
				t.Errorf("expected %v, got %v", test.size, size)
			}
		})
	}
}

func TestNumericEncoder_CanEncode(t *testing.T) {
	tests := []struct {
		content  string
		expected bool
	}{
		{
			content:  "123",
			expected: true,
		},
		{
			content:  "abc",
			expected: false,
		},
		{
			content:  "ア",
			expected: false,
		},
	}

	enc := &NumericEncoder{}

	for _, test := range tests {
		t.Run(test.content, func(t *testing.T) {
			canEncode := enc.CanEncode(test.content)
			if canEncode != test.expected {
				t.Errorf("expected %v, got %v", test.expected, canEncode)
			}
		})
	}
}

func TestNumericEncoder_Mode(t *testing.T) {
	enc := &NumericEncoder{}
	if enc.Mode() != EncodingModeNumeric {
		t.Errorf("expected %v, got %v", EncodingModeNumeric, enc.Mode())
	}
}
