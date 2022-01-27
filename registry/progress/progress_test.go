package progress_test

import (
	"testing"

	"github.com/gojek/optimus-extension-valor/model"
	"github.com/gojek/optimus-extension-valor/registry/progress"

	"github.com/stretchr/testify/suite"
)

type FactorySuite struct {
	suite.Suite
}

func (f *FactorySuite) TestRegister() {
	f.Run("should return error if fn is nil", func() {
		factory := progress.NewFactory()
		_type := "progressive"
		var fn model.NewProgress = nil

		actualErr := factory.Register(_type, fn)

		f.NotNil(actualErr)
	})

	f.Run("should return error fn is already registered", func() {
		factory := progress.NewFactory()
		_type := "progressive"
		var fn model.NewProgress = func(name string, total int) model.Progress {
			return nil
		}
		factory.Register(_type, fn)

		actualErr := factory.Register(_type, fn)

		f.NotNil(actualErr)
	})

	f.Run("should return nil if no error is found", func() {
		factory := progress.NewFactory()
		_type := "progressive"
		var fn model.NewProgress = func(name string, total int) model.Progress {
			return nil
		}

		actualErr := factory.Register(_type, fn)

		f.Nil(actualErr)
	})
}

func (f *FactorySuite) TestGet() {
	f.Run("should return nil and error type is not found", func() {
		factory := progress.NewFactory()
		_type := "progressive"
		var fn model.NewProgress = func(name string, total int) model.Progress {
			return nil
		}
		factory.Register(_type, fn)

		actualFn, actualErr := factory.Get("file")

		f.Nil(actualFn)
		f.NotNil(actualErr)
	})

	f.Run("should return fn and nil type is found found", func() {
		factory := progress.NewFactory()
		_type := "progressive"
		var fn model.NewProgress = func(name string, total int) model.Progress {
			return nil
		}
		factory.Register(_type, fn)

		actualFn, actualErr := factory.Get(_type)

		f.NotNil(actualFn)
		f.Nil(actualErr)
	})
}

func TestReaderFactorySuite(t *testing.T) {
	suite.Run(t, &FactorySuite{})
}
