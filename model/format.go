package model

// Format is a type to format input from one input to another
type Format func([]byte) ([]byte, error)
