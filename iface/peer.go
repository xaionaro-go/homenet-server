package iface

import (
	"net"

	"golang.org/x/crypto/ed25519"
)

type Peer interface {
	Object

	GetID() string
	GetIntAlias() uint32
	SetIntAlias(uint32)
	GetPublicKey() ed25519.PublicKey
	GetHost() net.IP
}
