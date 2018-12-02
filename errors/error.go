package errors

import (
	"github.com/xaionaro-go/errors"
)

type anyError struct {
	message string
}

func (e anyError) Error() string {
	return e.message
}

func Wrap(err error) errors.SmartError {
	return errors.Wrap(err)
}
