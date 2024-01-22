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
)

func ApplyMask(data [][]Cell, maskType int) {
	for idx, row := range data {
		for jdx, cell := range row {
			if cell.Type != CellTypeData {
				continue
			}

			data[idx][jdx].Value = maskFuncs[maskType](idx, jdx) != cell.Value
		}
	}
}

func CalculatePenaltyRule1(data [][]Cell) int {
	penalty := 0

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
	for jdx := 0; jdx < len(data); jdx++ {
		value := false
		for idx := 0; idx < len(data); idx++ {
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

func CalculatePenaltyRule2(data [][]Cell) int {
	penalty := 0

	size := len(data)

	for idx := 0; idx < size-1; idx++ {
		for jdx := 0; jdx < size-1; jdx++ {
			val := data[idx][jdx].Value
			if val == data[idx][jdx+1].Value && val == data[idx+1][jdx].Value && val == data[idx+1][jdx+1].Value {
				penalty += 3
			}
		}
	}

	return penalty
}

func CalculatePenaltyRule3(data [][]Cell) int {
	penalty := 0

	size := len(data)

	// Horizontal
	for idx := 0; idx < size; idx++ {
		for jdx := 0; jdx < size-10; jdx++ {
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
	for jdx := 0; jdx < size; jdx++ {
		for idx := 0; idx < size-10; idx++ {
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

func CalculatePenaltyRule4(data [][]Cell) int {
	darkCount := 0
	for _, row := range data {
		for _, cell := range row {
			if cell.Value {
				darkCount++
			}
		}
	}

	percentage := int(float64(darkCount) / float64(len(data)*len(data)) * 100)

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

func CalculatePenalty(data [][]Cell) int {
	rule1 := CalculatePenaltyRule1(data)
	rule2 := CalculatePenaltyRule2(data)
	rule3 := CalculatePenaltyRule3(data)
	rule4 := CalculatePenaltyRule4(data)

	return rule1 + rule2 + rule3 + rule4
}

func DetermineBestMask(data [][]Cell, errorCorrectionLevel ErrorCorrectionLevel) int {
	minPenalty := 1<<31 - 1
	bestMask := 0
	for maskType := 0; maskType < 8; maskType++ {
		copyData := make([][]Cell, len(data))
		for idx, row := range data {
			copyData[idx] = make([]Cell, len(row))
			copy(copyData[idx], row)
		}

		FillFormatBlock(copyData, errorCorrectionLevel, maskType)
		ApplyMask(copyData, maskType)
		penalty := CalculatePenalty(copyData)
		if penalty < minPenalty {
			minPenalty = penalty
			bestMask = maskType
		}
	}

	return bestMask
}
