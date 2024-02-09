package qrcode

import (
	"fmt"
	"io"

	"github.com/verte-zerg/qrcode/encode"
)

const (
	DEFAULT_SCALE = 4
)

type QRCode struct {
	// Content
	Content string
	// Options
	options *QRCodeOptions

	data [][]Cell
}

type QRCodeOptions struct {
	// Encoding is the encoding mode.
	// Default: calculated based on the content (can undestand only numeric, alphanumeric, latin1 and kanji).
	Mode encode.EncodingMode
	// Level is the error correction level.
	// Default: ErrorCorrectionLevelLow.
	ErrorLevel ErrorCorrectionLevel
	// Version is the version of the QR Code.
	// Default: calculated based on the content.
	Version int
}

type QRCodeOptionsMultiMode struct {
	// Level is the error correction level.
	// Default: ErrorCorrectionLevelLow.
	ErrorLevel ErrorCorrectionLevel

	// Version is the version of the QR Code.
	// Default: calculated based on the content.
	Version int

	// Enable micro QR code
	// Default: false
	MicroQR bool
}

func CreateMultiMode(blocks []*encode.EncodeBlock, options *QRCodeOptionsMultiMode) (*QRCode, error) {
	var err error
	if options == nil {
		options = &QRCodeOptionsMultiMode{}
	}
	qrCodeOptions := &QRCodeOptions{
		ErrorLevel: options.ErrorLevel,
	}

	version := options.Version

	if version == 0 {
		version, err = CalculateMinVersion(blocks, qrCodeOptions.ErrorLevel, options.MicroQR)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate min version: %w", err)
		}
	}
	qrCodeOptions.Version = version

	buf, err := GetBytesData(blocks, options.ErrorLevel, version)
	if err != nil {
		return nil, fmt.Errorf("failed to get bytes data: %w", err)
	}

	data := GenerateField(buf, version, options.ErrorLevel)

	return &QRCode{
		options: qrCodeOptions,
		data:    data,
	}, nil
}

func Create(content string, options *QRCodeOptions) (*QRCode, error) {
	if options == nil {
		options = &QRCodeOptions{}
	}

	mode := options.Mode
	if mode == 0 {
		var err error
		mode, err = encode.GetEncodingMode(content)
		if err != nil {
			return nil, fmt.Errorf("failed to get encoding mode: %w", err)
		}
	}

	encodeBlock := &encode.EncodeBlock{
		Mode: mode,
		Data: content,
	}
	return CreateMultiMode([]*encode.EncodeBlock{encodeBlock}, &QRCodeOptionsMultiMode{
		ErrorLevel: options.ErrorLevel,
		Version:    options.Version,
	})
}

func (qr *QRCode) Plot(writer io.Writer) error {
	return Plot(qr.data, writer, DEFAULT_SCALE)
}
