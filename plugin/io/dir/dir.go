package dir

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/gojek/optimus-extension-valor/model"
	"github.com/gojek/optimus-extension-valor/registry/io"
)

const _type = "dir"

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

	wg := &sync.WaitGroup{}
	mtx := &sync.Mutex{}

	outputError := make(model.Error)
	for _, data := range dataList {
		wg.Add(1)

		go func(w *sync.WaitGroup, m *sync.Mutex, dt *model.Data) {
			defer w.Done()

			key := fmt.Sprintf("%s [%s]", defaultErrKey, dt.Path)
			dirPath, _ := path.Split(dt.Path)
			if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
				m.Lock()
				outputError[key] = err
				m.Unlock()
				return
			}
			if err := ioutil.WriteFile(dt.Path, dt.Content, os.ModePerm); err != nil {
				m.Lock()
				outputError[key] = err
				m.Unlock()
				return
			}
		}(wg, mtx, data)
	}
	wg.Wait()
	if len(outputError) > 0 {
		return outputError
	}
	return nil
}

// ReadAll reads all files in a directory
func (d *Dir) ReadAll() ([]*model.Data, model.Error) {
	const defaultErrKey = "ReadAll"
	if err := d.validate(); err != nil {
		return nil, err
	}
	dirPath := d.getPath()
	filePaths, err := d.getFilePaths(dirPath, d.filterPath)
	if err != nil {
		return nil, err
	}

	wg := &sync.WaitGroup{}
	mtx := &sync.Mutex{}

	outputData := make([]*model.Data, len(filePaths))
	outputError := make(model.Error)
	for i, path := range filePaths {
		wg.Add(1)

		go func(idx int, w *sync.WaitGroup, m *sync.Mutex, pt string) {
			defer w.Done()

			key := fmt.Sprintf("%s [%s]", defaultErrKey, pt)
			content, err := ioutil.ReadFile(pt)
			if err != nil {
				m.Lock()
				outputError[key] = err
				m.Unlock()
			} else {
				data, err := d.postProcess(pt, content)
				if err != nil {
					m.Lock()
					outputError[key] = err
					m.Unlock()
				} else {
					m.Lock()
					outputData[idx] = data
					m.Unlock()
				}
			}
		}(i, wg, mtx, path)
	}
	wg.Wait()
	if len(outputError) > 0 {
		return nil, outputError
	}
	return outputData, nil
}

// ReadOne reads the first file in a directory
func (d *Dir) ReadOne() (*model.Data, model.Error) {
	const defaultErrKey = "ReadOne"
	if err := d.validate(); err != nil {
		return nil, err
	}
	dirPath := d.getPath()
	filePaths, pathErr := d.getFilePaths(dirPath, d.filterPath)
	if pathErr != nil {
		return nil, pathErr
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
