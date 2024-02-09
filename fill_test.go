package qrcode

import "testing"

func TestPositionNext(t *testing.T) {
	tests := []struct {
		pos  Position
		next Position
	}{
		{Position{10, 5, 11, -1, false, false}, Position{9, 5, 11, -1, false, false}},
		{Position{9, 5, 11, -1, false, false}, Position{10, 4, 11, -1, false, false}},
		{Position{10, 5, 11, 1, false, false}, Position{9, 5, 11, 1, false, false}},
		{Position{9, 5, 11, 1, false, false}, Position{10, 6, 11, 1, false, false}},
		{Position{9, 0, 11, -1, false, false}, Position{8, 0, 11, 1, false, false}},
		{Position{9, 10, 11, 1, false, false}, Position{8, 10, 11, -1, false, false}},
		{Position{0, 5, 11, -1, false, false}, Position{0, 4, 11, -1, false, false}},
		{Position{0, 5, 11, 1, false, false}, Position{0, 6, 11, 1, false, false}},
		{Position{7, 0, 11, -1, false, false}, Position{5, 0, 11, 1, false, true}},
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
		size := GetSize(test.version)
		if size != test.size {
			t.Errorf("Expected %d, got %d", test.size, size)
		}
	}
}
