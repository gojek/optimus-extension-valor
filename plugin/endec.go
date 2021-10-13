package plugin

// Decoder is a type to decode a raw data into a specified type output
type Decoder func([]byte, interface{}) error

// Encoder is a type to encode a specifid type input into an output raw data
type Encoder func(interface{}) ([]byte, error)
