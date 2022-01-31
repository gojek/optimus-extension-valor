package formatter_test

import (
	"testing"

	"github.com/gojek/optimus-extension-valor/model"
	"github.com/gojek/optimus-extension-valor/registry/formatter"

	"github.com/stretchr/testify/suite"
)

type FactorySuite struct {
	suite.Suite
}

func (f *FactorySuite) TestRegister() {
	f.Run("should return error if fn is nil", func() {
		factory := formatter.NewFactory()
		src := "json"
		dest := "yaml"
		var fn model.Format = nil

		actualErr := factory.Register(src, dest, fn)

		f.NotNil(actualErr)
	})

	f.Run("should return error fn is already registered", func() {
		factory := formatter.NewFactory()
		src := "json"
		dest := "yaml"
		var fn model.Format = func(b []byte) ([]byte, error) {
			return nil, nil
		}
		factory.Register(src, dest, fn)

		actualErr := factory.Register(src, dest, fn)

		f.NotNil(actualErr)
	})

	f.Run("should return nil if no error is found", func() {
		factory := formatter.NewFactory()
		src := "json"
		dest := "yaml"
		var fn model.Format = func(b []byte) ([]byte, error) {
			return nil, nil
		}

		actualErr := factory.Register(src, dest, fn)

		f.Nil(actualErr)
	})
}

func (f *FactorySuite) TestGet() {
	f.Run("should return nil and error type is not found", func() {
		factory := formatter.NewFactory()
		src := "json"
		dest := "yaml"
		var fn model.Format = func(b []byte) ([]byte, error) {
			return nil, nil
		}
		factory.Register(src, dest, fn)

		actualFn, actualErr := factory.Get(src, src)

		f.Nil(actualFn)
		f.NotNil(actualErr)
	})

	f.Run("should return fn and nil type is found found", func() {
		factory := formatter.NewFactory()
		src := "json"
		dest := "yaml"
		var fn model.Format = func(b []byte) ([]byte, error) {
			return nil, nil
		}
		factory.Register(src, dest, fn)

		actualFn, actualErr := factory.Get(src, dest)

		f.NotNil(actualFn)
		f.Nil(actualErr)
	})
}

func TestFactorySuite(t *testing.T) {
	suite.Run(t, &FactorySuite{})
}
