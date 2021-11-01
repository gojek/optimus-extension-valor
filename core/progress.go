package core

import (
	"strings"

	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
)

const (
	maxNameLength = 21
	defaultWidth  = 64
)

// Progress defines how a progress bar should be executed
type Progress struct {
	progress *mpb.Progress
	bar      *mpb.Bar
}

// Increment increments the progress
func (p *Progress) Increment() {
	p.bar.Increment()
}

// Wait waits for all bar to complete
func (p *Progress) Wait() {
	p.bar.SetTotal(0, true)
	p.progress.Wait()
}

// NewProgress initializes a progress bar
func NewProgress(name string, total int) *Progress {
	name = standardize(name)

	progress := mpb.New()
	bar := progress.Add(int64(total),
		mpb.NewBarFiller(mpb.BarStyle().Lbound("╢").Filler("█").Tip("▌").Padding("░").Rbound("╟")),
		mpb.BarWidth(defaultWidth),
		mpb.PrependDecorators(
			decor.Name(name, decor.WCSyncSpace),
			decor.Spinner(nil, decor.WCSyncSpace),
		),
		mpb.AppendDecorators(
			decor.Percentage(decor.WCSyncSpace),
			decor.OnComplete(decor.Name("processed"), "done"),
		),
	)
	return &Progress{
		progress: progress,
		bar:      bar,
	}
}

func standardize(input string) string {
	if len(input) > maxNameLength {
		input = input[:maxNameLength] + "..."
	}
	if len(input) < maxNameLength+3 {
		input = input + strings.Repeat(" ", maxNameLength+3-len(input))
	}
	return input
}
