package progress

import (
	"fmt"
	"strings"
	"time"

	"github.com/gojek/optimus-extension-valor/model"
	"github.com/gojek/optimus-extension-valor/registry/progress"
)

const iterativeType = "iterative"

// Iterative defines how a iterative progress bar should be executed
type Iterative struct {
	name  string
	total int

	currentCounter     int
	previousPercentage int

	finished  bool
	startTime time.Time
}

// Increment increments the progress
func (s *Iterative) Increment() {
	if s.currentCounter >= s.total || s.finished {
		return
	}
	s.currentCounter++

	currentPercentage := 100 * s.currentCounter / s.total
	if currentPercentage > s.previousPercentage {
		s.previousPercentage = currentPercentage

		currentLength := len(fmt.Sprintf("%d", s.currentCounter))
		totalLength := len(fmt.Sprintf("%d", s.total))

		current := fmt.Sprintf("%s%d", strings.Repeat(" ", totalLength-currentLength), s.currentCounter)
		fmt.Printf("%s: %s/%d [%3d%%]\n",
			s.name, current, s.total, currentPercentage,
		)
	}
}

// Wait finishes the progress
func (s *Iterative) Wait() {
	if !s.finished {
		fmt.Printf("total elapsed: %v\n", time.Now().Sub(s.startTime))
	}
	s.finished = true
}

// NewIterative initializes an iterative progress
func NewIterative(name string, total int) *Iterative {
	return &Iterative{
		name:      name,
		total:     total,
		startTime: time.Now(),
	}
}

func init() {
	err := progress.Progresses.Register(iterativeType, func(name string, total int) model.Progress {
		return NewIterative(name, total)
	})
	if err != nil {
		panic(err)
	}
}
