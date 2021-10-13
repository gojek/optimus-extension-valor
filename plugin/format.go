package plugin

// Formatter is a type to format input from one input to another
type Formatter func([]byte) ([]byte, error)
