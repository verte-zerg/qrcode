package qrcode

var (
	maskFuncs = []func(row, col int) bool{
		func(row, col int) bool { return (row+col)%2 == 0 },
		func(row, col int) bool { return row%2 == 0 },
		func(row, col int) bool { return col%3 == 0 },
		func(row, col int) bool { return (row+col)%3 == 0 },
		func(row, col int) bool { return (row/2+col/3)%2 == 0 },
		func(row, col int) bool { return (row*col)%2+(row*col)%3 == 0 },
		func(row, col int) bool { return ((row*col)%2+(row*col)%3)%2 == 0 },
		func(row, col int) bool { return ((row*col)%3+(row+col)%2)%2 == 0 },
	}

	microToNormalMask = []struct {
		microMask  int
		normalMask int
	}{
		{0, 1},
		{1, 4},
		{2, 6},
		{3, 7},
	}
)

func applyMask(data [][]Cell, maskType int) {
	for idx, row := range data {
		for jdx, cell := range row {
			if cell.Type != CellTypeData {
				continue
			}

			data[idx][jdx].Value = maskFuncs[maskType](idx, jdx) != cell.Value
		}
	}
}

func calculatePenaltyRule1(data [][]Cell) int {
	penalty := 0

	rows := len(data)
	if rows == 0 {
		return 0
	}
	cols := len(data[0])

	// Horizontal
	length := 0
	for _, row := range data {
		value := false
		for _, cell := range row {
			if length == 0 {
				value = cell.Value
				length = 1
				continue
			}

			if cell.Value == value {
				length++
			} else {
				if length >= 5 {
					penalty += length - 2
				}
				length = 1
				value = cell.Value
			}
		}

		if length >= 5 {
			penalty += length - 2
		}
		length = 0
	}

	// Vertical
	for jdx := 0; jdx < cols; jdx++ {
		value := false
		for idx := 0; idx < rows; idx++ {
			if length == 0 {
				value = data[idx][jdx].Value
				length = 1
				continue
			}

			if data[idx][jdx].Value == value {
				length++
			} else {
				if length >= 5 {
					penalty += length - 2
				}
				length = 1
				value = data[idx][jdx].Value
			}
		}

		if length >= 5 {
			penalty += length - 2
		}
		length = 0
	}

	return penalty
}

func calculatePenaltyRule2(data [][]Cell) int {
	penalty := 0

	rows := len(data)
	if rows == 0 {
		return 0
	}
	cols := len(data[0])

	for idx := 0; idx < rows-1; idx++ {
		for jdx := 0; jdx < cols-1; jdx++ {
			val := data[idx][jdx].Value
			if val == data[idx][jdx+1].Value && val == data[idx+1][jdx].Value && val == data[idx+1][jdx+1].Value {
				penalty += 3
			}
		}
	}

	return penalty
}

func calculatePenaltyRule3(data [][]Cell) int {
	penalty := 0

	rows := len(data)
	if rows == 0 {
		return 0
	}
	cols := len(data[0])

	// Horizontal
	for idx := 0; idx < rows; idx++ {
		for jdx := 0; jdx < cols-10; jdx++ {
			// dark-light-dark-dark-dark-light-dark-light-light-light-light
			if data[idx][jdx].Value && !data[idx][jdx+1].Value && data[idx][jdx+2].Value && data[idx][jdx+3].Value && data[idx][jdx+4].Value && !data[idx][jdx+5].Value && data[idx][jdx+6].Value && !data[idx][jdx+7].Value && !data[idx][jdx+8].Value && !data[idx][jdx+9].Value && !data[idx][jdx+10].Value {
				penalty += 40
			}

			// reverse order: light-light-light-light-dark-light-dark-dark-dark-light-dark
			if !data[idx][jdx].Value && !data[idx][jdx+1].Value && !data[idx][jdx+2].Value && !data[idx][jdx+3].Value && data[idx][jdx+4].Value && !data[idx][jdx+5].Value && data[idx][jdx+6].Value && data[idx][jdx+7].Value && data[idx][jdx+8].Value && !data[idx][jdx+9].Value && data[idx][jdx+10].Value {
				penalty += 40
			}
		}
	}

	// Vertical
	for jdx := 0; jdx < cols; jdx++ {
		for idx := 0; idx < rows-10; idx++ {
			// dark-light-dark-dark-dark-light-dark-light-light-light-light
			if data[idx][jdx].Value && !data[idx+1][jdx].Value && data[idx+2][jdx].Value && data[idx+3][jdx].Value && data[idx+4][jdx].Value && !data[idx+5][jdx].Value && data[idx+6][jdx].Value && !data[idx+7][jdx].Value && !data[idx+8][jdx].Value && !data[idx+9][jdx].Value && !data[idx+10][jdx].Value {
				penalty += 40
			}

			// reverse order: light-light-light-light-dark-light-dark-dark-dark-light-dark
			if !data[idx][jdx].Value && !data[idx+1][jdx].Value && !data[idx+2][jdx].Value && !data[idx+3][jdx].Value && data[idx+4][jdx].Value && !data[idx+5][jdx].Value && data[idx+6][jdx].Value && data[idx+7][jdx].Value && data[idx+8][jdx].Value && !data[idx+9][jdx].Value && data[idx+10][jdx].Value {
				penalty += 40
			}
		}
	}

	return penalty
}

func calculatePenaltyRule4(data [][]Cell) int {
	rows := len(data)
	if rows == 0 {
		return 0
	}
	cols := len(data[0])

	darkCount := 0
	for _, row := range data {
		for _, cell := range row {
			if cell.Value {
				darkCount++
			}
		}
	}

	percentage := int(float64(darkCount) / float64(rows*cols) * 100)

	if percentage == 50 {
		return 0
	}

	// 50% - 5%
	if percentage > 50 {
		return ((percentage/5)*5 - 50) * 2
	} else {
		return ((50 - percentage) / 5) * 5 * 2
	}
}

func calculatePenalty(data [][]Cell) int {
	rule1 := calculatePenaltyRule1(data)
	rule2 := calculatePenaltyRule2(data)
	rule3 := calculatePenaltyRule3(data)
	rule4 := calculatePenaltyRule4(data)

	return rule1 + rule2 + rule3 + rule4
}

func determineBestMask(data [][]Cell, errorCorrectionLevel ErrorCorrectionLevel) int {
	minPenalty := 1<<31 - 1
	bestMask := 0
	for maskType := 0; maskType < 8; maskType++ {
		fillFormatBlock(data, errorCorrectionLevel, maskType)
		applyMask(data, maskType)
		penalty := calculatePenalty(data)
		if penalty < minPenalty {
			minPenalty = penalty
			bestMask = maskType
		}
		applyMask(data, maskType)
	}

	return bestMask
}

func determineBestMaskMicro(data [][]Cell, version int, errorCorrectionLevel ErrorCorrectionLevel) int {
	minPenalty := 1<<31 - 1
	bestMask := 0
	for _, maskMap := range microToNormalMask {
		fillFormatBlockMicro(data, version, errorCorrectionLevel, maskMap.microMask)
		applyMask(data, maskMap.normalMask)
		penalty := calculatePenalty(data)
		if penalty < minPenalty {
			minPenalty = penalty
			bestMask = maskMap.microMask
		}
		applyMask(data, maskMap.normalMask)
	}

	return bestMask
}
