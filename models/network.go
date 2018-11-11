package models

import (
	"database/sql"

	"github.com/Sirupsen/logrus"

	"github.com/gin-gonic/gin"

	"github.com/xaionaro-go/atomicmap"

	"github.com/xaionaro-go/homenet-server/errors"
	"github.com/xaionaro-go/homenet-server/helpers"
	"github.com/xaionaro-go/homenet-server/storage"
)

type network struct {
	// this is supposed to be private (non-changable directly from an outside code) but serializable variables. So they're prefixed with "XxX_" to remind users to do not access them directly
	XxX_ID           string `json:"id"`
	XxX_PasswordHash string `json:"password_hash"`

	peers Peers
}

type NetworkT = network

func NewNetwork(id string) (*network, error) {
	oldNet, _ := Network().Get(id)
	if oldNet != nil {
		logrus.Errorf("A network with the same ID is already exists")
		return nil, errors.NewAlreadyExists(oldNet)
	}
	newNet := &network{
		XxX_ID: id,
	}
	Network().Set(newNet)
	return newNet, nil
}

func (net *network) GetPasswordHash() string {
	return net.XxX_PasswordHash
}

func (net *network) GetID() string {
	return net.XxX_ID
}

func (net *network) IGetID() interface{} {
	return net.GetID()
}

func (net *network) SetPasswordHash(newPasswordHash string) *network {
	// Yes we double-hashing the password:
	// - The first time on the client side.
	// - The second time on our side.
	net.XxX_PasswordHash = string(helpers.Hash([]byte(newPasswordHash)))
	return net
}

func (net *network) CheckPasswordHash(newPasswordHash string) bool {
	return helpers.CheckHash([]byte(net.XxX_PasswordHash), []byte(newPasswordHash))
}

func (net *network) GetPeers() Peers {
	return net.peers
}

func (net *network) GetPeerByID(peerID string) *peer {
	for _, peer := range net.peers {
		if peer.GetID() == peerID {
			return peer
		}
	}
	return nil
}

func (net *network) SaveToDisk() error {
	storage.GetSavingQueue() <- &[]network{*net}[0] // copy and send to the queue
	return nil
}

func (net *network) appendPeerIfNotExists(peer *peer) bool {
	oldPeer := net.GetPeerByID(peer.GetID())
	if oldPeer != nil {
		if oldPeer != peer { // two different peers with the same ID? Shouldn't be so
			logrus.Errorf("There already exsists another peer with the same ID. Cannot add to the network.")
		}
		return false
	}

	net.peers = append(net.peers, peer)
	return true
}

func (net *network) RemovePeer(peer *peer) bool {
	if len(net.peers) == 0 {
		return false
	}
	if len(net.peers) == 1 {
		if net.peers[0] == peer {
			net.peers = Peers{}
			return true
		} else {
			return false
		}
	}

	found := false
	newPeers := make(Peers, 0, len(net.peers)-1)
	for _, oldPeer := range net.peers {
		if oldPeer != peer {
			found = true
			continue
		}

		newPeers = append(newPeers, oldPeer)
	}

	if !found {
		return false
	}

	net.peers = newPeers
	return true
}

func (net *network) RemovePeerByID(peerID string) bool {
	peer := net.GetPeerByID(peerID)
	if peer == nil {
		return false
	}
	result := net.RemovePeer(peer)
	if result != true {
		logrus.Errorf("(*network).RemovePeerByID(): Shouldn't happened")
	}
	return result
}

func SetCTXNetwork(ctx *gin.Context, net *network) {
	ctx.Set("network", net)
}

func GetCTXNetwork(ctx *gin.Context) *network {
	net, _ := ctx.Get("network")
	return net.(*network)
}

type networkModel struct {
	network
}

var networkMap = storage.Get(&network{})

func Network() networkModel {
	return networkModel{}
}

func (m networkModel) Get(id string) (*network, error) {
	net, err := networkMap.Get(id)
	switch err {
	case nil:
		return net.(*network), nil
	case atomicmap.NotFound:
		return nil, sql.ErrNoRows
	default:
		return nil, err
	}
}

func (m networkModel) Set(net *network) {
	networkMap.Set(net.GetID(), net)
}
