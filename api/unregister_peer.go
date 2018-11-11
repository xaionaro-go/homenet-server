package api

import (
	"fmt"
)

type unregisterPeerAnswer struct {
	answerCommon
}

func (api *api) UnregisterPeer(networkID, peerID string) (int, error) {
	var answer unregisterPeerAnswer
	statusCode, err := api.DELETE(&answer, fmt.Sprintf("%s/peers/%s", networkID, peerID))
	return statusCode, err
}
