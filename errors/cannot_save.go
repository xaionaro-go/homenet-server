package errors

import (
	"fmt"

	"github.com/xaionaro-go/homenet-server/iface"
)

type CannotSave struct {
	InternalError
}

func NewCannotSave(obj iface.Object, err error) (result CannotSave) {
	result.message = fmt.Sprintf(`cannot save %T "%v": %s`, obj, obj.IGetID(), err.Error())
	return
}
