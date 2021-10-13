package plugin

import "github.com/gojek/optimus-extension-valor/model"

// Reader is a contract to read from a single source
type Reader interface {
	Read() (*model.Data, error)
	Next() bool
}

// Writer is a contract to write data
type Writer interface {
	Write(...*model.Data) error
}
