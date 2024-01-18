package Netpbm

import (
	"bufio"
	"fmt"
	"log"
	"math"
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

/*func main() {
	// Example usage
	pgm, err := ReadPGM("testImages\\pgm\\testP2.pgm")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	pgm.Rotate90CW()
	fmt.Println(pgm.data)
}*/

// ReadPGM reads a PGM image from a file and returns a struct that represents the image.
func ReadPGM(filename string) (*PGM, error) {
	var dimension string

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	if error := scanner.Err(); error != nil {
		log.Fatalln(error)
		return nil, error
	}

	scanner.Scan()

	magicNumber := scanner.Text()

	for scanner.Scan() {
		if scanner.Text()[0] == '#' {
			continue
		}
		break

	}

	dimension = scanner.Text()
	res := strings.Split(dimension, " ")
	height, _ := strconv.Atoi(res[1])
	width, _ := strconv.Atoi(res[0])

	scanner.Scan()

	max, _ := strconv.Atoi(scanner.Text())
	data := make([][]uint8, height)
	for i := range data {
		data[i] = make([]uint8, width)
	}
	if magicNumber == "P2" {
		for i := 0; i < height; i++ {
			scanner.Scan()
			line := scanner.Text()
			casebyte := strings.Fields(line)

			for j := 0; j < width; j++ {
				caseint, _ := strconv.Atoi(casebyte[j])
				data[i][j] = uint8(caseint)

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

func (pgm *PGM) Size() (int, int) {
	return pgm.width, pgm.height
}

func (pgm *PGM) At(x, y int) uint8 {
	return pgm.data[x][y]
}

func (pgm *PGM) Set(x, y int, value uint8) {
	pgm.data[x][y] = value
}

func (pgm *PGM) Save(filename string) error {
	fichier, err := os.Create(filename)
	if err != nil {
		fmt.Println("Erreur lors de la création du fichier")
		return err
	}

	_, err = fichier.WriteString(fmt.Sprintf("%v\n%v %v\n%v\n", pgm.magicNumber, pgm.width, pgm.height, pgm.max))
	if err != nil {
		fmt.Println("Erreur lors de l'écriture dans le fichier")
		return err
	}
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			var pixel uint8
			pixel = pgm.data[i][j]
			_, err = fichier.WriteString(fmt.Sprintf("%v ", pixel))
			if err != nil {
				fmt.Println("Erreur lors de l'écriture dans le fichier")
				return err
			}
		}
		_, err = fichier.WriteString(fmt.Sprintf("\n"))
	}
	fichier.Close()
	return nil
}

func (pgm *PGM) Invert() {
	for i := 0; i < pgm.width; i++ {
		for j := 0; j < pgm.height; j++ {
			pgm.data[i][j] = uint8(pgm.max) - pgm.data[i][j]
		}
	}
}

func (pgm *PGM) Flip() {
	for i := 0; i < pgm.height; i++ {
		for j, k := 0, pgm.width-1; j < k; j, k = j+1, k-1 {
			pgm.data[i][j], pgm.data[i][k] = pgm.data[i][k], pgm.data[i][j]
		}
	}
}

func (pgm *PGM) Flop() {
	for i, j := 0, pgm.height-1; i < j; i, j = i+1, j-1 {
		pgm.data[i], pgm.data[j] = pgm.data[j], pgm.data[i]
	}
}

func (pgm *PGM) SetMagicNumber(magicNumber string) {
	pgm.magicNumber = magicNumber
}

func (pgm *PGM) SetMaxValue(maxValue uint8) {
	pgm.max = maxValue
	for i := range pgm.data {
		for j := range pgm.data[i] {
			pgm.data[i][j] = uint8(math.Round(float64(pgm.data[i][j]) / float64(pgm.max) * 2))
			if pgm.data[i][j] == 4 {
				pgm.data[i][j] += 1
			}
		}
	}
}

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

	// Faire tourner les valeurs des pixels
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			rotated.data[j][rotated.height-i-1] = pgm.data[i][j]
		}
	}

	// Mettre à jour l'image PGM originale
	pgm.data, pgm.width, pgm.height = rotated.data, rotated.width, rotated.height
}

func (pgm *PGM) ToPBM() *PBM {

}
