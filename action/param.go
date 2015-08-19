package action

import "errors"

var (
	// ErrEmptyInput ...
	ErrEmptyInput = errors.New("relay/action: empty input")
)

// Param ...
type Param interface {
	IsSet() bool
	IsValid() bool
	IsOptional() bool
	Default() interface{}
	Validate(string) error
	Value() interface{}
}

// Params ...
type Params map[string]Param
