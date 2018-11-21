package errors

import (
	"fmt"

	"github.com/xaionaro-go/errors"

	"github.com/xaionaro-go/homenet-server/iface"
)

type TooManyPeers struct {
	BadRequest
}

func NewTooManyPeers(net iface.Network) errors.SmartError {
	result := TooManyPeers{}
	result.message = fmt.Sprintf(`too many peers within network %v: %v (of %v)`, net.GetID(), net.PeersCount(), net.PeersLimit())
	return errors.Wrap(&result)
}
