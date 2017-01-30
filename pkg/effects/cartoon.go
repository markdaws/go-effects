package effects

import (
	"image"
	"runtime"
)

// CTOpts options to pass to the Cartoon effect
type CTOpts struct {
	// BlurKernelSize is the gaussian blur kernel size. You might need to blur
	// the original input image to reduce the amount of noise you get in the edge
	// detection phase. Set to 0 to skip blur, otherwise the number must be an
	// odd number, the bigger the number the more blur
	BlurKernelSize int

	// EdgeThreshold is a number between 0 and 255 that specifies a cutoff point to
	// determine if an intensity change is an edge. Make smaller to include more details
	// as edges
	EdgeThreshold int

	// OilFilterSize specifies how bold the simulated strokes will be when turning the
	// style towards a painting, something around 5,10,15 should work well
	OilFilterSize int

	// OilLevels is the number of levels that the oil painting style will bucket colors in
	// to. Larger number to get more detail.
	OilLevels int

	// DebugPath is not empty is assumed to be a path where intermediate debug files can
	// be written to, such as the gaussian blured image and the sobel edge detection. This
	// can be useful for tweaking parameters
	DebugPath string
}

type cartoon struct {
	opts CTOpts
}

// Apply runs the image through the cartoon filter
func (c *cartoon) Apply(img *Image, numRoutines int) (*Image, error) {
	if numRoutines == 0 {
		numRoutines = runtime.GOMAXPROCS(0)
	}

	pipeline := Pipeline{}
	if c.opts.BlurKernelSize > 0 {
		pipeline.Add(NewGaussian(c.opts.BlurKernelSize, 1), nil)
	}
	pipeline.Add(NewGrayscale(GSLUMINOSITY), nil)
	pipeline.Add(NewSobel(c.opts.EdgeThreshold, false), nil)
	edgeImg, err := pipeline.Run(img, numRoutines)
	if err != nil {
		return nil, err
	}
	edgePix := edgeImg.img.Pix
	pf := func(ri, x, y, offset, inStride int, inPix, outPix []uint8) {

		r := inPix[offset]
		g := inPix[offset+1]
		b := inPix[offset+2]

		rEdge := edgePix[offset]
		if rEdge == 255 {
			r = 0
			b = 0
			g = 0
		}

		outPix[offset] = r
		outPix[offset+1] = g
		outPix[offset+2] = b
		outPix[offset+3] = 255
	}

	oil := NewOilPainting(c.opts.OilFilterSize, c.opts.OilLevels)
	oilImg, err := oil.Apply(img, numRoutines)
	if err != nil {
		return nil, err
	}

	out := &Image{
		img: image.NewRGBA(image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: img.Width, Y: img.Height},
		}),
		Width:  img.Width,
		Height: img.Height,

		// Have to take in to account pixels are lost in some of the effects around the edges,
		// so so only have the area where the two rections intersect from the edge detection and
		// the oil painting effect
		Bounds: oilImg.Bounds.Intersect(edgeImg.Bounds),
	}
	runParallel(numRoutines, oilImg, out.Bounds, out, pf, 0)
	return out, nil
}

// NewCartoon returns an effect that renders images as if they are drawn like a cartoon.
// It works by rendering the input image using the OilPainting effect, then drawing lines
// ontop of the image based on the Sobel edge detection method. You will probably have to
// play with the opts values to get a good result. Some starting values are:
// BlurKernelSize: 21
// EdgeThreshold: 40
// OilFilterSize: 15
// OilLevels: 15
func NewCartoon(opts CTOpts) Effect {
	return &cartoon{
		opts: opts,
	}
}
