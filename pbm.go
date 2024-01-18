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

/*func main() {
	ReadPBM("Sans_titreascii.pbm")
	Save("dougito.pbm")

}*/

func ReadPBM(filename string) (*PBM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var magicNumber string

	// Fonction pour lire la ligne suivante non commentée
	readNextLine := func() (string, error) {
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			// Ignorer les lignes vides ou les lignes commençant par "#"
			if line != "" && !strings.HasPrefix(line, "#") {
				return line, nil
			}
		}
		return "", scanner.Err()
	}

	// Lire la première ligne non commentée pour obtenir le nombre magique
	if magicNumber, err = readNextLine(); err != nil {
		return nil, err
	}

	if magicNumber != "P1" && magicNumber != "P4" {
		return nil, errors.New("unsupported file type")
	}

	// Lire les dimensions
	dimensions, err := readNextLine()
	if err != nil {
		return nil, err
	}

	dimTokens := strings.Fields(dimensions)
	if len(dimTokens) != 2 {
		return nil, errors.New("invalid image dimensions")
	}

	width, _ := strconv.Atoi(dimTokens[0])
	height, _ := strconv.Atoi(dimTokens[1])

	var data [][]bool

	// If the image is empty, initialize data with an empty slice
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
		} else if magicNumber == "P4" {
			// Calculate the number of padding bits
			paddingBits := (8 - width%8) % 8

			// Calculate the number of bytes per row, considering padding
			bytesPerRow := (width + paddingBits + 7) / 8

			// Create a buffer to read binary data
			buffer := make([]byte, bytesPerRow)
			for i := 0; i < height; i++ {
				_, err := file.Read(buffer)
				if err != nil {
					return nil, err
				}

				// Process the bits from the buffer
				for j := 0; j < width; j++ {
					// Get the byte containing the bit
					byteIndex := j / 8
					bitIndex := 7 - (j % 8)
					bit := (buffer[byteIndex] >> bitIndex) & 1
					data[i][j] = bit == 1
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

func (pbm *PBM) Size() (int, int) {
	return pbm.width, pbm.height
}

func (pbm *PBM) At(x, y int) bool {
	return pbm.data[x][y]
}
func (pbm *PBM) Set(x, y int, value bool) {
	pbm.data[x][y] = value
}

func (pbm *PBM) Save(filename string) error {
	fichier, err := os.Create(filename)
	if err != nil {
		fmt.Println("Erreur lors de la création du fichier")
		return err
	}

	_, err = fichier.WriteString(fmt.Sprintf("%v\n%v %v\n", pbm.magicNumber, pbm.width, pbm.height))
	if err != nil {
		fmt.Println("Erreur lors de l'écriture dans le fichier")
		return err
	}

	for i := 0; i < pbm.height; i++ {
		for j := 0; j < pbm.width; j++ {
			var pixel int
			if pbm.data[i][j] {
				pixel = 1
			} else {
				pixel = 0
			}
			_, err = fichier.WriteString(fmt.Sprintf("%v", pixel))
			if err != nil {
				fmt.Println("Erreur lors de l'écriture dans le fichier")
				return err
			}
		}
		_, err = fichier.WriteString("\n")
		if err != nil {
			fmt.Println("Erreur lors de l'écriture dans le fichier")
			return err
		}
	}

	return nil
}

func (pbm *PBM) Invert() {
	for i := 0; i < pbm.height; i++ {
		for j := 0; j < pbm.width; j++ {
			if pbm.data[i][j] {
				pbm.data[i][j] = false
			} else {
				pbm.data[i][j] = true
			}
		}
	}
}
func (pbm *PBM) Flip() {
	for i := 0; i < pbm.height; i++ {
		for j, k := 0, pbm.width-1; j < k; j, k = j+1, k-1 {
			pbm.data[i][j], pbm.data[i][k] = pbm.data[i][k], pbm.data[i][j]
		}
	}
}
func (pbm *PBM) Flop() {
	for i, j := 0, pbm.height-1; i < j; i, j = i+1, j-1 {
		pbm.data[i], pbm.data[j] = pbm.data[j], pbm.data[i]
	}
}

func (pbm *PBM) SetMagicNumber(magicNumber string) {
	pbm.magicNumber = magicNumber
}
