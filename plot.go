package qrcode

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
)

// plotRectangle fills a rectangle in the image with the given color
func plotRectangle(img *image.RGBA, x, y, size int, clr color.Color) {
	for idx := x * size; idx < (x+1)*size; idx++ {
		for jdx := y * size; jdx < (y+1)*size; jdx++ {
			img.Set(idx, jdx, clr)
		}
	}
}

// plot creates a PNG image from the given data and writes it to the writer
func plot(data [][]Cell, writer io.Writer, scale int) error {
	img := image.NewRGBA(image.Rect(0, 0, len(data)*scale, len(data)*scale))
	for y, row := range data {
		for x, cell := range row {
			clr := image.White
			if cell.Value {
				clr = image.Black
			}
			plotRectangle(img, x, y, scale, clr)
		}
	}

	err := png.Encode(writer, img)

	if err != nil {
		return fmt.Errorf("failed to encode png: %w", err)
	}

	return nil
}
