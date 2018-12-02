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

	msgMapStd := msgMap.ToSTDMap()
	msgMapSerializable := map[string]models.NegotiationMessageT{}
	for k, v := range msgMapStd {
		msgMapSerializable[k.(string)] = v.(models.NegotiationMessageT)
	}

	returnSuccess(ctx, msgMapSerializable)
	return
}
