package qrcode

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
)

type OutputFormat string

const (
	// Output formats
	PNG  = OutputFormat("png")
	JPEG = OutputFormat("jpeg")
	GIF  = OutputFormat("gif")
)

type MarkerType string

const (
	// Marker types
	Square         = MarkerType("square")
	Circle         = MarkerType("circle")
	SmallCircle    = MarkerType("smallCircle")
	HorizontalLine = MarkerType("horizontalLine")
	VerticalLine   = MarkerType("verticalLine")
)

// drawRectangle fills a rectangle in the image with the given options.BlackColor
func drawRectangle(img *image.RGBA, idx, idy int, options *PlotOptions) {
	for x := idx * options.Scale; x < (idx+1)*options.Scale; x++ {
		for y := idy * options.Scale; y < (idy+1)*options.Scale; y++ {
			img.Set(x+options.Border, y+options.Border, options.BlackColor)
		}
	}
}

// drawPartialRectangle fills a rectangle in the image with the given options.BlackColor
// but only specific part of the rectangle
func drawPartialRectangle(img *image.RGBA, idx, idy int, options *PlotOptions, up, bottom, left, right bool) {
	for x := idx * options.Scale; x < (idx+1)*options.Scale; x++ {
		for y := idy * options.Scale; y < (idy+1)*options.Scale; y++ {
			if up && y > idy*options.Scale+options.Scale/2 {
				continue
			}

			if bottom && y <= idy*options.Scale+options.Scale/2 {
				continue
			}

			if left && x > idx*options.Scale+options.Scale/2 {
				continue
			}

			if right && x <= idx*options.Scale+options.Scale/2 {
				continue
			}

			img.Set(x+options.Border, y+options.Border, options.BlackColor)
		}
	}
}

// drawCircle fills a circle in the image with the given options.BlackColor
func drawCircle(img *image.RGBA, idx, idy int, options *PlotOptions) {
	x_from := idx*options.Scale - 1
	x_to := (idx + 1) * options.Scale
	x_center := float32(x_to+x_from) / 2

	y_from := idy*options.Scale - 1
	y_to := (idy + 1) * options.Scale
	y_center := float32(y_to+y_from) / 2

	radius := float32(options.Scale) / 2

	for idx := x_from; idx < x_to; idx++ {
		for jdx := y_from; jdx < y_to; jdx++ {
			pos_x := float32(idx) - x_center
			pos_y := float32(jdx) - y_center

			if pos_x*pos_x+pos_y*pos_y <= radius*radius {
				img.Set(idx+options.Border, jdx+options.Border, options.BlackColor)
			}
		}
	}
}

// plotSquadMarker fills a rectangle in the image with the given color
func plotSquadMarkers(img *image.RGBA, data [][]Cell, options *PlotOptions) {
	size := len(data)

	for idx := 0; idx < size; idx++ {
		for idy := 0; idy < size; idy++ {
			if data[idx][idy].Value {
				drawRectangle(img, idx, idy, options)
			}
		}
	}
}

// plotHorizontalLineMarkers fills a horizontal line in the image with the given color
func plotHorizontalLineMarkers(img *image.RGBA, data [][]Cell, options *PlotOptions) {
	size := len(data)

	for idx := 0; idx < size; idx++ {
		for idy := 0; idy < size; idy++ {
			if !data[idx][idy].Value {
				continue
			}

			// if the cell is a search pattern, keep it as a rectangle
			if data[idx][idy].Type == CellTypeSearchPattern {
				drawRectangle(img, idx, idy, options)
				continue
			}

			// if the cell is the first or the last in the line, draw a rounded rectangle
			if (idx == 0 || !data[idx-1][idy].Value) || (idx == size-1 || !data[idx+1][idy].Value) {
				if idx < size-1 && data[idx+1][idy].Value {
					drawPartialRectangle(img, idx, idy, options, false, false, false, true)
				} else if idx > 0 && data[idx-1][idy].Value {
					drawPartialRectangle(img, idx, idy, options, false, false, true, false)
				}

				drawCircle(img, idx, idy, options)
				continue
			}

			// if the cell is in the middle of the line, draw a rectangle
			drawRectangle(img, idx, idy, options)
		}
	}
}

// plotVerticalLineMarkers fills a vertical line in the image with the given color
func plotVerticalLineMarkers(img *image.RGBA, data [][]Cell, options *PlotOptions) {
	size := len(data)

	for idx := 0; idx < size; idx++ {
		for idy := 0; idy < size; idy++ {
			if !data[idx][idy].Value {
				continue
			}

			// if the cell is a search pattern, keep it as a rectangle
			if data[idx][idy].Type == CellTypeSearchPattern {
				drawRectangle(img, idx, idy, options)
				continue
			}

			// if the cell is the first or the last in the line, draw a rounded rectangle
			if (idy == 0 || !data[idx][idy-1].Value) || (idy == size-1 || !data[idx][idy+1].Value) {
				if idy < size-1 && data[idx][idy+1].Value {
					drawPartialRectangle(img, idx, idy, options, false, true, false, false)
				} else if idy > 0 && data[idx][idy-1].Value {
					drawPartialRectangle(img, idx, idy, options, true, false, false, false)
				}

				drawCircle(img, idx, idy, options)
				continue
			}

			// if the cell is in the middle of the line, draw a rectangle
			drawRectangle(img, idx, idy, options)
		}
	}
}

// plotCircleCustomMarker fills a circle in the image with the given color
func plotCircleCustomMarkers(img *image.RGBA, data [][]Cell, options *PlotOptions, divider float32) {
	size := len(data)

	for idx := 0; idx < size; idx++ {
		for idy := 0; idy < size; idy++ {
			if !data[idx][idy].Value {
				continue
			}

			if data[idx][idy].Type == CellTypeSearchPattern {
				drawRectangle(img, idx, idy, options)
				continue
			}

			drawCircle(img, idx, idy, options)
		}
	}
}

// plotCircleMarker fills a circle in the image with the given color
var plotCircleMarkers = func(img *image.RGBA, data [][]Cell, options *PlotOptions) {
	plotCircleCustomMarkers(img, data, options, 2)
}

// plotSmallCircleMarker fills a small circle in the image with the given color
var plotSmallCircleMarkers = func(img *image.RGBA, data [][]Cell, options *PlotOptions) {
	plotCircleCustomMarkers(img, data, options, 3)
}

// getPlotMarkerFunc returns the plot marker function for the given marker type
func getPlotMarkerFunc(markerType MarkerType) (func(*image.RGBA, [][]Cell, *PlotOptions), error) {
	switch markerType {
	case Square:
		return plotSquadMarkers, nil
	case Circle:
		return plotCircleMarkers, nil
	case SmallCircle:
		return plotSmallCircleMarkers, nil
	case HorizontalLine:
		return plotHorizontalLineMarkers, nil
	case VerticalLine:
		return plotVerticalLineMarkers, nil
	default:
		return nil, fmt.Errorf("unsupported marker type: %s", markerType)
	}
}

// plotMarkers plots the markers on the image
func plotMarkers(data [][]Cell, img *image.RGBA, options *PlotOptions) error {
	markerFunc, err := getPlotMarkerFunc(options.MarkerType)

	if err != nil {
		return fmt.Errorf("failed to get plot marker function: %w", err)
	}

	imgSize := len(data)*options.Scale + 2*options.Border

	// fill the image with white
	for idx := 0; idx < imgSize; idx++ {
		for jdx := 0; jdx < imgSize; jdx++ {
			img.Set(idx, jdx, options.WhiteColor)
		}
	}

	markerFunc(img, data, options)
	return nil
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

// writeImage writes the image to the writer
func writeImage(img *image.RGBA, writer io.Writer, format OutputFormat) error {
	switch format {
	case PNG:
		return png.Encode(writer, img)
	case JPEG:
		return jpeg.Encode(writer, img, nil)
	case GIF:
		return gif.Encode(writer, img, nil)
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}

// plot creates a image from the given data and writes it to the writer
func plot(data [][]Cell, writer io.Writer, options *PlotOptions) error {
	if options == nil {
		return fmt.Errorf("options is nil")
	}

	imgSize := len(data)*options.Scale + 2*options.Border
	img := image.NewRGBA(image.Rect(0, 0, imgSize, imgSize))

	err := plotMarkers(data, img, options)
	if err != nil {
		return fmt.Errorf("failed to plot markers: %w", err)
	}

	if options.Border > 0 {
		plotBorder(img, options.Border, options.WhiteColor)
	}

	err = writeImage(img, writer, options.OutputFormat)

	if err != nil {
		return fmt.Errorf("failed to encode png: %w", err)
	}

	return nil
}
