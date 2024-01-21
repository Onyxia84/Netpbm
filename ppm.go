package Netpbm

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// Pixel represents a color pixel with red (R), green (G), and blue (B) components.
type Pixel struct {
	R, G, B uint8
}

// Point represents a 2D point with X and Y coordinates.
type Point struct {
	X, Y int
}

// PPM represents a Portable Pixmap image.
type PPM struct {
	data          [][]Pixel
	width, height int
	magicNumber   string
	max           uint8
}

// ReadPPM reads a PPM image from a file and returns a structure representing the image.
func ReadPPM(filename string) (*PPM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	line := scanner.Text()
	line = strings.TrimSpace(line)
	if line != "P3" && line != "P6" {
		return nil, fmt.Errorf("Not a Portable Pixmap file: bad magic number %s", line)
	}
	magicNumber := line

	// Read dimensions
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "#") {
			continue
		}
		break
	}

	dimension := strings.Fields(scanner.Text())
	width, _ := strconv.Atoi(dimension[0])
	height, _ := strconv.Atoi(dimension[1])

	scanner.Scan()
	maxValue, _ := strconv.Atoi(scanner.Text())

	// Read pixel data
	var ppm *PPM
	if magicNumber == "P3" {
		data := make([][]Pixel, height)
		for i := range data {
			data[i] = make([]Pixel, width)
		}

		// Read pixel data for P3 format
		for i := 0; i < height; i++ {
			scanner.Scan()
			line := scanner.Text()
			values := strings.Fields(line)
			for j := 0; j < width; j++ {
				if j == 0 {
					r, _ := strconv.Atoi(values[0])
					g, _ := strconv.Atoi(values[1])
					b, _ := strconv.Atoi(values[2])
					data[i][j] = Pixel{uint8(r), uint8(g), uint8(b)}
				} else {
					r, _ := strconv.Atoi(values[3*j])
					g, _ := strconv.Atoi(values[3*j+1])
					b, _ := strconv.Atoi(values[3*j+2])
					data[i][j] = Pixel{uint8(r), uint8(g), uint8(b)}
				}
			}
		}

		ppm = &PPM{
			data:        data,
			width:       width,
			height:      height,
			magicNumber: magicNumber,
			max:         uint8(maxValue),
		}
		fmt.Printf("%+v\n", PPM{data, width, height, magicNumber, uint8(maxValue)})
	}
	return ppm, nil
}

// Size returns the width and height of the PPM image.
func (ppm *PPM) Size() (int, int) {
	return ppm.width, ppm.height
}

// At returns the color pixel at the specified coordinates (x, y).
func (ppm *PPM) At(x, y int) Pixel {
	return ppm.data[y][x]
}

// Set sets the color pixel at the specified coordinates (x, y).
func (ppm *PPM) Set(x, y int, value Pixel) {
	ppm.data[y][x] = value
}

// Save saves the PPM image to a file and returns an error if there was a problem.
func (ppm *PPM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the PPM header to the file
	_, err = fmt.Fprintf(file, "%s\n%d %d\n%d\n", ppm.magicNumber, ppm.width, ppm.height, ppm.max)
	if err != nil {
		return err
	}

	// Write pixel data to the file
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			pixel := ppm.data[i][j]
			_, err := fmt.Fprintf(file, "%d %d %d ", pixel.R, pixel.G, pixel.B)
			if err != nil {
				return err
			}
		}
		_, err := fmt.Fprintln(file) // New line after each row of pixels
		if err != nil {
			return err
		}
	}

	return nil
}

// Invert inverts the colors of the PPM image.
func (ppm *PPM) Invert() {
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			pixel := ppm.data[i][j]
			invertedPixel := Pixel{
				R: ppm.max - pixel.R,
				G: ppm.max - pixel.G,
				B: ppm.max - pixel.B,
			}
			ppm.data[i][j] = invertedPixel
		}
	}
}

// Flip flips the PPM image horizontally.
func (ppm *PPM) Flip() {
	for i := 0; i < ppm.height; i++ {
		for j, k := 0, ppm.width-1; j < k; j, k = j+1, k-1 {
			ppm.data[i][j], ppm.data[i][k] = ppm.data[i][k], ppm.data[i][j]
		}
	}
}

// Flop flips the PPM image vertically.
func (ppm *PPM) Flop() {
	for i, j := 0, ppm.height-1; i < j; i, j = i+1, j-1 {
		ppm.data[i], ppm.data[j] = ppm.data[j], ppm.data[i]
	}
}

// SetMagicNumber sets the magic number of the PPM image.
func (ppm *PPM) SetMagicNumber(magicNumber string) {
	ppm.magicNumber = magicNumber
}

// SetMaxValue sets the maximum color value of the PPM image.
func (ppm *PPM) SetMaxValue(maxValue uint8) {
	// Check if the new maximum value is different from the current value
	if maxValue == ppm.max {
		return // No need to make changes if the maximum value is the same
	}

	// Calculate the proportionality factor to adjust pixel values
	scaleFactor := float64(maxValue) / float64(ppm.max)

	// Adjust pixel data based on the new maximum value
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			pixel := ppm.data[i][j]
			adjustedPixel := Pixel{
				R: uint8(float64(pixel.R) * scaleFactor),
				G: uint8(float64(pixel.G) * scaleFactor),
				B: uint8(float64(pixel.B) * scaleFactor),
			}
			ppm.data[i][j] = adjustedPixel
		}
	}

	// Update the new maximum value
	ppm.max = maxValue
}

// Rotate90CW rotates the PPM image 90 degrees clockwise.
func (ppm *PPM) Rotate90CW() {
	// Create a new matrix to store the rotated pixels
	rotatedData := make([][]Pixel, ppm.width)
	for i := range rotatedData {
		rotatedData[i] = make([]Pixel, ppm.height)
	}

	// Fill the new matrix with the rotated pixels
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			rotatedData[j][ppm.height-i-1] = ppm.data[i][j]
		}
	}

	// Update the image data with the rotated pixels
	ppm.data = rotatedData

	// Update the image dimensions (width and height)
	ppm.width, ppm.height = ppm.height, ppm.width
}

// ToPGM converts the PPM image to a PGM image.
func (ppm *PPM) ToPGM() *PGM {
	// Create a new matrix for PGM data
	pgmData := make([][]uint8, ppm.height)
	for i := range pgmData {
		pgmData[i] = make([]uint8, ppm.width)
	}

	// Convert PPM pixels to PGM grayscale levels
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			pixel := ppm.data[i][j]
			averageValue := uint8((uint32(pixel.R) + uint32(pixel.G) + uint32(pixel.B)) / 3)
			pgmData[i][j] = averageValue
		}
	}

	// Create a new instance of the PGM structure
	pgm := &PGM{
		data:        pgmData,
		width:       ppm.width,
		height:      ppm.height,
		magicNumber: "P2",
		max:         ppm.max,
	}

	return pgm
}

// ToPBM converts the PPM image to a PBM image.
func (ppm *PPM) ToPBM() *PBM {
	// Convert the PPM image to PGM
	pgm := ppm.ToPGM()

	// Convert the PGM image to PBM
	var data [][]bool
	data = make([][]bool, pgm.height)
	for i := range data {
		data[i] = make([]bool, pgm.width)
	}

	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			if pgm.data[i][j] > pgm.max/2 {
				data[i][j] = true
			} else {
				data[i][j] = false
			}
		}
	}

	// Create a new instance of the PBM structure
	pbm := &PBM{
		data:        data,
		width:       pgm.width,
		height:      pgm.height,
		magicNumber: "P1",
	}

	return pbm
}

func (ppm *PPM) DrawLine(p1, p2 Point, color Pixel) {
	// Handle points outside the image bounds
	if p1.X < 0 || p1.X >= ppm.width || p1.Y < 0 || p1.Y >= ppm.height {
		// Find the intersection point with the image bounds
		if p1.X < 0 {
			p1.X = 0
		} else if p1.X >= ppm.width {
			p1.X = ppm.width - 1
		}

		if p1.Y < 0 {
			p1.Y = 0
		} else if p1.Y >= ppm.height {
			p1.Y = ppm.height - 1
		}
	}

	dx := p2.X - p1.X
	dy := p2.Y - p1.Y

	// Determine the direction of the line
	var sx, sy int
	if dx > 0 {
		sx = 1
	} else {
		sx = -1
		dx = -dx
	}
	if dy > 0 {
		sy = 1
	} else {
		sy = -1
		dy = -dy
	}

	err := dx - dy

	// Draw the line
	for {
		// Check if the current point is within the image bounds
		if p1.X >= 0 && p1.X < ppm.width && p1.Y >= 0 && p1.Y < ppm.height {
			ppm.data[p1.Y][p1.X] = color
		}

		if p1.X == p2.X && p1.Y == p2.Y {
			break
		}

		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			p1.X += sx
		}
		if e2 < dx {
			err += dx
			p1.Y += sy
		}
	}
}

func (ppm *PPM) DrawRectangle(p1 Point, width, height int, color Pixel) {
	p2 := Point{p1.X + width, p1.Y}
	p3 := Point{p1.X + width, p1.Y + height}
	p4 := Point{p1.X, p1.Y + height}

	// Adjust dimensions if starting point is outside the image bounds
	if p1.X < 0 {
		width += p1.X
		p1.X = 0
	}
	if p1.Y < 0 {
		height += p1.Y
		p1.Y = 0
	}
	// Adjust dimensions if ending point is outside the image bounds
	if p1.X+width > ppm.width {
		width = ppm.width - p1.X
		ppm.DrawLine(p4, p1, color)
	} else {
		ppm.DrawLine(p4, p1, color)
		ppm.DrawLine(p2, p3, color)
	}
	if p1.Y+height > ppm.height {
		height = ppm.height - p1.Y
		ppm.DrawLine(p1, p2, color)
	} else {
		ppm.DrawLine(p1, p2, color)
		ppm.DrawLine(p3, p4, color)
	}

	// Check if the rectangle dimensions are now valid
	if width <= 0 || height <= 0 {
		// Invalid dimensions, do nothing
		return
	}
}

func (ppm *PPM) DrawFilledRectangle(p1 Point, width, height int, color Pixel) {
	ppm.DrawRectangle(p1, width, height, color)

	// Fill the rectangle with color
	for i := 0; i < ppm.height; i++ {
		var positions []int
		var numberPoints int
		for j := 0; j < ppm.width; j++ {
			if ppm.data[i][j] == color {
				numberPoints += 1
				positions = append(positions, j)
			}
		}
		if numberPoints > 1 {
			// Fill the pixels between the first and last points in the row
			for k := positions[0] + 1; k < positions[len(positions)-1]; k++ {
				ppm.data[i][k] = color
			}
		}
		// Handle the case where the rectangle exceeds image dimensions
		if height > ppm.height && width > ppm.width {
			// Fill the entire row if the rectangle is larger than the image
			for k := 0; k < ppm.width; k++ {
				ppm.data[i][k] = color
			}
		}
	}
}

func (ppm *PPM) DrawCircle(center Point, radius int, color Pixel) {
	// Draw the circumference of the circle
	for x := 0; x < ppm.height; x++ {
		for y := 0; y < ppm.width; y++ {
			dx := float64(x) - float64(center.X)
			dy := float64(y) - float64(center.Y)
			distance := math.Sqrt(dx*dx + dy*dy)

			// Check if the current pixel is on the circumference
			if math.Abs(distance-float64(radius)) < 1.0 && distance < float64(radius) {
				ppm.Set(x, y, color)
			}
		}
	}
	// Draw additional points on the axes to complete the circle
	ppm.Set(center.X-(radius-1), center.Y, color)
	ppm.Set(center.X+(radius-1), center.Y, color)
	ppm.Set(center.X, center.Y+(radius-1), color)
	ppm.Set(center.X, center.Y-(radius-1), color)
}

func (ppm *PPM) DrawFilledCircle(center Point, radius int, color Pixel) {
	// Draw the filled circle
	ppm.DrawCircle(center, radius, color)

	// Fill the pixels between the first and last points in each row
	for i := 0; i < ppm.height; i++ {
		var positions []int
		var numberPoints int
		for j := 0; j < ppm.width; j++ {
			if ppm.data[i][j] == color {
				numberPoints += 1
				positions = append(positions, j)
			}
		}
		if numberPoints > 1 {
			for k := positions[0] + 1; k < positions[len(positions)-1]; k++ {
				ppm.data[i][k] = color
			}
		}
	}
}

func (ppm *PPM) DrawTriangle(p1, p2, p3 Point, color Pixel) {
	// Draw the three sides of the triangle
	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p2, p3, color)
	ppm.DrawLine(p3, p1, color)
}

func (ppm *PPM) DrawFilledTriangle(p1, p2, p3 Point, color Pixel) {
	// Draw the filled triangle
	ppm.DrawTriangle(p1, p2, p3, color)

	// Fill the pixels between the first and last points in each row
	for i := 0; i < ppm.height; i++ {
		var positions []int
		var numberPoints int
		for j := 0; j < ppm.width; j++ {
			if ppm.data[i][j] == color {
				numberPoints += 1
				positions = append(positions, j)
			}
		}
		if numberPoints > 1 {
			for k := positions[0] + 1; k < positions[len(positions)-1]; k++ {
				ppm.data[i][k] = color
			}
		}
	}
}

func (ppm *PPM) DrawPolygon(points []Point, color Pixel) {
	numPoints := len(points)
	if numPoints < 3 {
		// A polygon must have at least 3 vertices
		return
	}

	// Draw lines between consecutive points to form the polygon
	for i := 0; i < numPoints-1; i++ {
		ppm.DrawLine(points[i], points[i+1], color)
	}

	// Draw the last line connecting the last and first points to close the polygon
	ppm.DrawLine(points[numPoints-1], points[0], color)
}

func (ppm *PPM) DrawFilledPolygon(points []Point, color Pixel) {
	// Draw the filled polygon
	ppm.DrawPolygon(points, color)

	// Fill the pixels between the first and last points in each row
	for i := 0; i < ppm.height; i++ {
		var positions []int
		var numberPoints int
		for j := 0; j < ppm.width; j++ {
			if ppm.data[i][j] == color {
				numberPoints += 1
				positions = append(positions, j)
			}
		}
		if numberPoints > 1 {
			for k := positions[0] + 1; k < positions[len(positions)-1]; k++ {
				ppm.data[i][k] = color
			}
		}
	}
}

func (ppm *PPM) DrawPerlinNoise(color1 Pixel, color2 Pixel) {
	// Function to generate Perlin noise value for given coordinates
	generatePerlinNoise := func(x, y float64) float64 {
		// You can customize this function based on your requirements
		// This is a simplified Perlin noise function
		return math.Sin(x*0.1) + math.Sin(y*0.1)
	}

	// Function to perform linear interpolation between two colors
	lerpColor := func(color1 Pixel, color2 Pixel, t float64) Pixel {
		clampedT := math.Max(0, math.Min(1, t))
		r := uint8(float64(color1.R)*(1-clampedT) + float64(color2.R)*clampedT)
		g := uint8(float64(color1.G)*(1-clampedT) + float64(color2.G)*clampedT)
		b := uint8(float64(color1.B)*(1-clampedT) + float64(color2.B)*clampedT)
		return Pixel{r, g, b}
	}

	// Iterate over each pixel in the image
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			// Generate Perlin noise value for the current pixel
			noiseValue := generatePerlinNoise(float64(x)/50, float64(y)/50)

			// Map the noise value to the color range between color1 and color2
			lerpedColor := lerpColor(color1, color2, (noiseValue+1)/2)

			// Set the pixel color in the image
			ppm.Set(x, y, lerpedColor)
		}
	}
}
