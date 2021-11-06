package model

// Reader is a contract to read from a single source
type Reader interface {
	ReadOne() (*Data, Error)
	ReadAll() ([]*Data, Error)
}

// Writer is a contract to write data
type Writer interface {
	Write(...*Data) Error
}

// GetPath gets the path
type GetPath func() string

// FilterPath filters the path with True value will be processed
type FilterPath func(string) bool

// PostProcess post processes data
type PostProcess func(path string, content []byte) (*Data, Error)
