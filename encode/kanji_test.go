package encode

import "testing"

func TestKanjiEncoder_Encode(t *testing.T) {
	tests := []struct {
		content  string
		expected []ValueBlock
	}{
		{
			content: "亜",
			expected: []ValueBlock{
				{Bits: 13, Value: 1439},
			},
		},
		{
			content: "亜亜",
			expected: []ValueBlock{
				{Bits: 13, Value: 1439},
				{Bits: 13, Value: 1439},
			},
		},
	}

	enc := &KanjiEncoder{}

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

	// Cannot encode to ShiftJIS
	t.Run("cannot encode", func(t *testing.T) {
		queue := make(chan ValueBlock, 100)
		err := enc.Encode("Ä", queue)
		if err == nil {
			t.Fatal("expected error")
		}
	})

	// The content is not Kanji, but it can be encoded to ShiftJIS
	t.Run("not Kanji", func(t *testing.T) {
		queue := make(chan ValueBlock, 100)
		err := enc.Encode("123", queue)
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestKanjiEncoder_CanEncode(t *testing.T) {
	tests := []struct {
		content  string
		expected bool
	}{
		{
			content:  "亜",
			expected: true,
		},
		{
			content:  "あア",
			expected: true,
		},
		{
			content:  "abc",
			expected: false,
		},
		{
			content:  "123",
			expected: false,
		},
	}

	enc := &KanjiEncoder{}

	for _, test := range tests {
		t.Run(test.content, func(t *testing.T) {
			canEncode := enc.CanEncode(test.content)
			if canEncode != test.expected {
				t.Errorf("expected %v, got %v", test.expected, canEncode)
			}
		})
	}
}

func TestKanjiEncoder_Size(t *testing.T) {
	tests := []struct {
		content  string
		expected int
	}{
		{
			content:  "亜",
			expected: 13,
		},
		{
			content:  "あア",
			expected: 26,
		},
	}

	enc := &KanjiEncoder{}

	for _, test := range tests {
		t.Run(test.content, func(t *testing.T) {
			size := enc.Size(test.content)
			if size != test.expected {
				t.Errorf("expected %v, got %v", test.expected, size)
			}
		})
	}
}

func TestKanjiEncoder_Mode(t *testing.T) {
	enc := &KanjiEncoder{}
	if enc.Mode() != EncodingModeKanji {
		t.Errorf("expected %v, got %v", EncodingModeKanji, enc.Mode())
	}
}
