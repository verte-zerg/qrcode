package qrcode

import (
	"fmt"
	"io"

	"github.com/verte-zerg/qrcode/encode"
)

const (
	// DEFAULT_SCALE is the default scale for the QR Code image.
	// The image will be len(data) * DEFAULT_SCALE x len(data) * DEFAULT_SCALE pixels.
	DEFAULT_SCALE = 4
)

const (
	// MicroQR versions
	M1 = -1
	M2 = -2
	M3 = -3
	M4 = -4
)

// QRCode is a struct that represents a QR Code.
type QRCode struct {
	// Content
	Content string
	// Options
	options *QRCodeOptions

	// Data
	Data [][]Cell
}

// QRCodeOptions is a struct that represents the options for the QR Code.
type QRCodeOptions struct {
	// Encoding is the encoding mode.
	// Default: calculated based on the content (can undestand only numeric, alphanumeric, byte, kanji or utf-8 with ECI)
	Mode encode.EncodingMode

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

// QRCodeOptionsMultiMode is a struct that represents the options for building multi-mode QR Codes.
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

// CreateMultiMode creates a QR Code with multiple modes.
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
		version, err = calculateMinVersion(blocks, qrCodeOptions.ErrorLevel, options.MicroQR)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate min version: %w", err)
		}
	}
	qrCodeOptions.Version = version

	buf, err := getBytesData(blocks, options.ErrorLevel, version)
	if err != nil {
		return nil, fmt.Errorf("failed to get bytes data: %w", err)
	}

	data := generateField(buf, version, options.ErrorLevel)

	return &QRCode{
		options: qrCodeOptions,
		Data:    data,
	}, nil
}

// Create creates a QR Code with the given content and options.
func Create(content string, options *QRCodeOptions) (*QRCode, error) {
	if options == nil {
		options = &QRCodeOptions{}
	}

	encodeBlock := &encode.EncodeBlock{
		Mode: options.Mode,
		Data: content,
	}
	if encodeBlock.Mode == 0 {
		var err error
		encodeBlock.Mode, err = encode.GetEncodingMode(content)

		// If the content is not valid for any mode, use UTF-8 with ECI
		if err != nil {
			encodeBlock.Mode = encode.EncodingModeECI
			encodeBlock.SubMode = encode.EncodingModeByte
			encodeBlock.AssignmentNumber = 26
		}
	}

	return CreateMultiMode([]*encode.EncodeBlock{encodeBlock}, &QRCodeOptionsMultiMode{
		ErrorLevel: options.ErrorLevel,
		Version:    options.Version,
		MicroQR:    options.MicroQR,
	})
}

// Plot plots the QR Code to the given writer.
func (qr *QRCode) Plot(writer io.Writer) error {
	return plot(qr.Data, writer, DEFAULT_SCALE)
}
