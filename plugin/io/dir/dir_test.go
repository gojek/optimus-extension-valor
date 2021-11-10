package dir_test

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/gojek/optimus-extension-valor/model"
	"github.com/gojek/optimus-extension-valor/plugin/io/dir"

	"github.com/stretchr/testify/suite"
)

const (
	defaulDirName   = "./out"
	defaultFileName = "test.yaml"
	defaultContent  = "message"
)

type DirSuite struct {
	suite.Suite
}

func (d *DirSuite) SetupSuite() {
	if err := os.MkdirAll(defaulDirName, os.ModePerm); err != nil {
		panic(err)
	}
	filePath := path.Join(defaulDirName, defaultFileName)
	if err := ioutil.WriteFile(filePath, []byte(defaultContent), os.ModePerm); err != nil {
		panic(err)
	}
}

func (d *DirSuite) TestWrite() {
	d.Run("should return error if data list is empty", func() {
		var dataList []*model.Data
		writer := dir.NewWriter()

		actualErr := writer.Write(dataList...)

		d.NotNil(actualErr)
	})

	d.Run("should return error if there's error during write", func() {
		dataList := []*model.Data{
			{
				Content: []byte(defaultContent),
				Path:    path.Join(defaulDirName, defaultFileName, "test"),
			},
			{
				Content: []byte(defaultContent),
				Path:    defaulDirName,
			},
		}
		writer := dir.NewWriter()

		actualErr := writer.Write(dataList...)

		d.NotNil(actualErr)
	})

	d.Run("should return error if there's error during write", func() {
		dataList := []*model.Data{
			{
				Content: []byte(defaultContent),
				Path:    path.Join(defaulDirName, defaultFileName),
			},
		}
		writer := dir.NewWriter()

		actualErr := writer.Write(dataList...)

		d.Nil(actualErr)
	})
}

func (d *DirSuite) TestReadAll() {
	d.Run("should return error if getPath nil", func() {
		var getPath model.GetPath = nil
		var filterPath model.FilterPath = func(s string) bool {
			return true
		}
		var postProcess model.PostProcess = func(path string, content []byte) (*model.Data, model.Error) {
			return &model.Data{
				Content: content,
				Path:    path,
			}, nil
		}
		reader := dir.NewReader(getPath, filterPath, postProcess)

		actualData, actualErr := reader.ReadAll()

		d.Nil(actualData)
		d.NotNil(actualErr)
	})

	d.Run("should return error if postProcess nil", func() {
		var getPath model.GetPath = func() string {
			return defaulDirName
		}
		var filterPath model.FilterPath = func(s string) bool {
			return true
		}
		var postProcess model.PostProcess = nil
		reader := dir.NewReader(getPath, filterPath, postProcess)

		actualData, actualErr := reader.ReadAll()

		d.Nil(actualData)
		d.NotNil(actualErr)
	})

	d.Run("should return error if no file paths found nil", func() {
		var getPath model.GetPath = func() string {
			return defaulDirName
		}
		var filterPath model.FilterPath = func(s string) bool {
			return false
		}
		var postProcess model.PostProcess = func(path string, content []byte) (*model.Data, model.Error) {
			return &model.Data{
				Content: content,
				Path:    path,
			}, nil
		}
		reader := dir.NewReader(getPath, filterPath, postProcess)

		actualData, actualErr := reader.ReadAll()

		d.Nil(actualData)
		d.NotNil(actualErr)
	})

	d.Run("should return error if error is found when post process", func() {
		var getPath model.GetPath = func() string {
			return defaulDirName
		}
		var filterPath model.FilterPath = func(s string) bool {
			return true
		}
		var postProcess model.PostProcess = func(path string, content []byte) (*model.Data, model.Error) {
			return nil, model.BuildError("test", errors.New("test error"))
		}
		reader := dir.NewReader(getPath, filterPath, postProcess)

		actualData, actualErr := reader.ReadAll()

		d.Nil(actualData)
		d.NotNil(actualErr)
	})

	d.Run("should return value if no error is found", func() {
		var getPath model.GetPath = func() string {
			return defaulDirName
		}
		var filterPath model.FilterPath = func(s string) bool {
			return true
		}
		var postProcess model.PostProcess = func(path string, content []byte) (*model.Data, model.Error) {
			return &model.Data{
				Content: content,
				Path:    path,
			}, nil
		}
		reader := dir.NewReader(getPath, filterPath, postProcess)

		actualData, actualErr := reader.ReadAll()

		d.NotNil(actualData)
		d.Nil(actualErr)
	})
}

func (d *DirSuite) TestReadOne() {
	d.Run("should return error if getPath nil", func() {
		var getPath model.GetPath = nil
		var filterPath model.FilterPath = func(s string) bool {
			return true
		}
		var postProcess model.PostProcess = func(path string, content []byte) (*model.Data, model.Error) {
			return &model.Data{
				Content: content,
				Path:    path,
			}, nil
		}
		reader := dir.NewReader(getPath, filterPath, postProcess)

		actualData, actualErr := reader.ReadOne()

		d.Nil(actualData)
		d.NotNil(actualErr)
	})

	d.Run("should return error if postProcess nil", func() {
		var getPath model.GetPath = func() string {
			return defaulDirName
		}
		var filterPath model.FilterPath = func(s string) bool {
			return true
		}
		var postProcess model.PostProcess = nil
		reader := dir.NewReader(getPath, filterPath, postProcess)

		actualData, actualErr := reader.ReadOne()

		d.Nil(actualData)
		d.NotNil(actualErr)
	})

	d.Run("should return error if no file paths found nil", func() {
		var getPath model.GetPath = func() string {
			return defaulDirName
		}
		var filterPath model.FilterPath = func(s string) bool {
			return false
		}
		var postProcess model.PostProcess = func(path string, content []byte) (*model.Data, model.Error) {
			return &model.Data{
				Content: content,
				Path:    path,
			}, nil
		}
		reader := dir.NewReader(getPath, filterPath, postProcess)

		actualData, actualErr := reader.ReadOne()

		d.Nil(actualData)
		d.NotNil(actualErr)
	})

	d.Run("should return error if error is found when post process", func() {
		var getPath model.GetPath = func() string {
			return defaulDirName
		}
		var filterPath model.FilterPath = func(s string) bool {
			return true
		}
		var postProcess model.PostProcess = func(path string, content []byte) (*model.Data, model.Error) {
			return nil, model.BuildError("test", errors.New("test error"))
		}
		reader := dir.NewReader(getPath, filterPath, postProcess)

		actualData, actualErr := reader.ReadOne()

		d.Nil(actualData)
		d.NotNil(actualErr)
	})

	d.Run("should return value if no error is found", func() {
		var getPath model.GetPath = func() string {
			return defaulDirName
		}
		var filterPath model.FilterPath = func(s string) bool {
			return true
		}
		var postProcess model.PostProcess = func(path string, content []byte) (*model.Data, model.Error) {
			return &model.Data{
				Content: content,
				Path:    path,
			}, nil
		}
		reader := dir.NewReader(getPath, filterPath, postProcess)

		actualData, actualErr := reader.ReadOne()

		d.NotNil(actualData)
		d.Nil(actualErr)
	})
}

func (d *DirSuite) TearDownSuite() {
	if err := os.RemoveAll(defaulDirName); err != nil {
		panic(err)
	}
}

func TestDirSuite(t *testing.T) {
	suite.Run(t, &DirSuite{})
}
