package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/markdaws/go-effects/pkg/effects"
)

func main() {
	effect := flag.String("effect", "", "The name of the effect to apply. Values are 'oil'")
	flag.Parse()

	switch *effect {
	case "oil":
		if len(flag.Args()) != 4 {
			fmt.Println("The oil effect requires 4 args, filterSize, levels, input file, output file\n")
			fmt.Println("Sample usage: goeffect -effect=oil 5 30 mypic.jpg mypic-oil.jpg\n")
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

	img, err := effects.LoadImage(flag.Arg(2))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	switch *effect {
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

		oilImg, err := effects.OilPainting(img, 0, filterSize, levels)
		if err != nil {
			fmt.Println("Failed to apply effect:", err)
			os.Exit(1)
		}

		err = oilImg.SaveAsJPG(flag.Arg(3), 90)
		if err != nil {
			fmt.Println("Failed to save modified image:", err)
			os.Exit(1)
		}
	}
}
