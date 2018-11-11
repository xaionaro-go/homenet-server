package methods

import (
	"net"

	"github.com/gin-gonic/gin"

	"github.com/xaionaro-go/homenet-server/errors"
	"github.com/xaionaro-go/homenet-server/models"
)

func RegisterPeer(ctx *gin.Context) {
	peerID := ctx.Param("id")
	address := ctx.Request.RemoteAddr

	network := models.GetCTXNetwork(ctx)
	peer := network.GetPeerByID(peerID)
	if peer == nil {
		peer = models.NewPeer(peerID)
	}

	peer.SetAddress(net.ParseIP(address))
	peer.SetNetwork(network)

	if err := peer.Save(); err != nil {
		returnError(ctx, errors.NewCannotSave(peer, err))
		return
	}

	returnSuccess(ctx, peer)
	return
}
