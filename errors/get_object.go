package errors

import (
	"database/sql"
	"fmt"

	"github.com/xaionaro-go/homenet-server/iface"
)

type GetObjectInvalidRequest struct {
	BadRequest
}

type GetObjectNotFound struct {
	BadRequest
}

func (e GetObjectNotFound) IsNotFound() {
}

func NewGetObjectNotFound(model iface.Model, id interface{}, additionalInfo ...interface{}) (result GetObjectNotFound) {
	result.message = fmt.Sprintf(`%T "%v" not found (%v)`, model, id, additionalInfo)
	return result
}

type GetObjectInternalError struct {
	InternalError
}

func NewGetObject(model iface.Model, id interface{}, err error) error {
	switch err {
	case sql.ErrNoRows:
		return NewGetObjectNotFound(model, id)
	default:
		result := GetObjectInternalError{}
		result.message = fmt.Sprintf(`unable to get %T "%v": %s`, model, id, err.Error())
		return result
	}
}
