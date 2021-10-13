package dir

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gojek/optimus-extension-valor/model"
	"github.com/gojek/optimus-extension-valor/plugin"
	"github.com/gojek/optimus-extension-valor/registry/io"
)

const _type = "dir"

// Dir represents directory operation
type Dir struct {
	dirPath  string
	metadata map[string]string

	populated bool

	filePathIndx int
	filePaths    []string
}

func (d *Dir) Write(dataList ...*model.Data) error {
	if len(dataList) == 0 {
		return errors.New("input is empty")
	}
	const pathKey = "path"
	const extensionKey = "extension"
	var missingDirKeys []string
	if d.metadata[pathKey] == "" {
		missingDirKeys = append(missingDirKeys, pathKey)
	}
	if d.metadata[extensionKey] == "" {
		missingDirKeys = append(missingDirKeys, extensionKey)
	}
	if len(missingDirKeys) > 0 {
		return fmt.Errorf("[%s] are empty in directory metadata", strings.Join(missingDirKeys, ", "))
	}
	for _, data := range dataList {
		if data.Metadata[pathKey] == "" {
			return fmt.Errorf("[%s] are empty in file metadata", pathKey)
		}
		separator := "/"
		dirPath := d.dirPath
		filePath := data.Path
		extension := d.metadata[extensionKey]
		path := path.Join(dirPath, fmt.Sprintf("%s.%s", filePath, extension))

		splitFilePath := strings.Split(path, separator)
		targetDirPath := strings.Join(splitFilePath[:len(splitFilePath)-1], separator)
		if err := os.MkdirAll(targetDirPath, os.ModePerm); err != nil {
			return err
		}
		if err := ioutil.WriteFile(path, data.Content, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func (d *Dir) Read() (*model.Data, error) {
	if !d.populated {
		if err := d.populate(); err != nil {
			return nil, err
		}
	}
	if !d.Next() {
		return nil, errors.New("no available Next item")
	}
	path := d.filePaths[d.filePathIndx]
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	d.filePathIndx++
	return &model.Data{
		Path:     path,
		Content:  content,
		Metadata: d.generateNewMetadata(),
	}, nil
}

func (d *Dir) generateNewMetadata() map[string]string {
	output := make(map[string]string)
	targetKey := "type"
	targetValue := "file"
	for key, value := range d.metadata {
		if key == targetKey {
			value = targetValue
		}
		output[key] = value
	}
	return output
}

func (d *Dir) populate() error {
	var filePaths []string
	err := filepath.Walk(d.dirPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			filePaths = append(filePaths, path)
		}
		return nil
	})
	if err != nil {
		return err
	}
	d.filePaths = filePaths
	d.populated = true
	return nil
}

// Next checks whether there is a next item to Read or not
func (d *Dir) Next() bool {
	if !d.populated {
		return true
	}
	return d.filePathIndx < len(d.filePaths)
}

// New initializes Dir based on path
func New(dirPath string, metadata map[string]string) *Dir {
	return &Dir{
		dirPath:  dirPath,
		metadata: metadata,
	}
}

func init() {
	io.Readers.Register(_type, func(path string, metadata map[string]string) plugin.Reader {
		return New(path, metadata)
	})
	io.Writers.Register(_type, func(path string, metadata map[string]string) plugin.Writer {
		return New(path, metadata)
	})
}
