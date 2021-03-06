package model

// Decode is a type to decode a raw data into a specified type output
type Decode func([]byte, interface{}) error

// Encode is a type to encode a specifid type input into an output raw data
type Encode func(interface{}) ([]byte, error)
