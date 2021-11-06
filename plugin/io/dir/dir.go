package dir

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/gojek/optimus-extension-valor/model"
	"github.com/gojek/optimus-extension-valor/registry/io"
)

const _type = "dir"

// ContentWrapper is a wrapper for loading content via channel
type ContentWrapper struct {
	Content []byte
	Error   model.Error
}

// Dir represents directory operation
type Dir struct {
	getPath     model.GetPath
	filterPath  model.FilterPath
	postProcess model.PostProcess
}

func (d *Dir) Write(dataList ...*model.Data) model.Error {
	const defaultErrKey = "Write"
	if len(dataList) == 0 {
		return model.BuildError(defaultErrKey, errors.New("data is empty"))
	}
	writeChans := make([]chan error, len(dataList))
	for i, data := range dataList {
		ch := make(chan error)
		go func(c chan error, dt *model.Data) {
			dirPath, _ := path.Split(dt.Path)
			if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
				c <- err
				return
			}
			if err := ioutil.WriteFile(dirPath, dt.Content, os.ModePerm); err != nil {
				c <- err
				return
			}
			c <- nil
		}(ch, data)
		writeChans[i] = ch
	}
	output := make(model.Error)
	for i, ch := range writeChans {
		err := <-ch
		if err != nil {
			key := fmt.Sprintf("%s [%s]", defaultErrKey, dataList[i].Path)
			output[key] = err
		}
	}
	if len(output) == 0 {
		output = nil
	}
	return output
}

// ReadAll reads all files in a directory
func (d *Dir) ReadAll() ([]*model.Data, model.Error) {
	const defaultErrKey = "ReadAll"
	if err := d.validate(); err != nil {
		return nil, model.BuildError(defaultErrKey, err)
	}
	dirPath := d.getPath()
	filePaths, err := d.getFilePaths(dirPath, d.filterPath)
	if err != nil {
		return nil, model.BuildError(defaultErrKey, err)
	}
	readChans := d.dispatchRead(filePaths)
	outputData := make([]*model.Data, len(readChans))
	outputError := make(model.Error)
	for i, ch := range readChans {
		path := filePaths[i]
		key := fmt.Sprintf("%s [%s]", defaultErrKey, path)
		readWrapper := <-ch
		if readWrapper.Error != nil {
			outputError[key] = readWrapper.Error
			continue
		}
		data, err := d.postProcess(path, readWrapper.Content)
		if err != nil {
			outputError[key] = err
		}
		outputData[i] = data
	}
	if len(outputError) > 0 {
		return nil, outputError
	}
	return outputData, nil
}

func (d *Dir) dispatchRead(filePaths []string) []chan *ContentWrapper {
	const defaultErrKey = "dispatchRead"
	readChans := make([]chan *ContentWrapper, len(filePaths))
	for i, path := range filePaths {
		ch := make(chan *ContentWrapper)
		go func(c chan *ContentWrapper, p string) {
			content, err := ioutil.ReadFile(p)
			if err != nil {
				c <- &ContentWrapper{
					Error: model.BuildError(defaultErrKey, err),
				}
			} else {
				c <- &ContentWrapper{
					Content: content,
				}
			}
		}(ch, path)
		readChans[i] = ch
	}
	return readChans
}

// ReadOne reads the first file in a directory
func (d *Dir) ReadOne() (*model.Data, model.Error) {
	const defaultErrKey = "ReadOne"
	if err := d.validate(); err != nil {
		return nil, model.BuildError(defaultErrKey, err)
	}
	dirPath := d.getPath()
	filePaths, pathErr := d.getFilePaths(dirPath, d.filterPath)
	if pathErr != nil {
		return nil, model.BuildError(defaultErrKey, pathErr)
	}
	content, readErr := ioutil.ReadFile(filePaths[0])
	if readErr != nil {
		return nil, model.BuildError(defaultErrKey, readErr)
	}
	return d.postProcess(filePaths[0], content)
}

func (d *Dir) getFilePaths(dirPath string, filterPath model.FilterPath) ([]string, model.Error) {
	const defaultErrKey = "getFilePaths"
	var output []string
	err := filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if filterPath == nil || filterPath(path) {
				output = append(output, path)
			}
		}
		return nil
	})
	if err != nil {
		return nil, model.BuildError(defaultErrKey, err)
	}
	if len(output) == 0 {
		return nil, model.BuildError(defaultErrKey, errors.New("no file path is found based on filter"))
	}
	return output, nil
}

func (d *Dir) validate() model.Error {
	const defaultErrKey = "validate"
	if d.getPath == nil {
		return model.BuildError(defaultErrKey, errors.New("getPath is nil"))
	}
	if d.postProcess == nil {
		return model.BuildError(defaultErrKey, errors.New("postProcess is nil"))
	}
	return nil
}

// NewReader initializes dir Reader
func NewReader(
	getPath model.GetPath,
	filterPath model.FilterPath,
	postProcess model.PostProcess,
) *Dir {
	return &Dir{
		getPath:     getPath,
		filterPath:  filterPath,
		postProcess: postProcess,
	}
}

// NewWriter initializes dir Writer
func NewWriter() *Dir {
	return &Dir{}
}

func init() {
	err := io.Readers.Register(_type, func(
		getPath model.GetPath,
		filterPath model.FilterPath,
		postProcess model.PostProcess,
	) model.Reader {
		return NewReader(getPath, filterPath, postProcess)
	})
	if err != nil {
		panic(err)
	}
	err = io.Writers.Register(_type, NewWriter())
	if err != nil {
		panic(err)
	}
}
