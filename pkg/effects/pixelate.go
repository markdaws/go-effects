package effects

import (
	"fmt"
	"image"
	"runtime"
)

func Pixelate(img *Image, numRoutines, blockSize int) (*Image, error) {
	if numRoutines == 0 {
		numRoutines = runtime.GOMAXPROCS(0)
	}

	if img.Bounds().Dx()%blockSize != 0 ||
		img.Bounds().Dy()%blockSize != 0 {
		return nil, fmt.Errorf("blockSize must divide exactly into the width and the height of the input image")
	}

	out := &Image{img: image.NewRGBA(img.img.Bounds())}

	nBlocksX := img.Bounds().Dx() / blockSize
	nBlocksY := img.Bounds().Dy() / blockSize
	nBlocks := nBlocksX * nBlocksY

	blocksR := make([]int, nBlocks)
	blocksG := make([]int, nBlocks)
	blocksB := make([]int, nBlocks)
	pixelsPerBlock := blockSize * blockSize

	pfCalc := func(ri, x, y, offset, inStride int, inPix, outPix []uint8) {
		r := inPix[offset]
		g := inPix[offset+1]
		b := inPix[offset+2]

		blockIndex := (y/blockSize)*nBlocksX + (x / blockSize)
		blocksR[blockIndex] += int(r)
		blocksG[blockIndex] += int(g)
		blocksB[blockIndex] += int(b)
	}

	pfSet := func(ri, x, y, offset, inStride int, inPix, outPix []uint8) {
		blockIndex := (y/blockSize)*nBlocksX + (x / blockSize)

		outPix[offset] = uint8(blocksR[blockIndex])
		outPix[offset+1] = uint8(blocksG[blockIndex])
		outPix[offset+2] = uint8(blocksB[blockIndex])
		outPix[offset+3] = 255
	}

	var pixelsPerRoutine int

	// Make sure the goroutines process on a block boundary
	pixelsPerRoutine = ((img.Bounds().Dx() / numRoutines) / blockSize) * blockSize
	runParallel(numRoutines, img, img.Bounds(), out, pfCalc, pixelsPerRoutine)

	// Divide by number of pixels
	for i := 0; i < nBlocks; i++ {
		blocksR[i] /= pixelsPerBlock
		blocksG[i] /= pixelsPerBlock
		blocksB[i] /= pixelsPerBlock
	}

	runParallel(numRoutines, img, img.Bounds(), out, pfSet, pixelsPerRoutine)
	return out, nil
}
