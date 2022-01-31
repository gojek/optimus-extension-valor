package model

// Reader is a contract to read from a source
type Reader interface {
	Read() (*Data, error)
}

// Writer is a contract to write data
type Writer interface {
	Write(*Data) error
}

// GetPath gets the path
type GetPath func() string

// PostProcess post processes data
type PostProcess func(path string, content []byte) (*Data, error)
