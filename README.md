# go-effects
Parallelized image manipulation effects, written in Go. All sample images used in this project are from: https://unsplash.com

## Overview
This library provides basic image effects, running on multiple goroutines to parallelize the processing. You can include the library directly in your code via gihutb.com/markdaws/go-effect/pkg/effects or there is a command line app that you can use to process photos in cmd/goeffects, to run:

```bash
go install ./cmd/goeffects 
goeffects --help
```

## Usage
Take a look at pkg/effects/effects_test.go for examples of how to use this library

##Package
github.com/markdaws/go-effects/pkg/effects

##Docs
See [godoc](https://godoc.org/github.com/markdaws/go-effects/pkg/effects)


##Oil Painting
This effect takes an input image and renders it styled as an oil painting. The boldness of the stroke and the range of the palette can be modified.

###Original Image
![](examples/mountain.jpg)
###Modified Image (filterSize:5, levels:30)
![](examples/mountain-oil-15-30.jpg)


##Cartoon
This effect renders an image as if it were drawn as a cartoon. The effect is achieved by rendering the image as an oil paiting and running edge detection on the input image, drawing the edges on top of the oil painting rendering.

###Original Image
![](examples/turtle.jpg)
###Modified Image (cartoon)
![](examples/turtle-cartoon.png)


##Pixelate
Renders the input image as a pixelated image, based on the blockSize specified by the caller.

###Original Image
![](examples/mountain.jpg)
###Modified Image (pixelate, blockSize: 10)
![](examples/mountain-pixelate-10.png)


##Pencil
Given an input image returns an image rendered as if it was drawn by pencil. This is just an inverted sobel image, you can specify a blur factor to reduce noise

###Original Image
![](examples/houses.jpg)
###Modified Image (Pencil, blurFactor: 5)
![](examples/houses-pencil.jpg)


##Sobel
Given an input image returns an image containing edge gradients values, based on the Sobel operator.  By default the pixel r,g,b values all contain the gradient intensity, but if you supply a threshold value to the function, then if the gradient intensity is >= threshold the pixel value will be 255 and if it is less than it will be 0.  This way you can set some threshold and use this for edge detection.

###Original Image
![](examples/turtle.jpg)
###Modified Image (Sobel)
![](examples/turtle-sobel.png)


##Gaussian
Applies a Gaussian blur to the input image. You can specify the kernelSize, larger values equate to more blurring and sigma, larger values give more weighting to pixels further from the target pixel.  Same values would be 11, 1 for example. The kernelSize must be an odd number.

###Original Image
![](examples/face.jpg)
###Modified Image (Gaussian)
![](examples/face-gaussian.png)


##Grayscale
Given an input image returns a grayscale version. Three algorithms are available, lightness (average of the max and min rgb value of a pixel), average (the average of the r,g,b values), luminosity (a weighted average of the rgb values based on how humans perceive color).

The luminence algorithm generally gives the best results.
###Original Image
![](examples/cabin.jpg)
###Modified Image (luminosity)
![](examples/cabin-gray-luminosity.png)
