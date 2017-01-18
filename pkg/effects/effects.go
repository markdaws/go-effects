package effects

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
)

type Image struct {
	img *image.RGBA
}

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

func OilPainting(img *Image, filterSize, levels int, cropOutput bool) (*Image, error) {
	w := img.img.Bounds().Dx()
	h := img.img.Bounds().Dy()

	// We lose some of the image due to the filter, if we want to crop then we will not write
	// the empty pixels
	bounds := img.img.Bounds()
	var out *image.RGBA
	if cropOutput {
		bounds = image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: bounds.Max.X - filterSize, Y: bounds.Max.Y - filterSize},
		}
	}
	out = image.NewRGBA(bounds)

	levels = levels - 1
	filterOffset := (filterSize - 1) / 2

	iBin := make([]int, levels+1)
	rBin := make([]int, levels+1)
	gBin := make([]int, levels+1)
	bBin := make([]int, levels+1)

	for y := filterOffset; y < h-filterOffset; y++ {
		for x := filterOffset; x < w-filterOffset; x++ {

			var maxIntensity int
			var maxIndex int

			reset(iBin)
			reset(rBin)
			reset(gBin)
			reset(bBin)

			for fy := -filterOffset; fy <= filterOffset; fy++ {
				for fx := -filterOffset; fx <= filterOffset; fx++ {

					p := img.img.At(x+fx, y+fy).(color.RGBA)

					ci := int(roundToInt32((float64(p.R+p.G+p.B) / 3.0 * float64(levels)) / 255.0))
					iBin[ci] += 1
					rBin[ci] += int(p.R)
					gBin[ci] += int(p.G)
					bBin[ci] += int(p.B)

					if iBin[ci] > maxIntensity {
						maxIntensity = iBin[ci]
						maxIndex = ci
					}
				}
			}

			r := rBin[maxIndex] / maxIntensity
			g := gBin[maxIndex] / maxIntensity
			b := bBin[maxIndex] / maxIntensity

			if cropOutput {
				out.Set(x-filterOffset, y-filterOffset, color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255})
			} else {
				out.Set(x, y, color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255})
			}
		}
	}

	return &Image{img: out}, nil
}

func roundToInt32(a float64) int32 {
	if a < 0 {
		return int32(a - 0.5)
	}
	return int32(a + 0.5)
}

func reset(s []int) {
	for i, _ := range s {
		s[i] = 0
	}
}
