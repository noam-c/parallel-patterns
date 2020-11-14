# Julia Set Painter

A single-threaded Go program that draws the Julia Set and outputs an image file.

While currently single-threaded, we can make this painter program multithreaded by employing the Producer/Consumer pattern.
Your *Producer* thread will be responsible for breaking down the work to be done into a set of tasks (one task for each pixel that needs to be drawn). After that, the Producer will output the tasks to a shared *channel*. *Consumer* threads will then take tasks from the channel and draw them. Each Consumer will loop through the channel and do work until there is no more work left to do.

## Project Instructions
All work will be done in `createImage()` in `main.go`. For ease of startup, the channel is already provided in the starting code.

1. Move the consumer part of `createImage()` into a new thread (i.e. `go func() {...}()`).
2. Copy and paste this consumer code or call it in a loop so that 4 threads are created in total.
3. Create a `sync.WaitGroup`, and set its count to 4.
4. At the end of each consumer thread (just after the `range` loop), call `Done()` on the `WaitGroup` object.
5. Call `Wait()` on the `WaitGroup` object at the end of the `createImage()` function.
6. Remove the second parameter (`imgWidth*imgHeight`) from the channel's `make` call since we don't need such a huge channel anymore.
7. Build and run the code! It should draw the same image twice as fast!

## Further Tweaks
* You can change the size of the image via the `ImageWidth` and `ImageHeight` constants at the top of `main.go`.
* You can change the color scheme of the image by changing the list of colors in `createColorPalette()` in `main.go`.
* Alter the call to `runtime.GOMAXPROCS` if your computer has more than 2 cores -- it should improve the speed of the program.