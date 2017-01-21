package effects

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"runtime"
	"sync"
)

// Image wrapper around internal pixels
type Image struct {
	img *image.RGBA
}

// Bounds returns the bounds of the pixels in the image
func (i *Image) Bounds() image.Rectangle {
	return i.img.Bounds()
}

// SaveAsJPG saves the image as a JPG. quality is between 1 and 100, 100 being best
func (i *Image) SaveAsJPG(path string, quality int) error {
	toImg, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create image: %s, %s", path, err)
	}
	err = jpeg.Encode(toImg, i.img, &jpeg.Options{Quality: quality})
	toImg.Close()
	if err != nil {
		return fmt.Errorf("failed to save image: %s, %s", path, err)
	}
	return nil
}

// SaveAsPNG saves the image as a PNG
func (i *Image) SaveAsPNG(path string) error {
	toImg, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create image: %s, %s", path, err)
	}

	err = png.Encode(toImg, i.img)
	toImg.Close()
	if err != nil {
		return fmt.Errorf("failed to save image: %s, %s", path, err)
	}
	return nil
}

// LoadImage loads the specified image from disk. Supported file types are png and jpg
func LoadImage(path string) (*Image, error) {
	srcReader, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read input image: %s, %s", path, err)
	}

	img, _, err := image.Decode(srcReader)
	srcReader.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to decode image on load: %s, %s", path, err)
	}

	outImg := image.NewRGBA(img.Bounds())
	draw.Draw(outImg, img.Bounds(), img, image.Point{}, draw.Over)

	return &Image{img: outImg}, nil
}

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

type pixelFunc func(ri, x, y, offset, inStride int, inPix, outPix []uint8)

func runParallel(numRoutines int, inImg *Image, inBounds image.Rectangle, outImg *Image, pf pixelFunc) {
	w := inBounds.Dx()
	h := inBounds.Dy()

	minX := inBounds.Min.X
	minY := inBounds.Min.Y
	stride := inImg.img.Stride
	start := minY*stride + minX*4
	inPix := inImg.img.Pix
	outPix := outImg.img.Pix

	wg := sync.WaitGroup{}
	xOffset := minX
	widthPerRoutine := w / numRoutines

	for r := 0; r < numRoutines; r++ {
		wg.Add(1)

		go func(ri, xStart, yStart, width, height int) {
			for x := xStart; x < xStart+width; x++ {
				for y := yStart; y < yStart+height; y++ {
					offset := start + y*stride + x*4
					pf(ri, x, y, offset, stride, inPix, outPix)
				}
			}
			wg.Done()
		}(r, xOffset, minY, widthPerRoutine, h)

		xOffset += widthPerRoutine
		if r == numRoutines-1 {
			widthPerRoutine = w - xOffset
		}
	}
	wg.Wait()
}

// Grayscale renders the input image as a grayscale image. numRoutines specifies how many
// goroutines should be used to process the image in parallel, use 0 to let the library decide
func Grayscale(img *Image, numRoutines int, algo GSAlgo) (*Image, error) {

	if numRoutines == 0 {
		numRoutines = runtime.GOMAXPROCS(0)
	}

	pf := func(ri, x, y, offset, inStride int, inPix, outPix []uint8) {
		var r, g, b uint8 = inPix[offset], inPix[offset+1], inPix[offset+2]
		switch algo {
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

	out := &Image{img: image.NewRGBA(img.img.Bounds())}
	runParallel(numRoutines, img, img.Bounds(), out, pf)
	return out, nil
}

// Sobel - the input image should be a grayscale image, the output will be a version of
// the input image with the Sobel edge detector applied to it.
func Sobel(img *Image, numRoutines int) (*Image, error) {
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

// OilPainting renders the input image as if it was painted like an oil painting. numRoutines specifies how many
// goroutines should be used to process the image in parallel, use 0 to let the library decide. filterSize specifies
// how bold the image should look, larger numbers equate to larger strokes, levels specifies how many buckets colors
// will be grouped in to, start with values 5,30 to see how that works.
func OilPainting(img *Image, numRoutines, filterSize, levels int) (*Image, error) {
	out := &Image{img: image.NewRGBA(img.img.Bounds())}
	levels = levels - 1
	filterOffset := (filterSize - 1) / 2

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

	inBounds := image.Rectangle{
		Min: image.Point{X: filterOffset, Y: filterOffset},
		Max: image.Point{X: img.Bounds().Dx() - 2*filterOffset, Y: img.Bounds().Dy() - 2*filterOffset},
	}

	runParallel(numRoutines, img, inBounds, out, pf)
	return out, nil
}

func roundToInt32(a float64) int32 {
	if a < 0 {
		return int32(a - 0.5)
	}
	return int32(a + 0.5)
}

func reset(s []int) {
	for i := range s {
		s[i] = 0
	}
}
