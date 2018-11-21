package errors

import (
	"fmt"

	"github.com/xaionaro-go/errors"

	"github.com/xaionaro-go/homenet-server/iface"
)

type PeerIDIsBusy struct {
	AlreadyExists
}

func NewPeerIDIsBusy(peerID string, net iface.Network) errors.SmartError {
	result := PeerIDIsBusy{}
	result.message = fmt.Sprintf(`peer ID %v is already busy in network %v`, peerID, net.GetID())
	return errors.Wrap(&result)
}
