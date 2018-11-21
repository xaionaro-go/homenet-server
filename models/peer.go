package models

import (
	"net"
	//"strconv"

	"github.com/Sirupsen/logrus"
)

type peer struct {
	// this is supposed to be private (non-changable directly from an outside code) but serializable variables. So they're prefixed with "XxX_" to remind users to do not access them directly
	XxX_ID       string `json:"id"`
	XxX_Name     string `json:"name"`
	XxX_Host     net.IP `json:"host"`
	XxX_IntAlias uint32 `json:"int_alias"` // this value is used to generate a persistent IP-address

	network *network
}

type PeerT = peer

func NewPeer(id string) *peer {
	return &peer{
		XxX_ID: id,
	}
}

type Peers []*peer

func (p *peer) GetID() string {
	return p.XxX_ID
}
func (p *peer) IGetID() interface{} {
	return p.GetID()
}
func (p *peer) GetIntAlias() uint32 {
	return p.XxX_IntAlias
}
func (p *peer) SetIntAlias(newIntAlias uint32) {
	p.XxX_IntAlias = newIntAlias
}

func (p *peer) SetName(name string) {
	p.XxX_Name = name
}

func (p *peer) SetNetwork(net *network) error {
	p.network = net

	_, err := net.AppendPeerIfNotExists(p)
	return err
}

func (p *peer) SetAddressByString(address string) {
	//host, portString, err := net.SplitHostPort(address)
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		logrus.Panicf("This shouldn't happened #0 (%v): %v", address, err)
	}
	p.XxX_Host = net.ParseIP(host)
	/*port64, err := strconv.ParseInt(portString, 10, 17)
	if err != nil {
		logrus.Panicf("This shouldn't happened #1 (%v): %v", portString, err)
	}
	p.XxX_Port = uint16(port64)*/
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
