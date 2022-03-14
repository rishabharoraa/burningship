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

func ComputeColor(num int, maxIterations int) (uint8, uint8, uint8) {
	if num == maxIterations {
		return 0, 0, 0
	}
	shade := 255 - uint8(num*2%255)
	return shade, shade, shade
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
