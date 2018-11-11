package models

import (
	"net"
)

type peer struct {
	id      string
	address net.IP
	network *network
}

type PeerT = peer

func NewPeer(id string) *peer {
	return &peer{
		id: id,
	}
}

type Peers []*peer

func (p *peer) GetID() string {
	return p.id
}
func (p *peer) IGetID() interface{} {
	return p.GetID()
}

func (p *peer) SetNetwork(net *network) {
	p.network = net
	net.appendPeerIfNotExists(p)
}

func (p *peer) SetAddress(address net.IP) {
	p.address = address
}

func (p *peer) Save() error {
	// We don't support peer saving, yet.
	// And it's not so required, actuallity. Clients will resend peers information every 10-60 seconds,
	// so the information will be restored very soon after a restart.
	return nil
}

func (p *peer) Remove() {
	p.network.RemovePeer(p)
}

type peerModel struct {
	peer
}

func Peer() peerModel {
	return peerModel{}
}
