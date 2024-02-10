package qrcode

import (
	"fmt"
	"math"

	"github.com/verte-zerg/qrcode/encode"
)

// Level of error correction
// Low - 7%
// Medium - 15%
// Quartile - 25%
// High - 30%
type ErrorCorrectionLevel int

const (
	ErrorCorrectionLevelLow ErrorCorrectionLevel = iota
	ErrorCorrectionLevelMedium
	ErrorCorrectionLevelQuartile
	ErrorCorrectionLevelHigh
)

var ErrContentTooLong = fmt.Errorf("content is too long")

// isVersionEnough checks if the given version can contain the data
func isVersionEnough(encodeBlocks []*encode.EncodeBlock, version int, dataSize int, ecl ErrorCorrectionLevel) (bool, error) {
	prefixBits := 0

	for _, block := range encodeBlocks {
		lengthBits, err := block.GetLengthBits(version)
		if err != nil {
			return false, fmt.Errorf("failed to get length bits: %w", err)
		}
		prefixBits += lengthBits + block.GetModeBits(version)
	}

	size := int(math.Ceil(float64(dataSize+prefixBits) / 8.0))

	var dataCodewords int
	if version < 0 {
		codewords := microCodewordsCount[-version]
		errorCodewords := microErrorCorrectionCodeWords[-version][ecl]

		if errorCodewords == 0 {
			return false, fmt.Errorf("unsupported error correction level: %v", ecl)
		}
		dataCodewords = codewords - errorCodewords
	} else {
		dataCodewords = codewordsCount[version] - errorCorrectionCodeWords[version][ecl]
	}

	return size <= dataCodewords, nil
}

// calculateMinVersion returns the minimum version for the given content, encoding mode, and error correction level.
// Alghorithm: iterate over versions from 1 to 40 (from M1 to M4 for MicroQR) and return the first version that can contain the content.
func calculateMinVersion(encodeBlocks []*encode.EncodeBlock, ecl ErrorCorrectionLevel, microQR bool) (int, error) {
	dataSize := 0
	for _, block := range encodeBlocks {
		blockSize, err := block.CalculateDataBitsCount()
		if err != nil {
			return 0, fmt.Errorf("failed to calculate data bits count: %w", err)
		}

		dataSize += blockSize
	}

	start, end, step := 1, 40, 1
	if microQR {
		start, end, step = -1, -4, -1
	}

	// For Normal QRs: 1 to 40
	// For Micro QRs: -1 to -4
	for version := start; version != end+step; version += step {
		ok, _ := isVersionEnough(encodeBlocks, version, dataSize, ecl)
		// if err != nil {
		// 	return 0, fmt.Errorf("failed to check version: %w", err)
		// }

		if ok {
			return version, nil
		}
	}

	return 0, ErrContentTooLong
}

// rearrangeDataBlocks rearranges the data blocks according to the QR code specification.
// When the QR code is split into data blocks, the data stream should be rearranged.
func rearrangeDataBlocks(data []byte, version int, errorLevel ErrorCorrectionLevel) []byte {
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
			blocksData = append(blocksData, data[dataIdx:dataIdx+block.DataCodewords])
			dataIdx += block.DataCodewords
		}
	}

	var buf []byte

	maxBlockSize := 0

	for i := 0; i < len(blocksData); i++ {
		if len(blocksData[i]) > maxBlockSize {
			maxBlockSize = len(blocksData[i])
		}
	}

	for i := 0; i < maxBlockSize; i++ {
		for j := 0; j < len(blocksData); j++ {
			if i < len(blocksData[j]) {
				buf = append(buf, blocksData[j][i])
			}
		}
	}

	return buf
}

// fillTerminator fills the data with terminator and padding bits based on the QR code specification.
func fillTerminator(data []byte, remainedBits int, version int, errorLevel ErrorCorrectionLevel) []byte {
	var availableCodewords int
	terminatorBits := 4

	// Micro QR Codes
	if version < 0 {
		availableCodewords = microCodewordsCount[-version] - microErrorCorrectionCodeWords[-version][errorLevel]
		terminatorBits = -version*2 + 1
	} else {
		availableCodewords = codewordsCount[version] - errorCorrectionCodeWords[version][errorLevel]
	}

	if remainedBits < terminatorBits && len(data) < availableCodewords {
		data = append(data, 0)
	}

	// TODO: refactor
	if remainedBits == 0 && terminatorBits == 9 && len(data) < availableCodewords {
		data = append(data, 0)
	}

	var terminator byte = 0b11101100
	hasEmptyCodewords := false
	for len(data) < availableCodewords {
		hasEmptyCodewords = true
		data = append(data, terminator)
		if terminator == 0b11101100 {
			terminator = 0b00010001
		} else {
			terminator = 0b11101100
		}
	}

	if hasEmptyCodewords && (version == -1 || version == -3) {
		data[len(data)-1] = 0
	}

	return data
}

// getBytesData returns the byte array for the given content, encoding mode, error correction level, and version.
func getBytesData(blocks []*encode.EncodeBlock, errorLevel ErrorCorrectionLevel, version int) ([]byte, error) {
	allBits := 0

	queue := make(chan encode.ValueBlock, 100)
	result := make(chan []byte)

	go encode.GenerateData(queue, result)

	for _, block := range blocks {
		blockBits, err := block.Encode(version, queue)
		if err != nil {
			return nil, fmt.Errorf("failed to encode data: %w", err)
		}

		allBits += blockBits
	}

	close(queue)
	data := <-result

	// add terminator
	remainedBits := len(data)*8 - allBits
	data = fillTerminator(data, remainedBits, version, errorLevel)

	errorData := getEDCData(data, version, errorLevel)
	data = rearrangeDataBlocks(data, version, errorLevel)
	data = append(data, errorData...)

	return data, nil
}

var (
	// Number of codewords for Micro QR Code version
	microCodewordsCount = [5]int{
		0, // added for shift start index to 1
		5, 10, 17, 24,
	}

	// Number of codewords for each version
	codewordsCount = [41]int{
		0, // added for shift start index to 1
		26, 44, 70, 100, 134, 172, 196, 242, 292, 346,
		404, 466, 532, 581, 655, 733, 815, 901, 991, 1085,
		1156, 1258, 1364, 1474, 1588, 1706, 1828, 1921, 2051, 2185,
		2323, 2465, 2611, 2761, 2876, 3034, 3196, 3362, 3532, 3706,
	}
)
