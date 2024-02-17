package encode

import "testing"

func TestEciEncoder_Encode(t *testing.T) {
	tests := []struct {
		assignmentNumber uint
		content          string
		expected         []ValueBlock
	}{
		{
			assignmentNumber: 26,
			content:          "ABÏア",
			expected: []ValueBlock{
				{Bits: 8, Value: 65},
				{Bits: 8, Value: 66},
				{Bits: 8, Value: 195},
				{Bits: 8, Value: 143},
				{Bits: 8, Value: 227},
				{Bits: 8, Value: 130},
				{Bits: 8, Value: 162},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.content, func(t *testing.T) {
			eci := &eciEncoder{
				AssignmentNumber: test.assignmentNumber,
				DataMode:         EncodingModeByte,
			}

			queue := make(chan ValueBlock, 100)
			err := eci.Encode(test.content, queue)
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
}

func TestEciEncoder_CanEncode(t *testing.T) {
	tests := []struct {
		assignmentNumber uint
		content          string
		expected         bool
	}{
		{
			assignmentNumber: 26,
			content:          "ABÏア",
			expected:         true,
		},
		{
			assignmentNumber: 5,
			content:          "ABÏアÄ",
			expected:         false,
		},
	}

	for _, test := range tests {
		t.Run(test.content, func(t *testing.T) {
			eci := &eciEncoder{
				AssignmentNumber: test.assignmentNumber,
				DataMode:         EncodingModeByte,
			}

			if eci.CanEncode(test.content) != test.expected {
				t.Errorf("expected %v, got %v", test.expected, !test.expected)
			}
		})
	}

	// Invalid assignment number
	t.Run("invalid assignment number", func(t *testing.T) {
		eci := &eciEncoder{
			AssignmentNumber: 100,
			DataMode:         EncodingModeByte,
		}

		if eci.CanEncode("abc") {
			t.Error("expected false")
		}
	})

	// Invalid content
	t.Run("invalid content", func(t *testing.T) {
		eci := &eciEncoder{
			AssignmentNumber: 4,
			DataMode:         EncodingModeByte,
		}

		if eci.CanEncode("АБВГД") {
			t.Error("expected false")
		}
	})
}

func TestEciEncoder_Size(t *testing.T) {
	tests := []struct {
		assignmentNumber uint
		content          string
		expected         int
	}{
		{
			assignmentNumber: 26,
			content:          "ABÏア",
			expected:         56,
		},
		{
			assignmentNumber: 7,
			content:          "АБВГД",
			expected:         40,
		},
	}

	for _, test := range tests {
		t.Run(test.content, func(t *testing.T) {
			eci := &eciEncoder{
				AssignmentNumber: test.assignmentNumber,
				DataMode:         EncodingModeByte,
			}

			size := eci.Size(test.content)
			if size != test.expected {
				t.Errorf("expected %v, got %v", test.expected, size)
			}
		})
	}

	// Invalid assignment number
	t.Run("invalid assignment number", func(t *testing.T) {
		eci := &eciEncoder{
			AssignmentNumber: 100,
			DataMode:         EncodingModeByte,
		}

		if eci.Size("abc") != 0 {
			t.Error("expected 0")
		}
	})

	// Invalid content
	t.Run("invalid content", func(t *testing.T) {
		eci := &eciEncoder{
			AssignmentNumber: 4,
			DataMode:         EncodingModeByte,
		}

		if eci.Size("АБВГД") != 0 {
			t.Error("expected 0")
		}
	})
}

func TestEciEncoder_Mode(t *testing.T) {
	eci := &eciEncoder{}
	if eci.Mode() != EncodingModeECI {
		t.Errorf("expected %v, got %v", EncodingModeECI, eci.Mode())
	}
}
