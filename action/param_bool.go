package action

import "strconv"

// ParamBool ...
type ParamBool struct {
	isOptional bool
	d          bool
	value      bool
	isValid    bool
	isSet      bool
}

// NewParamBool ...
func NewParamBool(d, isOptional bool) *ParamBool {
	return &ParamBool{
		d:          d,
		isOptional: isOptional,
	}
}

// Validate ...
func (pb *ParamBool) Validate(input string) error {
	if input == "" {
		pb.isSet = true
		pb.isValid = true
		pb.value = true

		return nil
	}

	v, err := strconv.ParseBool(input)
	if err != nil {
		pb.isSet = true

		return err
	}

	pb.isSet = true
	pb.isValid = true
	pb.value = v

	return nil
}

// IsSet ...
func (pb *ParamBool) IsSet() bool {
	return pb.isSet
}

// IsValid ...
func (pb *ParamBool) IsValid() bool {
	return pb.isValid
}

// IsOptional ...
func (pb *ParamBool) IsOptional() bool {
	return pb.isOptional
}

// Default ...
func (pb *ParamBool) Default() interface{} {
	return pb.d
}

// Value ...
func (pb *ParamBool) Value() interface{} {
	if pb.isSet {
		return pb.value
	}

	return pb.Default()
}
