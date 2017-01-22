package effects

import (
	"image"
	"path"
	"runtime"
)

// CTOptions options to pass to the Cartoon effect
type CTOptions struct {
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

// Cartoon renders the image stylized as a cartoon. It works by rendering the input image
// using the OilPainting effect, then drawing lines ontop of hte image based on the Sobel
// edge detection method. You will probably have to play with the opts values to get a good
// result. Some starting values are:
// BlurKernelSize: 21
// EdgeThreshold: 40
// OilFilterSize: 15
// OilLevels: 15
func Cartoon(img *Image, numRoutines int, opts CTOptions) (*Image, error) {
	if numRoutines == 0 {
		numRoutines = runtime.GOMAXPROCS(0)
	}

	inImg := img
	if opts.BlurKernelSize > 0 {
		var err error
		inImg, err = Gaussian(img, numRoutines, opts.BlurKernelSize, 1)
		if err != nil {
			return nil, err
		}

		if opts.DebugPath != "" {
			err = inImg.Save(path.Join(opts.DebugPath, "cartoon-gaussian.jpg"))
			if err != nil {
				return nil, err
			}
		}
	}

	grayImg, err := Grayscale(inImg, numRoutines, GSLUMINOSITY)
	if err != nil {
		return nil, err
	}

	edgeImg, err := Sobel(grayImg, numRoutines, opts.EdgeThreshold)
	if err != nil {
		return nil, err
	}
	if opts.DebugPath != "" {
		err = edgeImg.Save(path.Join(opts.DebugPath, "cartoon-edge.jpg"))
		if err != nil {
			return nil, err
		}
	}

	oilImg, err := OilPainting(img, numRoutines, opts.OilFilterSize, opts.OilLevels)
	if err != nil {
		return nil, err
	}
	if opts.DebugPath != "" {
		err = oilImg.Save(path.Join(opts.DebugPath, "cartoon-oil.jpg"))
		if err != nil {
			return nil, err
		}
	}

	out := &Image{img: image.NewRGBA(img.img.Bounds())}
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

	inBounds := img.Bounds()
	runParallel(numRoutines, oilImg, inBounds, out, pf, -1)
	return out, nil
}
