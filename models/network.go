package models

import (
	"database/sql"
	"sync"

	"github.com/Sirupsen/logrus"

	"github.com/gin-gonic/gin"

	"github.com/xaionaro-go/atomicmap"

	"github.com/xaionaro-go/homenet-server/errors"
	"github.com/xaionaro-go/homenet-server/helpers"
	"github.com/xaionaro-go/homenet-server/storage"
)

type networkInternals struct {
	// this is supposed to be private (non-changable directly from an outside code) but serializable variables. So they're prefixed with "XxX_" to remind users to do not access them directly
	XxX_ID           string `json:"id"`
	XxX_PasswordHash []byte `json:"password_hash"`

	peers                  Peers
	intAliasMap            atomicmap.Map
	negotiationMessagesMap atomicmap.Map
}

type network struct {
	mutex sync.Mutex

	XxX_Internals networkInternals `json:"networkData"`
}

type NetworkT = network

func NewNetwork(id string) (*network, error) {
	oldNet, _ := Network().Get(id)
	if oldNet != nil {
		logrus.Errorf("A network with the same ID is already exists")
		return nil, errors.NewAlreadyExists(oldNet)
	}
	newNet := &network{
		XxX_Internals: networkInternals{
			XxX_ID: id,
		},
	}
	newNet.XxX_Internals.AfterRestore()
	Network().Set(newNet)
	return newNet, nil
}

func (net *networkInternals) AfterRestore() error {
	net.intAliasMap = atomicmap.New()
	net.negotiationMessagesMap = atomicmap.New()
	return nil
}
func (net *network) AfterRestore() error {
	return net.XxX_Internals.AfterRestore()
}

func (net *network) Lock(fn func(*networkInternals)) {
	net.mutex.Lock()
	defer func() {
		if r := recover(); r != nil {
			net.mutex.Unlock()
			logrus.Panicln(r)
		}
	}() // defer is slow, so to unlock the network earlier we try to unlock it outside of the defer func. However if we got a panic then we need to recover and to unlock the network.
	fn(&net.XxX_Internals)
	net.mutex.Unlock()
}

func (net *networkInternals) GetPasswordHash() []byte {
	return net.XxX_PasswordHash
}
func (net *network) GetPasswordHash() (result []byte) {
	net.Lock(func(net *networkInternals) {
		result = net.GetPasswordHash()
	})
	return
}

func (net *networkInternals) GetID() string {
	return net.XxX_ID
}
func (net *networkInternals) IGetID() interface{} {
	return net.GetID()
}
func (net *network) GetID() (result string) {
	net.Lock(func(net *networkInternals) {
		result = net.GetID()
	})
	return
}
func (net *network) IGetID() interface{} {
	return net.GetID()
}

func (net *networkInternals) SetPasswordHash(newPasswordHash []byte) *networkInternals {
	// Yes we double-hashing the password:
	// - The first time on the client side.
	// - The second time on our side.
	net.XxX_PasswordHash = helpers.Hash(newPasswordHash)
	return net
}
func (net *network) SetPasswordHash(newPasswordHash []byte) *network {
	net.Lock(func(net *networkInternals) {
		net.SetPasswordHash(newPasswordHash)
	})
	return net
}

func (net *networkInternals) CheckPasswordHash(passwordHash []byte) bool {
	return helpers.CheckHash(net.GetPasswordHash(), passwordHash)
}
func (net *network) CheckPasswordHash(passwordHash []byte) (result bool) {
	net.Lock(func(net *networkInternals) {
		result = net.CheckPasswordHash(passwordHash)
	})
	return
}

func (net *networkInternals) GetPeers() Peers {
	return net.peers.Copy()
}
func (net *network) GetPeers() (result Peers) {
	net.Lock(func(net *networkInternals) {
		result = net.GetPeers()
	})
	return
}

func (net *networkInternals) GetPeerByID(peerID string) *peer {
	for _, peer := range net.peers {
		if peer.GetID() == peerID {
			return peer
		}
	}
	return nil
}
func (net *network) GetPeerByID(peerID string) (result *peer) {
	net.Lock(func(net *networkInternals) {
		result = net.GetPeerByID(peerID)
	})
	return
}

func (net *network) SaveToDisk() error {
	net.Lock(func(netInternals *networkInternals) {
		storage.GetSavingQueue() <- &[]network{*net}[0] // copy and send to the queue
	})
	return nil
}

func (net *networkInternals) PeersLimit() uint32 {
	return (1 << 8) - 5
}
func (net *network) PeersLimit() uint32 {
	return net.PeersLimit()
}

func (net *networkInternals) PeersCount() uint32 {
	return uint32(len(net.GetPeers()))
}
func (net *network) peersCount() (result uint32) {
	net.Lock(func(net *networkInternals) {
		result = net.PeersCount()
	})
	return
}

func (net *networkInternals) GetNegotiationMessagesMap(peerIDTo string) atomicmap.Map {
	r, _ := net.negotiationMessagesMap.Get(peerIDTo)
	if r == nil {
		return nil
	}
	return r.(atomicmap.Map)
}
func (net *network) GetNegotiationMessagesMap(peerIDTo string) atomicmap.Map {
	return net.XxX_Internals.GetNegotiationMessagesMap(peerIDTo)
}

func (net *networkInternals) AppendPeerIfNotExists(peer *peer) (bool, error) {
	if net.PeersCount() >= net.PeersLimit() {
		return false, errors.NewTooManyPeers(net).Wrap()
	}

	oldPeer := net.GetPeerByID(peer.GetID())
	if oldPeer != nil {
		if oldPeer != peer { // two different peers with the same ID? Shouldn't be so
			return false, errors.NewPeerIDIsBusy(peer.GetID(), net).Wrap()
		}
		return false, nil // already exists, not adding, but no error
	}

	if peer.GetIntAlias() == 0 {
		freeIntAlias := net.FindFreePeerIntAlias()
		if freeIntAlias == 0 {
			return false, errors.NewIntAliasIsBusy(0, net).Wrap()
		}
		peer.SetIntAlias(freeIntAlias)
	} else {
		anotherPeer := net.GetPeerByIntAlias(peer.GetIntAlias())
		if anotherPeer != nil {
			return false, errors.NewIntAliasIsBusy(peer.GetIntAlias(), net).Wrap()
		}
	}

	net.peers = append(net.peers, peer)
	net.intAliasMap.Set(peer.GetIntAlias(), peer)
	net.negotiationMessagesMap.Set(peer.GetID(), atomicmap.New())
	return true, nil
}

func (net *network) AppendPeerIfNotExists(peer *peer) (result0 bool, result1 error) {
	net.Lock(func(net *networkInternals) {
		result0, result1 = net.AppendPeerIfNotExists(peer)
	})
	return
}

func (net *networkInternals) FindFreePeerIntAlias() uint32 {
	for intAlias := uint32(1); intAlias < net.PeersLimit(); intAlias++ {
		if net.GetPeerByIntAlias(intAlias) == nil {
			return intAlias
		}
	}
	return 0 // not found
}
func (net *network) FindFreePeerIntAlias() (result uint32) {
	net.Lock(func(net *networkInternals) {
		result = net.FindFreePeerIntAlias()
	})
	return
}

func (net *networkInternals) GetPeerByIntAlias(intAlias uint32) *peer {
	if p, err := net.intAliasMap.Get(intAlias); err == nil {
		return p.(*peer)
	}
	return nil
}
func (net *network) GetPeerByIntAlias(intAlias uint32) (result *peer) {
	net.Lock(func(net *networkInternals) {
		result = net.GetPeerByIntAlias(intAlias)
	})
	return
}

func (net *networkInternals) RemovePeer(peer *peer) bool {
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
func (net *network) RemovePeer(peer *peer) (result bool) {
	net.Lock(func(net *networkInternals) {
		result = net.RemovePeer(peer)
	})
	return
}

func (net *networkInternals) RemovePeerByID(peerID string) bool {
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
func (net *network) RemovePeerByID(peerID string) (result bool) {
	net.Lock(func(net *networkInternals) {
		result = net.RemovePeerByID(peerID)
	})
	return
}

func (net *networkInternals) SetNegotiationMessage(msg *NegotiationMessage) error {
	mI, _ := net.negotiationMessagesMap.Get(msg.PeerIDTo)
	if mI == nil {
		return errors.NewGetObjectNotFound(&peer{}, msg.PeerIDTo, "networkID:"+net.GetID(), "network's internals")
	}
	m := mI.(atomicmap.Map)

	m.Set(msg.PeerIDFrom, *msg)
	return nil
}

func (net *network) SetNegotiationMessage(msg *NegotiationMessage) error {
	return net.XxX_Internals.SetNegotiationMessage(msg)
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
