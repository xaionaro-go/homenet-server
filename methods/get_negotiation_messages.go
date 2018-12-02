package methods

import (
	"github.com/gin-gonic/gin"

	"github.com/xaionaro-go/homenet-server/errors"
	"github.com/xaionaro-go/homenet-server/models"
)

func GetNegotiationMessages(ctx *gin.Context) {
	peerIDTo := ctx.Param("peer_id_to")

	network := models.GetCTXNetwork(ctx)

	msgMap := network.GetNegotiationMessagesMap(peerIDTo)
	if msgMap == nil {
		returnError(ctx, errors.NewGetObjectNotFound(models.NegotiationMessage(), peerIDTo, network))
		return
	}

	returnSuccess(ctx, msgMap.ToSTDMap())
	return
}
