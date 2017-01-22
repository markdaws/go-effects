package effects

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path"
	"strings"
	"sync"
)

//TODO: Block effect
//TODO: Pointilize effect
//TODO: Stained Glass

// Image wrapper around internal pixels
type Image struct {
	img *image.RGBA
}

// Bounds returns the bounds of the pixels in the image
func (i *Image) Bounds() image.Rectangle {
	return i.img.Bounds()
}

// Save saves the image as the file type defined by the extension in the path e.g. ,jpg or .png
func (i *Image) Save(outPath string) error {
	ext := strings.ToLower(path.Ext(outPath))

	switch path.Ext(outPath) {
	case ".jpg", ".jpeg":
		return i.SaveAsJPG(outPath, 90)
	case ".png":
		return i.SaveAsPNG(outPath)
	default:
		return fmt.Errorf("unsupported file type: %s", ext)
	}
	return nil
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
