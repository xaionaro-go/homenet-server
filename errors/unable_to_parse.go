package errors

import (
	"fmt"

	"github.com/xaionaro-go/errors"
)

type UnableToParse struct {
	BadRequest
}

func NewUnableToParse(err error) errors.SmartError {
	result := UnableToParse{}
	result.message = fmt.Sprintf(`unable to parse: %v`, err)
	return errors.Wrap(&result)
}
