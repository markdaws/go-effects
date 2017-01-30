package effects

import (
	"fmt"
	"image"
	"math"
	"runtime"
)

type gaussian struct {
	kernelSize int
	sigma      float64
}

// NewGaussian is an effect that applies a gaussian blur to the image
func NewGaussian(kernelSize int, sigma float64) Effect {
	return &gaussian{
		kernelSize: kernelSize,
		sigma:      sigma,
	}
}

func (g *gaussian) Apply(img *Image, numRoutines int) (*Image, error) {
	if !isOddInt(g.kernelSize) {
		return nil, fmt.Errorf("kernel size must be odd")
	}

	if numRoutines == 0 {
		numRoutines = runtime.GOMAXPROCS(0)
	}

	kernel := gaussianKernel(g.kernelSize, g.sigma)
	kernelOffset := (g.kernelSize - 1) / 2
	pf := func(ri, x, y, offset, inStride int, inPix, outPix []uint8) {
		var gr, gb, gg float64
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

	out := &Image{
		img: image.NewRGBA(image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: img.Width, Y: img.Height},
		}),
		Width:  img.Width,
		Height: img.Height,
		Bounds: Rect{
			X:      img.Bounds.X + kernelOffset,
			Y:      img.Bounds.Y + kernelOffset,
			Width:  img.Bounds.Width - 2*kernelOffset,
			Height: img.Bounds.Height - 2*kernelOffset,
		},
	}

	runParallel(numRoutines, img, out.Bounds, out, pf, 0)
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
