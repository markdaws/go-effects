package effects_test

import (
	"testing"

	"github.com/markdaws/go-effects/pkg/effects"
	"github.com/stretchr/testify/require"
)

const cabinPath = "../../test/cabin.jpg"

func TestOilPainting(t *testing.T) {
	img, err := effects.LoadImage(cabinPath)
	require.Nil(t, err)
	require.NotNil(t, img)

	oilImg, err := effects.OilPainting(img, 5, 30, true)
	require.Nil(t, err)
	require.NotNil(t, oilImg)

	err = oilImg.SaveAsJPG("../../test/cabin-oil.jpg", 90)
	require.Nil(t, err)

	err = oilImg.SaveAsPNG("../../test/cabin-oil.png")
	require.Nil(t, err)
}
