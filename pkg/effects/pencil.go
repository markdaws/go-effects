package effects

import (
	"fmt"
	"runtime"
)

// Pencil renders the input image as if it was drawn in pencil. It is simply
// an inverted Sobel image. You can specify the blurFactor, a value that must
// be odd, to blur the input image to get rid of the noise. This is the gaussian
// kernel size, larger numbers blur more but can significantly increase processing time.
func Pencil(img *Image, numRoutines int, blurFactor int) (*Image, error) {
	if !isOddInt(blurFactor) {
		return nil, fmt.Errorf("blurFactor must be odd")
	}

	if numRoutines == 0 {
		numRoutines = runtime.GOMAXPROCS(0)
	}

	inImg := img
	if blurFactor != 0 {
		var err error
		inImg, err = Gaussian(img, numRoutines, blurFactor, 1)
		if err != nil {
			return nil, err
		}
	}
	out, err := Sobel(inImg, numRoutines, -1, true)
	return out, err
}
