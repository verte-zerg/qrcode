package qrcode

import (
	"fmt"
	"image"
	"image/color"
	"io"

	"qrcode/encode"
)

const (
	// DEFAULT_SCALE is the default scale for the QR Code image.
	DEFAULT_SCALE = 4

	// DEFAULT_BORDER is the default border for the QR Code image.
	DEFAULT_BORDER = 0

	// DEFAULT_OUTPUT_FORMAT is the default output format for the QR Code image.
	DEFAULT_OUTPUT_FORMAT = PNG

	// DEFAULT_MARKER_TYPE is the default marker type for the QR Code image.
	DEFAULT_MARKER_TYPE = Square
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

type PlotOptions struct {
	// Scale is the scale for the QR Code image (in pixels).
	// The image will be len(data) * Scale x len(data) * Scale pixels.
	// Default: 4.
	Scale int

	// Border is the border for the QR Code image (in pixels).
	// Default: 0.
	Border int

	// OutputFormat is the format of the output image.
	// Default: PNG.
	OutputFormat OutputFormat

	// MarkerType is the type of marker to use for the QR Code.
	// Default: Square.
	MarkerType MarkerType

	// WhiteColor is the color to use for the white cells.
	// Default: image.White.
	WhiteColor color.Color

	// BlackColor is the color to use for the black cells.
	// Default: image.Black.
	BlackColor color.Color
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
			encodeBlock.AssignmentNumber = encode.UTF8
		}
	}

	return CreateMultiMode([]*encode.EncodeBlock{encodeBlock}, &QRCodeOptionsMultiMode{
		ErrorLevel: options.ErrorLevel,
		Version:    options.Version,
		MicroQR:    options.MicroQR,
	})
}

// Plot plots the QR Code to the given writer with the given options.
func (qr *QRCode) Plot(writer io.Writer, options *PlotOptions) error {
	if options == nil {
		options = &PlotOptions{}
	}

	if options.Scale == 0 {
		options.Scale = DEFAULT_SCALE
	}

	if options.OutputFormat == "" {
		options.OutputFormat = DEFAULT_OUTPUT_FORMAT
	}

	if options.MarkerType == "" {
		options.MarkerType = DEFAULT_MARKER_TYPE
	}

	if options.WhiteColor == nil {
		options.WhiteColor = image.White
	}

	if options.BlackColor == nil {
		options.BlackColor = image.Black
	}

	return plot(qr.Data, writer, options)
}
