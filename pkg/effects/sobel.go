package effects

import (
	"image"
	"math"
	"runtime"
)

type sobel struct {
	threshold int
	invert    bool
}

// NewSobel the input image should be a grayscale image, the output will be a version of
// the input image with the Sobel edge detector applied to it. A value of -1 for threshold
// will return an image whos rgb values are the sobel intensity values, if 0 <= threshold <= 255
// then the rgb values will be 255 if the intensity is >= threshold and 0 if the intensity
// is < threshold
func NewSobel(threshold int, invert bool) Effect {
	return &sobel{
		threshold: threshold,
		invert:    invert,
	}
}

func (s *sobel) Apply(img *Image, numRoutines int) (*Image, error) {
	if numRoutines == 0 {
		numRoutines = runtime.GOMAXPROCS(0)
	}

	sobelX := [][]int{
		{-1, 0, 1},
		{-2, 0, 2},
		{-1, 0, 1},
	}
	sobelY := [][]int{
		{-1, -2, -1},
		{0, 0, 0},
		{1, 2, 1},
	}

	pf := func(ri, x, y, offset, inStride int, inPix, outPix []uint8) {
		var px, py int
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				pOffset := offset + (dx*4 + dy*inStride)
				r := int(inPix[pOffset])
				px += sobelX[dx+1][dy+1] * r
				py += sobelY[dx+1][dy+1] * r
			}
		}

		val := uint8(math.Sqrt(float64(px*px + py*py)))
		if s.threshold != -1 {
			if val >= uint8(s.threshold) {
				val = 255
			} else {
				val = 0
			}
		}

		if s.invert {
			val = 255 - val
		}
		outPix[offset] = val
		outPix[offset+1] = val
		outPix[offset+2] = val
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
			X:      img.Bounds.X + 1,
			Y:      img.Bounds.Y + 1,
			Width:  img.Bounds.Width - 2,
			Height: img.Bounds.Height - 2,
		},
	}

	runParallel(numRoutines, img, out.Bounds, out, pf, 0)
	return out, nil
}
