package main

import (
	"fmt"
	"image-slicer/utils"
	"os"
	"time"

	_ "image/jpeg"
	_ "image/png"
)

const input = "./src"
const output = "./dist"

func init() {
	utils.ClearDir(output)

	fmt.Println("Output directory was cleared!")
	fmt.Println()
}

var progressChannels []chan int

func main() {
	rows := utils.GetUserInput("Enter amount of horizontal rows: ")
	cols := utils.GetUserInput("Enter amount of vertical columns: ")

	fmt.Println()

	start := time.Now()

	inputDir, err := os.ReadDir(input)
	if err != nil {
		fmt.Println("Error reading input directory: ", err)
		return
	}

	if len(inputDir) == 0 {
		fmt.Println("Input directory is empty")
		return
	}

	doneChannels := make([]chan string, len(inputDir))
	errorChannels := make([]chan error, len(inputDir))
	progressChannels = append(progressChannels, make(chan int))

	for i, image := range inputDir {
		doneChannels[i] = make(chan string, 1)
		errorChannels[i] = make(chan error, 1)

		if !image.IsDir() {

			imageName := image.Name()
			pathToImage := fmt.Sprintf("%s/%s", input, imageName)

			go utils.SliceImage(pathToImage, output, rows, cols, doneChannels[i], errorChannels[i], progressChannels[i])
		}
	}

	for i := range inputDir {
		select {
		case done := <-doneChannels[i]:
			if done != "" {
				fmt.Println("Image " + inputDir[i].Name() + " was processed successfully!")
				fmt.Println(done)
			}

		case err := <-errorChannels[i]:
			if err != nil {
				fmt.Println("Image" + inputDir[i].Name() + " was processed with error!")
				fmt.Println(err)
			}

		case progress := <-progressChannels[i]:
			fmt.Printf("Progress for image %s: %d%%\n", inputDir[i].Name(), progress)
		}

	}

	duration := time.Since(start)

	var consoleMessage string

	if len(inputDir) > 1 {
		consoleMessage = "Your images have been successfully sliced!"
	} else {
		consoleMessage = "Your image has been successfully sliced!"
	}

	fmt.Println()
	fmt.Println(consoleMessage)
	fmt.Printf("Total processing time: %.2f seconds", duration.Seconds())

	time.Sleep(time.Second)
}
