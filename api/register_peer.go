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
	statusCode, err := api.PUT(&answer, fmt.Sprintf("%s/peers/%s", networkID, peerID))
	return statusCode, answer.Result, err
}
