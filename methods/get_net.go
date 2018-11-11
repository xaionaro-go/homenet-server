package methods

import (
	"github.com/gin-gonic/gin"

	"github.com/xaionaro-go/homenet-server/models"
)

func GetNet(ctx *gin.Context) {
	returnSuccess(ctx, models.GetCTXNetwork(ctx))
	return
}
