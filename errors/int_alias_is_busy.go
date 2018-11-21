package errors

import (
	"fmt"

	"github.com/xaionaro-go/errors"

	"github.com/xaionaro-go/homenet-server/iface"
)

type IntAliasIsBusy struct {
	AlreadyExists
}

func NewIntAliasIsBusy(intAlias uint32, net iface.Network) errors.SmartError {
	result := IntAliasIsBusy{}
	result.message = fmt.Sprintf(`integer alias %v is already busy in network %v`, intAlias, net.GetID())
	return errors.Wrap(&result)
}
