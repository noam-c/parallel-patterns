package main

import (
	"math"

	"image/color"
	"image/draw"
)

// Stores offset and zoom in case we want to draw a specific portion of the image
type camera struct {
	offsetX int
	offsetY int
	zoom    float64
}

// MandelbrotPainter draws a Mandelbrot Set to an image canvas
type MandelbrotPainter struct {
	camera
	colorPalette []color.RGBA
}

// NewMandelbrotPainter creates a new Mandelbrot set maker, with no camera zoom
func NewMandelbrotPainter(colorPalette []color.RGBA) MandelbrotPainter {
	return MandelbrotPainter{
		camera: camera{
			offsetX: 0,
			offsetY: 0,
			zoom:    1.0,
		},
		colorPalette: colorPalette,
	}
}

// SetCamera changes the offset and zoom of the camera
func (m *MandelbrotPainter) SetCamera(offsetX, offsetY int, zoom float64) {
	m.camera.offsetX = offsetX
	m.camera.offsetY = offsetY
	m.camera.zoom = zoom
}

func blendColors(c1, c2 color.RGBA, t float64) color.RGBA {
	return color.RGBA{
		R: uint8(float64(c1.R)*(1.0-t) + float64(c2.R)*t),
		G: uint8(float64(c1.G)*(1.0-t) + float64(c2.G)*t),
		B: uint8(float64(c1.B)*(1.0-t) + float64(c2.B)*t),
	}
}

// DrawPixel draws the given pixel of the Mandelbrot set
// Algorithms from Wikipedia: https://en.wikipedia.org/wiki/Plotting_algorithms_for_the_Mandelbrot_set
func (m *MandelbrotPainter) DrawPixel(img draw.Image, pixelX, pixelY, w, h int) {
	fullWidth := m.zoom * float64(w)
	fullHeight := m.zoom * float64(h)
	// scaled x coordinate of pixel (scaled to lie in the Mandelbrot X scale (-2.5, 1))
	scaledX := ((float64(pixelX) + float64(m.offsetX)*m.zoom) * 3.5 / float64(fullWidth)) - 2.5
	// scaled y coordinate of pixel (scaled to lie in the Mandelbrot Y scale (-1, 1))
	scaledY := ((float64(pixelY) + float64(m.offsetY)*m.zoom) * 2.0 / float64(fullHeight)) - 1.0

	x := 0.0
	y := 0.0
	i := 0
	max := 1000
	for (x*x+y*y) <= 2*2 && i < max {
		xtemp := x*x - y*y + scaledX
		y = 2*x*y + scaledY
		x = xtemp
		i++
	}

	c := color.RGBA{}
	// Smooth coloring algorithm from:
	// https://en.wikipedia.org/wiki/Plotting_algorithms_for_the_Mandelbrot_set#Continuous_(smooth)_coloring
	if i > 1 && i < max {
		paletteIndex := float64(i)
		// sqrt of inner term removed using log simplification rules.
		zn := math.Log(x*x+y*y) / 2
		nu := math.Log(zn/math.Log(2)) / math.Log(2)
		// Rearranging the potential function.
		// Dividing zn by log(2) instead of log(N = 1<<8)
		// because we want the entire palette to range from the
		// center to radius 2, NOT our bailout radius.
		paletteIndex = paletteIndex + 1 - nu

		color1 := m.colorPalette[int(paletteIndex)%len(m.colorPalette)]
		color2 := m.colorPalette[(int(paletteIndex)+1)%len(m.colorPalette)]
		_, frac := math.Modf(paletteIndex)
		c = blendColors(color1, color2, frac)
	}

	img.Set(pixelX, pixelY, c)
}
