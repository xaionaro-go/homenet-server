package methods

import (
	"github.com/gin-gonic/gin"

	"github.com/xaionaro-go/homenet-server/errors"
	"github.com/xaionaro-go/homenet-server/models"
)

func SetNegotiationMessage(ctx *gin.Context) {
	peerIDTo := ctx.Param("peeer_id_to")
	peerIDFrom := ctx.Param("peer_id_from")

	network := models.GetCTXNetwork(ctx)

	negotiationMessage := models.NewNegotiationMessage()

	ctx.BindJSON(negotiationMessage)
	negotiationMessage.NetworkID = network.GetID()
	negotiationMessage.PeerIDFrom = peerIDFrom
	negotiationMessage.PeerIDTo = peerIDTo

	if err := network.SetNegotiationMessage(negotiationMessage); err != nil {
		returnError(ctx, errors.NewCannotSave(negotiationMessage, err))
		return
	}

	returnSuccess(ctx, negotiationMessage)
	return
}
