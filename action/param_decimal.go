package action

import "strconv"

// ParamDecimal ...
type ParamDecimal struct {
	d          int64
	isOptional bool
	value      int64
	isValid    bool
	isSet      bool
}

// NewParamDecimal ...
func NewParamDecimal(d int64, io bool) *ParamDecimal {
	return &ParamDecimal{
		d:          d,
		isOptional: io,
	}
}

// Validate ...
func (pd *ParamDecimal) Validate(input string) error {
	if input == "" {
		return ErrEmptyInput
	}

	v, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		pd.isSet = true
		return err
	}

	pd.isSet = true
	pd.isValid = true
	pd.value = v

	return nil
}

// IsSet ...
func (pd *ParamDecimal) IsSet() bool {
	return pd.isSet
}

// IsValid ...
func (pd *ParamDecimal) IsValid() bool {
	return pd.isValid
}

// IsOptional ...
func (pd *ParamDecimal) IsOptional() bool {
	return pd.isOptional
}

// Default ...
func (pd *ParamDecimal) Default() interface{} {
	return pd.d
}

// Value ...
func (pd *ParamDecimal) Value() interface{} {
	if pd.isSet {
		return pd.value
	}

	return pd.Default()
}
