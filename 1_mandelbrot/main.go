package main

import (
	"fmt"
	"image"
	"os"
	"runtime"
	"time"

	"image/color"
	"image/draw"
	"image/jpeg"
)

// Size of the picture to create. Larger images will be higher-resolution but take longer to make.
const (
	// ImageWidth is the width of the picture to draw.
	ImageWidth = 4000
	// ImageHeight is the height of the picture to draw.
	ImageHeight = 3000
)

// Color palette for the Mandelbrot set.
// Change this function to play with the colors of the image!
func createColorPalette() []color.RGBA {
	return []color.RGBA{
		color.RGBA{R: 66, G: 30, B: 15},
		color.RGBA{R: 25, G: 7, B: 26},
		color.RGBA{R: 9, G: 1, B: 47},
		color.RGBA{R: 4, G: 4, B: 73},
		color.RGBA{R: 0, G: 7, B: 100},
		color.RGBA{R: 12, G: 44, B: 138},
		color.RGBA{R: 24, G: 82, B: 177},
		color.RGBA{R: 57, G: 125, B: 209},
		color.RGBA{R: 134, G: 181, B: 229},
		color.RGBA{R: 211, G: 236, B: 248},
		color.RGBA{R: 241, G: 233, B: 191},
		color.RGBA{R: 248, G: 201, B: 95},
		color.RGBA{R: 255, G: 170, B: 0},
		color.RGBA{R: 204, G: 128, B: 0},
		color.RGBA{R: 153, G: 87, B: 0},
		color.RGBA{R: 106, G: 52, B: 3},
	}
}

// fillInRows paints the given rows of an image with the Mandelbrot Set
func fillInRows(painter MandelbrotPainter, img draw.Image, imgWidth, imgHeight, firstRow, lastRow int) {
	for x := 0; x < imgWidth; x++ {
		for y := firstRow; y < lastRow; y++ {
			painter.DrawPixel(img, x, y, imgWidth, imgHeight)
		}
	}
}

func createImage(imgWidth, imgHeight int) image.Image {
	colorPalette := createColorPalette()
	painter := NewMandelbrotPainter(colorPalette)
	//painter.SetCamera(imgWidth/3, imgHeight/4, 5.0) // Uncomment to draw a zoomed in part of the image!

	// 1. Split the fillInRows calls into four separate fillInRows calls -- split the workload across four calls.
	// 2. Create a sync.WaitGroup object, and set its count to 4 (we'll make four threads later).
	// 3. After each fillInRows call, call Done() on the WaitGroup object to signal that another fourth of the work was completed.
	// 4. After all of the fillInRows calls, Wait() on the WaitGroup object.
	// 5. Now, wrap each fillInRows/Done pair in a go func() { ... }() call
	// 6. Build and run the code! It should draw the same image twice as fast!

	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	fillInRows(painter, img, imgWidth, imgHeight, 0, imgHeight)

	return img
}

func main() {
	// Go only runs one thread at a time by default -- this call makes it use
	// more cores at once.
	// NOTE: Increasing GOMAXPROCS only helps if:
	// 1. Your code is actually multithreaded.
	// 2. You have multiple cores.
	// Many computers have 4 cores, so after adding concurrency, try changing it
	// to 4 and see what happens to your runtime!
	runtime.GOMAXPROCS(2)

	// Start timing the process by saving the current time.
	startTime := time.Now()

	// Make the image in memory.
	finalImage := createImage(ImageWidth, ImageHeight)

	// Write the image to a file.
	finalFile, _ := os.Create("output.jpg")
	jpeg.Encode(finalFile, finalImage, &jpeg.Options{Quality: 100})
	finalFile.Close()

	// Finish timing the process by calculating how much time passed.
	fmt.Println("Complete in", time.Since(startTime).Milliseconds(), "ms")
}
