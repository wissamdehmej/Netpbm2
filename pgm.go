package Netpbm2

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// PGM represents a PGM image
type PGM struct {
	data          [][]uint8
	width, height int
	magicNumber   string
	max           uint8
}

// ReadPGM reads a PGM image from a file and returns a struct that represents the image.
func ReadPGM(filename string) (*PGM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	// Read the magic number
	magicNumber, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading magic number: %v", err)
	}
	magicNumber = strings.TrimSpace(magicNumber)
	if magicNumber != "P2" && magicNumber != "P5" {
		return nil, fmt.Errorf("invalid magic number: %s", magicNumber)
	}

	// Read dimensions
	dimensions, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading dimensions: %v", err)
	}
	var width, height int
	_, err = fmt.Sscanf(strings.TrimSpace(dimensions), "%d %d", &width, &height)
	if err != nil {
		return nil, fmt.Errorf("invalid dimensions: %v", err)
	}

	// Read max value
	var max int
	_, err = fmt.Fscanf(reader, "%d\n", &max)
	if err != nil {
		return nil, fmt.Errorf("error reading max value: %v", err)
	}

	data := make([][]uint8, height)

	if magicNumber == "P2" {
		// Read P2 format (ASCII)
		for y := 0; y < height; y++ {
			data[y] = make([]uint8, width)
			for x := 0; x < width; x++ {
				var pixelValue int
				_, err := fmt.Fscanf(reader, "%d", &pixelValue)
				if err != nil {
					return nil, fmt.Errorf("error reading pixel data at position (%d, %d): %v", x, y, err)
				}
				data[y][x] = uint8(pixelValue)
			}
			// Read end of line
			_, _ = reader.ReadString('\n')
		}
	} else if magicNumber == "P5" {
		// Read P5 format (binary)
		for y := 0; y < height; y++ {
			data[y] = make([]uint8, width)
			err := binary.Read(reader, binary.BigEndian, &data[y])
			if err != nil {
				return nil, fmt.Errorf("error reading pixel data: %v", err)
			}
		}
	}

	return &PGM{data, width, height, magicNumber, uint8(max)}, nil
}

// Size returns the width and height of the PGM image.
func (pgm *PGM) Size() (int, int) {
	return pgm.width, pgm.height
}

// At returns the value of the pixel at the specified (x, y) coordinates.
func (pgm *PGM) At(x, y int) uint8 {
	if x >= 0 && x < pgm.width && y >= 0 && y < pgm.height {
		return pgm.data[y][x]
	}
	// You can choose how to handle out-of-bounds access, for example, return 0 or another default value.
	return 0
}

// Set sets the value of the pixel at the specified (x, y) coordinates.
func (pgm *PGM) Set(x, y int, value uint8) {
	if x >= 0 && x < pgm.width && y >= 0 && y < pgm.height {
		pgm.data[y][x] = value
	}
	// You can choose how to handle out-of-bounds access, for example, do nothing or log a message.
}

// Save saves the PGM image to a file in the same format as the original image.
func (pgm *PGM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	fmt.Fprint(writer, pgm.magicNumber+"\n")
	fmt.Fprintf(writer, "%d %d\n", pgm.width, pgm.height)
	fmt.Fprintf(writer, "%d\n", pgm.max)
	writer.Flush()
	if pgm.magicNumber == "P2" {
		for y, row := range pgm.data {
			for i, pixel := range row {
				xtra := " "
				if i == len(row)-1 {
					xtra = ""
				}
				//Here i convert uint8 to an int in order to finally convert it to a string
				fmt.Fprint(writer, strconv.Itoa(int(pixel))+xtra)
			}
			if y != len(pgm.data)-1 {
				fmt.Fprintln(writer, "")
			}
		}
		writer.Flush()
	} else if pgm.magicNumber == "P5" {
		for _, row := range pgm.data {
			for _, pixel := range row {
				//We can simply convert it to []byte
				_, err = file.Write([]byte{pixel})
				if err != nil {
					return fmt.Errorf("error writing pixel data: %v", err)
				}
			}
		}
	}
	return nil
}

// Invert inverts the colors of the PGM image.
func (pgm *PGM) Invert() {
	maxUint8 := uint8(pgm.max)

	for y := 0; y < pgm.height; y++ {
		for x := 0; x < pgm.width; x++ {
			// Inversion de la valeur du pixel
			pgm.data[y][x] = maxUint8 - pgm.data[y][x]
		}
	}
}

// Flip flips the PGM image horizontally.
func (pgm *PGM) Flip() {
	for y := 0; y < pgm.height; y++ {
		left := 0
		right := pgm.width - 1

		// Inverser les valeurs des pixels de gauche à droite
		for left < right {
			pgm.data[y][left], pgm.data[y][right] = pgm.data[y][right], pgm.data[y][left]
			left++
			right--
		}
	}
}

// Flop flops the PGM image vertically.
func (pgm *PGM) Flop() {
	top := 0
	bottom := pgm.height - 1

	for top < bottom {
		// Inverser les valeurs des pixels de haut en bas pour chaque colonne
		for x := 0; x < pgm.width; x++ {
			pgm.data[top][x], pgm.data[bottom][x] = pgm.data[bottom][x], pgm.data[top][x]
		}

		top++
		bottom--
	}
}

// SetMagicNumber sets the magic number of the PGM image.
func (pgm *PGM) SetMagicNumber(magicNumber string) {
	pgm.magicNumber = magicNumber
}

// SetMaxValue sets the max value of the PGM image.
func (pgm *PGM) SetMaxValue(maxValue uint8) {
	for y, _ := range pgm.data {
		for x, _ := range pgm.data[y] {
			prevValue := pgm.data[y][x]
			// Effectuez la multiplication avant la division et convertissez le type après la multiplication
			newValue := uint8((uint(prevValue) * 5) / uint(maxValue))
			pgm.data[y][x] = newValue
		}
	}
	// Mettez à jour la valeur maximale dans la structure PGM
	pgm.max = uint8(maxValue)
}

// Rotate90CW rotates the PGM image 90° clockwise.
func (pgm *PGM) Rotate90CW() {
	// Créer une nouvelle image avec les dimensions inversées
	rotatedData := make([][]uint8, pgm.width)
	for x := range rotatedData {
		rotatedData[x] = make([]uint8, pgm.height)
	}

	// Remplir la nouvelle image en effectuant la rotation
	for y := 0; y < pgm.height; y++ {
		for x := 0; x < pgm.width; x++ {
			rotatedData[x][pgm.height-y-1] = pgm.data[y][x]
		}
	}

	// Mettre à jour les dimensions de l'image actuelle
	pgm.width, pgm.height = pgm.height, pgm.width

	// Mettre à jour les données avec la nouvelle image pivotée
	pgm.data = rotatedData
}

// ToPBM converts the PGM image to PBM.
func (pgm *PGM) ToPBM() *PBM {
	// Créer une nouvelle instance de la struct PBM
	pbmInstance := &PBM{
		data:        make([][]bool, pgm.height),
		width:       pgm.width,
		height:      pgm.height,
		magicNumber: "P1",
	}

	// Remplir les données de la struct PBM en fonction des valeurs de l'image PGM
	for y := 0; y < pgm.height; y++ {
		pbmInstance.data[y] = make([]bool, pgm.width)
		for x := 0; x < pgm.width; x++ {
			// Convertir la valeur du pixel en bool (noir ou blanc)
			pbmInstance.data[y][x] = pgm.data[y][x] > uint8(pgm.max/2)
		}
	}

	return pbmInstance
}