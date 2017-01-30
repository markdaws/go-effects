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
)

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
}

// saveAsJPG saves the image as a JPG. quality is between 1 and 100, 100 being best
func (i *Image) saveAsJPG(path string, quality int) error {
	toImg, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create image: %s, %s", path, err)
	}
	err = jpeg.Encode(toImg, i.img, &jpeg.Options{Quality: quality})
	if err != nil {
		return fmt.Errorf("failed to encode image: %s, %s", path, err)
	}
	err = toImg.Close()
	if err != nil {
		return fmt.Errorf("failed to close image: %s, %s", path, err)
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
	if err != nil {
		return fmt.Errorf("failed to encode image: %s, %s", path, err)
	}
	err = toImg.Close()
	if err != nil {
		return fmt.Errorf("failed to close image: %s, %s", path, err)
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
	if err != nil {
		return nil, fmt.Errorf("failed to decode image on load: %s, %s", path, err)
	}
	err = srcReader.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close image on load: %s, %s", path, err)
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
