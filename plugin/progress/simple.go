package progress

import (
	"fmt"
	"strings"
	"time"

	"github.com/gojek/optimus-extension-valor/model"
	"github.com/gojek/optimus-extension-valor/registry/progress"
)

const simpleType = "simple"

// Simple defines how a simple progress bar should be executed
type Simple struct {
	name  string
	total int

	currentCounter     int
	previousPercentage int

	finished  bool
	startTime time.Time
}

// Increment increments the progress
func (s *Simple) Increment() {
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
func (s *Simple) Wait() {
	if !s.finished {
		fmt.Printf("total elapsed: %v\n", time.Now().Sub(s.startTime))
	}
	s.finished = true
}

// NewSimple initializes a simple progress
func NewSimple(name string, total int) *Simple {
	return &Simple{
		name:      name,
		total:     total,
		startTime: time.Now(),
	}
}

func init() {
	err := progress.Progresses.Register(simpleType, func(name string, total int) model.Progress {
		return NewSimple(name, total)
	})
	if err != nil {
		panic(err)
	}
}
