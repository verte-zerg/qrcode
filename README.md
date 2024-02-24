# QRCode Package

This package is a implementation of QR code encoding in Go. It offers an idiomatic Go implementation of QR code encoding,
requires only an `golang.org/x/text` dependency, maintainable by the Go team.

## Features

- [x] all modes (numeric, alphanumeric, byte, kanji, eci)
- [x] Micro QR codes
- [x] export to PNG/JPEG/GIF
- [x] ECI (Extended Channel Interpretation)
- [x] requires only an `golang.org/x/text` dependency
- [x] covered by unit/functional tests

## Installation

```bash
go get github.com/verte-zerg/qrcode
```

## Usage

### Import
```go
import "github.com/verte-zerg/qrcode"
```

### Create QR code

```go
qr, err := qrcode.Create("https://example.com", nil)
if err != nil {
    panic(err)
}
```

### Plot QR code

```go
file, err := os.Create("qrcode.png")
if err != nil {
    panic(err)
}

defer file.Close()

if err := qr.Plot(file, nil); err != nil {
    fmt.Println(err)
}
```

### Options

The `qrcode.Create` function accepts an `Options` struct as a second argument. The `Options` struct has the following fields:

```go
type QRCodeOptions struct {
	// Encoding is the encoding mode.
	// Default: calculated based on the content (numeric, alphanumeric, byte, kanji or utf-8 with ECI)
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
```

You can specify the encoding mode, error correction level, version and enable micro QR code.

Supported encoding modes:
- `encode.EncodingModeNumeric`
- `encode.EncodingModeAlphanumeric`
- `encode.EncodingModeByte`
- `encode.EncodingModeKanji`
- `encode.EncodingModeECI`

Supported error correction levels:
- `ErrorCorrectionLevelLow`
- `ErrorCorrectionLevelMedium`
- `ErrorCorrectionLevelQuartile`
- `ErrorCorrectionLevelHigh`

If Micro QR code is enabled, the version must be between M1 and M4.
Micro QR code version constants:
- `M1` (-1)
- `M2` (-2)
- `M3` (-3)
- `M4` (-4)

If you want to use specific ECI mode, you can use `qrcode.CreateMultiMode` function. The function can build QR code with several blocks of data with different modes.

```go
import (
    "github.com/verte-zerg/qrcode"
    "github.com/verte-zerg/qrcode/encode"
)

func main() {
	encodeBlocks := []*encode.EncodeBlock{
		{
			Mode: encode.EncodingModeNumeric,
			Data: "1234567890",
		},
		{
			Mode:             encode.EncodingModeECI,
			Data:             "привет мир",
			SubMode:          encode.EncodingModeByte,  // The mode must be always equal to EncodingModeByte for ECI
			AssignmentNumber: encode.ISO8859_5, // cyrillic
		},
	}

	qr, err := CreateMultiMode(encodeBlocks, nil)
	if err != nil {
		panic(err)
	}

	file, err := os.Create("qrcode_mix.png")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	if err := qr.Plot(file, nil); err != nil {
		panic(err)
	}
}
```

The list of supported ECI assignments can be found in the `encode` package.

You can specify several plot options using the `PlotOptions` struct:

```go
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
}
```

Supported output formats:
- `PNG`
- `JPEG`
- `GIF`


## Functions

`Create(content string, options *QRCodeOptions) (*QRCode, error)` - creates a QR code with the specified content and options.
`CreateMultiMode(blocks []*encode.EncodeBlock, options *QRCodeOptionsMultiMode) (*QRCode, error)` - creates a QR code with the specified blocks of data and options.
`(qr *QRCode) Plot(writer io.Writer, options *PlotOptions) error` - plots the QR code with the specified options to the writer.

## Roadmap

### Features

- [ ] add predefined QR code types (vCard, WiFi, etc.)
- [x] support other image formats (JPEG, GIF, etc.)
- [ ] data optimization algorithm
- [ ] custom data encoding
- [ ] structured append codes
- [ ] custom colors
- [ ] different shapes for the markers
- [ ] support adding a logo to the QR code

## License

This package is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing

If you find a bug or want to contribute to the code or documentation, you can help by submitting an issue or a pull request.