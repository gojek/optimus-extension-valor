package progress

import (
	"strings"

	"github.com/gojek/optimus-extension-valor/model"
	"github.com/gojek/optimus-extension-valor/registry/progress"

	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
)

const verboseType = "verbose"

const (
	maxNameLength = 21
	defaultWidth  = 64
)

// Verbose defines how a verbose progress should be executed
type Verbose struct {
	progress *mpb.Progress
	bar      *mpb.Bar
}

// Increment increments the progress
func (v *Verbose) Increment() {
	v.bar.Increment()
}

// Wait finishes the progress
func (v *Verbose) Wait() {
	v.bar.SetTotal(0, true)
	v.progress.Wait()
}

// NewVerbose initializes a verbose progress
func NewVerbose(name string, total int) *Verbose {
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
			decor.Elapsed(decor.ET_STYLE_MMSS, decor.WCSyncSpace),
		),
	)
	return &Verbose{
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

func init() {
	err := progress.Progresses.Register(verboseType, func(name string, total int) model.Progress {
		return NewVerbose(name, total)
	})
	if err != nil {
		panic(err)
	}
}
