package effects

import (
	"image"
	"runtime"
)

type brightness struct {
	offset int
}

// Apply applies the brgihtness effect to the input image
func (br *brightness) Apply(img *Image, numRoutines int) (*Image, error) {
	if numRoutines == 0 {
		numRoutines = runtime.GOMAXPROCS(0)
	}

	pf := func(ri, x, y, offset, inStride int, inPix, outPix []uint8) {

		r := int(inPix[offset]) + br.offset
		g := int(inPix[offset+1]) + br.offset
		b := int(inPix[offset+2]) + br.offset
		a := inPix[offset+3]

		outPix[offset] = uint8(rangeInt(r, 0, 255))
		outPix[offset+1] = uint8(rangeInt(g, 0, 255))
		outPix[offset+2] = uint8(rangeInt(b, 0, 255))
		outPix[offset+3] = a
	}

	out := &Image{
		img: image.NewRGBA(image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: img.Width, Y: img.Height},
		}),
		Width:  img.Width,
		Height: img.Height,
		Bounds: img.Bounds,
	}
	runParallel(numRoutines, img, out.Bounds, out, pf, 0)
	return out, nil
}

// NewBrightness returns an effect that can lighten of darken an image. To lighten an image
// set offset as a positive value between 0 and 255, to darken, set it as a negative number
func NewBrightness(offset int) Effect {
	return &brightness{offset: offset}
}
