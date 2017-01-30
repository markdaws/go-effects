package effects

import (
	"image"
	"runtime"
)

type oilPainting struct {
	filterSize int
	levels     int
}

func (op *oilPainting) Apply(img *Image, numRoutines int) (*Image, error) {
	levels := op.levels - 1
	filterOffset := (op.filterSize - 1) / 2

	if numRoutines == 0 {
		numRoutines = runtime.GOMAXPROCS(0)
	}

	var iBin, rBin, gBin, bBin [][]int
	iBin = make([][]int, numRoutines)
	rBin = make([][]int, numRoutines)
	gBin = make([][]int, numRoutines)
	bBin = make([][]int, numRoutines)
	for ri := 0; ri < numRoutines; ri++ {
		iBin[ri] = make([]int, levels+1)
		rBin[ri] = make([]int, levels+1)
		gBin[ri] = make([]int, levels+1)
		bBin[ri] = make([]int, levels+1)
	}

	pf := func(ri, x, y, offset, inStride int, inPix, outPix []uint8) {
		reset(iBin[ri])
		reset(rBin[ri])
		reset(gBin[ri])
		reset(bBin[ri])

		var maxIntensity int
		var maxIndex int

		for fy := -filterOffset; fy <= filterOffset; fy++ {
			for fx := -filterOffset; fx <= filterOffset; fx++ {
				fOffset := offset + (fx*4 + fy*inStride)
				r := inPix[fOffset]
				g := inPix[fOffset+1]
				b := inPix[fOffset+2]
				ci := int(roundToInt32((float64(r+g+b) / 3.0 * float64(levels)) / 255.0))
				iBin[ri][ci]++
				rBin[ri][ci] += int(r)
				gBin[ri][ci] += int(g)
				bBin[ri][ci] += int(b)

				if iBin[ri][ci] > maxIntensity {
					maxIntensity = iBin[ri][ci]
					maxIndex = ci
				}
			}
		}

		outPix[offset] = uint8(rBin[ri][maxIndex] / maxIntensity)
		outPix[offset+1] = uint8(gBin[ri][maxIndex] / maxIntensity)
		outPix[offset+2] = uint8(bBin[ri][maxIndex] / maxIntensity)
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
			X:      img.Bounds.X + filterOffset,
			Y:      img.Bounds.Y + filterOffset,
			Width:  img.Bounds.Width - 2*filterOffset,
			Height: img.Bounds.Height - 2*filterOffset,
		},
	}
	runParallel(numRoutines, img, out.Bounds, out, pf, 0)
	return out, nil
}

// NewOilPainting renders the input image as if it was painted like an oil painting. numRoutines specifies how many
// goroutines should be used to process the image in parallel, use 0 to let the library decide. filterSize specifies
// how bold the image should look, larger numbers equate to larger strokes, levels specifies how many buckets colors
// will be grouped in to, start with values 5,30 to see how that works.
func NewOilPainting(filterSize, levels int) Effect {
	return &oilPainting{filterSize: filterSize, levels: levels}
}
