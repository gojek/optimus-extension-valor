package std_test

import (
	"testing"

	"github.com/gojek/optimus-extension-valor/model"
	"github.com/gojek/optimus-extension-valor/plugin/io/std"

	"github.com/stretchr/testify/suite"
)

type StdSuite struct {
	suite.Suite
}

func (f *StdSuite) TestWrite() {
	f.Run("should return error if data is nil", func() {
		reader := std.New(model.TreatmentInfo)
		var data *model.Data = nil

		actualErr := reader.Write(data)

		f.NotNil(actualErr)
	})

	f.Run("should return result of write", func() {
		data := &model.Data{
			Type:    "std",
			Path:    "./output/test.file",
			Content: []byte("test content"),
		}
		reader := std.New(model.TreatmentInfo)

		actualErr := reader.Write(data)

		f.Nil(actualErr)
	})
}

func TestFileSuite(t *testing.T) {
	suite.Run(t, &StdSuite{})
}
