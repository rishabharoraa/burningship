package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"math/cmplx"
	"os"

	"github.com/nfnt/resize"
)

const CENTER complex128 = -1.75 - 0.03i
const SIZE float64 = 0.06

const RESOLUTION int = 16384

const MAX_ITERATIONS int = 72

func CalculateStep(size float64, resolution int) float64 {
	return (2 * size / float64(resolution))
}

func ComplexIterateBurningShip(z, c complex128) complex128 {
	_z := complex(math.Abs(real(z)), math.Abs(imag(z)))
	return _z*_z + c
}

func CalculateIterations(x, y float64, maxIterations int) int {
	iters := 0
	z := 0 + 0i
	c := complex(x, y)
	for ; iters < maxIterations; iters++ {
		if cmplx.Abs(z) > 2 {
			return iters
		}
		z = ComplexIterateBurningShip(z, c)
	}
	return iters
}

func MapRange(input, oldMin, oldMax int, newMin, newMax uint8) uint8 {
	return uint8(((float64(input)-float64(oldMin))*(float64(newMax)-float64(newMin)))/(float64(oldMax)-float64(oldMin)) + float64(newMin))
}

func ApplyFilter(shade, fr, fg, fb uint8) (uint8, uint8, uint8) {

	_r := MapRange(int(shade), 0, 255, 114, 20)
	_g := MapRange(int(shade), 0, 255, 31, 28)
	_b := MapRange(int(shade), 0, 255, 9, 60)
	return _r, _g, _b
}

func ComputeColor(num int, maxIterations int) (uint8, uint8, uint8) {
	var fr, fg, fb uint8 = 64, 64, 64
	if num == maxIterations {
		return 0, 0, 0
	}
	if num > int(float64(maxIterations)*0.4) {
		return 159, 43, 12
	}
	shade := MapRange(num, 0, maxIterations, 255, 0)
	return ApplyFilter(shade, fr, fg, fb)
	// return uint8(255 - num*4), uint8(255 - num*4), uint8(255 - num*4)
}

func Paint(points [][]int, maxIterations int) {

	size := len(points)

	img := image.NewRGBA(
		image.Rectangle{
			image.Point{0, 0},
			image.Point{size, size},
		},
	)

	for y := 0; y < size-1; y++ {
		for x := 0; x < size-1; x++ {
			r, g, b := ComputeColor(points[y][x], maxIterations)
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}
	file, _ := os.Create("burningShip.png")
	resizedImg := resize.Resize(1024, 1024, img, resize.Lanczos3)
	png.Encode(file, resizedImg)
}

func Plot(center complex128, size float64, resolution int, maxIterations int) {

	var points [][]int

	step := CalculateStep(size, resolution)

	var startX float64 = real(center) - float64(size)
	var endX float64 = real(center) + float64(size)

	var startY float64 = imag(center) - float64(size)
	var endY float64 = imag(center) + float64(size)

	for y := startY; y < endY; y += step {
		var line []int
		for x := startX; x < endX; x += step {
			iters := CalculateIterations(x, y, maxIterations)
			line = append(line, iters)
		}
		points = append(points, line)
	}

	Paint(points, maxIterations)
}

func main() {
	Plot(CENTER, SIZE, RESOLUTION, MAX_ITERATIONS)
}
