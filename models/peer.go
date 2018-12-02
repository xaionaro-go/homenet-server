package models

import (
	"net"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/xaionaro-go/homenet-server/errors"
)

type peer struct {
	// this is supposed to be private (non-changable directly from an outside code) but serializable variables. So they're prefixed with "XxX_" to remind users to do not access them directly
	XxX_ID       string  `json:"id"`
	XxX_Name     *string `json:"name"`
	XxX_Host     *net.IP `json:"host"`
	XxX_IntAlias uint32  `json:"int_alias"` // this value is used to generate a persistent IP-address

	updatedAt *time.Time
	network   *network
}

type PeerT = peer

func NewPeer(id string) *peer {
	return &peer{
		XxX_ID:    id,
		updatedAt: &[]time.Time{time.Now()}[0],
	}
}

type Peers []*peer

func (p *peer) updated() {
	atomic.StorePointer((*unsafe.Pointer)((unsafe.Pointer)(&p.updatedAt)), (unsafe.Pointer)(&[]time.Time{time.Now()}[0]))
}

func (p *peer) GetID() string {
	return p.XxX_ID
}
func (p *peer) IGetID() interface{} {
	return p.GetID()
}
func (p *peer) GetIntAlias() uint32 {
	return atomic.LoadUint32(&p.XxX_IntAlias)
}
func (p *peer) SetIntAlias(newIntAlias uint32) {
	atomic.StoreUint32(&p.XxX_IntAlias, newIntAlias)
	p.updated()
}
func (p *peer) GetName() string {
	if p.XxX_Name == nil {
		return ""
	}
	return *(*string)(atomic.LoadPointer((*unsafe.Pointer)((unsafe.Pointer)(&p.XxX_Name))))
}

func (p *peer) SetName(name string) {
	atomic.StorePointer((*unsafe.Pointer)((unsafe.Pointer)(&p.XxX_Name)), (unsafe.Pointer)(&name))
	p.updated()
}

func (p *peer) SetNetwork(net *network) error {
	if p.network != nil {
		return errors.NewNotImplemented("Cannot update network of a peer").Wrap()
	}
	p.network = net

	_, err := net.AppendPeerIfNotExists(p)
	p.updated()
	return errors.Wrap(err)
}

func (p *peer) GetHost() net.IP {
	if p.XxX_Host == nil {
		panic("Not initialized peer")
	}

	return *(*net.IP)(atomic.LoadPointer((*unsafe.Pointer)((unsafe.Pointer)(&p.XxX_Host))))
}

func (p *peer) SetAddressByString(address string) {
	//host, portString, err := net.SplitHostPort(address)
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		//logrus.Panicf("This shouldn't happened #0 (%v): %v", address, err)
		host = address
	}
	atomic.StorePointer((*unsafe.Pointer)((unsafe.Pointer)(&p.XxX_Host)), (unsafe.Pointer)(&[]net.IP{net.ParseIP(host)}[0]))
	/*port64, err := strconv.ParseInt(portString, 10, 17)
	if err != nil {
		logrus.Panicf("This shouldn't happened #1 (%v): %v", portString, err)
	}
	p.XxX_Port = uint16(port64)*/
	p.updated()
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
