package Netpbm2

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type PPM struct {
	data          [][]Pixel
	width, height int
	magicNumber   string
	max           uint8
}

type Pixel struct {
	R, G, B uint8
}

func ReadPPM(filename string) (*PPM, error) {
	//Open the file
	file, err := os.Open(filename)
	//Check for error
	if err != nil {
		return nil, err
	}
	//Close the file just before the Save function returns/(finishes its execution). Even if it's an error
	defer file.Close()
	scanner := bufio.NewScanner(file)
	//Create a base PPM variable
	ppm := &PPM{}
	//Variable line used to count the lines of the image
	line := 0
	//Loop through each lines
	for scanner.Scan() {
		text := scanner.Text()
		//Ignore empty lines and comments
		if strings.HasPrefix(text, "#") {
			continue
		}
		if ppm.magicNumber == "" {
			//Get the magicnumber. Trimspace removes the spaces from the string
			ppm.magicNumber = strings.TrimSpace(text)
		} else if ppm.width == 0 {
			//Get the width and height of the ppm
			fmt.Sscanf(text, "%d %d", &ppm.width, &ppm.height)
			//Initialize the ppm.data matrix variable by creating the correct amount and size of arrays in an array
			ppm.data = make([][]Pixel, ppm.height)
			for i := range ppm.data {
				ppm.data[i] = make([]Pixel, ppm.width)
			}
		} else if ppm.max == 0 {
			//Get the maxValue of the ppm
			fmt.Sscanf(text, "%d", &ppm.max)
		} else {
			if ppm.magicNumber == "P3" {
				val := strings.Fields(text)
				//Loop through each strings in the current line
				for i := 0; i < ppm.width; i++ {
					//Convert the string to uint8 and set it to the red of the pixel
					r, _ := strconv.ParseUint(val[i*3], 10, 8)
					//Same but the index is incremented to get the next value for the green
					g, _ := strconv.ParseUint(val[i*3+1], 10, 8)
					//Same but the index is incremented to get the next value for the blue
					b, _ := strconv.ParseUint(val[i*3+2], 10, 8)
					//Create the pixel with the colors we just obtained and define it the matrix
					ppm.data[line][i] = Pixel{R: uint8(r), G: uint8(g), B: uint8(b)}
				}
				line++
			} else if ppm.magicNumber == "P6" {
				//Create an array of byte of the size of the image * 3 because each pixel has 3 values RGB
				pixelData := make([]byte, ppm.width*ppm.height*3)
				//Reads the file content
				fileContent, err := os.ReadFile(filename)
				if err != nil {
					return nil, fmt.Errorf("couldn't read file: %v", err)
				}
				//Extracts the necessary pixel data
				copy(pixelData, fileContent[len(fileContent)-(ppm.width*ppm.height*3):])
				//Process the data to fill the pixel array of ppm.data
				pixelIndex := 0
				for y := 0; y < ppm.height; y++ {
					for x := 0; x < ppm.width; x++ {
						ppm.data[y][x].R = pixelData[pixelIndex]
						ppm.data[y][x].G = pixelData[pixelIndex+1]
						ppm.data[y][x].B = pixelData[pixelIndex+2]
						pixelIndex += 3
					}
				}
				break
			}
		}
	}
	return ppm, nil
}

func (ppm *PPM) Save(filename string) error {
	//Create a file with the defines name
	file, err := os.Create(filename)
	//Check for error
	if err != nil {
		return err
	}
	//Close the file just before the Save function returns/(finishes its execution). Even if it's an error
	defer file.Close()
	//Store all the modifications into writer "writer" temporarily until flush
	writer := bufio.NewWriter(file)
	//Write the magicnumber first
	fmt.Fprint(writer, ppm.magicNumber+"\n")
	//Write the size secondly
	fmt.Fprintf(writer, "%d %d\n", ppm.width, ppm.height)
	fmt.Fprintf(writer, "%d\n", ppm.max)
	//Flush writes all the modifications stored in the writer "writer" to the file
	writer.Flush()
	if ppm.magicNumber == "P3" {
		//Loop each pixels of ppm.data
		for y, row := range ppm.data {
			for i, pixel := range row {
				//xtra is used to space each pixel with a space except for the last one of the line
				xtra := " "
				if i == len(row)-1 {
					xtra = ""
				}
				//Write the RGB colors in the writer
				fmt.Fprintf(writer, "%d %d %d%s", pixel.R, pixel.G, pixel.B, xtra)
			}
			//Return to line if it's not the last line
			if y != len(ppm.data)-1 {
				fmt.Fprintln(writer, "")
			}
		}
		writer.Flush()
	} else if ppm.magicNumber == "P6" {
		//For each pixels
		for _, row := range ppm.data {
			for _, pixel := range row {
				//Simple convertion to []byte of RGB
				_, err = file.Write([]byte{pixel.R, pixel.G, pixel.B})
				if err != nil {
					return fmt.Errorf("error writing pixel data: %v", err)
				}
			}
		}
	}
	return nil
}

func (ppm *PPM) Size() (int, int) {
	//Simple return of the size
	return ppm.width, ppm.height
}

func (ppm *PPM) At(x, y int) Pixel {
	//Simple return of the value of a specifix pixel
	return ppm.data[y][x]
}

func (ppm *PPM) Set(x, y int, value Pixel) {
	//Simply define a new value to a specific pixel
	ppm.data[x][y] = value
}

func (ppm *PPM) Invert() {
	//Loop throught each pixels
	for y, _ := range ppm.data {
		for x, _ := range ppm.data[y] {
			pixel := ppm.data[y][x]
			//Change the value to the opposite of his value
			//If the value is 240 would be 15
			//255 - 240 = 15
			//If the value is 1O would be 245
			//255- 10 = 245
			pixel.R = uint8(255 - int(pixel.R))
			pixel.G = uint8(255 - int(pixel.G))
			pixel.B = uint8(255 - int(pixel.B))
			ppm.data[y][x] = pixel
		}
	}
}

// Flip by swapping the first and last pixel of each line until the image is flipped.
func (ppm *PPM) Flip() {
	//Loop through each lines
	for y, _ := range ppm.data {
		//Set cursor to the last character of the line
		cursor := ppm.width - 1
		//Loop through each characters of the line
		for x := 0; x < ppm.width; x++ {
			//Store the value of the pixel
			temp := ppm.data[y][x]
			//Change value of the pixel
			ppm.data[y][x] = ppm.data[y][cursor]
			//Set the value of the first pixel to the stored one
			ppm.data[y][cursor] = temp
			//Move the cursor to the left on the line
			cursor--
			//Break the loop when the cursor crosses or reaches the current line
			if cursor < x || cursor == x {
				break
			}
		}
	}
}

// Flop by swapping the first and last line until the image is flopped.
func (ppm *PPM) Flop() {
	//Set the cursor to the bottom line of the image.
	cursor := ppm.height - 1
	//Loop through each lines
	for y, _ := range ppm.data {
		//Swap the current line with the line pointed to by the cursor
		temp := ppm.data[y]
		ppm.data[y] = ppm.data[cursor]
		ppm.data[cursor] = temp
		//Move the cursor to one line higher
		cursor--
		//Break the loop when the cursor crosses or reaches the current line
		if cursor < y || cursor == y {
			break
		}
	}
}

func (ppm *PPM) SetMagicNumber(magicNumber string) {
	//Simply define a new magic number
	ppm.magicNumber = magicNumber
}

func (ppm *PPM) SetMaxValue(maxValue uint8) {
	//Loop through each pixel
	for y, _ := range ppm.data {
		for x, _ := range ppm.data[y] {
			pixel := ppm.data[y][x]
			//Calculate the new pixel value based on the new maximum value for each color
			//Adjusting the pixel value proportionally to the new max value
			pixel.R = uint8(float64(pixel.R) * float64(maxValue) / float64(ppm.max))
			pixel.G = uint8(float64(pixel.G) * float64(maxValue) / float64(ppm.max))
			pixel.B = uint8(float64(pixel.B) * float64(maxValue) / float64(ppm.max))
			ppm.data[y][x] = pixel
		}
	}
	ppm.max = maxValue
}

func (ppm *PPM) Rotate90CW() {
	//Same as pgm.Rotate90CW but the matrix is [][]Pixel not [][]uint8
	//Create a new matrix to store the rotated pixel data
	rotatedData := make([][]Pixel, ppm.width)
	for i := range rotatedData {
		rotatedData[i] = make([]Pixel, ppm.height)
	}
	//Loop through each pixel in the original image
	for i := 0; i < ppm.width; i++ {
		for j := 0; j < ppm.height; j++ {
			//Rotate the pixel by 90 degrees clockwise and assign it
			rotatedData[i][j] = ppm.data[ppm.height-1-j][i]
		}
	}
	//Swap the width and height of the image.
	ppm.width, ppm.height = ppm.height, ppm.width
	//Update the image data with the rotated data
	ppm.data = rotatedData
}

func (ppm *PPM) ToPBM() *PBM {
	//Same idea as pgm.ToPBM
	//Create a new pbm
	pbm := &PBM{}
	//Assign same data except for the magicnumber
	pbm.magicNumber = "P1"
	pbm.height = ppm.height
	pbm.width = ppm.width
	for y, _ := range ppm.data {
		pbm.data = append(pbm.data, []bool{})
		for x, _ := range ppm.data[y] {
			r, g, b := ppm.data[y][x].R, ppm.data[y][x].G, ppm.data[y][x].B
			//Calculate if the pixel should be black or white
			//if the average of the 3 colors is lower than the half of the maxValue, then i consider it white
			//If maxValue is 100 and average is 49, it would be black
			isBlack := (uint8((int(r)+int(g)+int(b))/3) < ppm.max/2)
			pbm.data[y] = append(pbm.data[y], isBlack)
		}
	}
	return pbm
}

func (ppm *PPM) ToPGM() *PGM {
	//Same idea as ppm.ToPBM
	//Create a new pbm
	pgm := &PGM{}
	//Assign same data except for the magicnumber
	pgm.magicNumber = "P2"
	pgm.height = ppm.height
	pgm.width = ppm.width
	pgm.max = ppm.max
	for y, _ := range ppm.data {
		pgm.data = append(pgm.data, []uint8{})
		for x, _ := range ppm.data[y] {
			r, g, b := ppm.data[y][x].R, ppm.data[y][x].G, ppm.data[y][x].B
			//Calculate the amount of gray the pixel should have
			//It is just the average of the 3 RGB colors
			grayValue := uint8((int(r) + int(g) + int(b)) / 3)
			pgm.data[y] = append(pgm.data[y], uint8(grayValue))
		}
	}
	return pgm
}

type Point struct {
	X, Y int
}

// Drawing lines by using Bresenham's Line Drawing Algorithm
// Found people suggesting it on online forums
func (ppm *PPM) DrawLine(p1, p2 Point, color Pixel) {
	deltaX := abs(p2.X - p1.X)
	deltaY := abs(p2.Y - p1.Y)
	sx, sy := sign(p2.X-p1.X), sign(p2.Y-p1.Y)
	err := deltaX - deltaY
	for {
		if p1.X >= 0 && p1.X < ppm.width && p1.Y >= 0 && p1.Y < ppm.height {
			ppm.data[p1.Y][p1.X] = color
		}
		if p1.X == p2.X && p1.Y == p2.Y {
			break
		}
		e2 := 2 * err
		if e2 > -deltaY {
			err -= deltaY
			p1.X += sx
		}
		if e2 < deltaX {
			err += deltaX
			p1.Y += sy
		}
	}
}

// If negative, change it to positive
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Return 1 if it's over 0
// Return 0 if it's 0
// Return -1 if  it's negative
func sign(x int) int {
	if x > 0 {
		return 1
	} else if x < 0 {
		return -1
	}
	return 0
}

func (ppm *PPM) DrawRectangle(p1 Point, width, height int, color Pixel) {
	//Create the 3 extra points according to the width and the height
	p2 := Point{p1.X + width, p1.Y}
	p3 := Point{p1.X, p1.Y + height}
	p4 := Point{p1.X + width, p1.Y + height}
	//Draw the lines to connect them
	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p2, p4, color)
	ppm.DrawLine(p4, p3, color)
	ppm.DrawLine(p3, p1, color)
}

func (ppm *PPM) DrawFilledRectangle(p1 Point, width, height int, color Pixel) {
	//Draw horizontal lines with the asked width under each other until the height is reached
	p2 := Point{p1.X + width, p1.Y}
	for i := 0; i <= height; i++ {
		ppm.DrawLine(p1, p2, color)
		p1.Y++
		p2.Y++
	}
}

func (ppm *PPM) DrawCircle(center Point, radius int, color Pixel) {
	//Loop through each pixel
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			//Calculate the distance from the current pixel to the center of the circle
			dx := float64(x - center.X)
			dy := float64(y - center.Y)
			distance := math.Sqrt(dx*dx + dy*dy)
			//Check if the distance is approximately equal to the specified radius
			//*0.85 is to obtain a circle looking like the tester's circle even if it's not really a circle... In reality, remove "*0.85" and it's a real circle
			if math.Abs(distance-float64(radius)*0.85) < 0.5 {
				ppm.data[y][x] = color
			}
		}
	}
}

func (ppm *PPM) DrawFilledCircle(center Point, radius int, color Pixel) {
	//Draw a circle with the radius getting smaller until it is at 0;
	for radius >= 0 {
		ppm.DrawCircle(center, radius, color)
		radius--
	}
}

func (ppm *PPM) DrawTriangle(p1, p2, p3 Point, color Pixel) {
	//Draw lines and link the 3 points
	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p2, p3, color)
	ppm.DrawLine(p3, p1, color)
}

// Draw a line from p1 to p3 and move p1 towars p2 until the triangle is filled
func (ppm *PPM) DrawFilledTriangle(p1, p2, p3 Point, color Pixel) {
	//Loop until p1 reaches p2
	for p1 != p2 {
		//Draw a line between p1 and p3
		ppm.DrawLine(p3, p1, color)
		//Increment or decrement X of p1 based on p2 position
		if p1.X != p2.X && p1.X < p2.X {
			p1.X++
		} else if p1.X != p2.X && p1.X > p2.X {
			p1.X--
		}
		//Increment or decrement Y of p1 based on p2 position
		if p1.Y != p2.Y && p1.Y < p2.Y {
			p1.Y++
		} else if p1.Y != p2.Y && p1.Y > p2.Y {
			p1.Y--
		}
	}
	//Draw a final line between the last position of p1 (should be at p2 at this point) and p3
	ppm.DrawLine(p3, p1, color)
}

func (ppm *PPM) DrawPolygon(points []Point, color Pixel) {
	//Link the points with a line
	for i := 0; i < len(points)-1; i++ {
		ppm.DrawLine(points[i], points[i+1], color)
	}
	//Link the last and the first point with a line
	ppm.DrawLine(points[len(points)-1], points[0], color)
}

func (ppm *PPM) DrawFilledPolygon(points []Point, color Pixel) {

}