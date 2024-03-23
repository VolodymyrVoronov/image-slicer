package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"image/jpeg"
	"image/png"
)

type Coords struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func GetUserInput(message string) int {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(message)

	scanner.Scan()

	inputString := scanner.Text()
	inputInt, err := strconv.Atoi(inputString)

	if err != nil {
		fmt.Println("Error: ", err)
		fmt.Println("Default value: 1")

		return 1
	}

	if inputInt < 1 {
		fmt.Println("Error: input for rows/columns must be greater than 0")
		fmt.Println("Default value: 1")

		return 1
	}

	return inputInt
}

func GetImageFormat(file string) (string, error) {
	extension := filepath.Ext(file)

	if extension != ".png" && extension != ".jpg" && extension != ".jpeg" {
		return "", fmt.Errorf(fmt.Sprintln("Error: Format has to be png, jpg or jpeg."))
	}

	format := strings.ToLower(extension[1:])

	return format, nil
}

func GetImageName(file string) string {
	name := filepath.Base(file)
	name = strings.TrimSuffix(name, filepath.Ext(file))

	return name
}

func SliceImage(fileOriginal string, outPutDir string, rows int, cols int, doneChannel chan string, errorChannel chan error) {
	file, err := os.Open(fileOriginal)
	if err != nil {
		fmt.Println("Error while opening image: ", err)
		errorChannel <- err
		return
	}
	defer file.Close()

	imageFormat, err := GetImageFormat(fileOriginal)
	if err != nil {
		errorChannel <- err
		return
	}

	imageOriginal, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error while decoding image: ", err)
		errorChannel <- err
		return
	}

	bounds := imageOriginal.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	sliceWidth := width / cols
	sliceHeight := height / rows

	var slicedImageCoords []Coords

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			x0 := j * sliceWidth
			y0 := i * sliceHeight
			x1 := (j + 1) * sliceWidth
			y1 := (i + 1) * sliceHeight

			slice := imageOriginal.(interface {
				SubImage(r image.Rectangle) image.Image
			}).SubImage(image.Rect(x0, y0, x1, y1))

			slicedImageCoords = append(slicedImageCoords, Coords{X: x0, Y: y0})

			formatOfFile, err := GetImageFormat(fileOriginal)
			if err != nil {
				fmt.Println("Error while getting image format: ", err)
				errorChannel <- err
				return
			}

			fileName := fmt.Sprintf("%s-%d%d.%s", GetImageName(fileOriginal), i+1, j+1, formatOfFile)

			slicedImageFile, err := os.Create(filepath.Join(outPutDir, fileName))
			if err != nil {
				fmt.Println("Error while creating image: ", err)
				errorChannel <- err
				return
			}
			defer slicedImageFile.Close()

			if imageFormat == "png" {
				SaveInFormat(slicedImageFile, slice, "png")
			} else if imageFormat == "jpg" || imageFormat == "jpeg" {
				SaveInFormat(slicedImageFile, slice, "jpeg")
			}
		}
	}

	var output string

	output += "["

	for i, coord := range slicedImageCoords {
		output += fmt.Sprintf("{x: %d, y: %d}", coord.X, coord.Y)

		if i < len(slicedImageCoords)-1 {
			output += ", "
		}
	}

	output += "]"

	WriteDataToFileAsJSON(slicedImageCoords, filepath.Join(outPutDir, fmt.Sprintf("%s.json", GetImageName(fileOriginal))))

	doneChannel <- fmt.Sprintf("Coords of %s: %s \n", GetImageName(fileOriginal), output)
}

func SaveInFormat(w *os.File, m image.Image, format string) {
	switch format {
	case "png":
		png.Encode(w, m)

	case "jpeg":
		jpeg.Encode(w, m, &jpeg.Options{Quality: 100})
	}
}

func WriteDataToFileAsJSON(data interface{}, fileDir string) error {
	j, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = os.WriteFile(fileDir, j, 0644)
	if err != nil {
		return err
	}

	return nil
}

func ClearDir(dirPath string) error {
	dir, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	if len(dir) == 0 {
		fmt.Println("Output directory is empty! No need to clear it.")
		fmt.Println()
		return nil
	}

	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			err = os.Remove(path)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Output directory was cleared!")
	fmt.Println()

	return nil
}
