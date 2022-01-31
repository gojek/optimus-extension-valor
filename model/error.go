package model

import (
	"encoding/json"
	"fmt"
	"sync"
)

// Error defines how execution error is constructed
type Error struct {
	keyToValue map[string]interface{}

	initialized bool
	mtx         *sync.Mutex
}

// Add adds a new error based on a specified key
func (e *Error) Add(key string, value interface{}) {
	if !e.initialized {
		e.keyToValue = make(map[string]interface{})
		e.mtx = &sync.Mutex{}
		e.initialized = true
	}
	e.mtx.Lock()
	e.keyToValue[key] = value
	e.mtx.Unlock()
}

// Error returns the summary of error
func (e *Error) Error() string {
	var output string
	if len(e.keyToValue) > 0 {
		var key string
		for k := range e.keyToValue {
			key = k
			break
		}
		output = fmt.Sprintf("error with key [%s]", key)
		if len(e.keyToValue) > 1 {
			output += fmt.Sprintf(" and %d others", len(e.keyToValue)-1)
		}
	}
	return output
}

// JSON returns the complete error message representation
func (e *Error) JSON() []byte {
	mapError := e.buildMap()
	output, err := json.MarshalIndent(mapError, "", " ")
	if err != nil {
		return []byte(err.Error())
	}
	return output
}

// Length returns the number of errors stored so far
func (e *Error) Length() int {
	return len(e.keyToValue)
}

func (e *Error) buildMap() map[string]interface{} {
	output := make(map[string]interface{})
	for key, value := range e.keyToValue {
		if customErr, ok := value.(Error); ok {
			mV := customErr.buildMap()
			output[key] = mV
		} else {
			if err, ok := value.(error); ok {
				output[key] = err.Error()
			} else {
				output[key] = value
			}
		}
	}
	return output
}
