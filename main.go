package main

import (
	"fmt"
	"image"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/fogleman/gg"
)

func main() {
	// File must be JPEG/JPG -> if PNG, change to LoadPNG
	file := "parrot.jpeg"
	refImg, _ := gg.LoadJPG(file)
	fmt.Println("Image Loaded")
	context := gg.NewContextForImage(refImg)
	bounds := refImg.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	fmt.Println("Image Width:", width, "Image Height:", height)

	// Sets the scale to create rectangles at a proper proportion to the image
	wScale := int(math.Round(float64(width) / 100))
	hScale := int(math.Round(float64(height) / 100))
	fmt.Println("Width Scale:", wScale, "Height Scale:", hScale)

	fmt.Println("Drawing Started")
	rand.Seed(time.Now().UnixNano())
	w, h := 1, 1
	for y := 0; y < height; y += h {
		h = rand.Intn(hScale) + hScale
		for x := 0; x < width; x += w {
			w = rand.Intn(wScale) + wScale
			context.Push()
			if x+w > width {
				w = width - x
			}
			if y+h > height {
				h = height - y
			}
			context.DrawRectangle(float64(x), float64(y), float64(w), float64(h))
			r, g, b := getAverageColor(x, y, x+w, y+h, refImg)

			// Rectangle is slightly translucent to be able to see some finer detail
			context.SetRGBA255(r, g, b, 235)
			context.Fill()
			context.Pop()
		}
	}
	fmt.Println("Drawing Ended")
	newFile := strings.Split(file, ".")
	// Saves the file under the original name + filtered as a PNG
	context.SavePNG(newFile[0] + "_filtered.png")
	fmt.Println("File Saved")
}

// Returns the average red, green, and blue colors within a rectangle of the image
func getAverageColor(x0, y0, x1, y1 int, img image.Image) (int, int, int) {
	r := make([]uint32, (y1-y0)*(x1-x0))
	g := make([]uint32, (y1-y0)*(x1-x0))
	b := make([]uint32, (y1-y0)*(x1-x0))
	idx := 0

	// Get all colors in the range and add them to their slices
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			color := img.At(x, y)
			red, green, blue, _ := color.RGBA()
			r[idx], g[idx], b[idx] = red, green, blue
			idx += 1
		}
	}

	// Get the averages for each color slice
	dR := getAverage(r)
	dG := getAverage(g)
	dB := getAverage(b)
	return dR, dG, dB
}

// Takes a slice of uint32 and returns an average color in an int
func getAverage(slice []uint32) int {
	var sum uint32 = 0
	for i := 0; i < len(slice); i++ {
		sum += (slice[i])
	}
	count := len(slice)
	// Take the square root to get the true average
	avg := math.Sqrt(float64(sum / uint32(count)))

	return int(avg)
}
