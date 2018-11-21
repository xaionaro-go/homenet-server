package iface

type Peer interface {
	Object

	GetID() string
	GetIntAlias() uint32
	SetIntAlias(uint32)
}
