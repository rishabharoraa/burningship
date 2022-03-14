package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"math/cmplx"
	"os"
)

const CENTER complex128 = -1.75 - 0.03i
const SIZE float64 = 0.06

const RESOLUTION int = 2048

const MAX_ITERATIONS int = 2048

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

func MapRange(input, oldMin, oldMax, newMin, newMax uint8) uint8 {
	return uint8(((float64(input)-float64(oldMin))*(float64(newMax)-float64(newMin)))/(float64(oldMax)-float64(oldMin)) + float64(newMin))
}

func ApplyFilter(shade, fr, fg, fb uint8) (uint8, uint8, uint8) {

	_r := MapRange(shade, 0, 255, fr, 255)
	_g := MapRange(shade, 0, 255, fg, 255)
	_b := MapRange(shade, 0, 255, fb, 255)
	return (_r*2 + 3) % 255, (_g*2 + 3) % 255, (_b*2 + 3) % 255
}

func ComputeColor(num int, maxIterations int) (uint8, uint8, uint8) {
	var fr, fg, fb uint8 = 0, 0, 0
	if num == maxIterations {
		return fr, fg, fb
	}
	shade := 255 - uint8(num*2%255)
	return ApplyFilter(shade, fr, fg, fb)
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
	png.Encode(file, img)
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
