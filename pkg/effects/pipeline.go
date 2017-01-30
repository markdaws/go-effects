package effects

// Pipeline allows multiple effects to be composed together easily
type Pipeline struct {
	effects []item
}

type item struct {
	effect   Effect
	callback func(*Image)
}

// Add adds an effect to the pipeline
func (p *Pipeline) Add(e Effect, callback func(*Image)) {
	p.effects = append(p.effects, item{effect: e, callback: callback})
}

// Run executes all of the effects in the order they were passed to the Add function
// on the input image and returns the results.
func (p *Pipeline) Run(img *Image, numRoutines int) (*Image, error) {
	currentImg := img
	for _, item := range p.effects {
		outImg, err := item.effect.Apply(currentImg, numRoutines)
		if err != nil {
			return nil, err
		}
		if item.callback != nil {
			item.callback(outImg)
		}
		currentImg = outImg
	}
	return currentImg, nil
}
