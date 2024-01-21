package Netpbm2

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Pixel struct {
	R, G, B uint8
}

type PPM struct {
	data          [][]Pixel
	width, height int
	magicNumber   string
	max           uint8
}

func ReadPPM(filename string) (*PPM, error) {
	//Open the file
	file, err := os.Open(filename)
	//Check the potentiel error
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	ppm := &PPM{}
	line := 0
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "#") {
			continue
		}
		if ppm.magicNumber == "" {
			ppm.magicNumber = strings.TrimSpace(text)
		} else if ppm.width == 0 {
			fmt.Sscanf(text, "%d %d", &ppm.width, &ppm.height)
			ppm.data = make([][]Pixel, ppm.height)
			for i := range ppm.data {
				ppm.data[i] = make([]Pixel, ppm.width)
			}
		} else if ppm.max == 0 {
			fmt.Sscanf(text, "%d", &ppm.max)
		} else {
			if ppm.magicNumber == "P3" {
				// If the PPM image format is P3 (ASCII), parse the color values from the current line.
				val := strings.Fields(text)
				// Iterate through each set of color values in the current line.
				for i := 0; i < ppm.width; i++ {
					// Convert the string to uint8 and assign it to the red component of the pixel.
					r, _ := strconv.ParseUint(val[i*3], 10, 8)
					// Increment the index to obtain the next value for the green component.
					g, _ := strconv.ParseUint(val[i*3+1], 10, 8)
					// Increment the index again to get the next value for the blue component.
					b, _ := strconv.ParseUint(val[i*3+2], 10, 8)
					// Create a pixel with the obtained color values and store it in the data matrix.
					ppm.data[line][i] = Pixel{R: uint8(r), G: uint8(g), B: uint8(b)}
				}
				// Move to the next line in the PPM image data matrix.
				line++

			} else if ppm.magicNumber == "P6" {
				//Create an array of byte of the size of the image * 3 because each pixel has 3 values RGB
				pixelData := make([]byte, ppm.width*ppm.height*3)
				fileContent, err := os.ReadFile(filename)
				if err != nil {
					return nil, fmt.Errorf("couldn't read file: %v", err)
				}
				// Copy the pixel data from the end of the file content to the pixelData slice.
				copy(pixelData, fileContent[len(fileContent)-(ppm.width*ppm.height*3):])
				// Initialize an index variable to keep track of the current position in the pixelData slice.
				pixelIndex := 0
				// Iterate through each row in the PPM image.
				for y := 0; y < ppm.height; y++ {
					// Iterate through each pixel in the current row.
					for x := 0; x < ppm.width; x++ {
						// Set the red component of the pixel at position (x, y) from the pixelData slice.
						ppm.data[y][x].R = pixelData[pixelIndex]
						// Increment the index to get the green component of the pixel.
						ppm.data[y][x].G = pixelData[pixelIndex+1]
						// Increment the index again to get the blue component of the pixel.
						ppm.data[y][x].B = pixelData[pixelIndex+2]
						// Move the index to the next set of three values (next pixel).
						pixelIndex += 3
					}
				}
				break
			}
		}
	}
	return ppm, nil
}

// Return size
func (ppm *PPM) Size() (int, int) {
	return ppm.width, ppm.height
}

// Return value a pixel
func (ppm *PPM) At(x, y int) Pixel {
	return ppm.data[y][x]
}

// Define a new value pixel
func (ppm *PPM) Set(x, y int, value Pixel) {
	if x >= 0 && x < ppm.width && y >= 0 && y < ppm.height {
		ppm.data[y][x] = value
	} else {
		fmt.Println("Error: Coordinates out of bounds.")
	}
}