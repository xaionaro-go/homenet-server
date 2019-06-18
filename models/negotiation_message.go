package models

import (
	"net"
	"time"
)

type Protocol string

const (
	ProtocolUDP = Protocol("udp")
)

type NegotiationMessage struct {
	NetworkID               string    `json:"network_id,omitempty"`
	PeerIDFrom              string    `json:"peer_id_from,omitempty"`
	PeerIDTo                string    `json:"peer_id_to,omitempty"`
	Protocol                Protocol  `json:"protocol"`
	SourcePort              uint16    `json:"source_port"`
	LocalAddress            net.IP    `json:"local_address"`
	RequireReverseDirection bool      `json:"require_reverse_direction"`
	CreatedAt               time.Time `json:"created_at,omitempty"`
}

func NewNegotiationMessage() *NegotiationMessage {
	return &NegotiationMessage{}
}

func (msg *NegotiationMessage) GetID() string {
	return msg.NetworkID + "/" + msg.PeerIDTo + "/" + msg.PeerIDFrom
}
func (msg *NegotiationMessage) IGetID() interface{} {
	return msg.GetID()
}
