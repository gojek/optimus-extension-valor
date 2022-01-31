package progress

import (
	"strings"

	"github.com/gojek/optimus-extension-valor/model"
	"github.com/gojek/optimus-extension-valor/registry/progress"

	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
)

const progressiveType = "progressive"

const (
	maxNameLength = 21
	defaultWidth  = 64
)

// Progressive defines how a progressive progress should be executed
type Progressive struct {
	progress *mpb.Progress
	bar      *mpb.Bar
}

// Increase inceases the progress by the specified num
func (v *Progressive) Increase(num int) {
	v.bar.IncrBy(num)
}

// Wait finishes the progress
func (v *Progressive) Wait() {
	v.bar.SetTotal(0, true)
	v.progress.Wait()
}

// NewProgressive initializes a progressive progress
func NewProgressive(name string, total int) *Progressive {
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
	return &Progressive{
		progress: progress,
		bar:      bar,
	}
}

func standardize(input string) string {
	extender := "..."
	if len(input) > maxNameLength {
		input = input[:maxNameLength] + extender
	}
	limit := maxNameLength + len(extender)
	if len(input) < limit {
		input = input + strings.Repeat(" ", limit-len(input))
	}
	return input
}

func init() {
	err := progress.Progresses.Register(progressiveType, func(name string, total int) model.Progress {
		return NewProgressive(name, total)
	})
	if err != nil {
		panic(err)
	}
}
