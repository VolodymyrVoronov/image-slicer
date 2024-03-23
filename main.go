package main

import (
	"fmt"
	"image-slicer/utils"
	"io/fs"
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

func main() {
	rows := utils.GetUserInput("Enter amount of horizontal rows: ")
	cols := utils.GetUserInput("Enter amount of vertical columns: ")

	progressChannel := make(chan int)

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

	go func() {
		totalProgress := 0
		totalImages := len(inputDir)

		for progress := range progressChannel {
			totalProgress += progress
			overallProgress := totalProgress * 100 / (totalImages * 100)

			fmt.Printf("Overall Progress: %d%%\n", overallProgress)

			if overallProgress >= 100 {
				close(progressChannel)
			}
		}
	}()

	progress := calculateProgress(doneChannels, inputDir)
	progressChannel <- progress

	for i, image := range inputDir {
		doneChannels[i] = make(chan string, 1)
		errorChannels[i] = make(chan error, 1)

		if !image.IsDir() {

			imageName := image.Name()
			pathToImage := fmt.Sprintf("%s/%s", input, imageName)

			func(i int) {
				utils.SliceImage(pathToImage, output, rows, cols, doneChannels[i], errorChannels[i], progressChannel)
			}(i)

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

func calculateProgress(doneChannels []chan string, inputDir []fs.DirEntry) int {
	processedImages := 0
	for _, doneChannel := range doneChannels {
		select {
		case <-doneChannel:
			processedImages++
		default:
		}
	}

	totalImages := len(inputDir)
	progress := processedImages * 100 / totalImages

	return progress
}
