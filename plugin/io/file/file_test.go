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
	defaultDirName  = "./out"
	defaultFileName = "test.yaml"
	defaultContent  = "message"
)

type FileSuite struct {
	suite.Suite
}

func (f *FileSuite) SetupSuite() {
	if err := os.MkdirAll(defaultDirName, os.ModePerm); err != nil {
		panic(err)
	}
	filePath := path.Join(defaultDirName, defaultFileName)
	if err := ioutil.WriteFile(filePath, []byte(defaultContent), os.ModePerm); err != nil {
		panic(err)
	}
}

func (f *FileSuite) TestRead() {
	f.Run("should return error if getPath nil", func() {
		var getPath model.GetPath = nil
		var postProcess model.PostProcess = func(path string, content []byte) (*model.Data, error) {
			return &model.Data{
				Content: content,
				Path:    path,
			}, nil
		}
		reader := file.New(getPath, postProcess)

		actualData, actualErr := reader.Read()

		f.Nil(actualData)
		f.NotNil(actualErr)
	})

	f.Run("should return error if postProcess nil", func() {
		var getPath model.GetPath = func() string {
			return defaultDirName
		}
		var postProcess model.PostProcess = nil
		reader := file.New(getPath, postProcess)

		actualData, actualErr := reader.Read()

		f.Nil(actualData)
		f.NotNil(actualErr)
	})

	f.Run("should return error if error is found when post process", func() {
		var getPath model.GetPath = func() string {
			return defaultDirName
		}
		var postProcess model.PostProcess = func(path string, content []byte) (*model.Data, error) {
			return nil, errors.New("test error")
		}
		reader := file.New(getPath, postProcess)

		actualData, actualErr := reader.Read()

		f.Nil(actualData)
		f.NotNil(actualErr)
	})

	f.Run("should return value if no error is found", func() {
		var getPath model.GetPath = func() string {
			return path.Join(defaultDirName, defaultFileName)
		}
		var postProcess model.PostProcess = func(path string, content []byte) (*model.Data, error) {
			return &model.Data{
				Content: content,
				Path:    path,
			}, nil
		}
		reader := file.New(getPath, postProcess)

		actualData, actualErr := reader.Read()

		f.NotNil(actualData)
		f.Nil(actualErr)
	})
}

func (f *FileSuite) TestWrite() {
	f.Run("should return error if data is nil", func() {
		reader := file.New(nil, nil)
		var data *model.Data = nil

		actualErr := reader.Write(data)

		f.NotNil(actualErr)
	})

	f.Run("should return result of write", func() {
		defer func() { os.RemoveAll("./output") }()
		data := &model.Data{
			Type:    "file",
			Path:    path.Join(defaultDirName, defaultFileName),
			Content: []byte(defaultContent),
		}
		reader := file.New(nil, nil)

		actualErr := reader.Write(data)

		f.Nil(actualErr)
	})
}

func (f *FileSuite) TearDownSuite() {
	if err := os.RemoveAll(defaultDirName); err != nil {
		panic(err)
	}
}

func TestFileSuite(t *testing.T) {
	suite.Run(t, &FileSuite{})
}
