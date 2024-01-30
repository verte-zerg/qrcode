package qrcode

import (
	"testing"
)

func TransposeMatrix(data [][]Cell) [][]Cell {
	if len(data) == 0 {
		return [][]Cell{}
	}

	transposedData := make([][]Cell, len(data[0]))
	for idx := range transposedData {
		transposedData[idx] = make([]Cell, len(data))

		for jdx := range transposedData[idx] {
			transposedData[idx][jdx] = data[jdx][idx]
		}
	}

	return transposedData
}

func CompareMatrix(data1 [][]Cell, data2 [][]Cell) bool {
	if len(data1) != len(data2) {
		return false
	}

	for idx, row := range data1 {
		if len(row) != len(data2[idx]) {
			return false
		}

		for jdx, cell := range row {
			if cell != data2[idx][jdx] {
				return false
			}
		}
	}

	return true
}

func TestCalculatePenaltyRule1(t *testing.T) {
	whiteCell := Cell{Type: CellTypeData, Value: false}
	blackCell := Cell{Type: CellTypeData, Value: true}

	tests := []struct {
		data     [][]Cell
		expected int
	}{
		{
			data:     [][]Cell{},
			expected: 0,
		},
		{
			data: [][]Cell{
				{whiteCell, whiteCell, whiteCell, whiteCell},
				{blackCell, blackCell, blackCell, blackCell},
			},
			expected: 0,
		},
		{
			data: [][]Cell{
				{whiteCell, whiteCell, whiteCell, whiteCell, whiteCell, blackCell},
			},
			expected: 3,
		},
		{
			data: [][]Cell{
				{whiteCell, whiteCell, whiteCell, whiteCell, whiteCell},
				{blackCell, blackCell, blackCell, blackCell, blackCell},
			},
			expected: 3 + 3,
		},
		{
			data: [][]Cell{
				{whiteCell, whiteCell, whiteCell, whiteCell, whiteCell, whiteCell},
			},
			expected: 4,
		},
	}

	for _, test := range tests {
		penalty := CalculatePenaltyRule1(test.data)

		if penalty != test.expected {
			t.Errorf("CalculatePenaltyRule1 failed, expected %d, got %d", test.expected, penalty)
		}

		transposedData := TransposeMatrix(test.data)
		penalty = CalculatePenaltyRule1(transposedData)

		if penalty != test.expected {
			t.Errorf("CalculatePenaltyRule1 failed, expected %d, got %d", test.expected, penalty)
		}
	}
}

func TestCalculatePenaltyRule2(t *testing.T) {
	whiteCell := Cell{Type: CellTypeData, Value: false}
	blackCell := Cell{Type: CellTypeData, Value: true}

	tests := []struct {
		data     [][]Cell
		expected int
	}{
		{
			data:     [][]Cell{},
			expected: 0,
		},
		{
			data: [][]Cell{
				{whiteCell, whiteCell, whiteCell, whiteCell},
			},
			expected: 0,
		},
		{
			data: [][]Cell{
				{whiteCell, whiteCell, whiteCell, whiteCell, whiteCell},
				{whiteCell, whiteCell, blackCell, blackCell, blackCell},
			},
			expected: 3,
		},
		{
			data: [][]Cell{
				{whiteCell, whiteCell, whiteCell},
				{whiteCell, whiteCell, whiteCell},
				{whiteCell, whiteCell, whiteCell},
			},
			expected: 12,
		},
		{
			data: [][]Cell{
				{whiteCell, whiteCell, whiteCell, blackCell, blackCell, blackCell},
				{whiteCell, whiteCell, whiteCell, blackCell, blackCell, blackCell},
				{whiteCell, whiteCell, whiteCell, blackCell, blackCell, blackCell},
			},
			expected: 24,
		},
	}

	for _, test := range tests {
		penalty := CalculatePenaltyRule2(test.data)

		if penalty != test.expected {
			t.Errorf("CalculatePenaltyRule2 failed, expected %d, got %d", test.expected, penalty)
		}

		transposedData := TransposeMatrix(test.data)

		penalty = CalculatePenaltyRule2(transposedData)

		if penalty != test.expected {
			t.Errorf("CalculatePenaltyRule2 failed, expected %d, got %d", test.expected, penalty)
		}
	}
}

func TestCalculatePenaltyRule3(t *testing.T) {
	whiteCell := Cell{Type: CellTypeData, Value: false}
	blackCell := Cell{Type: CellTypeData, Value: true}

	tests := []struct {
		data     [][]Cell
		expected int
	}{
		{
			data:     [][]Cell{},
			expected: 0,
		},
		{
			data: [][]Cell{
				{blackCell, whiteCell, blackCell, blackCell, blackCell, whiteCell, blackCell, whiteCell, whiteCell, whiteCell, whiteCell},
			},
			expected: 40,
		},
		{
			data: [][]Cell{
				{whiteCell, whiteCell, whiteCell, whiteCell, blackCell, whiteCell, blackCell, blackCell, blackCell, whiteCell, blackCell},
			},
			expected: 40,
		},
		{
			data: [][]Cell{
				{blackCell, whiteCell, blackCell, blackCell, blackCell, whiteCell, blackCell, whiteCell, whiteCell, whiteCell, whiteCell},
				{whiteCell, whiteCell, whiteCell, whiteCell, blackCell, whiteCell, blackCell, blackCell, blackCell, whiteCell, blackCell},
			},
			expected: 80,
		},
	}

	for _, test := range tests {
		penalty := CalculatePenaltyRule3(test.data)

		if penalty != test.expected {
			t.Errorf("CalculatePenaltyRule3 failed, expected %d, got %d", test.expected, penalty)
		}

		transposedData := TransposeMatrix(test.data)

		penalty = CalculatePenaltyRule3(transposedData)

		if penalty != test.expected {
			t.Errorf("CalculatePenaltyRule3 failed, expected %d, got %d", test.expected, penalty)
		}
	}
}

func TestCalculatePenaltyRule4(t *testing.T) {
	whiteCell := Cell{Type: CellTypeData, Value: false}
	blackCell := Cell{Type: CellTypeData, Value: true}

	tests := []struct {
		data     [][]Cell
		expected int
	}{
		{
			data:     [][]Cell{},
			expected: 0,
		},
		{
			data: [][]Cell{
				{blackCell},
			},
			expected: 100,
		},
		{
			data: [][]Cell{
				{blackCell, blackCell},
				{whiteCell, whiteCell},
			},
			expected: 0,
		},
		{
			data: [][]Cell{
				{blackCell, blackCell},
				{blackCell, blackCell},
				{blackCell, whiteCell},
			},
			expected: 60,
		},
		{
			data: [][]Cell{
				{whiteCell, whiteCell},
				{whiteCell, whiteCell},
				{whiteCell, blackCell},
			},
			expected: 60,
		},
	}

	for _, test := range tests {
		penalty := CalculatePenaltyRule4(test.data)

		if penalty != test.expected {
			t.Errorf("CalculatePenaltyRule4 failed, expected %d, got %d", test.expected, penalty)
		}

		transposedData := TransposeMatrix(test.data)

		penalty = CalculatePenaltyRule4(transposedData)

		if penalty != test.expected {
			t.Errorf("CalculatePenaltyRule4 failed, expected %d, got %d", test.expected, penalty)
		}
	}
}

func TestApplyMask(t *testing.T) {
	whiteCell := Cell{Type: CellTypeData, Value: false}
	blackCell := Cell{Type: CellTypeData, Value: true}

	tests := []struct {
		data     [][]Cell
		mask     int
		expected [][]Cell
	}{
		// Mask won't apply to not data cells
		{
			data: [][]Cell{
				{Cell{Type: CellTypeFormat, Value: false}},
			},
			mask: 0,
			expected: [][]Cell{
				{Cell{Type: CellTypeFormat, Value: false}},
			},
		},
		// Mask 0
		{
			data: [][]Cell{
				{whiteCell, whiteCell},
				{whiteCell, whiteCell},
			},
			mask: 0,
			expected: [][]Cell{
				{blackCell, whiteCell},
				{whiteCell, blackCell},
			},
		},
		// Mask 1
		{
			data: [][]Cell{
				{whiteCell, whiteCell},
				{whiteCell, whiteCell},
			},
			mask: 1,
			expected: [][]Cell{
				{blackCell, blackCell},
				{whiteCell, whiteCell},
			},
		},
		// Mask 2
		{
			data: [][]Cell{
				{whiteCell, whiteCell, whiteCell, whiteCell},
				{whiteCell, whiteCell, whiteCell, whiteCell},
			},
			mask: 2,
			expected: [][]Cell{
				{blackCell, whiteCell, whiteCell, blackCell},
				{blackCell, whiteCell, whiteCell, blackCell},
			},
		},
		// Mask 3
		{
			data: [][]Cell{
				{whiteCell, whiteCell, whiteCell, whiteCell},
				{whiteCell, whiteCell, whiteCell, whiteCell},
				{whiteCell, whiteCell, whiteCell, whiteCell},
				{whiteCell, whiteCell, whiteCell, whiteCell},
			},
			mask: 3,
			expected: [][]Cell{
				{blackCell, whiteCell, whiteCell, blackCell},
				{whiteCell, whiteCell, blackCell, whiteCell},
				{whiteCell, blackCell, whiteCell, whiteCell},
				{blackCell, whiteCell, whiteCell, blackCell},
			},
		},
		// Mask 4
		{
			data: [][]Cell{
				{whiteCell, whiteCell, whiteCell, whiteCell},
				{whiteCell, whiteCell, whiteCell, whiteCell},
				{whiteCell, whiteCell, whiteCell, whiteCell},
				{whiteCell, whiteCell, whiteCell, whiteCell},
			},
			mask: 4,
			expected: [][]Cell{
				{blackCell, blackCell, blackCell, whiteCell},
				{blackCell, blackCell, blackCell, whiteCell},
				{whiteCell, whiteCell, whiteCell, blackCell},
				{whiteCell, whiteCell, whiteCell, blackCell},
			},
		},
		// Mask 5
		{
			data: [][]Cell{
				{whiteCell, whiteCell, whiteCell, whiteCell, whiteCell},
				{whiteCell, whiteCell, whiteCell, whiteCell, whiteCell},
				{whiteCell, whiteCell, whiteCell, whiteCell, whiteCell},
				{whiteCell, whiteCell, whiteCell, whiteCell, whiteCell},
				{whiteCell, whiteCell, whiteCell, whiteCell, whiteCell},
			},
			mask: 5,
			expected: [][]Cell{
				{blackCell, blackCell, blackCell, blackCell, blackCell},
				{blackCell, whiteCell, whiteCell, whiteCell, whiteCell},
				{blackCell, whiteCell, whiteCell, blackCell, whiteCell},
				{blackCell, whiteCell, blackCell, whiteCell, blackCell},
				{blackCell, whiteCell, whiteCell, blackCell, whiteCell},
			},
		},
		// Mask 6
		{
			data: [][]Cell{
				{whiteCell, whiteCell, whiteCell, whiteCell},
				{whiteCell, whiteCell, whiteCell, whiteCell},
				{whiteCell, whiteCell, whiteCell, whiteCell},
				{whiteCell, whiteCell, whiteCell, whiteCell},
			},
			mask: 6,
			expected: [][]Cell{
				{blackCell, blackCell, blackCell, blackCell},
				{blackCell, blackCell, blackCell, whiteCell},
				{blackCell, blackCell, whiteCell, blackCell},
				{blackCell, whiteCell, blackCell, whiteCell},
			},
		},
		// Mask 7
		{
			data: [][]Cell{
				{whiteCell, whiteCell, whiteCell, whiteCell},
				{whiteCell, whiteCell, whiteCell, whiteCell},
				{whiteCell, whiteCell, whiteCell, whiteCell},
				{whiteCell, whiteCell, whiteCell, whiteCell},
			},
			mask: 7,
			expected: [][]Cell{
				{blackCell, whiteCell, blackCell, whiteCell},
				{whiteCell, whiteCell, whiteCell, blackCell},
				{blackCell, whiteCell, whiteCell, whiteCell},
				{whiteCell, blackCell, whiteCell, blackCell},
			},
		},
	}

	for _, test := range tests {
		ApplyMask(test.data, test.mask)

		if !CompareMatrix(test.data, test.expected) {
			t.Errorf("ApplyMask failed, expected %v, got %v", test.expected, test.data)
		}
	}

}
