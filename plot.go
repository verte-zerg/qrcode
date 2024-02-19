package qrcode

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
)

// plotRectangle fills a rectangle in the image with the given color
func plotRectangle(img *image.RGBA, x, y, size, shift int, clr color.Color) {
	for idx := x * size; idx < (x+1)*size; idx++ {
		for jdx := y * size; jdx < (y+1)*size; jdx++ {
			img.Set(idx+shift, jdx+shift, clr)
		}
	}
}

// plotBorder fills the border of the image with the given color
func plotBorder(img *image.RGBA, border int, clr color.Color) {
	size := img.Bounds().Size().X
	for idx := 0; idx < size; idx++ {
		for borderIdx := 0; borderIdx < border; borderIdx++ {
			img.Set(borderIdx, idx, clr)
			img.Set(idx, borderIdx, clr)
			img.Set(size-1-borderIdx, idx, clr)
			img.Set(idx, size-1-borderIdx, clr)
		}
	}
}

// plot creates a PNG image from the given data and writes it to the writer
func plot(data [][]Cell, writer io.Writer, scale, border int) error {
	imgSize := len(data)*scale + 2*border

	img := image.NewRGBA(image.Rect(0, 0, imgSize, imgSize))
	for idx, row := range data {
		for jdx, cell := range row {
			clr := image.White
			if cell.Value {
				clr = image.Black
			}
			plotRectangle(img, jdx, idx, scale, border, clr)
		}
	}

	// Draw border
	if border > 0 {
		plotBorder(img, border, image.White)
	}

	err := png.Encode(writer, img)

	if err != nil {
		return fmt.Errorf("failed to encode png: %w", err)
	}

	return nil
}
