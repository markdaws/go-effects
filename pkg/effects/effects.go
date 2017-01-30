package effects

import (
	"math"
	"sync"
)

// Effect interface for any effect type
type Effect interface {
	// Apply applies the effect to the input image and returns an output image
	Apply(img *Image, numRoutines int) (*Image, error)
}

type pixelFunc func(ri, x, y, offset, inStride int, inPix, outPix []uint8)

func runParallel(numRoutines int, inImg *Image, inBounds Rect, outImg *Image, pf pixelFunc, blockWidth int) {
	w := inBounds.Width
	h := inBounds.Height

	minX := inBounds.X
	minY := inBounds.Y

	stride := inImg.img.Stride
	inPix := inImg.img.Pix
	outPix := outImg.img.Pix

	wg := sync.WaitGroup{}
	xOffset := minX

	var widthPerRoutine int
	if blockWidth != 0 {
		widthPerRoutine = blockWidth
	} else {
		widthPerRoutine = w / numRoutines
	}

	for r := 0; r < numRoutines; r++ {
		wg.Add(1)

		if r == numRoutines-1 {
			widthPerRoutine = (minX + w) - xOffset
		}

		go func(ri, xStart, yStart, width, height int) {
			for x := xStart; x < xStart+width; x++ {
				for y := yStart; y < yStart+height; y++ {
					offset := y*stride + x*4
					pf(ri, x, y, offset, stride, inPix, outPix)
				}
			}
			wg.Done()
		}(r, xOffset, minY, widthPerRoutine, h)

		xOffset += widthPerRoutine
	}
	wg.Wait()
}

func roundToInt32(a float64) int32 {
	if a < 0 {
		return int32(a - 0.5)
	}
	return int32(a + 0.5)
}

func rangeInt(i, min, max int) int {
	return int(math.Min(math.Max(float64(i), float64(min)), float64(max)))
}

func isOddInt(i int) bool {
	return i%2 != 0
}

func reset(s []int) {
	for i := range s {
		s[i] = 0
	}
}
