package model

// Progress is a contract for process progress
type Progress interface {
	Increase(int)
	Wait()
}

// NewProgress is a function to initialize a progress
type NewProgress func(name string, total int) Progress
