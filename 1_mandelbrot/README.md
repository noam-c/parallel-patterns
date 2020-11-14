# Mandelbrot Set Painter

A single-threaded Go program that draws the Mandelbrot Set and outputs an image file.

While currently single-threaded, we can make this painter program multithreaded using threads and a barrier.
The `createImage()` function currently calls `fillInRows` once and passes in all the rows for one thread to draw. We will split this into four different calls and put each call in a separate thread.

## Project Instructions
All work will be done in `createImage()` in `main.go`.

1. Split the `fillInRows` calls into four separate `fillInRows` calls -- split the workload across four calls.
2. Create a `sync.WaitGroup` object, and set its count to 4 (we'll make four threads later).
3. After each `fillInRows` call, call `Done()` on the `WaitGroup` object to signal that another fourth of the work was completed.
4. After all of the `fillInRows` calls, `Wait()` on the `WaitGroup` object.
5. Now, wrap each `fillInRows`/`Done` pair in a `go func() { ... }()` call.
6. Build and run the code! It should draw the same image twice as fast!

## Further Tweaks
* You can change the size of the image via the `ImageWidth` and `ImageHeight` constants at the top of `main.go`.
* You can change the color scheme of the image by changing the list of colors in `createColorPalette()` in `main.go`.
* Alter the call to `runtime.GOMAXPROCS` if your computer has more than 2 cores -- it should improve the speed of the program.