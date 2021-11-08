package model

import (
	"encoding/json"
	"fmt"
)

// Error is error with in key-value formatted
type Error map[string]interface{}

// Error returns error that represent the field error
func (e Error) Error() string {
	var output string
	if len(e) > 0 {
		var key string
		for k := range e {
			key = k
			break
		}
		output = fmt.Sprintf("error with key [%s]", key)
		if len(e) > 1 {
			output = output + " " + fmt.Sprintf("and %d others", len(e)-1)
		}
	}
	return output
}

// JSON converts field error into its JSON representation
func (e Error) JSON() []byte {
	mapError := e.buildMap()
	output, err := json.MarshalIndent(mapError, "", " ")
	if err != nil {
		return []byte(err.Error())
	}
	return output
}

func (e Error) buildMap() map[string]interface{} {
	output := make(map[string]interface{})
	for key, value := range e {
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

// BuildError builds a new error based on key and value
func BuildError(key string, value interface{}) Error {
	return map[string]interface{}{
		key: value,
	}
}
