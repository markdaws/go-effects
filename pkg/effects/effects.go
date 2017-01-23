package effects

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"path"
	"strings"
	"sync"
)

//TODO: Pointilize effect
//TODO: Stained Glass

// Rect used for image bounds
type Rect struct {
	X      int
	Y      int
	Width  int
	Height int
}

func (r Rect) String() string {
	return fmt.Sprintf("X:%d, Y:%d, Width:%d, Height:%d", r.X, r.Y, r.Width, r.Height)
}
func (r Rect) Intersect(r2 Rect) Rect {
	x := math.Max(float64(r.X), float64(r2.X))
	num1 := math.Min(float64(r.X+r.Width), float64(r2.X+r2.Width))

	y := math.Max(float64(r.Y), float64(r2.Y))
	num2 := math.Min(float64(r.Y+r.Height), float64(r2.Y+r2.Height))
	if num1 >= x && num2 >= y {
		return Rect{X: int(x), Y: int(y), Width: int(num1 - x), Height: int(num2 - y)}
	} else {
		return Rect{}
	}
}
func (r Rect) IsEmpty() bool {
	return r.Width == 0 || r.Height == 0
}
func (r Rect) ToImageRect() image.Rectangle {
	return image.Rectangle{
		Min: image.Point{X: r.X, Y: r.Y},
		Max: image.Point{X: r.X + r.Width, Y: r.Y + r.Height},
	}
}

// Image wrapper around internal pixels
type Image struct {
	img    *image.RGBA
	Bounds Rect
	Width  int
	Height int
}

// SaveOpts specifies some save parameters that can be specified when saving
// an image
type SaveOpts struct {
	// JPEGCompression a value between 1 and 100, if 0 specified, defaults to 95.
	// Higher values are better quality. Only applicable if the file ends with a
	// .jpg or .jpeg extension
	JPEGCompression int

	// ClipToBounds if true only the image inside the region specified by the bounds is
	// saved. Sometimes aftere running an image effect you may have outer bands that are
	// dead pixels, setting this to true crops them out.
	ClipToBounds bool
}

// Save saves the image as the file type defined by the extension in the path e.g. ,jpg or .png
func (i *Image) Save(outPath string, opts SaveOpts) error {
	ext := strings.ToLower(path.Ext(outPath))

	final := i
	if opts.ClipToBounds {
		final = &Image{
			img: image.NewRGBA(image.Rectangle{
				Min: image.Point{X: 0, Y: 0},
				Max: image.Point{X: i.Bounds.Width, Y: i.Bounds.Height},
			}),
			Width:  i.Bounds.Width,
			Height: i.Bounds.Height,
			Bounds: Rect{X: 0, Y: 0, Width: i.Bounds.Width, Height: i.Bounds.Height},
		}

		draw.Draw(final.img, final.Bounds.ToImageRect(), i.img, i.Bounds.ToImageRect().Min, draw.Src)
	}

	switch path.Ext(outPath) {
	case ".jpg", ".jpeg":
		cmpLvl := opts.JPEGCompression
		if cmpLvl == 0 {
			cmpLvl = 95
		}
		return final.saveAsJPG(outPath, cmpLvl)
	case ".png":
		return final.saveAsPNG(outPath)
	default:
		return fmt.Errorf("unsupported file type: %s", ext)
	}
	return nil
}

// saveAsJPG saves the image as a JPG. quality is between 1 and 100, 100 being best
func (i *Image) saveAsJPG(path string, quality int) error {
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

// saveAsPNG saves the image as a PNG
func (i *Image) saveAsPNG(path string) error {
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

	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y
	return &Image{
		img:    outImg,
		Width:  w,
		Height: h,
		Bounds: Rect{X: 0, Y: 0, Width: w, Height: h},
	}, nil
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

func isOddInt(i int) bool {
	return i%2 != 0
}

func reset(s []int) {
	for i := range s {
		s[i] = 0
	}
}
