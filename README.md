# go-effects
Parallelized image manipulation effects, written in Go.

## Overview
This library provides basic image effects, running on multiple goroutines to parallelize the processing.

## Usage
Take a look at pkg/effects/effects_test.go for examples of how to use this library

##Package
github.com/markdaws/go-effects/pkg/effects

##Docs
See [godoc](https://godoc.org/github.com/markdaws/go-effects/pkg/effects)

##Effects
###Oil Painting
This effect takes an input image and renders it styled as an oil painting. The boldness of the stroke and the range of the palette can be modified.

####Original Image
![](examples/mountain.jpg)
####Modified Image (filterSize:5, levels:30)
![](examples/mountain-oil-15-30.jpg)

###Grayscale
Given an input image returns a grayscale version. Three algorithms are available, lightness (average of the max and min rgb value of a pixel), average (the average of the r,g,b values), luminosity (a weighted average of the rgb values based on how humans perceive color).

The luminence algorithm generally gives the best results.
####Original Image
![](examples/cabin.jpg)
####Modified Image (luminosity)
![](examples/cabin-gray-luminosity.png)
