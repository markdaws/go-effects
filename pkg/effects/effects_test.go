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