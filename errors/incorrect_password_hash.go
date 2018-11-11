package errors

import (
	"fmt"

	"github.com/xaionaro-go/homenet-server/iface"
)

type IncorrectPasswordHash struct {
	BadRequest
}

func NewIncorrectPasswordHash(obj iface.Object) (result IncorrectPasswordHash) {
	result.message = fmt.Sprintf(`incorrect password hash for %T "%v"`, obj, obj.IGetID())
	return
}
