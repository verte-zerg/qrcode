package qrcode

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"testing"
)

func TestPlotRectangle(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	step := 2
	data := [][]int{
		// x, y, left, right, top, bottom
		{0, 0, 0, 2, 0, 2},
		{1, 1, 2, 4, 2, 4},
		{2, 2, 4, 6, 4, 6},
		{3, 3, 6, 8, 6, 8},
		{4, 4, 8, 10, 8, 10},
	}

	whiteClr := color.RGBA{255, 255, 255, 255}
	for _, row := range data {
		x, y, left, right, top, bottom := row[0], row[1], row[2], row[3], row[4], row[5]
		plotRectangle(img, x, y, step, image.White)
		for idx := left; idx < right; idx++ {
			for jdx := top; jdx < bottom; jdx++ {
				clr := img.At(idx, jdx)
				if clr != whiteClr {
					t.Errorf("expected black pixel at (%v, %v), got %v", idx, jdx, img.At(idx, jdx))
				}
			}
		}
	}
}

type InvalidWriter struct{}

func (w *InvalidWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("invalid writer")
}

func TestPlot(t *testing.T) {
	// Valid test
	t.Run("valid", func(t *testing.T) {
		data := [][]Cell{
			{
				{Value: true},
				{Value: false},
				{Value: true},
			},
			{
				{Value: false},
				{Value: true},
				{Value: false},
			},
			{
				{Value: true},
				{Value: false},
				{Value: true},
			},
		}

		var buf bytes.Buffer
		err := plot(data, &buf, 1)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// Open the generated image
		img, _, err := image.Decode(&buf)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		clrWhite := color.RGBA{255, 255, 255, 255}
		clrBlack := color.RGBA{0, 0, 0, 255}
		for idx, row := range data {
			for jdx, cell := range row {
				clr := img.At(jdx, idx)
				if cell.Value && clr != clrBlack {
					t.Errorf("expected black pixel at (%v, %v), got %v", jdx, idx, clr)
				}
				if !cell.Value && clr != clrWhite {
					t.Errorf("expected white pixel at (%v, %v), got %v", jdx, idx, clr)
				}
			}
		}
	})

	// Invalid test
	t.Run("invalid", func(t *testing.T) {
		data := [][]Cell{}

		// Create a buffer that will always return an error when writing
		var buf InvalidWriter

		// Call the function and check for error
		err := plot(data, &buf, 1)
		if err == nil {
			t.Error("expected an error, but got nil")
		}
	})

}
