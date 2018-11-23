package middleware

import (
	"encoding/base64"

	"github.com/gin-gonic/gin"

	"github.com/xaionaro-go/homenet-server/errors"
	"github.com/xaionaro-go/homenet-server/models"
)

func GetNetwork() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		netID := ctx.Param("net")
		passwordHash := ctx.Param("password_hash")
		if passwordHash == "" {
			passwordHashB, _ := base64.StdEncoding.DecodeString(ctx.Request.Header.Get("X-Homenet-Accesshash"))
			passwordHash = string(passwordHashB)
		}

		net, err := models.Network().Get(netID)
		if net == nil || err != nil {
			returnError(ctx, errors.NewGetObject(models.Network(), netID, err))
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
