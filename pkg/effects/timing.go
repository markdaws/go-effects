package effects

import (
	"fmt"
	"time"
)

//TODO: Move this in to its own project

// Timing allows the caller to capture timing information through a codebase
type Timing struct {
	Labels        map[string]int
	runningLabels map[string]time.Time
}

// NewTiming returns a new Timing instance
func NewTiming() *Timing {
	return &Timing{
		Labels:        make(map[string]int),
		runningLabels: make(map[string]time.Time),
	}
}

// Time starts a timer with the specified label e.g. 'load-image'
func (t *Timing) Time(label string) {
	t.runningLabels[label] = time.Now()
}

// TimeEnd stops the timer with the specified label and returns the number of ms
// The value is also accessible at .Labels['load-image'] for example.
func (t *Timing) TimeEnd(label string) int {
	start, ok := t.runningLabels[label]
	if !ok {
		return -1
	}

	elapsed := int(time.Since(start).Nanoseconds() / 1000000)
	t.Labels[label] = elapsed
	delete(t.runningLabels, label)
	return elapsed
}

// String returns a list of all the labels and the amount of time allocated by each
func (t *Timing) String() string {
	var out string
	out += "TIMINGS ----------\n"
	for k, v := range t.Labels {
		out += fmt.Sprintf("%s : %dms\n", k, v)
	}
	out += "------------------\n"
	return out
}
