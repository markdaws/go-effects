package effects

import (
	"fmt"
	"image"
	"math"
	"runtime"
)

func Gaussian(img *Image, numRoutines, kernelSize int, sigma float64) (*Image, error) {
	if !isOddInt(kernelSize) {
		return nil, fmt.Errorf("kernel size must be odd")
	}

	if numRoutines == 0 {
		numRoutines = runtime.GOMAXPROCS(0)
	}

	out := &Image{img: image.NewRGBA(img.img.Bounds())}

	kernel := gaussianKernel(kernelSize, sigma)

	pf := func(ri, x, y, offset, inStride int, inPix, outPix []uint8) {
		var gr, gb, gg float64
		kernelOffset := (kernelSize - 1) / 2
		for dy := -kernelOffset; dy <= kernelOffset; dy++ {
			for dx := -kernelOffset; dx <= kernelOffset; dx++ {
				pOffset := offset + (dx*4 + dy*inStride)
				r := inPix[pOffset]
				g := inPix[pOffset+1]
				b := inPix[pOffset+2]

				scale := kernel[dx+kernelOffset][dy+kernelOffset]
				gr += scale * float64(r)
				gg += scale * float64(g)
				gb += scale * float64(b)
			}
		}

		outPix[offset] = uint8(gr)
		outPix[offset+1] = uint8(gg)
		outPix[offset+2] = uint8(gb)
		outPix[offset+3] = 255
	}

	kernelOffset := (kernelSize - 1) / 2
	inBounds := image.Rectangle{
		Min: image.Point{X: kernelOffset, Y: kernelOffset},
		Max: image.Point{X: img.Bounds().Dx() - 2*kernelOffset, Y: img.Bounds().Dy() - 2*kernelOffset},
	}

	runParallel(numRoutines, img, inBounds, out, pf, -1)
	return out, nil
}

func gaussianKernel(dimension int, sigma float64) [][]float64 {
	k := make([][]float64, dimension)
	sum := 0.0
	for x := 0; x < dimension; x++ {
		k[x] = make([]float64, dimension)
		for y := 0; y < dimension; y++ {
			k[x][y] = gaussianXY(x, y, sigma)
			sum += k[x][y]
		}
	}

	scale := 1.0 / sum
	for y := 0; y < dimension; y++ {
		for x := 0; x < dimension; x++ {
			k[x][y] *= scale
		}
	}

	return k
}

// expects x,y to be 0 at the center of the kernel
func gaussianXY(x, y int, sigma float64) float64 {
	return ((1.0 / (2 * math.Pi * sigma * sigma)) * math.E) - (float64(x*x+y*y) / (2 * sigma * sigma))
}
