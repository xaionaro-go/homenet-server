package api

import (
	"bytes"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/ed25519"

	"github.com/xaionaro-go/errors"

	"github.com/xaionaro-go/homenet-server/models"
)

type registerPeerAnswer struct {
	answerCommon
	Result models.PeerT
}

func (api *api) RegisterPeer(networkID, peerID, peerName string, publicKey ed25519.PublicKey) (int, *models.PeerT, error) {
	if len(networkID) == 0 {
		return 0, nil, networkIDCannotBeEmpty.Wrap()
	}

	if len(peerID) == 0 {
		return 0, nil, peerIDCannotBeEmpty.Wrap()
	}

	var answer registerPeerAnswer
	params := map[string]interface{}{}
	if peerName != "" {
		params["peer_name"] = peerName
	}

	var publicKeyBuf bytes.Buffer
	_, err := base64.NewEncoder(base64.URLEncoding, &publicKeyBuf).Write(publicKey[:])
	if err != nil {
		return 0, nil, invalidPublicKey.Wrap(err)
	}
	params["public_key"] = publicKeyBuf.String()

	statusCode, err := api.PUT(&answer, fmt.Sprintf("%s/peers/%s", networkID, peerID), nil, params)
	return statusCode, &answer.Result, errors.Wrap(err)
}
