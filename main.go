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
}

func main() {
	rows := utils.GetUserInput("Enter amount of horizontal slices: ")
	cols := utils.GetUserInput("Enter amount of vertical slices: ")

	fmt.Println()

	start := time.Now()

	inputDir, err := os.ReadDir(input)
	if err != nil {
		fmt.Println("Error reading input directory:", err)
		return
	}

	if len(inputDir) == 0 {
		fmt.Println("Source directory is empty")
		return
	}

	for _, image := range inputDir {
		if !image.IsDir() {

			imageName := image.Name()

			utils.SliceImage(fmt.Sprintf("%s/%s", input, imageName), output, rows, cols)
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
