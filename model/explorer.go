package model

// ExplorePath explores path from its root with filter
type ExplorePath func(root string, filter func(string) bool) ([]string, error)
