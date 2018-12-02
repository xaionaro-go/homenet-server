package models

import (
	"net"
	"time"

	"github.com/xaionaro-go/atomicmap"

	"github.com/xaionaro-go/homenet-server/errors"
)

type negotiationMessage struct {
	NetworkID    string    `json:"network_id"`
	PeerIDFrom   string    `json:"peer_id_from"`
	PeerIDTo     string    `json:"peer_id_to"`
	Protocol     string    `json:"protocol"`
	SourcePort   uint16    `json:"source_port"`
	LocalAddress net.IP    `json:"local_address"`
	CreatedAt    time.Time `json:"created_at"`
}

func NewNegotiationMessage(networkID, peerIDFrom, peerIDTo string) *negotiationMessage {
	return &negotiationMessage{
		NetworkID:  networkID,
		PeerIDFrom: peerIDFrom,
		PeerIDTo:   peerIDTo,
	}
}

var negotiationMessagesMap = atomicmap.New()

func negotiationMessagesCreateMapForPeer(p *peer) error {
	net := p.GetNetwork()
	key := net.GetID() + "/" + p.GetID()
	negotiationMessagesMap.Set(key, atomicmap.NewWithArgs(uint64(net.PeersLimit()), nil))
	return nil
}

func (msg *negotiationMessage) GetMapKey() string {
	return msg.NetworkID + "/" + msg.PeerIDTo
}

func (msg *negotiationMessage) GetID() string {
	return msg.GetMapKey() + "/" + msg.PeerIDFrom
}
func (msg *negotiationMessage) IGetID() interface{} {
	return msg.GetID()
}

func (msg *negotiationMessage) Save() error {
	mI, _ := negotiationMessagesMap.Get(msg.GetMapKey)
	if mI == nil {
		return errors.NewGetObjectNotFound(&peer{}, msg.PeerIDTo, "networkID:"+msg.NetworkID, "negotiationMessage storage")
	}
	m := mI.(atomicmap.Map)

	m.Set(msg.PeerIDFrom, *msg)
	return nil
}
