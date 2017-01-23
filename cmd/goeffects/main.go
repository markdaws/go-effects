package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/markdaws/go-effects/pkg/effects"
)

func main() {
	effect := flag.String("effect", "", "The name of the effect to apply. Values are 'oil|sobel|gaussian|cartoon|pixelate'")
	flag.Parse()

	switch *effect {
	case "oil":
		if len(flag.Args()) != 4 {
			fmt.Println("The oil effect requires 4 args, input path, output path, filterSize, levels\n")
			fmt.Println("Sample usage: goeffects -effect=oil mypic.jpg mypic-oil.jpg 5 30\n")
			flag.PrintDefaults()
			os.Exit(1)
		}
	case "sobel":
		if len(flag.Args()) != 3 {
			fmt.Println("The sobel effect requires 3 args, input path, output path, threshold\n")
			fmt.Println("Sample usage: goeffects -effect=sobel mypic.jpg mypic-sobel.jpg 100\n")
			flag.PrintDefaults()
			os.Exit(1)
		}
	case "gaussian":
		if len(flag.Args()) != 4 {
			fmt.Println("The gaussian effect requires 4 args, input path, output path, kernelSize, sigma\n")
			fmt.Println("Sample usage: goeffects -effect=gaussian mypic.jpg mypic-gaussian.jpg 9 1\n")
			flag.PrintDefaults()
			os.Exit(1)
		}
	case "cartoon":
		if len(flag.Args()) != 6 {
			fmt.Println("The cartoon effect requires 6 args, input path, output path, blurStrength, edgeThreshold, oilBoldness, oilLevels")
			fmt.Println("Sample usage: goeffects -effect=cartoon mypic.jpg mypic-cartoon.jpg 21 40 15 15")
			flag.PrintDefaults()
			os.Exit(1)
		}
	case "pixelate":
		if len(flag.Args()) != 3 {
			fmt.Println("The pixelate effect requires 3 args, input path, output path, block size")
			fmt.Println("Sample usage: goeffects -effect=pixelate mypic.jpg mypic-pixelate.jpg 12")
			flag.PrintDefaults()
			os.Exit(1)
		}
	case "":
		fmt.Println("The effect option is required\n")
		flag.PrintDefaults()
		os.Exit(1)

	default:
		fmt.Println("Unknown effect option value\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	var inPath, outPath string
	inPath = flag.Arg(0)
	outPath = flag.Arg(1)

	img, err := effects.LoadImage(inPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var outImg *effects.Image
	switch *effect {
	case "gaussian":
		kernelSize, err := strconv.Atoi(flag.Arg(2))
		if err != nil {
			fmt.Println("invalid kernelSize value")
			os.Exit(1)
		}
		sigma, err := strconv.ParseFloat(flag.Arg(3), 64)
		if err != nil {
			fmt.Println("invalid sigma value")
			os.Exit(1)
		}
		outImg, err = effects.Gaussian(img, 0, kernelSize, sigma)
		if err != nil {
			fmt.Println("Failed to apply effect:", err)
			os.Exit(1)
		}
	case "sobel":
		threshold, err := strconv.Atoi(flag.Arg(2))
		if err != nil {
			fmt.Println("invalid threshold value")
			os.Exit(1)
		}
		outImg, err = effects.Sobel(img, 0, threshold)
		if err != nil {
			fmt.Println("Failed to apply effect:", err)
			os.Exit(1)
		}
	case "oil":
		filterSize, err := strconv.Atoi(flag.Arg(0))
		if err != nil {
			fmt.Println("Invalid filterSize value:", err)
			os.Exit(1)
		}
		if filterSize <= 3 {
			fmt.Println("FilterSize must be at least 3")
			os.Exit(1)
		}

		levels, err := strconv.Atoi(flag.Arg(1))
		if err != nil {
			fmt.Println("Invalid levels value:", err)
			os.Exit(1)
		}
		if levels < 1 {
			fmt.Println("Levels must be at least 1")
			os.Exit(1)
		}

		outImg, err = effects.OilPainting(img, 0, filterSize, levels)
		if err != nil {
			fmt.Println("Failed to apply effect:", err)
			os.Exit(1)
		}
	case "cartoon":
		blurStrength, err := strconv.Atoi(flag.Arg(2))
		if err != nil {
			fmt.Println("Invalid blurStrength value")
			os.Exit(1)
		}
		edgeThreshold, err := strconv.Atoi(flag.Arg(3))
		if err != nil {
			fmt.Println("Invalid edgeThreshold value")
			os.Exit(1)
		}
		oilFilterSize, err := strconv.Atoi(flag.Arg(4))
		if err != nil {
			fmt.Println("Invalid oilFilterSize value")
			os.Exit(1)
		}
		oilLevels, err := strconv.Atoi(flag.Arg(5))
		if err != nil {
			fmt.Println("Invalid oilLevels value")
			os.Exit(1)
		}
		opts := effects.CTOpts{
			BlurKernelSize: blurStrength,
			EdgeThreshold:  edgeThreshold,
			OilFilterSize:  oilFilterSize,
			OilLevels:      oilLevels,
			DebugPath:      "",
		}
		outImg, err = effects.Cartoon(img, 0, opts)
		if err != nil {
			fmt.Println("Failed to apply effect:", err)
			os.Exit(1)
		}
	case "pixelate":
		blockSize, err := strconv.Atoi(flag.Arg(2))
		if err != nil {
			fmt.Println("Invalid blockSize value")
			os.Exit(1)
		}
		outImg, err = effects.Pixelate(img, 0, blockSize)
		if err != nil {
			fmt.Println("Failed to apply effect:", err)
			os.Exit(1)
		}
	}

	err = outImg.Save(outPath)
	if err != nil {
		fmt.Println("Failed to save modified image:", err)
		os.Exit(1)
	}

}
