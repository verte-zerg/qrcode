package qrcode

import "testing"

func TestPositionNext(t *testing.T) {
	tests := []struct {
		pos  position
		next position
	}{
		{position{10, 5, 11, -1, false, false}, position{9, 5, 11, -1, false, false}},
		{position{9, 5, 11, -1, false, false}, position{10, 4, 11, -1, false, false}},
		{position{10, 5, 11, 1, false, false}, position{9, 5, 11, 1, false, false}},
		{position{9, 5, 11, 1, false, false}, position{10, 6, 11, 1, false, false}},
		{position{9, 0, 11, -1, false, false}, position{8, 0, 11, 1, false, false}},
		{position{9, 10, 11, 1, false, false}, position{8, 10, 11, -1, false, false}},
		{position{0, 5, 11, -1, false, false}, position{0, 4, 11, -1, false, false}},
		{position{0, 5, 11, 1, false, false}, position{0, 6, 11, 1, false, false}},
		{position{7, 0, 11, -1, false, false}, position{5, 0, 11, 1, false, true}},
	}

	for _, test := range tests {
		test.pos.Next()
		if test.pos != test.next {
			t.Errorf("Expected %v, got %v", test.next, test.pos)
		}
	}
}

func TestGetSize(t *testing.T) {
	tests := []struct {
		version int
		size    int
	}{
		{-1, 11},
		{-4, 17},
		{1, 21},
		{40, 177},
	}

	for _, test := range tests {
		size := getSize(test.version)
		if size != test.size {
			t.Errorf("Expected %d, got %d", test.size, size)
		}
	}
}
