package io_test

import (
	"testing"

	"github.com/gojek/optimus-extension-valor/model"
	"github.com/gojek/optimus-extension-valor/registry/io"

	"github.com/stretchr/testify/suite"
)

type ReaderFactorySuite struct {
	suite.Suite
}

func (r *ReaderFactorySuite) TestRegister() {
	r.Run("should return error if fn is nil", func() {
		factory := io.NewReaderFactory()
		_type := "file"
		var fn io.ReaderFn = nil

		actualErr := factory.Register(_type, fn)

		r.NotNil(actualErr)
	})

	r.Run("should return error fn is already registered", func() {
		factory := io.NewReaderFactory()
		_type := "file"
		var fn io.ReaderFn = func(getPath model.GetPath, postProcess model.PostProcess) model.Reader {
			return nil
		}
		factory.Register(_type, fn)

		actualErr := factory.Register(_type, fn)

		r.NotNil(actualErr)
	})

	r.Run("should return nil if no error is found", func() {
		factory := io.NewReaderFactory()
		_type := "file"
		var fn io.ReaderFn = func(getPath model.GetPath, postProcess model.PostProcess) model.Reader {
			return nil
		}

		actualErr := factory.Register(_type, fn)

		r.Nil(actualErr)
	})
}

func (r *ReaderFactorySuite) TestGet() {
	r.Run("should return nil and error type is not found", func() {
		factory := io.NewReaderFactory()
		_type := "file"
		var fn io.ReaderFn = func(getPath model.GetPath, postProcess model.PostProcess) model.Reader {
			return nil
		}
		factory.Register(_type, fn)

		actualFn, actualErr := factory.Get("dir")

		r.Nil(actualFn)
		r.NotNil(actualErr)
	})

	r.Run("should return fn and nil type is found found", func() {
		factory := io.NewReaderFactory()
		_type := "file"
		var fn io.ReaderFn = func(getPath model.GetPath, postProcess model.PostProcess) model.Reader {
			return nil
		}
		factory.Register(_type, fn)

		actualFn, actualErr := factory.Get(_type)

		r.NotNil(actualFn)
		r.Nil(actualErr)
	})
}

func TestReaderFactorySuite(t *testing.T) {
	suite.Run(t, &ReaderFactorySuite{})
}
