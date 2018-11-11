package methods

import (
	"github.com/gin-gonic/gin"

	"github.com/xaionaro-go/homenet-server/errors"
	"github.com/xaionaro-go/homenet-server/models"
)

func RegisterPeer(ctx *gin.Context) {
	peerID := ctx.Param("id")
	address := ctx.Request.RemoteAddr
	peerName := ctx.Param("name")

	network := models.GetCTXNetwork(ctx)
	peer := network.GetPeerByID(peerID)
	if peer == nil {
		peer = models.NewPeer(peerID)
	}

	peer.SetAddressByString(address)
	peer.SetNetwork(network)
	if peerName != "" {
		peer.SetName(peerName)
	}

	if err := peer.Save(); err != nil {
		returnError(ctx, errors.NewCannotSave(peer, err))
		return
	}

	returnSuccess(ctx, peer)
	return
}
