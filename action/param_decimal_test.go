package action

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewParamDecimalInt64(t *testing.T) {
	success := map[string]int64{
		"123123": 123123,
		"123":    123,
		"012":    12,
		"1":      1,
		"0":      0,
	}

	for input, expected := range success {
		p := NewParamDecimal(0, false)
		err := p.Validate(input)

		if assert.NoError(t, err) {
			assert.Equal(t, p.Value().(int64), expected)
		}
	}

	fail := map[string]string{
		"1231asd23": "strconv.ParseInt: parsing \"1231asd23\": invalid syntax",
		"-":         "strconv.ParseInt: parsing \"-\": invalid syntax",
		"":          ErrEmptyInput.Error(),
	}

	for input, expected := range fail {
		p := NewParamDecimal(0, false)
		err := p.Validate(input)

		assert.EqualError(t, err, expected)
	}
}
