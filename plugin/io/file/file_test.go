package file_test

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/gojek/optimus-extension-valor/model"
	"github.com/gojek/optimus-extension-valor/plugin/io/file"

	"github.com/stretchr/testify/suite"
)

const (
	defaulDirName   = "./out"
	defaultFileName = "test.yaml"
	defaultContent  = "message"
)

type FileSuite struct {
	suite.Suite
}

func (f *FileSuite) SetupSuite() {
	if err := os.MkdirAll(defaulDirName, os.ModePerm); err != nil {
		panic(err)
	}
	filePath := path.Join(defaulDirName, defaultFileName)
	if err := ioutil.WriteFile(filePath, []byte(defaultContent), os.ModePerm); err != nil {
		panic(err)
	}
}

func (f *FileSuite) TestReadAll() {
	f.Run("should return error if getPath nil", func() {
		var getPath model.GetPath = nil
		var postProcess model.PostProcess = func(path string, content []byte) (*model.Data, model.Error) {
			return &model.Data{
				Content: content,
				Path:    path,
			}, nil
		}
		reader := file.New(getPath, postProcess)

		actualData, actualErr := reader.ReadAll()

		f.Nil(actualData)
		f.NotNil(actualErr)
	})

	f.Run("should return error if postProcess nil", func() {
		var getPath model.GetPath = func() string {
			return defaulDirName
		}
		var postProcess model.PostProcess = nil
		reader := file.New(getPath, postProcess)

		actualData, actualErr := reader.ReadAll()

		f.Nil(actualData)
		f.NotNil(actualErr)
	})

	f.Run("should return error if error is found when post process", func() {
		var getPath model.GetPath = func() string {
			return defaulDirName
		}
		var postProcess model.PostProcess = func(path string, content []byte) (*model.Data, model.Error) {
			return nil, model.BuildError("test", errors.New("test error"))
		}
		reader := file.New(getPath, postProcess)

		actualData, actualErr := reader.ReadAll()

		f.Nil(actualData)
		f.NotNil(actualErr)
	})

	f.Run("should return value if no error is found", func() {
		var getPath model.GetPath = func() string {
			return path.Join(defaulDirName, defaultFileName)
		}
		var postProcess model.PostProcess = func(path string, content []byte) (*model.Data, model.Error) {
			return &model.Data{
				Content: content,
				Path:    path,
			}, nil
		}
		reader := file.New(getPath, postProcess)

		actualData, actualErr := reader.ReadAll()

		f.NotNil(actualData)
		f.Nil(actualErr)
	})
}

func (f *FileSuite) TestReadOne() {
	f.Run("should return error if getPath nil", func() {
		var getPath model.GetPath = nil
		var postProcess model.PostProcess = func(path string, content []byte) (*model.Data, model.Error) {
			return &model.Data{
				Content: content,
				Path:    path,
			}, nil
		}
		reader := file.New(getPath, postProcess)

		actualData, actualErr := reader.ReadOne()

		f.Nil(actualData)
		f.NotNil(actualErr)
	})

	f.Run("should return error if postProcess nil", func() {
		var getPath model.GetPath = func() string {
			return defaulDirName
		}
		var postProcess model.PostProcess = nil
		reader := file.New(getPath, postProcess)

		actualData, actualErr := reader.ReadOne()

		f.Nil(actualData)
		f.NotNil(actualErr)
	})

	f.Run("should return error if error is found when post process", func() {
		var getPath model.GetPath = func() string {
			return defaulDirName
		}
		var postProcess model.PostProcess = func(path string, content []byte) (*model.Data, model.Error) {
			return nil, model.BuildError("test", errors.New("test error"))
		}
		reader := file.New(getPath, postProcess)

		actualData, actualErr := reader.ReadOne()

		f.Nil(actualData)
		f.NotNil(actualErr)
	})

	f.Run("should return value if no error is found", func() {
		var getPath model.GetPath = func() string {
			return path.Join(defaulDirName, defaultFileName)
		}
		var postProcess model.PostProcess = func(path string, content []byte) (*model.Data, model.Error) {
			return &model.Data{
				Content: content,
				Path:    path,
			}, nil
		}
		reader := file.New(getPath, postProcess)

		actualData, actualErr := reader.ReadOne()

		f.NotNil(actualData)
		f.Nil(actualErr)
	})
}

func (f *FileSuite) TearDownSuite() {
	if err := os.RemoveAll(defaulDirName); err != nil {
		panic(err)
	}
}

func TestFileSuite(t *testing.T) {
	suite.Run(t, &FileSuite{})
}
