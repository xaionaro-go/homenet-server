package methods

import (
	"encoding/base64"

	"github.com/gin-gonic/gin"
	"github.com/xaionaro-go/homenet-server/errors"
	"github.com/xaionaro-go/homenet-server/models"
)

func RegisterNet(ctx *gin.Context) {
	netID := ctx.Param("net")

	network, err := models.NewNetwork(netID)
	if err != nil {
		returnError(ctx, err)
		return
	}

	passwordHash := []byte(ctx.Param("password_hash"))
	if len(passwordHash) == 0 {
		passwordHash, _ = base64.StdEncoding.DecodeString(ctx.Request.Header.Get("X-Homenet-Accesshash"))
	}

	network.SetPasswordHash(passwordHash)
	if err := network.SaveToDisk(); err != nil {
		returnError(ctx, errors.NewCannotSave(network, err))
		return
	}

	returnSuccess(ctx, network)
	return
}
