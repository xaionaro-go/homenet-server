package api

import (
	"fmt"

	"github.com/xaionaro-go/errors"
)

type unregisterPeerAnswer struct {
	answerCommon
}

func (api *api) UnregisterPeer(networkID, peerID string) (int, error) {
	if len(networkID) == 0 {
		return 0, networkIDCannotBeEmpty.Wrap()
	}

	if len(peerID) == 0 {
		return 0, peerIDCannotBeEmpty.Wrap()
	}

	var answer unregisterPeerAnswer
	statusCode, err := api.DELETE(&answer, fmt.Sprintf("%s/peers/%s", networkID, peerID))
	return statusCode, errors.Wrap(err)
}
