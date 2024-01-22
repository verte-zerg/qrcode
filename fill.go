package qrcode

import "fmt"

type CellType int

const (
	CellTypeData CellType = iota
	CellTypeFormat
	CellTypeVersion
	CellTypeAlignmentPattern
	CellTypeSearchPattern
	CellTypeSyncPattern
	CellTypeDelimiter
)

type VersionBlockInfo struct {
	ErrorCorrectionLevel ErrorCorrectionLevel
	MaskPattern          int
}

var (
	SearchPattern = [7][7]bool{
		{true, true, true, true, true, true, true},
		{true, false, false, false, false, false, true},
		{true, false, true, true, true, false, true},
		{true, false, true, true, true, false, true},
		{true, false, true, true, true, false, true},
		{true, false, false, false, false, false, true},
		{true, true, true, true, true, true, true},
	}

	AlignmentPattern = [5][5]bool{
		{true, true, true, true, true},
		{true, false, false, false, true},
		{true, false, true, false, true},
		{true, false, false, false, true},
		{true, true, true, true, true},
	}

	AlignmentPositions = [41][]int{
		{}, // added for shift versions, as first version has index 1
		{},
		{6, 18},
		{6, 22},
		{6, 26},
		{6, 30},
		{6, 34},
		{6, 22, 38},
		{6, 24, 42},
		{6, 26, 46},
		{6, 28, 50},
		{6, 30, 54},
		{6, 32, 58},
		{6, 34, 62},
		{6, 26, 46, 66},
		{6, 26, 48, 70},
		{6, 26, 50, 74},
		{6, 30, 54, 78},
		{6, 30, 56, 82},
		{6, 30, 58, 86},
		{6, 34, 62, 90},
		{6, 28, 50, 72, 94},
		{6, 26, 50, 74, 98},
		{6, 30, 54, 78, 102},
		{6, 28, 54, 80, 106},
		{6, 32, 58, 84, 110},
		{6, 30, 58, 86, 114},
		{6, 34, 62, 90, 118},
		{6, 26, 50, 74, 98, 122},
		{6, 30, 54, 78, 102, 126},
		{6, 26, 52, 78, 104, 130},
		{6, 30, 56, 82, 108, 134},
		{6, 34, 60, 86, 112, 138},
		{6, 30, 58, 86, 114, 142},
		{6, 34, 62, 90, 118, 146},
		{6, 30, 54, 78, 102, 126, 150},
		{6, 24, 50, 76, 102, 128, 154},
		{6, 28, 54, 80, 106, 132, 158},
		{6, 32, 58, 84, 110, 136, 162},
		{6, 26, 54, 82, 110, 138, 166},
		{6, 30, 58, 86, 114, 142, 170},
	}

	// VersionValues is a table of version information for each version.
	// The first index is the version number 7, as the version block starts at version 7.
	VersionValues = [34][18]byte{
		{0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 1, 0, 0, 1, 0, 1, 0, 0},
		{0, 0, 1, 0, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 1, 1, 0, 0},
		{0, 0, 1, 0, 0, 1, 1, 0, 1, 0, 1, 0, 0, 1, 1, 0, 0, 1},
		{0, 0, 1, 0, 1, 0, 0, 1, 0, 0, 1, 1, 0, 1, 0, 0, 1, 1},
		{0, 0, 1, 0, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 0, 1, 1, 0},
		{0, 0, 1, 1, 0, 0, 0, 1, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0},
		{0, 0, 1, 1, 0, 1, 1, 0, 0, 0, 0, 1, 0, 0, 0, 1, 1, 1},
		{0, 0, 1, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 1, 1, 0, 1},
		{0, 0, 1, 1, 1, 1, 1, 0, 0, 1, 0, 0, 1, 0, 1, 0, 0, 0},
		{0, 1, 0, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 1, 1, 0, 0, 0},
		{0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 1, 0, 1, 1, 1, 0, 1},
		{0, 1, 0, 0, 1, 0, 1, 0, 1, 0, 0, 0, 0, 1, 0, 1, 1, 1},
		{0, 1, 0, 0, 1, 1, 0, 1, 0, 1, 0, 0, 1, 1, 0, 0, 1, 0},
		{0, 1, 0, 1, 0, 0, 1, 0, 0, 1, 1, 0, 1, 0, 0, 1, 1, 0},
		{0, 1, 0, 1, 0, 1, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0, 1, 1},
		{0, 1, 0, 1, 1, 0, 1, 0, 0, 0, 1, 1, 0, 0, 1, 0, 0, 1},
		{0, 1, 0, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 0},
		{0, 1, 1, 0, 0, 0, 1, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 0},
		{0, 1, 1, 0, 0, 1, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 1},
		{0, 1, 1, 0, 1, 0, 1, 1, 1, 1, 1, 0, 1, 0, 1, 0, 1, 1},
		{0, 1, 1, 0, 1, 1, 0, 0, 0, 0, 1, 0, 0, 0, 1, 1, 1, 0},
		{0, 1, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 1, 1, 0, 1, 0},
		{0, 1, 1, 1, 0, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1},
		{0, 1, 1, 1, 1, 0, 1, 1, 0, 1, 0, 1, 1, 1, 0, 1, 0, 1},
		{0, 1, 1, 1, 1, 1, 0, 0, 1, 0, 0, 1, 0, 1, 0, 0, 0, 0},
		{1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 1, 0, 1, 0, 1, 0, 1},
		{1, 0, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 1, 1, 0, 0, 0, 0},
		{1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 1, 0, 1, 1, 1, 0, 1, 0},
		{1, 0, 0, 0, 1, 1, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1},
		{1, 0, 0, 1, 0, 0, 1, 0, 1, 1, 0, 0, 0, 0, 1, 0, 1, 1},
		{1, 0, 0, 1, 0, 1, 0, 1, 0, 0, 0, 0, 1, 0, 1, 1, 1, 0},
		{1, 0, 0, 1, 1, 0, 1, 0, 1, 0, 0, 1, 1, 0, 0, 1, 0, 0},
		{1, 0, 0, 1, 1, 1, 0, 1, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1},
		{1, 0, 1, 0, 0, 0, 1, 1, 0, 0, 0, 1, 1, 0, 1, 0, 0, 1},
	}

	// FormatValues is a table of format information for each error correction level and mask pattern.
	// Structure: [error correction level][mask pattern][15 bits]
	FormatValues = [4][8][15]byte{
		{
			{1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 0, 0, 1, 0, 0},
			{1, 1, 1, 0, 0, 1, 0, 1, 1, 1, 1, 0, 0, 1, 1},
			{1, 1, 1, 1, 1, 0, 1, 1, 0, 1, 0, 1, 0, 1, 0},
			{1, 1, 1, 1, 0, 0, 0, 1, 0, 0, 1, 1, 1, 0, 1},
			{1, 1, 0, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 1, 1},
			{1, 1, 0, 0, 0, 1, 1, 0, 0, 0, 1, 1, 0, 0, 0},
			{1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1},
			{1, 1, 0, 1, 0, 0, 1, 0, 1, 1, 1, 0, 1, 1, 0},
		},
		{
			{1, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0},
			{1, 0, 1, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 1},
			{1, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 0, 0},
			{1, 0, 1, 1, 0, 1, 1, 0, 1, 0, 0, 1, 0, 1, 1},
			{1, 0, 0, 0, 1, 0, 1, 1, 1, 1, 1, 1, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 1, 1, 1, 0},
			{1, 0, 0, 1, 1, 1, 1, 1, 0, 0, 1, 0, 1, 1, 1},
			{1, 0, 0, 1, 0, 1, 0, 1, 0, 1, 0, 0, 0, 0, 0},
		},
		{
			{0, 1, 1, 0, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1},
			{0, 1, 1, 0, 0, 0, 0, 0, 1, 1, 0, 1, 0, 0, 0},
			{0, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 0, 0, 0, 1},
			{0, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0},
			{0, 1, 0, 0, 1, 0, 0, 1, 0, 1, 1, 0, 1, 0, 0},
			{0, 1, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 1, 1},
			{0, 1, 0, 1, 1, 1, 0, 1, 1, 0, 1, 1, 0, 1, 0},
			{0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 0, 1},
		},
		{
			{0, 0, 1, 0, 1, 1, 0, 1, 0, 0, 0, 1, 0, 0, 1},
			{0, 0, 1, 0, 0, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0},
			{0, 0, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 1, 1, 1},
			{0, 0, 1, 1, 0, 0, 1, 1, 1, 0, 1, 0, 0, 0, 0},
			{0, 0, 0, 0, 1, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0},
			{0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 1, 0, 1, 0, 1},
			{0, 0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 1, 1, 0, 0},
			{0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 1, 1, 0, 1, 1},
		},
	}
)

type Cell struct {
	Value bool
	Type  CellType
}

type Position struct {
	X, Y        int
	Size        int
	Direction   int
	evenReverse bool
}

func (p *Position) Next() {
	even := (p.X%2 == 0) != p.evenReverse

	if even {
		if p.X != 0 {
			p.X--
		} else {
			p.Y += p.Direction
		}
	} else {
		p.X++
		p.Y += p.Direction

		if p.Y == p.Size || p.Y == -1 {
			p.Direction *= -1
			p.Y += p.Direction
			p.X -= 2
		}
	}

	if p.X == 6 {
		p.X--
		p.evenReverse = true
	}
}

func GetSize(version int) int {
	return 17 + 4*version
}

func FillSearchPattern(field [][]Cell) {
	size := len(field)
	for i := 0; i < 7; i++ {
		for j := 0; j < 7; j++ {
			value := SearchPattern[i][j]
			field[i][j] = Cell{Value: value, Type: CellTypeSearchPattern}
			field[i][size-1-j] = Cell{Value: value, Type: CellTypeSearchPattern}
			field[size-1-i][j] = Cell{Value: value, Type: CellTypeSearchPattern}
		}
	}

	for i := 0; i < 8; i++ {
		// top-left
		field[i][7] = Cell{Value: false, Type: CellTypeSearchPattern}
		field[7][i] = Cell{Value: false, Type: CellTypeSearchPattern}

		// top-right
		field[i][size-8] = Cell{Value: false, Type: CellTypeSearchPattern}
		field[7][size-1-i] = Cell{Value: false, Type: CellTypeSearchPattern}

		// bottom-left
		field[size-8][i] = Cell{Value: false, Type: CellTypeSearchPattern}
		field[size-1-i][7] = Cell{Value: false, Type: CellTypeSearchPattern}
	}
}

func FillSyncPattern(field [][]Cell) {
	size := len(field)
	value := true
	for i := 8; i < size-8; i++ {
		field[6][i] = Cell{Value: value, Type: CellTypeSyncPattern}
		field[i][6] = Cell{Value: value, Type: CellTypeSyncPattern}
		value = !value
	}
}

func FillAlignmentPattern(field [][]Cell, version int) {
	positions := AlignmentPositions[version]
	for idx, i := range positions {
		for jdx, j := range positions {
			if (idx == 0 && jdx == 0) || (idx == 0 && jdx == len(positions)-1) || (idx == len(positions)-1 && jdx == 0) {
				continue
			}
			for k := 0; k < 5; k++ {
				for l := 0; l < 5; l++ {
					field[i+k-2][j+l-2] = Cell{Value: AlignmentPattern[k][l], Type: CellTypeAlignmentPattern}
				}
			}
		}
	}
}

func FillVersionBlock(field [][]Cell, version int) {
	size := len(field)
	if version < 7 {
		return
	}
	values := VersionValues[version-7]
	for i := 0; i < 6; i++ {
		for j := 0; j < 3; j++ {
			field[size-11+j][i] = Cell{Value: values[i*3+j] == 1, Type: CellTypeVersion}
			field[5-i][size-9-j] = Cell{Value: values[i*3+j] == 1, Type: CellTypeVersion}
		}
	}
}

func FillFormatBlock(field [][]Cell, errorCorrectionLevel ErrorCorrectionLevel, maskPattern int) {
	size := len(field)
	values := FormatValues[errorCorrectionLevel][maskPattern]
	// values := make([]byte, 15)
	for i := 0; i < 8; i++ {
		shiftDelimiter := 0
		if i > 5 {
			shiftDelimiter = 1
		}
		// up -> down, on the right side of the top-left search pattern
		field[i+shiftDelimiter][8] = Cell{Value: values[14-i] == 1, Type: CellTypeFormat}

		// left -> right, on the bottom side of the top-left search pattern
		field[8][i+shiftDelimiter] = Cell{Value: values[i] == 1, Type: CellTypeFormat}

		// right -> left, on the bottom side of the top-right search pattern
		field[8][size-1-i] = Cell{Value: values[14-i] == 1, Type: CellTypeFormat}

		// down -> up, on the right side of the bottom-left search pattern
		field[size-1-i][8] = Cell{Value: values[i] == 1, Type: CellTypeFormat}
	}

	// always dark module
	field[size-8][8] = Cell{Value: true, Type: CellTypeFormat}
}

func FillEmptyFormatBlock(field [][]Cell) {
	size := len(field)
	for i := 0; i < 8; i++ {
		shiftDelimiter := 0
		if i > 5 {
			shiftDelimiter = 1
		}
		// up -> down, on the right side of the top-left search pattern
		field[i+shiftDelimiter][8] = Cell{Value: false, Type: CellTypeFormat}

		// left -> right, on the bottom side of the top-left search pattern
		field[8][i+shiftDelimiter] = Cell{Value: false, Type: CellTypeFormat}

		// right -> left, on the bottom side of the top-right search pattern
		field[8][size-1-i] = Cell{Value: false, Type: CellTypeFormat}

		// down -> up, on the right side of the bottom-left search pattern
		field[size-1-i][8] = Cell{Value: false, Type: CellTypeFormat}
	}

	// always dark module
	field[size-8][8] = Cell{Value: true, Type: CellTypeFormat}
}

func FillDataBlock(field [][]Cell, data []byte) {
	size := len(field)
	bitIdx := 0
	pos := Position{X: size - 1, Y: size - 1, Size: size, Direction: -1}
	for byteIdx := 0; byteIdx < len(data); {
		if field[pos.Y][pos.X].Type == CellTypeData {
			field[pos.Y][pos.X] = Cell{Value: data[byteIdx]&(1<<uint(7-bitIdx)) != 0, Type: CellTypeData}
			bitIdx++
			if bitIdx > 7 {
				bitIdx = 0
				byteIdx++
			}
		}
		pos.Next()
	}

	fmt.Println(pos)
}

func GenerateField(data []byte, version int, errorCorrectionLevel ErrorCorrectionLevel) [][]Cell {
	size := GetSize(version)
	field := make([][]Cell, size)
	for i := range field {
		field[i] = make([]Cell, size)
	}

	FillSearchPattern(field)
	FillSyncPattern(field)
	FillAlignmentPattern(field, version)
	FillVersionBlock(field, version)
	FillEmptyFormatBlock(field)
	FillDataBlock(field, data)

	bestMask := DetermineBestMask(field, errorCorrectionLevel)
	FillFormatBlock(field, errorCorrectionLevel, bestMask)
	ApplyMask(field, bestMask)

	return field
}
