package api

import (
	"fmt"

	"github.com/xaionaro-go/homenet-server/models"
)

type registerPeerAnswer struct {
	answerCommon
	Result models.PeerT
}

func (api *api) RegisterPeer(networkID, peerID, peerName string) (int, models.PeerT, error) {
	var answer registerPeerAnswer
	params := map[string]interface{}{}
	if peerName != "" {
		params["peer_name"] = peerName
	}
	statusCode, err := api.PUT(&answer, fmt.Sprintf("%s/peers/%s", networkID, peerID), params)
	return statusCode, answer.Result, err
}
