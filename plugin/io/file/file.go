package file

import (
	"errors"
	"io/ioutil"
	"os"
	"path"

	"github.com/gojek/optimus-extension-valor/model"
	"github.com/gojek/optimus-extension-valor/registry/io"
)

const _type = "file"

// File represents file operation
type File struct {
	getPath     model.GetPath
	postProcess model.PostProcess
}

// Read reads file from a path
func (f *File) Read() (*model.Data, error) {
	if err := f.validate(); err != nil {
		return nil, err
	}
	path := f.getPath()
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return f.postProcess(path, content)
}

func (f *File) validate() error {
	if f.getPath == nil {
		return errors.New("getPath is nil")
	}
	if f.postProcess == nil {
		return errors.New("postProcess is nil")
	}
	return nil
}

// Write writes data to destination
func (f *File) Write(data *model.Data) error {
	if data == nil {
		return errors.New("data is nil")
	}
	dirPath, _ := path.Split(data.Path)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return err
	}
	return ioutil.WriteFile(data.Path, data.Content, os.ModePerm)
}

// New initializes File based on path
func New(getPath model.GetPath, postProcess model.PostProcess) *File {
	return &File{
		getPath:     getPath,
		postProcess: postProcess,
	}
}

func init() {
	if err := io.Readers.Register(_type,
		func(getPath model.GetPath, postProcess model.PostProcess) model.Reader {
			return New(getPath, postProcess)
		},
	); err != nil {
		panic(err)
	}
	if err := io.Writers.Register(_type,
		func(treatment model.OutputTreatment) model.Writer {
			return New(nil, nil)
		},
	); err != nil {
		panic(err)
	}
}
