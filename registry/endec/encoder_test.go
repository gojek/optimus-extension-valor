package endec_test

import (
	"testing"

	"github.com/gojek/optimus-extension-valor/model"
	"github.com/gojek/optimus-extension-valor/registry/endec"

	"github.com/stretchr/testify/suite"
)

type EncodeFactorySuite struct {
	suite.Suite
}

func (e *EncodeFactorySuite) TestRegister() {
	e.Run("should return error if fn is nil", func() {
		factory := endec.NewEncodeFactory()
		format := "yaml"
		var fn model.Encode = nil

		actualErr := factory.Register(format, fn)

		e.NotNil(actualErr)
	})

	e.Run("should return error fn is already registered", func() {
		factory := endec.NewEncodeFactory()
		format := "yaml"
		var fn model.Encode = func(i interface{}) ([]byte, model.Error) {
			return nil, nil
		}
		factory.Register(format, fn)

		actualErr := factory.Register(format, fn)

		e.NotNil(actualErr)
	})

	e.Run("should return nil if no error is found", func() {
		factory := endec.NewEncodeFactory()
		format := "yaml"
		var fn model.Encode = func(i interface{}) ([]byte, model.Error) {
			return nil, nil
		}

		actualErr := factory.Register(format, fn)

		e.Nil(actualErr)
	})
}

func (e *EncodeFactorySuite) TestGet() {
	e.Run("should return nil and error type is not found", func() {
		factory := endec.NewEncodeFactory()
		yamlFormat := "yaml"
		jsonFormat := "json"
		var fn model.Encode = func(i interface{}) ([]byte, model.Error) {
			return nil, nil
		}
		factory.Register(yamlFormat, fn)

		actualFn, actualErr := factory.Get(jsonFormat)

		e.Nil(actualFn)
		e.NotNil(actualErr)
	})

	e.Run("should return fn and nil type is found found", func() {
		factory := endec.NewEncodeFactory()
		yamlFormat := "yaml"
		var fn model.Encode = func(i interface{}) ([]byte, model.Error) {
			return nil, nil
		}
		factory.Register(yamlFormat, fn)

		actualFn, actualErr := factory.Get(yamlFormat)

		e.NotNil(actualFn)
		e.Nil(actualErr)
	})
}

func TestEncodeFactorySuite(t *testing.T) {
	suite.Run(t, &EncodeFactorySuite{})
}
