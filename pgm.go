package Netpbm

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type PGM struct {
	data          [][]uint8
	width, height int
	magicNumber   string
	max           uint8
}

// Function to read a PGM image.
func ReadPGM(filename string) (*PGM, error) {
	var dimension string

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}
	// Close the file at the end of the function.
	defer file.Close()

	scanner := bufio.NewScanner(file)

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
		return nil, err
	}

	// Read the magic number
	scanner.Scan()
	magicNumber := scanner.Text()

	// Skip comments
	for scanner.Scan() {
		if scanner.Text()[0] == '#' {
			continue
		}
		break
	}

	// Read dimensions
	dimension = scanner.Text()
	res := strings.Split(dimension, " ")
	height, _ := strconv.Atoi(res[1])
	width, _ := strconv.Atoi(res[0])

	// Read maximum pixel value
	scanner.Scan()
	max, _ := strconv.Atoi(scanner.Text())
	data := make([][]uint8, height)
	for i := range data {
		data[i] = make([]uint8, width)
	}

	// Read pixel values
	if magicNumber == "P2" {
		for i := 0; i < height; i++ {
			scanner.Scan()
			line := scanner.Text()
			caseBytes := strings.Fields(line)

			for j := 0; j < width; j++ {
				caseInt, _ := strconv.Atoi(caseBytes[j])
				data[i][j] = uint8(caseInt)
			}
		}
	}

	return &PGM{
		data:        data,
		width:       width,
		height:      height,
		magicNumber: magicNumber,
		max:         uint8(max),
	}, nil
}

// Function that returns the width and height of the PGM image.
func (pgm *PGM) Size() (int, int) {
	return pgm.width, pgm.height
}

// Function that returns the value of a pixel at the specified coordinates.
func (pgm *PGM) At(x, y int) uint8 {
	return pgm.data[x][y]
}

// Function that changes the value of a pixel at the specified coordinates.
func (pgm *PGM) Set(x, y int, value uint8) {
	pgm.data[x][y] = value
}

// Function that saves a PGM image.
func (pgm *PGM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file")
		return err
	}
	// Close the file at the end of the function.
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%v\n%v %v\n%v\n", pgm.magicNumber, pgm.width, pgm.height, pgm.max))
	if err != nil {
		fmt.Println("Error writing to file")
		return err
	}

	// Write pixel values
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			var pixel uint8
			pixel = pgm.data[i][j]
			_, err = file.WriteString(fmt.Sprintf("%v ", pixel))
			if err != nil {
				fmt.Println("Error writing to file")
				return err
			}
		}
		_, err = file.WriteString(fmt.Sprintf("\n"))
	}

	return nil
}

// Function that inverts the colors of the image.
func (pgm *PGM) Invert() {
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			pgm.data[i][j] = pgm.max - pgm.data[i][j]
		}
	}
}

// Function that horizontally flips the image.
func (pgm *PGM) Flip() {
	for i := 0; i < pgm.height; i++ {
		for j, k := 0, pgm.width-1; j < k; j, k = j+1, k-1 {
			pgm.data[i][j], pgm.data[i][k] = pgm.data[i][k], pgm.data[i][j]
		}
	}
}

// Function that vertically flips the image.
func (pgm *PGM) Flop() {
	for i, j := 0, pgm.height-1; i < j; i, j = i+1, j-1 {
		pgm.data[i], pgm.data[j] = pgm.data[j], pgm.data[i]
	}
}

// Function that changes the magic number of the image.
func (pgm *PGM) SetMagicNumber(magicNumber string) {
	pgm.magicNumber = magicNumber
}

// Function that sets a new maximum value for pixel intensity.
func (pgm *PGM) SetMaxValue(maxValue uint8) {
	// Set the multiplier
	multiplier := float64(maxValue) / float64(pgm.max)
	// Update the maximum value
	pgm.max = maxValue

	// Update pixel values with the new maximum value
	for i := range pgm.data {
		for j := range pgm.data[i] {
			// Modify the value of each pixel proportionally
			pgm.data[i][j] = uint8(float64(pgm.data[i][j]) * multiplier)
		}
	}
}

// Function that rotates the image 90 degrees clockwise.
func (pgm *PGM) Rotate90CW() {
	rotated := PGM{
		data:        make([][]uint8, pgm.width),
		width:       pgm.height,
		height:      pgm.width,
		magicNumber: pgm.magicNumber,
		max:         pgm.max,
	}

	for i := range rotated.data {
		rotated.data[i] = make([]uint8, rotated.width)
	}

	// Rotate pixel values
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			rotated.data[j][rotated.height-i-1] = pgm.data[i][j]
		}
	}

	// Update the original PGM image
	pgm.data, pgm.width, pgm.height = rotated.data, rotated.width, rotated.height
}

// Function that converts the PGM image to PBM format.
func (pgm *PGM) ToPBM() *PBM {
	var data [][]bool
	data = make([][]bool, pgm.height)
	for i := 0; i < pgm.height; i++ {
		data[i] = make([]bool, pgm.width)
	}

	// Convert pixel values to binary
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			if pgm.data[i][j] > pgm.max/2 {
				data[i][j] = true
			} else {
				data[i][j] = false
			}
		}
	}

	return &PBM{
		data:        data,
		width:       pgm.width,
		height:      pgm.height,
		magicNumber: "P1",
	}
}
