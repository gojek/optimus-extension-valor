package file

import (
	"errors"
	"io/ioutil"

	"github.com/gojek/optimus-extension-valor/model"
	"github.com/gojek/optimus-extension-valor/registry/io"
)

const _type = "file"

// File represents file operation
type File struct {
	getPath     model.GetPath
	postProcess model.PostProcess
}

// ReadOne reads one file from a path
func (f *File) ReadOne() (*model.Data, model.Error) {
	const defaultErrKey = "ReadOne"
	if err := f.validate(); err != nil {
		return nil, err
	}
	path := f.getPath()
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, model.BuildError(defaultErrKey, err)
	}
	return f.postProcess(path, content)
}

// ReadAll reads one file from a path but is returned as a slice
func (f *File) ReadAll() ([]*model.Data, model.Error) {
	const defaultErrKey = "ReadAll"
	if err := f.validate(); err != nil {
		return nil, err
	}
	data, err := f.ReadOne()
	if err != nil {
		return nil, err
	}
	return []*model.Data{data}, nil
}

func (f *File) validate() model.Error {
	const defaultErrKey = "validate"
	if f.getPath == nil {
		return model.BuildError(defaultErrKey, errors.New("getPath is nil"))
	}
	if f.postProcess == nil {
		return model.BuildError(defaultErrKey, errors.New("postProcess is nil"))
	}
	return nil
}

// New initializes File based on path
func New(getPath model.GetPath, postProcess model.PostProcess) *File {
	return &File{
		getPath:     getPath,
		postProcess: postProcess,
	}
}

func init() {
	err := io.Readers.Register(_type, func(
		getPath model.GetPath,
		filterPath model.FilterPath,
		postProcess model.PostProcess,
	) model.Reader {
		return New(getPath, postProcess)
	})
	if err != nil {
		panic(err)
	}
}
