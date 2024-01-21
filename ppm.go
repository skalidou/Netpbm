package Netpbm

import (
	"errors"
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

// ReadPPM reads a PPM image from a file and returns a struct that represents the image.
func ReadPPM(filename string) (*PPM, error) {
	var err error
	var magicNumber string = ""
	var width int
	var height int
	var maxval int
	var counter int
	var headersize int
	var splitfile []string
	file, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("error opening the file")
	}
	if strings.Contains(string(file), "\r") {
		splitfile = strings.SplitN(string(file), "\r\n", -1)
	} else {
		splitfile = strings.SplitN(string(file), "\n", -1)
	}
	for i := range splitfile {
		if strings.Contains(splitfile[i], "P3") {
			magicNumber = "P3"
		} else if strings.Contains(splitfile[i], "P6") {
			magicNumber = "P6"
		}
		if strings.HasPrefix(splitfile[i], "#") && maxval != 0 {
			headersize = counter
		}
		splitl := strings.SplitN(splitfile[i], " ", -1)
		if width == 0 && height == 0 && len(splitl) >= 2 {
			width, err = strconv.Atoi(splitl[0])
			if err != nil {
				fmt.Println("error reading widht")
			}
			height, err = strconv.Atoi(splitl[1])
			if err != nil {
				fmt.Println("error reading height")
			}
			headersize = counter
		}

		if maxval == 0 && width != 0 {
			maxval, err = strconv.Atoi(splitfile[i])
			headersize = counter
		}
		counter++
	}

	data := make([][]Pixel, height)

	for j := 0; j < height; j++ {
		data[j] = make([]Pixel, width)
	}
	var splitdata []string

	if counter > headersize {
		for i := 0; i < height; i++ {
			splitdata = strings.SplitN(splitfile[headersize+1+i], " ", -1)
			for j := 0; j < width*3; j += 3 {
				r, _ := strconv.Atoi(splitdata[j])
				if r > maxval {
					r = maxval
				}
				g, _ := strconv.Atoi(splitdata[j+1])
				if g > maxval {
					g = maxval
				}
				b, _ := strconv.Atoi(splitdata[j+2])
				if b > maxval {
					b = maxval
				}
				data[i][j/3] = Pixel{R: uint8(r), G: uint8(g), B: uint8(b)}
			}
		}
	}
	return &PPM{data: data, width: width, height: height, magicNumber: magicNumber, max: uint8(maxval)}, err
}

// Size returns the width and height of the image.
func (ppm *PPM) Size() (int, int) {
	return ppm.width, ppm.height
}

// At returns the value of the pixel at (x, y).
func (ppm *PPM) At(x, y int) Pixel {
	return ppm.data[y][x]
}

// Set sets the value of the pixel at (x, y).
func (ppm *PPM) Set(x, y int, value Pixel) {
	ppm.data[x][y] = value
}

// Save saves the PPM image to a file and returns an error if there was a problem.
func (ppm *PPM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		file.Close()
		return err
	}
	_, err = fmt.Fprintln(file, ppm.magicNumber)
	if err != nil {
		file.Close()
		return err
	}
	_, err = fmt.Fprintln(file, ppm.width, ppm.height)
	if err != nil {
		file.Close()
		return err
	}
	_, err = fmt.Fprintln(file, ppm.max)
	if err != nil {
		file.Close()
		return err
	}

	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			if ppm.data[y][x].R > ppm.max || ppm.data[y][x].G > ppm.max || ppm.data[y][x].B > ppm.max {
				errors.New("data value is too high")
			} else {
				fmt.Fprint(file, ppm.data[y][x].R, ppm.data[y][x].G, ppm.data[y][x].B, " ")
			}
		}
		fmt.Fprintln(file)
	}
	return err
}

// Invert inverts the colors of the PPM image.
func (ppm *PPM) Invert() {
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			r := ppm.max - ppm.data[i][j].R
			g := ppm.max - ppm.data[i][j].G
			b := ppm.max - ppm.data[i][j].B
			ppm.data[i][j] = Pixel{R: r, G: g, B: b}
		}
	}
}

// Flip flips the PPM image horizontally.
func (ppm *PPM) Flip() {
	for i := 0; i < len(ppm.data); i++ {
		for j := 0; j < len(ppm.data[i])/2; j++ {
			startdata := ppm.data[i][j]
			ppm.data[i][j] = ppm.data[i][len(ppm.data[i])-1-j]
			ppm.data[i][len(ppm.data[i])-1-j] = startdata
		}
	}
}

// Flop flops the PPM image vertically.
func (ppm *PPM) Flop() {
	for i := 0; i < len(ppm.data)/2; i++ {
		startdata := ppm.data[i]
		ppm.data[i] = ppm.data[len(ppm.data)-1-i]
		ppm.data[len(ppm.data)-1-i] = startdata
	}
}

// SetMagicNumber sets the magic number of the PPM image.
func (ppm *PPM) SetMagicNumber(magicNumber string) {
	ppm.magicNumber = magicNumber
}

// SetMaxValue sets the max value of the PPM image.
func (ppm *PPM) SetMaxValue(maxValue uint8) {
	oldMax := ppm.max
	ppm.max = maxValue
	var r, g, b uint8
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			r = uint8(float64(ppm.data[i][j].R) * float64(ppm.max) / float64(oldMax))
			g = uint8(float64(ppm.data[i][j].G) * float64(ppm.max) / float64(oldMax))
			b = uint8(float64(ppm.data[i][j].B) * float64(ppm.max) / float64(oldMax))
			ppm.data[i][j] = Pixel{R: r, G: g, B: b}
		}
	}
}

// Rotate90CW rotates the PPM image 90° clockwise.
func (ppm *PPM) Rotate90CW() {
	newData := make([][]Pixel, ppm.width)
	for i := 0; i < ppm.height; i++ {
		newData[i] = make([]Pixel, ppm.height)
		for j := 0; j < ppm.width; j++ {
			newData[i][j] = ppm.data[j][i]
		}
	}

	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width/2; j++ {
			prevdata := newData[i][j]
			newData[i][j] = newData[i][ppm.height-j-1]
			newData[i][ppm.height-j-1] = prevdata
		}
	}
	ppm.data = newData
}

// ToPGM converts the PPM image to PGM.
func (ppm *PPM) ToPGM() *pgm {
	var magicnb string
	data := make([][]uint8, ppm.height)
	if ppm.magicNumber == "P3" {
		magicnb = "P2"
	} else if ppm.magicNumber == "P6" {
		magicnb = "P5"
	}
	for i := 0; i < ppm.height; i++ {
		data[i] = make([]uint8, ppm.width)
		for j := 0; j < ppm.width; j++ {
			data[i][j] = uint8((int(ppm.data[i][j].R) + int(ppm.data[i][j].G) + int(ppm.data[i][j].B)) / 3)
		}
	}

	return &pgm{magicNumber: magicnb, data: data, width: ppm.width, height: ppm.height, max: ppm.max}
}

// ToPBM converts the PPM image to PBM.
func (ppm *PPM) ToPBM() *PBM {
	var magicnb string
	data := make([][]bool, ppm.height)
	if ppm.magicNumber == "P3" {
		magicnb = "P1"
	} else if ppm.magicNumber == "P6" {
		magicnb = "P4"
	}
	for i := 0; i < ppm.height; i++ {
		data[i] = make([]bool, ppm.width)
		for j := 0; j < ppm.width; j++ {
			data[i][j] = uint8((int(ppm.data[i][j].R)+int(ppm.data[i][j].G)+int(ppm.data[i][j].B))/3) < ppm.max/2
		}
	}
	return &PBM{magicNumber: magicnb, width: ppm.width, height: ppm.height, data: data}
}

type Point struct {
	X, Y int
}

// DrawLine draws a line between two points in the PPM image using the specified color.
func (ppm *PPM) DrawLine(p1, p2 Point, color Pixel) {
	// Calculate the differences in X and Y coordinates
	dx := float64(p2.X - p1.X)
	dy := float64(p2.Y - p1.Y)
	// Determine the number of steps based on the maximum difference in coordinates
	steps := int(math.Max(math.Abs(dx), math.Abs(dy)))
	// Calculate the incremental changes in X and Y coordinates
	xIncrement := dx / float64(steps)
	yIncrement := dy / float64(steps)
	// Initialize the starting coordinates
	x, y := float64(p1.X), float64(p1.Y)
	// Draw the line by setting the specified color at each point along the line
	for i := 0; i <= steps; i++ {
		ppm.Set(int(x), int(y), color)
		x += xIncrement
		y += yIncrement
	}
}

// DrawRectangle draws a rectangle.
func (ppm *PPM) DrawRectangle(p1 Point, width, height int, color Pixel) {
	if p1.Y+height < ppm.height {
		for y := p1.Y; y <= p1.Y+height; y++ {
			if p1.X+width < ppm.width {
				for x := p1.X; x <= p1.X+width; x++ {
					if x == p1.X || x == p1.X+width {
						ppm.data[y][x] = color
					} else if y == p1.Y || y == p1.Y+height {
						ppm.data[y][x] = color
					}
				}
			} else if p1.X+width >= ppm.width {
				for x := p1.X; x < ppm.width; x++ {
					if y == p1.Y || y == p1.Y+height {
						ppm.data[y][x] = color
					} else if x == p1.X || x == p1.X+width {
						ppm.data[y][x] = color
					}
				}
			}
		}
	} else if p1.Y+height >= ppm.height {
		for y := p1.Y; y < ppm.height; y++ {
			if p1.X+width < ppm.width {
				for x := p1.X; x <= p1.X+width; x++ {
					if x == p1.X || x == p1.X+width {
						ppm.data[y][x] = color
					} else if y == p1.Y || y == p1.Y+height {
						ppm.data[y][x] = color
					}
				}
			} else if p1.X+width >= ppm.width {
				for x := p1.X; x < ppm.width; x++ {
					if y == p1.Y || y == p1.Y+height {
						ppm.data[y][x] = color
					} else if x == p1.X || x == p1.X+width {
						ppm.data[y][x] = color
					}
				}
			}
		}
	}

}

// DrawFilledRectangle draws a filled rectangle.
func (ppm *PPM) DrawFilledRectangle(p1 Point, width, height int, color Pixel) {
	// If the height of the rectangle is within the limit of the file
	if p1.Y+height < ppm.height {
		for y := p1.Y; y < p1.Y+height; y++ {
			// If the width of the rectangle is within the data of the file
			if p1.X+width < ppm.width {
				for x := p1.X; x <= p1.X+width; x++ {
					ppm.data[y][x] = color
				}
				// If the width of the rectangle is within the data of the file
			} else if p1.X+width > ppm.width {
				for x := p1.X; x < ppm.width; x++ {
					ppm.data[y][x] = color
				}
			}
		}
		// If the height of the rectangle goes beyond the limit of the file
	} else if p1.Y+height > ppm.height {
		for y := p1.Y; y < ppm.height; y++ {
			// If the width of the rectangle is within the data of the file
			if p1.X+width < ppm.width {
				for x := p1.X; x <= p1.X+width; x++ {
					ppm.data[y][x] = color
				}
				// If the width of the rectangle is within the data of the file
			} else if p1.X+width > ppm.width {
				for x := p1.X; x < ppm.width; x++ {
					ppm.data[y][x] = color
				}
			}
		}
	}
}

// DrawCircle draws a circle.
func (ppm *PPM) DrawCircle(center Point, radius int, color Pixel) {
	// ...
}

// DrawFilledCircle draws a filled circle.
func (ppm *PPM) DrawFilledCircle(center Point, radius int, color Pixel) {
	// ...
}

// DrawTriangle draws a triangle.
func (ppm *PPM) DrawTriangle(p1, p2, p3 Point, color Pixel) {
	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p2, p3, color)
	ppm.DrawLine(p3, p1, color)
}

// DrawFilledTriangle draws a filled triangle.
func (ppm *PPM) DrawFilledTriangle(p1, p2, p3 Point, color Pixel) {
	// ...
}

// DrawPolygon draws a polygon.
func (ppm *PPM) DrawPolygon(points []Point, color Pixel) {
	for i := 0; i < len(points); i++ {
		nextIndex := (i + 1) % len(points)
		ppm.DrawLine(points[i], points[nextIndex], color)
	}
}

// DrawFilledPolygon draws a filled polygon.
func (ppm *PPM) DrawFilledPolygon(points []Point, color Pixel) {
	// ...
}
