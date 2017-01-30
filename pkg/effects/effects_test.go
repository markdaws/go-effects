package effects_test

import (
	"fmt"
	"testing"

	"github.com/markdaws/go-effects/pkg/effects"
	"github.com/stretchr/testify/require"
)

const cabinPath = "../../test/cabin.jpg"

func TestOilPainting(t *testing.T) {
	timing := effects.NewTiming()

	timing.Time("load")
	img, err := effects.LoadImage(cabinPath)
	timing.TimeEnd("load")
	require.Nil(t, err)
	require.NotNil(t, img)

	timing.Time("oil-serial")
	oil := effects.NewOilPainting(5, 30)
	oilImg, err := oil.Apply(img, 1)
	timing.TimeEnd("oil-serial")
	require.Nil(t, err)
	require.NotNil(t, oilImg)

	err = oilImg.Save("../../test/cabin-oil.jpg", effects.SaveOpts{ClipToBounds: true})
	require.Nil(t, err)

	err = oilImg.Save("../../test/cabin-oil.png", effects.SaveOpts{ClipToBounds: true})
	require.Nil(t, err)

	timing.Time("oil-parallel")
	oilImg, err = oil.Apply(img, 0)
	timing.TimeEnd("oil-parallel")
	require.Nil(t, err)
	require.NotNil(t, oilImg)

	err = oilImg.Save("../../test/cabin-parallel-oil.jpg", effects.SaveOpts{ClipToBounds: true})
	require.Nil(t, err)

	fmt.Println(img.Bounds)
	fmt.Println(timing)
}

func TestGrayscale(t *testing.T) {
	timing := effects.NewTiming()

	timing.Time("load")
	img, err := effects.LoadImage(cabinPath)
	timing.TimeEnd("load")
	require.Nil(t, err)
	require.NotNil(t, img)

	timing.Time("grayscale-average")
	gsAverage := effects.NewGrayscale(effects.GSAVERAGE)
	grayImg, err := gsAverage.Apply(img, 1)
	timing.TimeEnd("grayscale-average")
	require.Nil(t, err)
	require.NotNil(t, grayImg)
	err = grayImg.Save("../../test/cabin-gray-average.jpg", effects.SaveOpts{})
	require.Nil(t, err)

	timing.Time("grayscale-lightness")
	gsLightness := effects.NewGrayscale(effects.GSLIGHTNESS)
	grayImg, err = gsLightness.Apply(img, 1)
	timing.TimeEnd("grayscale-lightness")
	require.Nil(t, err)
	require.NotNil(t, grayImg)
	err = grayImg.Save("../../test/cabin-gray-lightness.jpg", effects.SaveOpts{})
	require.Nil(t, err)

	timing.Time("grayscale-luminosity")
	gsLuminosity := effects.NewGrayscale(effects.GSLUMINOSITY)
	grayImg, err = gsLuminosity.Apply(img, 1)
	timing.TimeEnd("grayscale-luminosity")
	require.Nil(t, err)
	require.NotNil(t, grayImg)
	err = grayImg.Save("../../test/cabin-gray-luminosity.jpg", effects.SaveOpts{})
	require.Nil(t, err)

	timing.Time("grayscale-parallel-luminosity")
	grayImg, err = gsLuminosity.Apply(img, 0)
	timing.TimeEnd("grayscale-parallel-luminosity")
	require.Nil(t, err)
	require.NotNil(t, grayImg)
	err = grayImg.Save("../../test/cabin-gray-parallel-luminosity.jpg", effects.SaveOpts{})
	require.Nil(t, err)

	fmt.Println(img.Bounds)
	fmt.Println(timing)
}

func TestSobel(t *testing.T) {
	timing := effects.NewTiming()

	timing.Time("load")
	img, err := effects.LoadImage("../../test/turtle.jpg")
	timing.TimeEnd("load")
	require.Nil(t, err)
	require.NotNil(t, img)

	timing.Time("grayscale-luminosity")
	grayscale := effects.NewGrayscale(effects.GSLUMINOSITY)
	grayImg, err := grayscale.Apply(img, 1)
	timing.TimeEnd("grayscale-luminosity")
	require.Nil(t, err)
	require.NotNil(t, grayImg)

	// The sobel image contains all of the intensity values, since we pass -1
	timing.Time("sobel")
	sobel := effects.NewSobel(-1, false)
	sobelImg, err := sobel.Apply(img, 0)
	require.Nil(t, err)
	require.NotNil(t, sobelImg)
	timing.TimeEnd("sobel")

	err = sobelImg.Save("../../test/turtle-sobel.jpg", effects.SaveOpts{ClipToBounds: true})
	require.Nil(t, err)

	// The sobel image contains pixels of value either 255 or 0, 255 if the sobel gradient is
	// >= threshold, 0 otherwise
	timing.Time("sobel-threshold-200")
	sobel = effects.NewSobel(200, false)
	sobelImg, err = sobel.Apply(img, 0)
	require.Nil(t, err)
	require.NotNil(t, sobelImg)
	timing.TimeEnd("sobel-threshold-200")

	err = sobelImg.Save("../../test/turtle-sobel-threshold-200.jpg", effects.SaveOpts{ClipToBounds: true})
	require.Nil(t, err)

	fmt.Println(img.Bounds)
	fmt.Println(timing)
}

func TestPencil(t *testing.T) {
	timing := effects.NewTiming()

	timing.Time("load")
	img, err := effects.LoadImage("../../test/houses.jpg")
	timing.TimeEnd("load")
	require.Nil(t, err)
	require.NotNil(t, img)

	timing.Time("grayscale-luminosity")
	gs := effects.NewGrayscale(effects.GSLUMINOSITY)
	grayImg, err := gs.Apply(img, 1)
	timing.TimeEnd("grayscale-luminosity")
	require.Nil(t, err)
	require.NotNil(t, grayImg)

	// Pencil is just the sobel inverted with no thresholding
	timing.Time("pencil")
	pencil := effects.NewPencil(5)
	pencilImg, err := pencil.Apply(img, 0)
	require.Nil(t, err)
	require.NotNil(t, pencilImg)
	timing.TimeEnd("pencil")

	err = pencilImg.Save("../../test/houses-pencil.jpg", effects.SaveOpts{ClipToBounds: true})
	require.Nil(t, err)

	fmt.Println(img.Bounds)
	fmt.Println(timing)
}

func TestGaussian(t *testing.T) {
	timing := effects.NewTiming()

	timing.Time("load")
	img, err := effects.LoadImage("../../test/face.jpg")
	timing.TimeEnd("load")
	require.Nil(t, err)
	require.NotNil(t, img)

	timing.Time("gaussian")
	effect := effects.NewGaussian(21, 1)
	gaussianImg, err := effect.Apply(img, 0)
	timing.TimeEnd("gaussian")
	require.Nil(t, err)
	require.NotNil(t, gaussianImg)
	err = gaussianImg.Save("../../test/face-gaussian.jpg", effects.SaveOpts{ClipToBounds: true})
	require.Nil(t, err)

	fmt.Println(img.Bounds)
	fmt.Println(timing)
}

func TestCartoon(t *testing.T) {
	timing := effects.NewTiming()

	timing.Time("load")
	img, err := effects.LoadImage("../../test/turtle.jpg")
	timing.TimeEnd("load")
	require.Nil(t, err)
	require.NotNil(t, img)

	timing.Time("cartoon")
	opts := effects.CTOpts{
		BlurKernelSize: 21,
		EdgeThreshold:  40,
		OilFilterSize:  20,
		OilLevels:      12,
		DebugPath:      "../../test",
	}
	effect := effects.NewCartoon(opts)
	cartoonImg, err := effect.Apply(img, 0)
	require.Nil(t, err)
	require.NotNil(t, cartoonImg)
	timing.TimeEnd("cartoon")

	err = cartoonImg.Save("../../test/turtle-cartoon.jpg", effects.SaveOpts{ClipToBounds: true})
	require.Nil(t, err)

	fmt.Println(img.Bounds)
	fmt.Println(timing)
}

func TestPixelate(t *testing.T) {
	timing := effects.NewTiming()

	timing.Time("load")
	img, err := effects.LoadImage("../../test/turtle.jpg")
	timing.TimeEnd("load")
	require.Nil(t, err)
	require.NotNil(t, img)

	timing.Time("pixelate")
	pixelate := effects.NewPixelate(20)
	pixelImg, err := pixelate.Apply(img, 0)
	require.Nil(t, err)
	require.NotNil(t, pixelImg)
	timing.TimeEnd("pixelate")

	err = pixelImg.Save("../../test/turtle-20-pixelate.jpg", effects.SaveOpts{})
	require.Nil(t, err)

	fmt.Println(img.Bounds)
	fmt.Println(timing)
}
