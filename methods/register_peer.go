package methods

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xaionaro-go/secureio"

	"github.com/xaionaro-go/homenet-server/errors"
	"github.com/xaionaro-go/homenet-server/models"
)

func RegisterPeer(ctx *gin.Context) {
	peerID := ctx.Param("id")
	peerName := ctx.Query("peer_name")
	publicKeyEncoded := ctx.Query("public_key")

	publicKeyDecoder := base64.NewDecoder(base64.URLEncoding, bytes.NewReader([]byte(publicKeyEncoded)))

	var publicKey [secureio.PublicKeySize]byte
	_, err := publicKeyDecoder.Read(publicKey[:])
	if err != nil {
		returnError(ctx, errors.NewUnableToParse(fmt.Errorf("invalid public key: %v", err)))
		return
	}

	address := ctx.Request.Header.Get("X-Forwarded-For")
	if len(address) == 0 {
		address = ctx.Request.RemoteAddr
	}

	network := models.GetCTXNetwork(ctx)
	peer := network.GetPeerByID(peerID)
	if peer == nil {
		peer = models.NewPeer(peerID)
		if err := peer.SetNetwork(network); err != nil {
			returnError(ctx, errors.NewCannotSave(peer, err))
			return
		}
	}

	peer.SetAddressByString(address)
	peer.SetPublicKey(publicKey)
	if peerName != "" {
		peer.SetName(peerName)
	}

	if err := peer.Save(); err != nil {
		returnError(ctx, errors.NewCannotSave(peer, err))
		return
	}

	returnSuccess(ctx, peer)
	return
}
