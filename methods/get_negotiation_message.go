package methods

import (
	"github.com/gin-gonic/gin"

	"github.com/xaionaro-go/homenet-server/errors"
	"github.com/xaionaro-go/homenet-server/models"
)

func GetNegotiationMessage(ctx *gin.Context) {
	peerIDTo := ctx.Param("peeer_id_to")
	peerIDFrom := ctx.Param("peer_id_from")

	network := models.GetCTXNetwork(ctx)

	msgMap := network.GetNegotiationMessagesMap(peerIDTo)
	if msgMap == nil {
		returnError(ctx, errors.NewGetObjectNotFound(models.NegotiationMessage(), peerIDTo, network))
		return
	}

	msg, _ := msgMap.Get(peerIDFrom)
	if msg == nil {
		returnError(ctx, errors.NewGetObjectNotFound(models.NegotiationMessage(), peerIDFrom, network, peerIDTo))
		return
	}

	returnSuccess(ctx, msg)
	return
}
