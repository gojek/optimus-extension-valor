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

func (w *WriterFactorySuite) TestRegister() {
	w.Run("should return error if fn is nil", func() {
		factory := io.NewWriterFactory()
		_type := "file"
		var writerFn io.WriterFn = nil

		actualErr := factory.Register(_type, writerFn)

		w.NotNil(actualErr)
	})

	w.Run("should return error fn is already registered", func() {
		factory := io.NewWriterFactory()
		_type := "file"
		var writerFn io.WriterFn = func(model.OutputTreatment) model.Writer {
			return &mocks.Writer{}
		}
		factory.Register(_type, writerFn)

		actualErr := factory.Register(_type, writerFn)

		w.NotNil(actualErr)
	})

	w.Run("should return nil if no error is found", func() {
		factory := io.NewWriterFactory()
		_type := "file"
		var writerFn io.WriterFn = func(model.OutputTreatment) model.Writer {
			return &mocks.Writer{}
		}

		actualErr := factory.Register(_type, writerFn)

		w.Nil(actualErr)
	})
}

func (w *WriterFactorySuite) TestGet() {
	w.Run("should return nil and error type is not found", func() {
		factory := io.NewWriterFactory()
		_type := "file"
		var writerFn io.WriterFn = func(model.OutputTreatment) model.Writer {
			return &mocks.Writer{}
		}
		factory.Register(_type, writerFn)
		factory.Register(_type, writerFn)

		actualWriter, actualErr := factory.Get("dir")

		w.Nil(actualWriter)
		w.NotNil(actualErr)
	})

	w.Run("should return fn and nil type is found found", func() {
		factory := io.NewWriterFactory()
		_type := "file"
		var writerFn io.WriterFn = func(model.OutputTreatment) model.Writer {
			return &mocks.Writer{}
		}
		factory.Register(_type, writerFn)
		factory.Register(_type, writerFn)

		actualWriter, actualErr := factory.Get(_type)

		w.NotNil(actualWriter)
		w.Nil(actualErr)
	})
}

func TestWriterFactorySuite(t *testing.T) {
	suite.Run(t, &WriterFactorySuite{})
}
