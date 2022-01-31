package explorer_test

import (
	"testing"

	"github.com/gojek/optimus-extension-valor/model"
	"github.com/gojek/optimus-extension-valor/registry/explorer"

	"github.com/stretchr/testify/suite"
)

type FactorySuite struct {
	suite.Suite
}

func (f *FactorySuite) TestRegister() {
	f.Run("should return error if fn is nil", func() {
		factory := explorer.NewFactory()
		_type := "file"
		var exploreFn model.ExplorePath = nil

		actualErr := factory.Register(_type, exploreFn)

		f.NotNil(actualErr)
	})

	f.Run("should return error fn is already registered", func() {
		factory := explorer.NewFactory()
		_type := "file"
		exploreFn := func(root string, filter func(string) bool) ([]string, error) {
			return nil, nil
		}
		factory.Register(_type, exploreFn)

		actualErr := factory.Register(_type, exploreFn)

		f.NotNil(actualErr)
	})

	f.Run("should return nil if no error is found", func() {
		factory := explorer.NewFactory()
		_type := "file"
		exploreFn := func(root string, filter func(string) bool) ([]string, error) {
			return nil, nil
		}

		actualErr := factory.Register(_type, exploreFn)

		f.Nil(actualErr)
	})
}

func (f *FactorySuite) TestGet() {
	f.Run("should return nil and error if _type is not found", func() {
		factory := explorer.NewFactory()
		_type := "file"

		actualFn, actualErr := factory.Get(_type)

		f.Nil(actualFn)
		f.NotNil(actualErr)
	})

	f.Run("should return fn and nil if _type is found", func() {
		factory := explorer.NewFactory()
		_type := "file"
		factory.Register(_type, func(root string, filter func(string) bool) ([]string, error) {
			return nil, nil
		})

		actualFn, actualErr := factory.Get(_type)

		f.NotNil(actualFn)
		f.Nil(actualErr)
	})
}

func TestWriterFactorySuite(t *testing.T) {
	suite.Run(t, &FactorySuite{})
}
