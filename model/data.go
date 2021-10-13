package model

// Data contains information with type, path, and raw content
type Data struct {
	Path     string
	Content  []byte
	Metadata map[string]string
}
