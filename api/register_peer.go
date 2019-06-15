package api

import (
	"fmt"

	"github.com/xaionaro-go/errors"

	"github.com/xaionaro-go/homenet-server/models"
)

type registerPeerAnswer struct {
	answerCommon
	Result models.PeerT
}

func (api *api) RegisterPeer(networkID, peerID, peerName string) (int, *models.PeerT, error) {
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
	statusCode, err := api.PUT(&answer, fmt.Sprintf("%s/peers/%s", networkID, peerID), nil, params)
	return statusCode, &answer.Result, errors.Wrap(err)
}
