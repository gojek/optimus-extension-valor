package model

// Evaluate evaluates snippet
type Evaluate func(name, snippet string) (string, error)
