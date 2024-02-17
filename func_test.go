package qrcode

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/verte-zerg/qrcode/encode"
)

func ValidateContent(content string, filename string, microQR bool) error {
	cmd := exec.Command("./Reader", filename)

	stdout, err := cmd.Output()
	if err != nil {
		return err
	}

	if cmd.ProcessState.ExitCode() != 0 {
		return errors.New("exit code is not 0")
	}

	output := string(stdout)

	if output != content {
		if strings.Contains(output, content) {
			return nil
		}
		return errors.New("content is not equal: " + output + " != " + content)
	}

	return nil
}

func ValidateContentRaw(data [][]Cell, content string) error {
	cmd := exec.Command("./RawReader")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	size := len(data)

	buf := make([]byte, size*(size+1)+1)

	trueValue := byte('X')
	falseValue := byte(' ')
	newLineValue := byte('\n')
	sharpValue := byte('#')

	i := 0
	for _, row := range data {
		for _, cell := range row {
			if cell.Value {
				buf[i] = trueValue
			} else {
				buf[i] = falseValue
			}
			i++
		}
		buf[i] = newLineValue
		i++
	}

	buf[i] = sharpValue

	stdin.Write(buf)

	stdin.Close()

	stdout, err := cmd.Output()
	if err != nil {
		return err
	}

	if cmd.ProcessState.ExitCode() != 0 {
		return errors.New("exit code is not 0")
	}

	output := string(stdout)

	if output != content {
		if strings.Contains(output, content) {
			return nil
		}
		return errors.New("content is not equal: " + output + " != " + content)
	}

	return nil
}

func TestPlotBase(t *testing.T) {
	content := "849"
	qr, err := Create(content, &QRCodeOptions{
		Mode:       encode.EncodingModeNumeric,
		ErrorLevel: ErrorCorrectionLevelMedium,
	})

	if err != nil {
		t.Error(err)
	}

	file, err := os.Create("test_new.png")
	if err != nil {
		t.Error(err)
	}

	defer file.Close()

	if err := qr.Plot(file); err != nil {
		t.Error(err)
	}

	err = ValidateContent(content, "test_new.png", false)
	if err != nil {
		t.Error(err)
	}
}

func TestPlotMicro(t *testing.T) {
	content := "8492314"
	qr, err := Create(
		content,
		&QRCodeOptions{
			ErrorLevel: ErrorCorrectionLevelLow,
			MicroQR:    true,
		},
	)

	if err != nil {
		t.Error(err)
	}

	file, err := os.Create("test_micro.png")
	if err != nil {
		t.Error(err)
	}

	defer file.Close()

	if err := qr.Plot(file); err != nil {
		t.Error(err)
	}

	err = ValidateContent(content, "test_micro.png", false)
	if err != nil {
		t.Error(err)
	}
}

func TestPlotUTF8(t *testing.T) {
	content := "asdfklj;ååß∂∆…¬å´œ¨®ˆπø∑´∆˚çå˜ß¬˚…¬√“‘ˆœ‘ø´®“π\\"
	qr, err := CreateMultiMode([]*encode.EncodeBlock{{
		Mode:             encode.EncodingModeECI,
		Data:             content,
		SubMode:          encode.EncodingModeByte,
		AssignmentNumber: 26,
	},
	}, &QRCodeOptionsMultiMode{
		ErrorLevel: ErrorCorrectionLevelMedium,
	},
	)

	if err != nil {
		t.Error(err)
	}

	file, err := os.Create("test_utf8.png")
	if err != nil {
		t.Error(err)
	}

	defer file.Close()

	if err := qr.Plot(file); err != nil {
		t.Error(err)
	}

	err = ValidateContent(content, "test_utf8.png", false)
	if err != nil {
		t.Error(err)
	}
}

func TestPlotCyrilic(t *testing.T) {
	content := "АВГДЕ"
	qr, err := CreateMultiMode([]*encode.EncodeBlock{{
		Mode:             encode.EncodingModeECI,
		Data:             content,
		SubMode:          encode.EncodingModeByte,
		AssignmentNumber: 7,
	},
	}, &QRCodeOptionsMultiMode{
		ErrorLevel: ErrorCorrectionLevelMedium,
	},
	)

	if err != nil {
		t.Error(err)
	}

	file, err := os.Create("test_cyrilic.png")
	if err != nil {
		t.Error(err)
	}

	defer file.Close()

	if err := qr.Plot(file); err != nil {
		t.Error(err)
	}

	err = ValidateContent(content, "test_cyrilic.png", false)
	if err != nil {
		t.Error(err)
	}
}

func TestPlotLong(t *testing.T) {
	// generate random content string with length 30
	var content string // = "https://www.qrcode.com/"
	// content := "['give you up','let you down','run around and desert you'].map(x=>'Never gonna '+x)"
	for i := 0; i < 1600; i++ {
		content += string(byte(rand.Intn(25) + 97))
	}

	qr, err := Create(content, &QRCodeOptions{
		Mode:       encode.EncodingModeByte,
		ErrorLevel: ErrorCorrectionLevelQuartile,
	})
	if err != nil {
		t.Error(err)
		return
	}

	file, err := os.Create("test_long.png")
	if err != nil {
		t.Error(err)
		return
	}

	defer file.Close()

	if err := qr.Plot(file); err != nil {
		t.Error(err)
	}
}

func TestPlotMixed(t *testing.T) {
	encodeBlocks := []*encode.EncodeBlock{
		{
			Mode: encode.EncodingModeNumeric,
			Data: "1234567890",
		},
		{
			Mode: encode.EncodingModeAlphaNumeric,
			Data: "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		},
		{
			Mode: encode.EncodingModeByte,
			Data: "abcdefghijklmnopqrstuvwxyz",
		},
		{
			Mode: encode.EncodingModeKanji,
			Data: "茗",
		},
	}

	qr, err := CreateMultiMode(encodeBlocks, nil)
	if err != nil {
		t.Error(err)
		return
	}

	file, err := os.Create("test_mix.png")
	if err != nil {
		t.Error(err)
		return
	}

	defer file.Close()

	if err := qr.Plot(file); err != nil {
		t.Error(err)
	}

	content := ""
	for _, block := range encodeBlocks {
		content += block.Data
	}

	err = ValidateContent(content, "test_mix.png", false)
	if err != nil {
		t.Error(err)
	}
}

type ComparingTestMode struct {
	mode     encode.EncodingMode
	modeName string
	from     int
	to       int
}

type ComparingTestErrorLevel struct {
	errorLevel     ErrorCorrectionLevel
	errorLevelName string
	testModes      []ComparingTestMode
}

type ComparingTest struct {
	version     int
	errorLevels []ComparingTestErrorLevel
}

func GenerateNumericContent(size int) string {
	content := make([]byte, size)
	numbers := []byte("0123456789")
	for i := 0; i < size; i++ {
		content[i] = numbers[rand.Intn(10)]
	}
	return string(content)
}

func GenerateAlphaNumericContent(size int) string {
	content := make([]byte, size)
	chars := []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ $%*+-./:")
	for i := 0; i < size; i++ {
		content[i] = chars[rand.Intn(45)]
	}
	return string(content)
}

func GenerateByteContent(size int) string {
	content := make([]rune, size)
	chars := []rune("!\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}")
	for i := 0; i < size; i++ {
		content[i] = chars[rand.Intn(93)]
	}
	return string(content)
}

func GenerateKanjiContent(size int) string {
	content := make([]rune, size)
	kanji := rune('茗')
	for i := 0; i < size; i++ {
		content[i] = kanji
	}
	return string(content)
}

func TestPlotMicroComparing(t *testing.T) {
	modes := []encode.EncodingMode{
		encode.EncodingModeNumeric,
		encode.EncodingModeAlphaNumeric,
		encode.EncodingModeByte,
		encode.EncodingModeKanji,
	}

	modesNames := []string{
		"numeric",
		"alphanumeric",
		"byte",
		"kanji",
	}

	errorLevels := []ErrorCorrectionLevel{
		ErrorCorrectionLevelLow,
		ErrorCorrectionLevelMedium,
		ErrorCorrectionLevelQuartile,
		ErrorCorrectionLevelHigh,
	}

	errorLevelsNames := []string{
		"low",
		"medium",
		"quartile",
		"high",
	}

	tests := make([]ComparingTest, 0)
	for i := 1; i <= 4; i++ {
		errorLevelsTests := make([]ComparingTestErrorLevel, 0)
		for j := 0; j < 4; j++ {
			modeTests := make([]ComparingTestMode, 0)
			for k := 0; k < 4; k++ {
				limit := ContentLengthLimitsMicro[i][j][k]
				if limit != 0 {
					modeTests = append(modeTests, ComparingTestMode{
						mode:     modes[k],
						modeName: modesNames[k],
						from:     1,
						to:       limit,
					})
				}
			}
			if len(modeTests) > 0 {
				errorLevelsTests = append(errorLevelsTests, ComparingTestErrorLevel{
					errorLevel:     errorLevels[j],
					errorLevelName: errorLevelsNames[j],
					testModes:      modeTests,
				})
			}
		}
		if len(errorLevelsTests) > 0 {
			tests = append(tests, ComparingTest{
				version:     -i,
				errorLevels: errorLevelsTests,
			})
		}
	}

	tempDir := t.TempDir()

	for _, test := range tests {
		testName := fmt.Sprintf("version_m%d", -test.version)
		t.Run(testName, func(t *testing.T) {
			version := test.version
			errorLevelTests := test.errorLevels
			t.Parallel()
			for _, errorLevelTest := range errorLevelTests {
				testName := fmt.Sprintf("error_level_%s", errorLevelTest.errorLevelName)
				t.Run(testName, func(t *testing.T) {
					errorLevel := errorLevelTest.errorLevel
					testModes := errorLevelTest.testModes
					t.Parallel()
					for _, testMode := range testModes {
						testName := fmt.Sprintf("mode_%v", testMode.modeName)
						t.Run(testName, func(t *testing.T) {
							mode := testMode.mode
							from, to := testMode.from, testMode.to
							t.Parallel()

							for content_size := from; content_size <= to; content_size++ {
								var content string
								switch mode {
								case encode.EncodingModeNumeric:
									content = GenerateNumericContent(content_size)
								case encode.EncodingModeAlphaNumeric:
									content = GenerateAlphaNumericContent(content_size)
								case encode.EncodingModeByte:
									content = GenerateByteContent(content_size)
								case encode.EncodingModeKanji:
									content = GenerateKanjiContent(content_size)
								}

								block := &encode.EncodeBlock{
									Mode: mode,
									Data: content,
								}

								qr, err := CreateMultiMode([]*encode.EncodeBlock{block}, &QRCodeOptionsMultiMode{
									ErrorLevel: errorLevel,
									Version:    version,
								})

								if err != nil && !errors.Is(err, ErrContentTooLong) {
									t.Error(err)
									continue
								}

								filename := fmt.Sprintf("test_micro_%v_%v_%v.png", version, errorLevel, mode)

								filename = filepath.Join(tempDir, filename)

								file, err := os.Create(filename)
								if err != nil {
									t.Error(err)
									continue
								}

								if err := qr.Plot(file); err != nil {
									t.Error(err)
									file.Close()
									continue
								}

								file.Close()

								err = ValidateContent(content, filename, true)
								if err != nil {
									t.Error(err)
								}
							}
						})
					}
				})
			}
		})
	}
}

func TestPlotComparing(t *testing.T) {
	modes := []encode.EncodingMode{
		encode.EncodingModeNumeric,
		encode.EncodingModeAlphaNumeric,
		encode.EncodingModeByte,
		encode.EncodingModeKanji,
	}

	modesNames := []string{
		"numeric",
		"alphanumeric",
		"byte",
		"kanji",
	}

	errorLevels := []ErrorCorrectionLevel{
		ErrorCorrectionLevelLow,
		ErrorCorrectionLevelMedium,
		ErrorCorrectionLevelQuartile,
		ErrorCorrectionLevelHigh,
	}

	errorLevelsNames := []string{
		"low",
		"medium",
		"quartile",
		"high",
	}

	tests := make([]ComparingTest, 0)
	previours_limits := ContentLengthLimits[0]
	for i := 1; i <= 40; i++ {
		errorLevelsTests := make([]ComparingTestErrorLevel, 0)
		for j := 0; j < 4; j++ {
			modeTests := make([]ComparingTestMode, 0)
			for k := 0; k < 4; k++ {
				limit := ContentLengthLimits[i][j][k]
				from := previours_limits[j][k] + 1

				if limit != 0 {
					modeTests = append(modeTests, ComparingTestMode{
						mode:     modes[k],
						modeName: modesNames[k],
						from:     from,
						to:       limit,
					})
				}
			}
			if len(modeTests) > 0 {
				errorLevelsTests = append(errorLevelsTests, ComparingTestErrorLevel{
					errorLevel:     errorLevels[j],
					errorLevelName: errorLevelsNames[j],
					testModes:      modeTests,
				})
			}
		}
		if len(errorLevelsTests) > 0 {
			tests = append(tests, ComparingTest{
				version:     i,
				errorLevels: errorLevelsTests,
			})
		}
	}

	// f, err := os.Create("cpu.pprof")
	// if err != nil {
	// 	panic(err)
	// }
	// defer f.Close()

	// pprof.StartCPUProfile(f)
	// defer pprof.StopCPUProfile()

	// P, P, P -> 1362s
	// N, P, P -> 2132s
	// P, N, P -> 1287s
	// P, P, N -> 1277s

	for _, test := range tests {
		testName := fmt.Sprintf("version_%d", test.version)
		if test.version > 5 {
			return
		}
		t.Run(testName, func(t *testing.T) {
			version := test.version
			errorLevelTests := test.errorLevels
			// t.Parallel()
			for _, errorLevelTest := range errorLevelTests {
				testName := fmt.Sprintf("errors_level_%s", errorLevelTest.errorLevelName)
				t.Run(testName, func(t *testing.T) {
					errorLevel := errorLevelTest.errorLevel
					testModes := errorLevelTest.testModes
					t.Parallel()
					for _, testMode := range testModes {
						testName := fmt.Sprintf("mode_%v", testMode.modeName)
						t.Run(testName, func(t *testing.T) {
							mode := testMode.mode
							from, to := testMode.from, testMode.to
							// t.Parallel()

							for content_size := from; content_size <= to; content_size++ {
								var content string
								switch mode {
								case encode.EncodingModeNumeric:
									content = GenerateNumericContent(content_size)
								case encode.EncodingModeAlphaNumeric:
									content = GenerateAlphaNumericContent(content_size)
								case encode.EncodingModeByte:
									content = GenerateByteContent(content_size)
								case encode.EncodingModeKanji:
									content = GenerateKanjiContent(content_size)
								}

								block := &encode.EncodeBlock{
									Mode: mode,
									Data: content,
								}

								qr, err := CreateMultiMode([]*encode.EncodeBlock{block}, &QRCodeOptionsMultiMode{
									ErrorLevel: errorLevel,
									Version:    version,
								})

								if err != nil && !errors.Is(err, ErrContentTooLong) {
									t.Error(err)
									continue
								}

								err = ValidateContentRaw(qr.Data, content)
								if err != nil {
									t.Error(err)
								}
							}
						})
					}
				})
			}
		})
	}
}
