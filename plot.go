package qrcode

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
)

func PlotRectangle(img *image.RGBA, x, y, size int, clr color.Color) {
	for idx := x * size; idx < (x+1)*size; idx++ {
		for jdx := y * size; jdx < (y+1)*size; jdx++ {
			img.Set(idx, jdx, clr)
		}
	}
}

func Plot(data [][]Cell, writer io.Writer) error {
	scale := 4
	img := image.NewRGBA(image.Rect(0, 0, len(data)*scale, len(data)*scale))
	for y, row := range data {
		for x, cell := range row {
			clr := image.White
			if cell.Value {
				clr = image.Black
			}
			PlotRectangle(img, x, y, scale, clr)
		}
	}

	err := png.Encode(writer, img)

	if err != nil {
		return fmt.Errorf("failed to encode png: %w", err)
	}

	return nil
}
