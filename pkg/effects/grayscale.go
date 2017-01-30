package effects

import (
	"image"
	"math"
	"runtime"
)

// GSAlgo the type of algorithm to use when converting an image to it's grayscale equivalent
type GSAlgo int

const (
	// GSLIGHTNESS is the average of the min and max r,g,b value
	GSLIGHTNESS GSAlgo = iota

	// GSAVERAGE is the average of the r,g,b values of each pixel
	GSAVERAGE

	// GSLUMINOSITY used a weighting for r,g,b based on how the human eye perceives colors
	GSLUMINOSITY
)

type grayscale struct {
	algo GSAlgo
}

func (gs *grayscale) Apply(img *Image, numRoutines int) (*Image, error) {
	if numRoutines == 0 {
		numRoutines = runtime.GOMAXPROCS(0)
	}

	pf := func(ri, x, y, offset, inStride int, inPix, outPix []uint8) {
		var r, g, b uint8 = inPix[offset], inPix[offset+1], inPix[offset+2]
		switch gs.algo {
		case GSLIGHTNESS:
			max := math.Max(math.Max(float64(r), float64(g)), float64(b))
			min := math.Max(math.Min(float64(r), float64(g)), float64(b))
			r = uint8(max + min/2)
			g = r
			b = r
		case GSAVERAGE:
			r = (r + g + b) / 3
			g = r
			b = r
		case GSLUMINOSITY:
			r = uint8(0.21*float64(r) + 0.72*float64(g) + 0.07*float64(b))
			g = r
			b = r
		}
		outPix[offset] = r
		outPix[offset+1] = g
		outPix[offset+2] = b
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
			X:      img.Bounds.X,
			Y:      img.Bounds.Y,
			Width:  img.Bounds.Width,
			Height: img.Bounds.Height,
		},
	}

	runParallel(numRoutines, img, out.Bounds, out, pf, 0)
	return out, nil
}

// NewGrayscale renders the input image as a grayscale image. numRoutines specifies how many
// goroutines should be used to process the image in parallel, use 0 to let the library decide
func NewGrayscale(algo GSAlgo) Effect {
	return &grayscale{algo: algo}
}
