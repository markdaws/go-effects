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
	oilImg, err := effects.OilPainting(img, 1, 5, 30)
	timing.TimeEnd("oil-serial")
	require.Nil(t, err)
	require.NotNil(t, oilImg)

	err = oilImg.Save("../../test/cabin-oil.jpg", effects.SaveOpts{ClipToBounds: true})
	require.Nil(t, err)

	err = oilImg.Save("../../test/cabin-oil.png", effects.SaveOpts{ClipToBounds: true})
	require.Nil(t, err)

	timing.Time("oil-parallel")
	oilImg, err = effects.OilPainting(img, 0, 5, 30)
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
	grayImg, err := effects.Grayscale(img, 1, effects.GSAVERAGE)
	timing.TimeEnd("grayscale-average")
	require.Nil(t, err)
	require.NotNil(t, grayImg)
	err = grayImg.Save("../../test/cabin-gray-average.jpg", effects.SaveOpts{})
	require.Nil(t, err)

	timing.Time("grayscale-lightness")
	grayImg, err = effects.Grayscale(img, 1, effects.GSLIGHTNESS)
	timing.TimeEnd("grayscale-lightness")
	require.Nil(t, err)
	require.NotNil(t, grayImg)
	err = grayImg.Save("../../test/cabin-gray-lightness.jpg", effects.SaveOpts{})
	require.Nil(t, err)

	timing.Time("grayscale-luminosity")
	grayImg, err = effects.Grayscale(img, 1, effects.GSLUMINOSITY)
	timing.TimeEnd("grayscale-luminosity")
	require.Nil(t, err)
	require.NotNil(t, grayImg)
	err = grayImg.Save("../../test/cabin-gray-luminosity.jpg", effects.SaveOpts{})
	require.Nil(t, err)

	timing.Time("grayscale-parallel-luminosity")
	grayImg, err = effects.Grayscale(img, 0, effects.GSLUMINOSITY)
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
	grayImg, err := effects.Grayscale(img, 1, effects.GSLUMINOSITY)
	timing.TimeEnd("grayscale-luminosity")
	require.Nil(t, err)
	require.NotNil(t, grayImg)

	// The sobel image contains all of the intensity values, since we pass -1
	timing.Time("sobel")
	sobelImg, err := effects.Sobel(img, 0, -1)
	require.Nil(t, err)
	require.NotNil(t, sobelImg)
	timing.TimeEnd("sobel")

	err = sobelImg.Save("../../test/turtle-sobel.jpg", effects.SaveOpts{ClipToBounds: true})
	require.Nil(t, err)

	// The sobel image contains pixels of value either 255 or 0, 255 if the sobel gradient is
	// >= threshold, 0 otherwise
	timing.Time("sobel-threshold-200")
	sobelImg, err = effects.Sobel(img, 0, 200)
	require.Nil(t, err)
	require.NotNil(t, sobelImg)
	timing.TimeEnd("sobel-threshold-200")

	err = sobelImg.Save("../../test/turtle-sobel-threshold-200.jpg", effects.SaveOpts{ClipToBounds: true})
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
	gaussianImg, err := effects.Gaussian(img, 0, 21, 1)
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
	cartoonImg, err := effects.Cartoon(img, 0, opts)
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
	pixelImg, err := effects.Pixelate(img, 0, 20)
	require.Nil(t, err)
	require.NotNil(t, pixelImg)
	timing.TimeEnd("pixelate")

	err = pixelImg.Save("../../test/turtle-20-pixelate.jpg", effects.SaveOpts{})
	require.Nil(t, err)

	fmt.Println(img.Bounds)
	fmt.Println(timing)
}
