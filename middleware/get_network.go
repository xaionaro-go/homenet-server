package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/xaionaro-go/homenet-server/errors"
	"github.com/xaionaro-go/homenet-server/models"
)

func GetNetwork() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		netID := ctx.Param("net")
		passwordHash := ctx.Param("password_hash")

		net, err := models.Network().Get(netID)
		if err != nil {
			returnError(ctx, errors.NewGetObject(models.Network, netID, err))
			return
		}

		if !net.CheckPasswordHash(passwordHash) {
			returnError(ctx, errors.NewIncorrectPasswordHash(net))
			return
		}

		models.SetCTXNetwork(ctx, net)
		ctx.Next()
	}
}
