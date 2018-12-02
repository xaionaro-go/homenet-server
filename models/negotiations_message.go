package models

import (
	"net"
	"time"
)

type negotiationsMessage struct {
	PeerIDFrom   string    `json:"peer_id_from"`
	PeerIDTo     string    `json:"peer_id_to"`
	Protocol     string    `json:"protocol"`
	SourcePort   uint16    `json:"source_port"`
	LocalAddress net.IP    `json:"local_address"`
	CreatedAt    time.Time `json:"created_at"`
}
