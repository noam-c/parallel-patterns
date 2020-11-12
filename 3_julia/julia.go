package main

import (
	"image/color"
	"image/draw"
	"math"
)

// Stores offset and zoom in case we want to draw a specific portion of the image
type camera struct {
	offsetX int
	offsetY int
	zoom    float64
}

// JuliaPainter draws a Julia Set to an image canvas
type JuliaPainter struct {
	camera
	colorPalette []color.RGBA
}

// NewJuliaPainter creates a new Julia set maker, with no camera zoom
func NewJuliaPainter(colorPalette []color.RGBA) JuliaPainter {
	return JuliaPainter{
		camera: camera{
			offsetX: 0,
			offsetY: 0,
			zoom:    1.0,
		},
		colorPalette: colorPalette,
	}
}

// SetCamera changes the offset and zoom of the camera
func (m *JuliaPainter) SetCamera(offsetX, offsetY int, zoom float64) {
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

// DrawPixel draws the given pixel of the Julia set
// Algorithms from Wikipedia: https://en.wikipedia.org/wiki/Julia_set
func (m *JuliaPainter) DrawPixel(img draw.Image, pixelX, pixelY, w, h int) {
	const cX = -0.7
	const cY = 0.27015
	const r = 3.0

	fullWidth := m.zoom * float64(w)
	fullHeight := m.zoom * float64(h)
	// scaled x coordinate of pixel (scaled to lie in the Julia radius)
	scaledX := ((float64(pixelX) + float64(m.offsetX)*m.zoom) * r / float64(fullWidth)) - r/2.0
	// scaled y coordinate of pixel (scaled to lie in the Julia radius)
	scaledY := ((float64(pixelY) + float64(m.offsetY)*m.zoom) * r / float64(fullHeight)) - r/2.0

	i := 0
	const max = 10000

	for (scaledX*scaledX+scaledY*scaledY) < r*r && i < max {
		xTemp := scaledX*scaledX - scaledY*scaledY
		scaledY = 2*scaledX*scaledY + cY
		scaledX = xTemp + cX

		i++
	}

	c := color.RGBA{}
	// Smooth coloring/blending algorithm adapted from Mandelbrot Set algorithms:
	// https://en.wikipedia.org/wiki/Plotting_algorithms_for_the_Mandelbrot_set#Continuous_(smooth)_coloring
	if i < max {
		paletteIndex := float64(i)

		// My (noam-c) twist to the algorithm to stay on the same color multiple times to reduce banding
		paletteIndex = paletteIndex / 10.0

		// sqrt of inner term removed using log simplification rules.
		zn := math.Log(scaledX*scaledX+scaledY*scaledY) / 2
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
