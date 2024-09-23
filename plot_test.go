package qrcode

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"testing"
)

var (
	outputFormatsColor = map[OutputFormat]struct {
		white, black color.Color
	}{
		PNG:  {color.RGBA{255, 255, 255, 255}, color.RGBA{0, 0, 0, 255}},
		JPEG: {color.YCbCr{255, 128, 128}, color.YCbCr{0, 128, 128}},
		GIF:  {color.RGBA{255, 255, 255, 255}, color.RGBA{0, 0, 0, 255}},
	}
)

func TestPlotRectangle(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 6, 6))
	scale := 2
	data := [][]int{
		// left, right, top, bottom
		{0, 2, 0, 2},
		{2, 4, 2, 4},
		{4, 6, 4, 6},
	}

	wC := Cell{Value: false}
	bC := Cell{Value: true}

	cells := [][]Cell{
		{bC, bC, wC, wC, wC, wC},
		{bC, bC, wC, wC, wC, wC},
		{wC, wC, bC, bC, wC, wC},
		{wC, wC, bC, bC, wC, wC},
		{wC, wC, wC, wC, bC, bC},
		{wC, wC, wC, wC, bC, bC},
	}

	// fill the image with white
	for idx := 0; idx < 10; idx++ {
		for jdx := 0; jdx < 10; jdx++ {
			img.Set(idx, jdx, image.White)
		}
	}

	plotOptions := &PlotOptions{
		Scale:        scale,
		Border:       0,
		OutputFormat: PNG,
		MarkerType:   Square,
		WhiteColor:   image.White,
		BlackColor:   image.Black,
	}

	blackClr := color.RGBA{0, 0, 0, 255}
	plotSquadMarkers(img, cells, plotOptions)

	for _, row := range data {
		left, right, top, bottom := row[0], row[1], row[2], row[3]
		for idx := left; idx < right; idx++ {
			for jdx := top; jdx < bottom; jdx++ {
				clr := img.At(idx, jdx)
				if clr != blackClr {
					t.Errorf("expected black pixel at (%v, %v), got %v", idx, jdx, img.At(idx, jdx))
				}
			}
		}
	}
}

func TestPlotBorder(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 6, 6))
	border := 2
	colors := outputFormatsColor[PNG]
	whiteClr := colors.white
	blackClr := colors.black

	// Fill the image with black
	for idx := 0; idx < 6; idx++ {
		for jdx := 0; jdx < 6; jdx++ {
			img.Set(idx, jdx, blackClr)
		}
	}

	expected := [][]color.Color{
		{whiteClr, whiteClr, whiteClr, whiteClr, whiteClr, whiteClr},
		{whiteClr, whiteClr, whiteClr, whiteClr, whiteClr, whiteClr},
		{whiteClr, whiteClr, blackClr, blackClr, whiteClr, whiteClr},
		{whiteClr, whiteClr, blackClr, blackClr, whiteClr, whiteClr},
		{whiteClr, whiteClr, whiteClr, whiteClr, whiteClr, whiteClr},
		{whiteClr, whiteClr, whiteClr, whiteClr, whiteClr, whiteClr},
	}

	plotBorder(img, border, image.White)
	for idx := 0; idx < 6; idx++ {
		for jdx := 0; jdx < 6; jdx++ {
			if clr := img.At(idx, jdx); clr != expected[idx][jdx] {
				t.Errorf("expected %v at (%v, %v), got %v", expected[idx][jdx], idx, jdx, clr)
			}
		}
	}
}

type InvalidWriter struct{}

func (w *InvalidWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("invalid writer")
}

func TestPlot(t *testing.T) {
	plotOptions := &PlotOptions{
		Scale:        1,
		Border:       0,
		OutputFormat: PNG,
		MarkerType:   Square,
		WhiteColor:   image.White,
		BlackColor:   image.Black,
	}

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

		formats := []OutputFormat{PNG, GIF}
		for _, format := range formats {

			t.Run(string(format), func(t *testing.T) {
				var buf bytes.Buffer

				plotOptions.OutputFormat = format

				err := plot(data, &buf, plotOptions)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				// Open the generated image
				img, _, err := image.Decode(&buf)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				colors := outputFormatsColor[format]
				clrWhite := colors.white
				clrBlack := colors.black
				for idx, row := range data {
					for jdx, cell := range row {
						clr := img.At(jdx, idx)
						if cell.Value && clr != clrBlack {
							t.Errorf("expected black pixel at (%v, %v), expected %v, got %v", jdx, idx, clrBlack, clr)
						}
						if !cell.Value && clr != clrWhite {
							t.Errorf("expected white pixel at (%v, %v), expected %v, got %v", jdx, idx, clrWhite, clr)
						}
					}
				}
			})
		}
	})

	// Invalid test
	t.Run("invalid", func(t *testing.T) {
		data := [][]Cell{}

		// Create a buffer that will always return an error when writing
		var buf InvalidWriter

		// Call the function and check for error
		err := plot(data, &buf, plotOptions)
		if err == nil {
			t.Error("expected an error, but got nil")
		}
	})

}
