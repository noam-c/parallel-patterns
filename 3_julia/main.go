package main

import (
	"fmt"
	"image"
	"os"
	"runtime"
	"sync"
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
	// OutputFile is the name of the JPEG image file to output.
	OutputFile = "output.jpg"
)

// Color palette for the Julia set.
// Change this function to play with the colors of the image!
func createColorPalette() []color.RGBA {
	return []color.RGBA{
		{R: 0, G: 0, B: 51},
		{R: 0, G: 0, B: 77},
		{R: 0, G: 0, B: 102},
		{R: 26, G: 26, B: 127},
		{R: 51, G: 26, B: 153},
		{R: 77, G: 26, B: 153},
		{R: 77, G: 51, B: 153},
		{R: 102, G: 51, B: 127},
		{R: 102, G: 51, B: 127},
		{R: 127, G: 77, B: 127},
		{R: 153, G: 77, B: 127},
		{R: 153, G: 77, B: 127},
		{R: 189, G: 102, B: 102},
		{R: 204, G: 102, B: 102},
		{R: 204, G: 102, B: 102},
		{R: 230, G: 127, B: 77},
		{R: 230, G: 127, B: 51},
		{R: 230, G: 127, B: 51},
		{R: 230, G: 153, B: 51},
		{R: 230, G: 153, B: 51},
		{R: 230, G: 153, B: 51},
		{R: 230, G: 189, B: 51},
		{R: 230, G: 189, B: 77},
	}
}

type task struct {
	X int
	Y int
}

func doTask(painter JuliaPainter, img draw.Image, imgWidth, imgHeight int, t task) {
	painter.DrawPixel(img, t.X, t.Y, imgWidth, imgHeight)
}

func createImage(imgWidth, imgHeight int) image.Image {
	colorPalette := createColorPalette()
	painter := NewJuliaPainter(colorPalette)
	//painter.SetCamera(imgWidth/5, imgHeight/3, 3.0) // Uncomment to draw a zoomed in part of the image!

	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	// Create a channel and make sure there's room in it for all the pixels we
	// have. If we had background threads taking work off the channel as we put
	// work on, we wouldn't need to specify such a large size. But until we're
	// multithreaded, this is the only way to make this work.
	// Also, channels are SPECIFICALLY for multithreaded work -- for
	// single-threaded code, an array or a Queue data structure would suffice.
	workQueue := make(chan task, 20)

	wg := sync.WaitGroup{}

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			// Until we reach the channel's end, perform all the tasks that we find.
			// This is the "consumer" part of the code and should be run on a background thread.
			// And then, after that, copy-paste to make a total of 4 background consumer threads.
			// We also should set up our consumers BEFORE we start the producer, so
			// there is a consumer ready as soon as the producer creates some work.
			for task := range workQueue {
				doTask(painter, img, imgWidth, imgHeight, task)
			}
			wg.Done()
		}()
	}

	// This is the "producer" part of the code and may or may not run on a
	// background thread. It sets up each job to be done and then puts it in the
	// work queue to give one or more consumers tasks to do.
	go func() {
		wg.Add(1)
		for x := 0; x < imgWidth; x++ {
			for y := 0; y < imgHeight; y++ {
				workQueue <- task{x, y}
			}
		}
		close(workQueue)
		wg.Done()
	}()

	/*
		1. Move the consumer part of createImage() into a new thread (i.e. go func() {...}()).
		2. Copy and paste this consumer code or call it in a loop so that 4 threads are created in total.
		3. Create a sync.WaitGroup, and set its count to 4.
		4. At the end of each consumer thread (just after the range loop), call Done() on the WaitGroup object.
		5. Call Wait() on the WaitGroup object at the end of the createImage() function.
		6. Move the consumer work above the producer code section
		7. Remove the second parameter (imgWidth*imgHeight) from the channel's make call since we don't need such a huge channel anymore.
		8. Build and run the code! It should draw the same image twice as fast!
	*/

	wg.Wait()
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
	finalFile, _ := os.Create(OutputFile)
	jpeg.Encode(finalFile, finalImage, &jpeg.Options{Quality: 100})
	finalFile.Close()

	// Finish timing the process by calculating how much time passed.
	fmt.Println("Complete in", time.Since(startTime).Milliseconds(), "ms")
	fmt.Println("Image created:", OutputFile)
}
