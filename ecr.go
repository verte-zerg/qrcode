package qrcode

func calculateEDCPoly(data []byte, codewords int) []byte {
	dataPoly := &polynomial{data}
	degree := codewords - len(data)
	dataPoly = dataPoly.IncreaseDegree(degree)
	edcPoly := dataPoly.Divide(generatorPolynomials[degree])
	return edcPoly.Coefficients
}

func getEDCData(data []byte, version int, errorLevel ErrorCorrectionLevel) []byte {
	var buf []byte

	var blocks []ecBlock
	if version < 0 {
		blocks = microErrorCorrectionBlocks[-version][errorLevel]
	} else {
		blocks = errorCorrectionBlocks[version][errorLevel]
	}

	var blocksData [][]byte
	dataIdx := 0
	for _, block := range blocks {
		for i := 0; i < block.Blocks; i++ {
			errorData := calculateEDCPoly(data[dataIdx:dataIdx+block.DataCodewords], block.TotalCodewords)
			blocksData = append(blocksData, errorData)
			dataIdx += block.DataCodewords
		}
	}

	maxBlockSize := 0

	for i := 0; i < len(blocksData); i++ {
		if len(blocksData[i]) > maxBlockSize {
			maxBlockSize = len(blocksData[i])
		}
	}

	for i := 0; i < maxBlockSize; i++ {
		for j := 0; j < len(blocksData); j++ {
			if len(blocksData[j]) > 0 {
				buf = append(buf, blocksData[j][0])
				blocksData[j] = blocksData[j][1:]
			}
		}
	}

	return buf
}

// Count of error correction code words for Micro QR Code version and error correction level
// Structure: [version][error correction level]
var microErrorCorrectionCodeWords = [5][4]int{
	{0, 0, 0, 0}, // added to shift the index by 1
	{2, 0, 0, 0},
	{5, 6, 0, 0},
	{6, 8, 0, 0},
	{8, 10, 14, 0},
}

// Count of error correction code words for each version and error correction level
// Structure: [version][error correction level]
var errorCorrectionCodeWords = [41][4]int{
	{0, 0, 0, 0}, // added to shift the index by 1
	{7, 10, 13, 17},
	{10, 16, 22, 28},
	{15, 26, 36, 44},
	{20, 36, 52, 64},
	{26, 48, 72, 88},
	{36, 64, 96, 112},
	{40, 72, 108, 130},
	{48, 88, 132, 156},
	{60, 110, 160, 192},
	{72, 130, 192, 224},
	{80, 150, 224, 264},
	{96, 176, 260, 308},
	{104, 198, 288, 352},
	{120, 216, 320, 384},
	{132, 240, 360, 432},
	{144, 280, 408, 480},
	{168, 308, 448, 532},
	{180, 338, 504, 588},
	{196, 364, 546, 650},
	{224, 416, 600, 700},
	{224, 442, 644, 750},
	{252, 476, 690, 816},
	{270, 504, 750, 900},
	{300, 560, 810, 960},
	{312, 588, 870, 1050},
	{336, 644, 952, 1110},
	{360, 700, 1020, 1200},
	{390, 728, 1050, 1260},
	{420, 784, 1140, 1350},
	{450, 812, 1200, 1440},
	{480, 868, 1290, 1530},
	{510, 924, 1350, 1620},
	{540, 980, 1440, 1710},
	{570, 1036, 1530, 1800},
	{570, 1064, 1590, 1890},
	{600, 1120, 1680, 1980},
	{630, 1204, 1770, 2100},
	{660, 1260, 1860, 2220},
	{720, 1316, 1950, 2310},
	{750, 1372, 2040, 2430},
}

type ecBlock struct {
	Blocks         int
	TotalCodewords int
	DataCodewords  int
	Ratio          int //TODO: rename
}

// ['give you up','let you down','run around and desert you'].map(x=>'Never gonna '+x)
// Never gonna give you upNever gonna let you downNever gonna run around and desert you

// Error correction blocks for Micro QR Code version and error correction level
// Structure: [version][error correction level][block]
//
// (5,3,0)Ь
// (10,5,1 )Ь
// (10,4,2)Ь (17,11,2)Ь (17,9,4)
// (24,16,3)Ь (24,14,5) (24,10,7)
var microErrorCorrectionBlocks = [5][4][]ecBlock{
	{{}, {}, {}, {}}, // added to shift the index by 1
	{{{1, 5, 3, 0}}, {}, {}, {}},
	{{{1, 10, 5, 1}}, {{1, 10, 4, 2}}, {}, {}},
	{{{1, 17, 11, 2}}, {{1, 17, 9, 4}}, {}, {}},
	{{{1, 24, 16, 3}}, {{1, 24, 14, 5}}, {{1, 24, 10, 7}}, {}},
}

// Error correction blocks for each version and error correction level
// Structure: [version][error correction level][block]
var errorCorrectionBlocks = [41][4][]ecBlock{
	{{}, {}, {}, {}}, // added to shift the index by 1
	{{{1, 26, 19, 2}}, {{1, 26, 16, 4}}, {{1, 26, 13, 6}}, {{1, 26, 9, 8}}},
	{{{1, 44, 34, 4}}, {{1, 44, 28, 8}}, {{1, 44, 22, 11}}, {{1, 44, 16, 14}}},
	{{{1, 70, 55, 7}}, {{1, 70, 44, 13}}, {{2, 35, 17, 9}}, {{2, 35, 13, 11}}},
	{{{1, 100, 80, 10}}, {{2, 50, 32, 9}}, {{2, 50, 24, 13}}, {{4, 25, 9, 8}}},
	{{{1, 134, 108, 13}}, {{2, 67, 43, 12}}, {{2, 33, 15, 9}, {2, 34, 16, 9}}, {{2, 33, 11, 11}, {2, 34, 12, 11}}},
	{{{2, 86, 68, 9}}, {{4, 43, 27, 8}}, {{4, 43, 19, 12}}, {{4, 43, 15, 14}}},
	{{{2, 98, 78, 10}}, {{4, 49, 31, 9}}, {{2, 32, 14, 9}, {4, 33, 15, 9}}, {{4, 39, 13, 13}, {1, 40, 14, 13}}},
	{{{2, 121, 97, 12}}, {{2, 60, 38, 11}, {2, 61, 39, 11}}, {{4, 40, 18, 11}, {2, 41, 19, 11}}, {{4, 40, 14, 13}, {2, 41, 15, 13}}},
	{{{2, 146, 116, 15}}, {{3, 58, 36, 11}, {2, 59, 37, 11}}, {{4, 36, 16, 10}, {4, 37, 17, 10}}, {{4, 36, 12, 12}, {4, 37, 13, 12}}},
	{{{2, 86, 68, 9}, {2, 87, 69, 9}}, {{4, 69, 43, 13}, {1, 70, 44, 13}}, {{6, 43, 19, 12}, {2, 44, 20, 12}}, {{6, 43, 15, 14}, {2, 44, 16, 14}}},
	{{{4, 101, 81, 10}}, {{1, 80, 50, 15}, {4, 81, 51, 15}}, {{4, 50, 22, 14}, {4, 51, 23, 14}}, {{3, 36, 12, 12}, {8, 37, 13, 12}}},
	{{{2, 116, 92, 12}, {2, 117, 93, 12}}, {{6, 58, 36, 11}, {2, 59, 37, 11}}, {{4, 46, 20, 13}, {6, 47, 21, 13}}, {{7, 42, 14, 14}, {4, 43, 15, 14}}},
	{{{4, 133, 107, 13}}, {{8, 59, 37, 11}, {1, 60, 38, 11}}, {{8, 44, 20, 12}, {4, 45, 21, 12}}, {{12, 33, 11, 11}, {4, 34, 12, 11}}},
	{{{3, 145, 115, 15}, {1, 146, 116, 15}}, {{4, 64, 40, 12}, {5, 65, 41, 12}}, {{11, 36, 16, 10}, {5, 37, 17, 10}}, {{11, 36, 12, 12}, {5, 37, 13, 12}}},
	{{{5, 109, 87, 11}, {1, 110, 88, 11}}, {{5, 65, 41, 12}, {5, 66, 42, 12}}, {{5, 54, 24, 15}, {7, 55, 25, 15}}, {{11, 36, 12, 12}, {7, 37, 13, 12}}},
	{{{5, 122, 98, 12}, {1, 123, 99, 12}}, {{7, 73, 45, 14}, {3, 74, 46, 14}}, {{15, 43, 19, 12}, {2, 44, 20, 12}}, {{3, 45, 15, 15}, {13, 46, 16, 15}}},
	{{{1, 135, 107, 14}, {5, 136, 108, 14}}, {{10, 74, 46, 14}, {1, 75, 47, 14}}, {{1, 50, 22, 14}, {15, 51, 23, 14}}, {{2, 42, 14, 14}, {17, 43, 15, 14}}},
	{{{5, 150, 120, 15}, {1, 151, 121, 15}}, {{9, 69, 43, 13}, {4, 70, 44, 13}}, {{17, 50, 22, 14}, {1, 51, 23, 14}}, {{2, 42, 14, 14}, {19, 43, 15, 14}}},
	{{{3, 141, 113, 14}, {4, 142, 114, 14}}, {{3, 70, 44, 13}, {11, 71, 45, 13}}, {{17, 47, 21, 13}, {4, 48, 22, 13}}, {{9, 39, 13, 13}, {16, 40, 14, 13}}},
	{{{3, 135, 107, 14}, {5, 136, 108, 14}}, {{3, 67, 41, 13}, {13, 68, 42, 13}}, {{15, 54, 24, 15}, {5, 55, 25, 15}}, {{15, 43, 15, 14}, {10, 44, 16, 14}}},
	{{{4, 144, 116, 14}, {4, 145, 117, 14}}, {{17, 68, 42, 13}}, {{17, 50, 22, 14}, {6, 51, 23, 14}}, {{19, 46, 16, 15}, {6, 47, 17, 15}}},
	{{{2, 139, 111, 14}, {7, 140, 112, 14}}, {{17, 74, 46, 14}}, {{7, 54, 24, 15}, {16, 55, 25, 15}}, {{34, 37, 13, 12}}},
	{{{4, 151, 121, 15}, {5, 152, 122, 15}}, {{4, 75, 47, 14}, {14, 76, 48, 14}}, {{11, 54, 24, 15}, {14, 55, 25, 15}}, {{16, 45, 15, 15}, {14, 46, 16, 15}}},
	{{{6, 147, 117, 15}, {4, 148, 118, 15}}, {{6, 73, 45, 14}, {14, 74, 46, 14}}, {{11, 54, 24, 15}, {16, 55, 25, 15}}, {{30, 46, 16, 15}, {2, 47, 17, 15}}},
	{{{8, 132, 106, 13}, {4, 133, 107, 13}}, {{8, 75, 47, 14}, {13, 76, 48, 14}}, {{7, 54, 24, 15}, {22, 55, 25, 15}}, {{22, 45, 15, 15}, {13, 46, 16, 15}}},
	{{{10, 142, 114, 14}, {2, 143, 115, 14}}, {{19, 74, 46, 14}, {4, 75, 47, 14}}, {{28, 50, 22, 14}, {6, 51, 23, 14}}, {{33, 46, 16, 15}, {4, 47, 17, 15}}},
	{{{8, 152, 122, 15}, {4, 153, 123, 15}}, {{22, 73, 45, 14}, {3, 74, 46, 14}}, {{8, 53, 23, 15}, {26, 54, 24, 15}}, {{12, 45, 15, 15}, {28, 46, 16, 15}}},
	{{{3, 147, 117, 15}, {10, 148, 118, 15}}, {{3, 73, 45, 14}, {23, 74, 46, 14}}, {{4, 54, 24, 15}, {31, 55, 25, 15}}, {{11, 45, 15, 15}, {31, 46, 16, 15}}},
	{{{7, 146, 116, 15}, {7, 147, 117, 15}}, {{21, 73, 45, 14}, {7, 74, 46, 14}}, {{1, 53, 23, 15}, {37, 54, 24, 15}}, {{19, 45, 15, 15}, {26, 46, 16, 15}}},
	{{{5, 145, 115, 15}, {10, 146, 116, 15}}, {{19, 75, 47, 14}, {10, 76, 48, 14}}, {{15, 54, 24, 15}, {25, 55, 25, 15}}, {{23, 45, 15, 15}, {25, 46, 16, 15}}},
	{{{13, 145, 115, 15}, {3, 146, 116, 15}}, {{2, 74, 46, 14}, {29, 75, 47, 14}}, {{42, 54, 24, 15}, {1, 55, 25, 15}}, {{23, 45, 15, 15}, {28, 46, 16, 15}}},
	{{{17, 145, 115, 15}}, {{10, 74, 46, 14}, {23, 75, 47, 14}}, {{10, 54, 24, 15}, {35, 55, 25, 15}}, {{19, 45, 15, 15}, {35, 46, 16, 15}}},
	{{{17, 145, 115, 15}, {1, 146, 116, 15}}, {{14, 74, 46, 14}, {21, 75, 47, 14}}, {{29, 54, 24, 15}, {19, 55, 25, 15}}, {{11, 45, 15, 15}, {46, 46, 16, 15}}},
	{{{13, 145, 115, 15}, {6, 146, 116, 15}}, {{14, 74, 46, 14}, {23, 75, 47, 14}}, {{44, 54, 24, 15}, {7, 55, 25, 15}}, {{59, 46, 16, 15}, {1, 47, 17, 15}}},
	{{{12, 151, 121, 15}, {7, 152, 122, 15}}, {{12, 75, 47, 14}, {26, 76, 48, 14}}, {{39, 54, 24, 15}, {14, 55, 25, 15}}, {{22, 45, 15, 15}, {41, 46, 16, 15}}},
	{{{6, 151, 121, 15}, {14, 152, 122, 15}}, {{6, 75, 47, 14}, {34, 76, 48, 14}}, {{46, 54, 24, 15}, {10, 55, 25, 15}}, {{2, 45, 15, 15}, {64, 46, 16, 15}}},
	{{{17, 152, 122, 15}, {4, 153, 123, 15}}, {{29, 74, 46, 14}, {14, 75, 47, 14}}, {{49, 54, 24, 15}, {10, 55, 25, 15}}, {{24, 45, 15, 15}, {46, 46, 16, 15}}},
	{{{4, 152, 122, 15}, {18, 153, 123, 15}}, {{13, 74, 46, 14}, {32, 75, 47, 14}}, {{48, 54, 24, 15}, {14, 55, 25, 15}}, {{42, 45, 15, 15}, {32, 46, 16, 15}}},
	{{{20, 147, 117, 15}, {4, 148, 118, 15}}, {{40, 75, 47, 14}, {7, 76, 48, 14}}, {{43, 54, 24, 15}, {22, 55, 25, 15}}, {{10, 45, 15, 15}, {67, 46, 16, 15}}},
	{{{19, 148, 118, 15}, {6, 149, 119, 15}}, {{18, 75, 47, 14}, {31, 76, 48, 14}}, {{34, 54, 24, 15}, {34, 55, 25, 15}}, {{20, 45, 15, 15}, {61, 46, 16, 15}}},
}
