package methods

import (
	"github.com/gin-gonic/gin"
	"github.com/xaionaro-go/homenet-server/errors"
	"github.com/xaionaro-go/homenet-server/models"
)

func UnregisterPeer(ctx *gin.Context) {
	peerID := ctx.Param("id")

	net := models.GetCTXNetwork(ctx)
	if !net.RemovePeerByID(peerID) {
		returnError(ctx, errors.NewGetObjectNotFound(models.Peer(), peerID, net))
		return
	}

	returnSuccess(ctx, nil)
	return
}
