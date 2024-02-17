package qrcode

type CellType int

// Cell type represents the type of a cell in the QR code matrix.
const (
	CellTypeData CellType = iota
	CellTypeFormat
	CellTypeVersion
	CellTypeAlignmentPattern
	CellTypeSearchPattern
	CellTypeSyncPattern
	CellTypeDelimiter
)

var (
	// searchPattern is a 7x7 search pattern for the QR code.
	searchPattern = [7][7]bool{
		{true, true, true, true, true, true, true},
		{true, false, false, false, false, false, true},
		{true, false, true, true, true, false, true},
		{true, false, true, true, true, false, true},
		{true, false, true, true, true, false, true},
		{true, false, false, false, false, false, true},
		{true, true, true, true, true, true, true},
	}

	// alignmentPattern is a 5x5 alignment pattern for the QR code.
	alignmentPattern = [5][5]bool{
		{true, true, true, true, true},
		{true, false, false, false, true},
		{true, false, true, false, true},
		{true, false, false, false, true},
		{true, true, true, true, true},
	}

	// alignmentPositions is a table of alignment pattern positions for each version.
	alignmentPositions = [41][]int{
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

	// versionValues is a table of version information for each version.
	// The zero index is the version number 7, as the version block starts from version 7.
	versionValues = [34][18]byte{
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

	// microFormatValues is a table of format information for each version, error correction level, and mask pattern.
	// Structure: [version][error correction level][mask pattern][15 bits]
	microFormatValues = [5][4][4][15]byte{
		{}, // added for shift versions, as first version has index 1
		{ // Version M1
			{
				{1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1},
				{1, 0, 0, 0, 0, 0, 1, 0, 1, 1, 1, 0, 0, 1, 0},
				{1, 0, 0, 1, 1, 1, 0, 0, 0, 1, 0, 1, 0, 1, 1},
				{1, 0, 0, 1, 0, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0},
			},
		},
		{ // Version M2
			{ // Error correction level L
				{1, 0, 1, 0, 1, 0, 1, 1, 0, 1, 0, 1, 1, 1, 0},
				{1, 0, 1, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0, 1},
				{1, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0},
				{1, 0, 1, 1, 0, 1, 0, 1, 1, 1, 1, 0, 1, 1, 1},
			},
			{ // Error correction level M
				{1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 0, 0, 1, 1},
				{1, 1, 0, 0, 0, 1, 0, 1, 0, 1, 0, 0, 1, 0, 0},
				{1, 1, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 0, 1},
				{1, 1, 0, 1, 0, 0, 0, 1, 1, 0, 0, 1, 0, 1, 0},
			},
		},
		{ // Version M3
			{ // Error correction level L
				{1, 1, 1, 0, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 0},
				{1, 1, 1, 0, 0, 1, 1, 0, 1, 0, 0, 1, 1, 1, 1},
				{1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 1, 0, 1, 1, 0},
				{1, 1, 1, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 1},
			},
			{ // Error correction level M
				{0, 0, 0, 0, 1, 1, 0, 1, 1, 0, 1, 1, 1, 1, 0},
				{0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 1, 0, 0, 1},
				{0, 0, 0, 1, 1, 0, 0, 1, 0, 1, 1, 0, 0, 0, 0},
				{0, 0, 0, 1, 0, 0, 1, 1, 0, 0, 0, 0, 1, 1, 1},
			},
		},
		{ // Version M4
			{ // Error correction level L
				{0, 0, 1, 0, 1, 1, 1, 0, 0, 1, 1, 0, 1, 0, 1},
				{0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0},
				{0, 0, 1, 1, 1, 0, 1, 0, 1, 0, 1, 1, 0, 1, 1},
				{0, 0, 1, 1, 0, 0, 0, 0, 1, 1, 0, 1, 1, 0, 0},
			},
			{ // Error correction level M
				{0, 1, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0},
				{0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1},
				{0, 1, 0, 1, 1, 1, 1, 0, 1, 1, 0, 0, 1, 1, 0},
				{0, 1, 0, 1, 0, 1, 0, 0, 1, 0, 1, 0, 0, 0, 1},
			},
			{ // Error correction level Q
				{0, 1, 1, 0, 1, 0, 0, 1, 1, 1, 0, 0, 0, 1, 1},
				{0, 1, 1, 0, 0, 0, 1, 1, 1, 0, 1, 0, 1, 0, 0},
				{0, 1, 1, 1, 1, 1, 0, 1, 0, 0, 0, 1, 1, 0, 1},
				{0, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 0, 1, 0},
			},
		},
	}

	// formatValues is a table of format information for each error correction level and mask pattern.
	// Structure: [error correction level][mask pattern][15 bits]
	formatValues = [4][8][15]byte{
		{ // Error correction level L
			{1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 0, 0, 1, 0, 0},
			{1, 1, 1, 0, 0, 1, 0, 1, 1, 1, 1, 0, 0, 1, 1},
			{1, 1, 1, 1, 1, 0, 1, 1, 0, 1, 0, 1, 0, 1, 0},
			{1, 1, 1, 1, 0, 0, 0, 1, 0, 0, 1, 1, 1, 0, 1},
			{1, 1, 0, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 1, 1},
			{1, 1, 0, 0, 0, 1, 1, 0, 0, 0, 1, 1, 0, 0, 0},
			{1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1},
			{1, 1, 0, 1, 0, 0, 1, 0, 1, 1, 1, 0, 1, 1, 0},
		},
		{ // Error correction level M
			{1, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0},
			{1, 0, 1, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 1},
			{1, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 0, 0},
			{1, 0, 1, 1, 0, 1, 1, 0, 1, 0, 0, 1, 0, 1, 1},
			{1, 0, 0, 0, 1, 0, 1, 1, 1, 1, 1, 1, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 1, 1, 1, 0},
			{1, 0, 0, 1, 1, 1, 1, 1, 0, 0, 1, 0, 1, 1, 1},
			{1, 0, 0, 1, 0, 1, 0, 1, 0, 1, 0, 0, 0, 0, 0},
		},
		{ // Error correction level Q
			{0, 1, 1, 0, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1},
			{0, 1, 1, 0, 0, 0, 0, 0, 1, 1, 0, 1, 0, 0, 0},
			{0, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 0, 0, 0, 1},
			{0, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0},
			{0, 1, 0, 0, 1, 0, 0, 1, 0, 1, 1, 0, 1, 0, 0},
			{0, 1, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 1, 1},
			{0, 1, 0, 1, 1, 1, 0, 1, 1, 0, 1, 1, 0, 1, 0},
			{0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 0, 1},
		},
		{ // Error correction level H
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

// Cell represents a cell in the QR code matrix with a value and a type.
type Cell struct {
	Value bool
	Type  CellType
}

// position is a iterator, which is used to fill the QR code matrix in a spiral pattern.
type position struct {
	X, Y        int
	Size        int
	Direction   int
	Micro       bool
	evenReverse bool
}

func (p *position) Next() {
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

	if p.X == 6 && !p.Micro {
		p.X--
		p.evenReverse = true
	}
}

// getSize returns the square size of the QR code matrix for the given version.
func getSize(version int) int {
	if version < 0 {
		return 9 + (-version)*2
	}
	return 17 + 4*version
}

// fillSearchPattern fills the search patterns in the QR code matrix.
func fillSearchPattern(version int, field [][]Cell) {
	size := len(field)
	for i := 0; i < 7; i++ {
		for j := 0; j < 7; j++ {
			value := searchPattern[i][j]
			field[i][j] = Cell{Value: value, Type: CellTypeSearchPattern}

			if version > 0 {
				field[i][size-1-j] = Cell{Value: value, Type: CellTypeSearchPattern}
				field[size-1-i][j] = Cell{Value: value, Type: CellTypeSearchPattern}
			}
		}
	}

	for i := 0; i < 8; i++ {
		// top-left
		field[i][7] = Cell{Value: false, Type: CellTypeSearchPattern}
		field[7][i] = Cell{Value: false, Type: CellTypeSearchPattern}

		if version > 0 {
			// top-right
			field[i][size-8] = Cell{Value: false, Type: CellTypeSearchPattern}
			field[7][size-1-i] = Cell{Value: false, Type: CellTypeSearchPattern}

			// bottom-left
			field[size-8][i] = Cell{Value: false, Type: CellTypeSearchPattern}
			field[size-1-i][7] = Cell{Value: false, Type: CellTypeSearchPattern}
		}
	}
}

// fillSyncPattern fills the sync patterns in the QR code matrix.
func fillSyncPattern(version int, field [][]Cell) {
	size := len(field)
	value := true

	row := 6
	col := 6
	emptySpace := 8
	if version < 0 {
		row = 0
		col = 0
		emptySpace = 0
	}

	for i := 8; i < size-emptySpace; i++ {
		field[row][i] = Cell{Value: value, Type: CellTypeSyncPattern}
		field[i][col] = Cell{Value: value, Type: CellTypeSyncPattern}
		value = !value
	}
}

// fillAlignmentPattern fills the alignment patterns in the QR code matrix.
func fillAlignmentPattern(field [][]Cell, version int) {
	positions := alignmentPositions[version]
	for idx, i := range positions {
		for jdx, j := range positions {
			if (idx == 0 && jdx == 0) || (idx == 0 && jdx == len(positions)-1) || (idx == len(positions)-1 && jdx == 0) {
				continue
			}
			for k := 0; k < 5; k++ {
				for l := 0; l < 5; l++ {
					field[i+k-2][j+l-2] = Cell{Value: alignmentPattern[k][l], Type: CellTypeAlignmentPattern}
				}
			}
		}
	}
}

// fillVersionBlock fills the version block in the QR code matrix.
func fillVersionBlock(field [][]Cell, version int) {
	size := len(field)
	if version < 7 {
		return
	}
	values := versionValues[version-7]
	for i := 0; i < 6; i++ {
		for j := 0; j < 3; j++ {
			field[size-11+j][i] = Cell{Value: values[i*3+j] == 1, Type: CellTypeVersion}
			field[5-i][size-9-j] = Cell{Value: values[i*3+j] == 1, Type: CellTypeVersion}
		}
	}
}

// fillFormatBlock fills the format block in the QR code matrix.
func fillFormatBlock(field [][]Cell, errorCorrectionLevel ErrorCorrectionLevel, maskPattern int) {
	size := len(field)
	values := formatValues[errorCorrectionLevel][maskPattern]
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

// fillEmptyFormatBlock set type CellTypeFormat to the format block cells.
// This is used for avoid fill the data cells over the format block cells.
// The format block will be filled later (after the data cells are filled and the best mask is determined).
func fillEmptyFormatBlock(field [][]Cell) {
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

// fillFormatBlockMicro fills the format block in the QR code matrix for micro QR codes.
func fillFormatBlockMicro(field [][]Cell, version int, errorCorrectionLevel ErrorCorrectionLevel, maskPattern int) {
	values := microFormatValues[-version][errorCorrectionLevel][maskPattern]
	for i := 0; i < 8; i++ {
		// up -> down, on the right side of the top-left search pattern
		field[i+1][8] = Cell{Value: values[14-i] == 1, Type: CellTypeFormat}

		// left -> right, on the bottom side of the top-left search pattern
		field[8][i+1] = Cell{Value: values[i] == 1, Type: CellTypeFormat}
	}
}

// fillEmptyFormatBlockMicro set type CellTypeFormat to the format block cells for micro QR codes.
// This is used for avoid fill the data cells over the format block cells.
// The format block will be filled later (after the data cells are filled and the best mask is determined).
func fillEmptyFormatBlockMicro(field [][]Cell) {
	for i := 0; i < 8; i++ {
		// up -> down, on the right side of the top-left search pattern
		field[i+1][8] = Cell{Value: false, Type: CellTypeFormat}

		// left -> right, on the bottom side of the top-left search pattern
		field[8][i+1] = Cell{Value: false, Type: CellTypeFormat}
	}
}

// fillDataBlock fills the data block in the QR code matrix.
func fillDataBlock(field [][]Cell, data []byte) {
	size := len(field)
	bitIdx := 0
	pos := position{X: size - 1, Y: size - 1, Size: size, Direction: -1}
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
}

// fillDataBlockMicro fills the data block in the QR code matrix for micro QR codes.
func fillDataBlockMicro(field [][]Cell, data []byte, version int, errorLevel ErrorCorrectionLevel) {
	size := len(field)
	bitIdx := 0
	pos := position{X: size - 1, Y: size - 1, Size: size, Direction: -1, Micro: true}
	for byteIdx := 0; byteIdx < len(data); {
		if field[pos.Y][pos.X].Type == CellTypeData {
			field[pos.Y][pos.X] = Cell{Value: data[byteIdx]&(1<<uint(7-bitIdx)) != 0, Type: CellTypeData}
			bitIdx++

			// For M1, M3L and M3M, the last data byte has 4 bits
			if bitIdx == 4 {
				if version == -1 && errorLevel == ErrorCorrectionLevelLow && byteIdx == 2 {
					bitIdx = 8
				} else if version == -3 && errorLevel == ErrorCorrectionLevelLow && byteIdx == 10 {
					bitIdx = 8
				} else if version == -3 && errorLevel == ErrorCorrectionLevelMedium && byteIdx == 8 {
					bitIdx = 8
				}
			}

			if bitIdx > 7 {
				bitIdx = 0
				byteIdx++
			}
		}
		pos.Next()
	}
}

// generateField creates a QR code matrix based on the given data, version and error correction level.
func generateField(data []byte, version int, errorCorrectionLevel ErrorCorrectionLevel) [][]Cell {
	if version < 0 {
		return generateMicroQRField(data, version, errorCorrectionLevel)
	}
	size := getSize(version)
	field := make([][]Cell, size)
	for i := range field {
		field[i] = make([]Cell, size)
	}

	fillSearchPattern(version, field)
	fillSyncPattern(version, field)
	fillAlignmentPattern(field, version)
	fillVersionBlock(field, version)
	fillEmptyFormatBlock(field)
	fillDataBlock(field, data)

	bestMask := determineBestMask(field, errorCorrectionLevel)
	fillFormatBlock(field, errorCorrectionLevel, bestMask)
	applyMask(field, bestMask)

	return field
}

// generateMicroQRField creates a micro QR code matrix based on the given data, version and error correction level.
func generateMicroQRField(data []byte, version int, errorCorrectionLevel ErrorCorrectionLevel) [][]Cell {
	size := getSize(version)
	field := make([][]Cell, size)
	for i := range field {
		field[i] = make([]Cell, size)
	}

	fillSearchPattern(version, field)
	fillSyncPattern(version, field)
	fillEmptyFormatBlockMicro(field)
	fillDataBlockMicro(field, data, version, errorCorrectionLevel)

	bestMask := determineBestMaskMicro(field, version, errorCorrectionLevel)
	fillFormatBlockMicro(field, version, errorCorrectionLevel, bestMask)
	applyMask(field, microToNormalMask[bestMask].normalMask)

	return field
}
