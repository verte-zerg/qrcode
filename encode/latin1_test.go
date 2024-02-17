package encode

import "testing"

func TestByteEncoder_Encode(t *testing.T) {
	tests := []struct {
		content string
		blocks  []ValueBlock
	}{
		{
			content: "ABC",
			blocks: []ValueBlock{
				{Bits: 8, Value: 65},
				{Bits: 8, Value: 66},
				{Bits: 8, Value: 67},
			},
		},
		{
			content: "Å",
			blocks: []ValueBlock{
				{Bits: 8, Value: 197},
			},
		},
		{
			content: "123",
			blocks: []ValueBlock{
				{Bits: 8, Value: 49},
				{Bits: 8, Value: 50},
				{Bits: 8, Value: 51},
			},
		},
	}

	enc := &byteEncoder{}

	for _, test := range tests {
		t.Run(test.content, func(t *testing.T) {
			queue := make(chan ValueBlock, 100)
			err := enc.Encode(test.content, queue)
			if err != nil {
				t.Fatal(err)
			}

			close(queue)

			i := 0
			for block := range queue {
				if i >= len(test.blocks) {
					t.Fatalf("unexpected block: %v", block)
				}

				if block != test.blocks[i] {
					t.Errorf("expected %v, got %v", test.blocks[i], block)
				}
				i++
			}
		})
	}

	// Invalid content
	t.Run("invalid content", func(t *testing.T) {
		queue := make(chan ValueBlock, 100)
		err := enc.Encode("ア", queue)
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestByteEncoder_CanEncode(t *testing.T) {
	tests := []struct {
		content  string
		expected bool
	}{
		{
			content:  "ABC",
			expected: true,
		},
		{
			content:  "Å",
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
			content:  "ア",
			expected: false,
		},
	}

	enc := &byteEncoder{}

	for _, test := range tests {
		t.Run(test.content, func(t *testing.T) {
			if enc.CanEncode(test.content) != test.expected {
				t.Errorf("expected %v, got %v", test.expected, !test.expected)
			}
		})
	}
}

func TestByteEncoder_Size(t *testing.T) {
	tests := []struct {
		content  string
		expected int
	}{
		{
			content:  "ABC",
			expected: 24,
		},
		{
			content:  "Å",
			expected: 8,
		},
		{
			content:  "123",
			expected: 24,
		},
	}

	enc := &byteEncoder{}

	for _, test := range tests {
		t.Run(test.content, func(t *testing.T) {
			if enc.Size(test.content) != test.expected {
				t.Errorf("expected %v, got %v", test.expected, enc.Size(test.content))
			}
		})
	}
}

func TestByteEncoder_Mode(t *testing.T) {
	enc := &byteEncoder{}
	if enc.Mode() != EncodingModeByte {
		t.Errorf("expected %v, got %v", EncodingModeByte, enc.Mode())
	}
}
