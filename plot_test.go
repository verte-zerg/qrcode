package qrcode

import (
	"errors"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"

	"github.com/verte-zerg/qrcode/encode"
)

func ValidateContent(content string, filename string) error {
	cmd := exec.Command("zbarimg", "-q", filename)

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

func TestPlot(t *testing.T) {
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

	err = ValidateContent(content, "test_new.png")
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
		Mode:       encode.EncodingModeLatin1,
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
	encodeBlocks := []encode.EncodeBlock{
		{
			Mode: encode.EncodingModeNumeric,
			Data: "1234567890",
		},
		{
			Mode: encode.EncodingModeAlphaNumeric,
			Data: "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		},
		{
			Mode: encode.EncodingModeLatin1,
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

	err = ValidateContent(content, "test_mix.png")
	if err != nil {
		t.Error(err)
	}
}

func TestPlotWithComparing(t *testing.T) {
	modes := []encode.EncodingMode{
		encode.EncodingModeNumeric,
		encode.EncodingModeAlphaNumeric,
		encode.EncodingModeLatin1,
		encode.EncodingModeKanji,
	}

	for _, mode := range modes {
		name := strconv.Itoa(int(mode))
		t.Run(name, func(t *testing.T) {
			mode := mode
			t.Parallel()
			for i := 0; i < 500; i++ {
				var content string

				size := i + 1
				for i := 0; i < size; i++ {
					if mode == encode.EncodingModeNumeric || mode == encode.EncodingModeAlphaNumeric {
						content += strconv.Itoa(rand.Intn(10))
					}

					if mode == encode.EncodingModeLatin1 {
						content += string(byte(rand.Intn(25) + 97))
					}

					if mode == encode.EncodingModeKanji {
						content += "茗"
					}
				}

				qr, err := Create(content, &QRCodeOptions{
					Mode:       mode,
					ErrorLevel: ErrorCorrectionLevelLow,
				})

				if err != nil && !errors.Is(err, ErrContentTooLong) {
					t.Error(err)
					continue
				}

				file, err := os.Create("test_" + name + ".png")
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

				err = ValidateContent(content, "test_"+name+".png")
				if err != nil {
					t.Error(err)
				}
			}
		})
	}
}
