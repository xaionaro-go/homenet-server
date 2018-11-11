package errors

import (
	"fmt"

	"github.com/xaionaro-go/homenet-server/iface"
)

type AlreadyExists struct {
	BadRequest
}

func NewAlreadyExists(obj iface.Object) (result AlreadyExists) {
	result.message = fmt.Sprintf(`%T "%v" is already exists`, obj, obj.IGetID())
	return
}
