package io_test

import (
	"testing"

	"github.com/gojek/optimus-extension-valor/mocks"
	"github.com/gojek/optimus-extension-valor/model"
	"github.com/gojek/optimus-extension-valor/registry/io"

	"github.com/stretchr/testify/suite"
)

type WriterFactorySuite struct {
	suite.Suite
}

func (r *WriterFactorySuite) TestRegister() {
	r.Run("should return error if fn is nil", func() {
		factory := io.NewWriterFactory()
		_type := "file"
		var writer model.Writer = nil

		actualErr := factory.Register(_type, writer)

		r.NotNil(actualErr)
	})

	r.Run("should return error fn is already registered", func() {
		factory := io.NewWriterFactory()
		_type := "file"
		var writer model.Writer = &mocks.Writer{}
		factory.Register(_type, writer)

		actualErr := factory.Register(_type, writer)

		r.NotNil(actualErr)
	})

	r.Run("should return nil if no error is found", func() {
		factory := io.NewWriterFactory()
		_type := "file"
		var writer model.Writer = &mocks.Writer{}

		actualErr := factory.Register(_type, writer)

		r.Nil(actualErr)
	})
}

func (r *WriterFactorySuite) TestGet() {
	r.Run("should return nil and error type is not found", func() {
		factory := io.NewWriterFactory()
		_type := "file"
		var writer model.Writer = &mocks.Writer{}
		factory.Register(_type, writer)

		actualWriter, actualErr := factory.Get("dir")

		r.Nil(actualWriter)
		r.NotNil(actualErr)
	})

	r.Run("should return fn and nil type is found found", func() {
		factory := io.NewWriterFactory()
		_type := "file"
		var writer model.Writer = &mocks.Writer{}
		factory.Register(_type, writer)

		actualWriter, actualErr := factory.Get(_type)

		r.NotNil(actualWriter)
		r.Nil(actualErr)
	})
}

func TestWriterFactorySuite(t *testing.T) {
	suite.Run(t, &WriterFactorySuite{})
}
