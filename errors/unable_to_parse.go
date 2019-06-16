package errors

import (
	"fmt"
)

type UnableToParse struct {
	BadRequest
}

func NewUnableToParse(err error) error {
	result := UnableToParse{}
	result.message = fmt.Sprintf(`unable to parse: %v`, err)
	return &result
}
