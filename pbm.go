package Netpbm

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PBM struct {
	data        [][]bool
	width       int
	height      int
	magicNumber string
}

// Function to read a PBM image.
func ReadPBM(filename string) (*PBM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var magicNumber string

	// Function to read each line one by one.
	readNextLine := func() (string, error) {
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			// Skip potential comments.
			if line != "" && !strings.HasPrefix(line, "#") {
				return line, nil
			}
		}
		return "", scanner.Err()
	}

	// Read the first uncommented line to get the magic number.
	if magicNumber, err = readNextLine(); err != nil {
		return nil, err
	}

	if magicNumber != "P1" && magicNumber != "P4" {
		return nil, errors.New("unsupported file type")
	}

	dimensions, err := readNextLine()
	if err != nil {
		return nil, err
	}
	// Read the second uncommented line to get width and height.
	dimTokens := strings.Fields(dimensions)
	if len(dimTokens) != 2 {
		return nil, errors.New("invalid image dimensions")
	}

	width, _ := strconv.Atoi(dimTokens[0])
	height, _ := strconv.Atoi(dimTokens[1])

	var data [][]bool

	// If the image is empty, initialize data with an empty slice.
	if width > 0 && height > 0 {
		data = make([][]bool, height)
		for i := range data {
			data[i] = make([]bool, width)
		}

		if magicNumber == "P1" {
			for i := 0; i < height; i++ {
				line, err := readNextLine()
				if err != nil {
					return nil, err
				}

				tokens := strings.Fields(line)
				for j, token := range tokens {
					pixel, err := strconv.Atoi(token)
					if err != nil {
						return nil, err
					}
					data[i][j] = pixel == 1
				}
			}
		}
	}

	return &PBM{
		data:        data,
		width:       width,
		height:      height,
		magicNumber: magicNumber,
	}, nil
}

// Function that returns the width and height of the PBM image.
func (pbm *PBM) Size() (int, int) {
	return pbm.width, pbm.height
}

// Function that returns the value of a pixel at the specified coordinates.
func (pbm *PBM) At(x, y int) bool {
	return pbm.data[x][y]
}

// Function that changes the value of a pixel at the specified coordinates.
func (pbm *PBM) Set(x, y int, value bool) {
	pbm.data[x][y] = value
}

// Function that saves a PBM image.
func (pbm *PBM) Save(filename string) error {
	fileSave, err := os.Create(filename)
	if err != nil {
		return err
	}
	// Close the file at the end of the function.
	defer fileSave.Close()

	fmt.Fprintf(fileSave, "%s\n%d %d\n", pbm.magicNumber, pbm.width, pbm.height)

	if pbm.magicNumber == "P1" {
		for _, row := range pbm.data {
			for _, pixel := range row {
				if pixel {
					fmt.Fprint(fileSave, "1 ")
				} else {
					fmt.Fprint(fileSave, "0 ")
				}
			}
			fmt.Fprintln(fileSave)
		}
	}
	return nil
}

// Function that inverts the colors of the image.
func (pbm *PBM) Invert() {
	for i := 0; i < pbm.height; i++ {
		for j := 0; j < pbm.width; j++ {
			pbm.data[i][j] = !pbm.data[i][j]
		}
	}
}

// Function that horizontally flips the image.
func (pbm *PBM) Flip() {
	for i := 0; i < pbm.height; i++ {
		for j, k := 0, pbm.width-1; j < k; j, k = j+1, k-1 {
			pbm.data[i][j], pbm.data[i][k] = pbm.data[i][k], pbm.data[i][j]
		}
	}
}

// Function that vertically flips the image.
func (pbm *PBM) Flop() {
	for i, j := 0, pbm.height-1; i < j; i, j = i+1, j-1 {
		pbm.data[i], pbm.data[j] = pbm.data[j], pbm.data[i]
	}
}

// Function that changes the magic number of the image.
func (pbm *PBM) SetMagicNumber(magicNumber string) {
	pbm.magicNumber = magicNumber
}
