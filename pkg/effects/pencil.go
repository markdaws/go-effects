package effects

import (
	"fmt"
	"runtime"
)

type pencil struct {
	blurFactor int
}

func (p *pencil) Apply(img *Image, numRoutines int) (*Image, error) {
	if !isOddInt(p.blurFactor) {
		return nil, fmt.Errorf("blurFactor must be odd")
	}

	if numRoutines == 0 {
		numRoutines = runtime.GOMAXPROCS(0)
	}

	inImg := img
	if p.blurFactor != 0 {
		var err error
		gaussian := NewGaussian(p.blurFactor, 1)
		inImg, err = gaussian.Apply(img, numRoutines)
		if err != nil {
			return nil, err
		}
	}
	sobel := NewSobel(-1, true)
	out, err := sobel.Apply(inImg, numRoutines)
	return out, err
}

// NewPencil renders the input image as if it was drawn in pencil. It is simply
// an inverted Sobel image. You can specify the blurFactor, a value that must
// be odd, to blur the input image to get rid of the noise. This is the gaussian
// kernel size, larger numbers blur more but can significantly increase processing time.
func NewPencil(blurFactor int) Effect {
	return &pencil{blurFactor: blurFactor}
}
