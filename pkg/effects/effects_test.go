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

	err = oilImg.SaveAsJPG("../../test/cabin-oil.jpg", 90)
	require.Nil(t, err)

	err = oilImg.SaveAsPNG("../../test/cabin-oil.png")
	require.Nil(t, err)

	timing.Time("oil-parallel")
	oilImg, err = effects.OilPainting(img, 0, 5, 30)
	timing.TimeEnd("oil-parallel")
	require.Nil(t, err)
	require.NotNil(t, oilImg)

	err = oilImg.SaveAsPNG("../../test/cabin-parallel-oil.png")
	require.Nil(t, err)

	fmt.Println(img.Bounds())
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
	err = grayImg.SaveAsPNG("../../test/cabin-gray-average.png")
	require.Nil(t, err)

	timing.Time("grayscale-lightness")
	grayImg, err = effects.Grayscale(img, 1, effects.GSLIGHTNESS)
	timing.TimeEnd("grayscale-lightness")
	require.Nil(t, err)
	require.NotNil(t, grayImg)
	err = grayImg.SaveAsPNG("../../test/cabin-gray-lightness.png")
	require.Nil(t, err)

	timing.Time("grayscale-luminosity")
	grayImg, err = effects.Grayscale(img, 1, effects.GSLUMINOSITY)
	timing.TimeEnd("grayscale-luminosity")
	require.Nil(t, err)
	require.NotNil(t, grayImg)
	err = grayImg.SaveAsPNG("../../test/cabin-gray-luminosity.png")
	require.Nil(t, err)

	timing.Time("grayscale-parallel-luminosity")
	grayImg, err = effects.Grayscale(img, 0, effects.GSLUMINOSITY)
	timing.TimeEnd("grayscale-parallel-luminosity")
	require.Nil(t, err)
	require.NotNil(t, grayImg)
	err = grayImg.SaveAsPNG("../../test/cabin-gray-parallel-luminosity.png")
	require.Nil(t, err)

	fmt.Println(img.Bounds())
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

	err = sobelImg.SaveAsPNG("../../test/turtle-sobel.png")
	require.Nil(t, err)

	// The sobel image contains pixels of value either 255 or 0, 255 if the sobel gradient is
	// >= threshold, 0 otherwise
	timing.Time("sobel-threshold-200")
	sobelImg, err = effects.Sobel(img, 0, 200)
	require.Nil(t, err)
	require.NotNil(t, sobelImg)
	timing.TimeEnd("sobel-threshold-200")

	err = sobelImg.SaveAsPNG("../../test/turtle-sobel-threshold-200.png")
	require.Nil(t, err)

	fmt.Println(img.Bounds())
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
	err = gaussianImg.SaveAsPNG("../../test/face-gaussian.png")
	require.Nil(t, err)

	fmt.Println(img.Bounds())
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
	opts := effects.CTOptions{
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

	err = cartoonImg.SaveAsPNG("../../test/turtle-cartoon.png")
	require.Nil(t, err)

	fmt.Println(img.Bounds())
	fmt.Println(timing)
}
