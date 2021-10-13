package file

import (
	"errors"
	"io/ioutil"

	"github.com/gojek/optimus-extension-valor/model"
	"github.com/gojek/optimus-extension-valor/plugin"
	"github.com/gojek/optimus-extension-valor/registry/io"
)

const _type = "file"

// File represents file operation
type File struct {
	read bool

	path     string
	metadata map[string]string
}

func (f *File) Read() (*model.Data, error) {
	if !f.Next() {
		return nil, errors.New("no available Next item")
	}
	content, err := ioutil.ReadFile(f.path)
	if err != nil {
		return nil, err
	}
	f.read = true
	return &model.Data{
		Path:     f.path,
		Content:  content,
		Metadata: f.metadata,
	}, nil
}

// Next checks whether there is an item to Read or not
func (f *File) Next() bool {
	return !f.read
}

// New initializes File based on path
func New(path string, metadata map[string]string) *File {
	return &File{
		path:     path,
		metadata: metadata,
	}
}

func init() {
	io.Readers.Register("file", func(path string, metadata map[string]string) plugin.Reader {
		return New(path, metadata)
	})
}
