package endec_test

import (
	"testing"

	"github.com/gojek/optimus-extension-valor/model"
	"github.com/gojek/optimus-extension-valor/registry/endec"

	"github.com/stretchr/testify/suite"
)

type DecodeFactorySuite struct {
	suite.Suite
}

func (d *DecodeFactorySuite) TestRegister() {
	d.Run("should return error if fn is nil", func() {
		factory := endec.NewDecodeFactory()
		format := "yaml"
		var fn model.Decode = nil

		actualErr := factory.Register(format, fn)

		d.NotNil(actualErr)
	})

	d.Run("should return error fn is already registered", func() {
		factory := endec.NewDecodeFactory()
		format := "yaml"
		var fn model.Decode = func(b []byte, i interface{}) error {
			return nil
		}
		factory.Register(format, fn)

		actualErr := factory.Register(format, fn)

		d.NotNil(actualErr)
	})

	d.Run("should return nil if no error is found", func() {
		factory := endec.NewDecodeFactory()
		format := "yaml"
		var fn model.Decode = func(b []byte, i interface{}) error {
			return nil
		}

		actualErr := factory.Register(format, fn)

		d.Nil(actualErr)
	})
}

func (d *DecodeFactorySuite) TestGet() {
	d.Run("should return nil and error type is not found", func() {
		factory := endec.NewDecodeFactory()
		yamlFormat := "yaml"
		jsonFormat := "json"
		var fn model.Decode = func(b []byte, i interface{}) error {
			return nil
		}
		factory.Register(yamlFormat, fn)

		actualFn, actualErr := factory.Get(jsonFormat)

		d.Nil(actualFn)
		d.NotNil(actualErr)
	})

	d.Run("should return fn and nil type is found found", func() {
		factory := endec.NewDecodeFactory()
		yamlFormat := "yaml"
		var fn model.Decode = func(b []byte, i interface{}) error {
			return nil
		}
		factory.Register(yamlFormat, fn)

		actualFn, actualErr := factory.Get(yamlFormat)

		d.NotNil(actualFn)
		d.Nil(actualErr)
	})
}

func TestDecoderFactorySuite(t *testing.T) {
	suite.Run(t, &DecodeFactorySuite{})
}
