package models

import (
	"net"
	"time"
)

type Protocol string

const (
	ProtocolUDP = Protocol("udp")
)

type negotiationMessage struct {
	NetworkID    string    `json:"network_id"`
	PeerIDFrom   string    `json:"peer_id_from"`
	PeerIDTo     string    `json:"peer_id_to"`
	Protocol     Protocol  `json:"protocol"`
	SourcePort   uint16    `json:"source_port"`
	LocalAddress net.IP    `json:"local_address"`
	CreatedAt    time.Time `json:"created_at"`
}

type NegotiationMessageT negotiationMessage

func NewNegotiationMessage() *negotiationMessage {
	return &negotiationMessage{}
}

func (msg *negotiationMessage) GetID() string {
	return msg.NetworkID + "/" + msg.PeerIDTo + "/" + msg.PeerIDFrom
}
func (msg *negotiationMessage) IGetID() interface{} {
	return msg.GetID()
}

type negotiationMessageModel struct {
	negotiationMessage
}

func NegotiationMessage() negotiationMessageModel {
	return negotiationMessageModel{}
}
