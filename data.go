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

// CalculateMinVersion returns the minimum version for the given content, encoding mode, and error correction level.
// Alghorithm: iterate over versions from 1 to 40 and return the first version that can contain the content.
// Content is the string to encode.
// Mode is one of the EncodingMode constants.
// Error correction level is one of the ErrorCorrectionLevel constants.
func CalculateMinVersion(encodeBlocks []*encode.EncodeBlock, ecl ErrorCorrectionLevel) (int, error) {
	dataSize := 0
	for _, block := range encodeBlocks {
		blockSize, err := block.CalculateDataBitsCount()
		if err != nil {
			return 0, fmt.Errorf("failed to calculate data bits count: %w", err)
		}

		dataSize += blockSize
	}

	for version := 1; version <= 40; version++ {
		prefixBits := 0

		for _, block := range encodeBlocks {
			lengthBits, err := block.GetLengthBits(version)
			if err != nil {
				return 0, fmt.Errorf("failed to get length bits: %w", err)
			}
			prefixBits += lengthBits + block.GetModeBits()
		}

		size := int(math.Ceil(float64(dataSize+prefixBits) / 8.0))
		dataCodewords := CodewordsCount[version] - ErrorCorrectionCodeWords[version][ecl]
		if size <= dataCodewords {
			return version, nil
		}
	}

	return 0, ErrContentTooLong
}

// RearrangeDataBlocks rearranges the data blocks according to the QR code specification.
// When the QR code is split into data blocks, the data stream should be rearranged.
func RearrangeDataBlocks(data []byte, version int, errorLevel ErrorCorrectionLevel) []byte {
	blocks := ErrorCorrectionBlocks[version][errorLevel]
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

// GetBytesData returns the byte array for the given content, encoding mode, error correction level, and version.
func GetBytesData(blocks []*encode.EncodeBlock, errorLevel ErrorCorrectionLevel, version int) ([]byte, error) {
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

	availableCodewords := CodewordsCount[version] - ErrorCorrectionCodeWords[version][errorLevel]
	if remainedBits < 4 && len(data) < availableCodewords {
		data = append(data, 0)
	}

	var terminator byte = 0b11101100
	for len(data) < availableCodewords {
		data = append(data, terminator)
		if terminator == 0b11101100 {
			terminator = 0b00010001
		} else {
			terminator = 0b11101100
		}
	}

	errorData := GetEDCData(data, version, errorLevel)
	data = RearrangeDataBlocks(data, version, errorLevel)
	data = append(data, errorData...)

	return data, nil

}

var (
	// Number of codewords for each version
	CodewordsCount = [41]int{
		0, // added for shift start index to 1
		26, 44, 70, 100, 134, 172, 196, 242, 292, 346,
		404, 466, 532, 581, 655, 733, 815, 901, 991, 1085,
		1156, 1258, 1364, 1474, 1588, 1706, 1828, 1921, 2051, 2185,
		2323, 2465, 2611, 2761, 2876, 3034, 3196, 3362, 3532, 3706,
	}
)
