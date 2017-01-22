package effects

import (
	"image"
	"math"
	"runtime"
)

// Sobel the input image should be a grayscale image, the output will be a version of
// the input image with the Sobel edge detector applied to it. A value of -1 for threshold
// will return an image whos rgb values are the sobel intensity values, if 0 <= threshold <= 255
// then the rgb values will be 255 if the intensity is >= threshold and 0 if the intensity
// is < threshold
func Sobel(img *Image, numRoutines, threshold int) (*Image, error) {
	if numRoutines == 0 {
		numRoutines = runtime.GOMAXPROCS(0)
	}

	out := &Image{img: image.NewRGBA(img.img.Bounds())}

	sobelX := [][]int{
		[]int{-1, 0, 1},
		[]int{-2, 0, 2},
		[]int{-1, 0, 1},
	}
	sobelY := [][]int{
		[]int{-1, -2, -1},
		[]int{0, 0, 0},
		[]int{1, 2, 1},
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
		if threshold != -1 {
			if val >= uint8(threshold) {
				val = 255
			} else {
				val = 0
			}
		}
		outPix[offset] = val
		outPix[offset+1] = val
		outPix[offset+2] = val
		outPix[offset+3] = 255
	}

	inBounds := image.Rectangle{
		Min: image.Point{X: 1, Y: 1},
		Max: image.Point{X: img.Bounds().Dx() - 2, Y: img.Bounds().Dy() - 2},
	}

	runParallel(numRoutines, img, inBounds, out, pf)
	return out, nil
}
