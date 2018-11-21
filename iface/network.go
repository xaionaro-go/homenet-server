package iface

type Network interface {
	Object

	GetID() string
	PeersLimit() uint32
	PeersCount() uint32
}
