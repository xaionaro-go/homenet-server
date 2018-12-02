package errors

import (
	"fmt"

	"github.com/xaionaro-go/errors"
)

type NotImplemented struct {
	BadRequest
}

func NewNotImplemented(description string) errors.SmartError {
	result := NotImplemented{}
	result.message = fmt.Sprintf(`not implemented, yet: %v`, description)
	return errors.Wrap(&result)
}
