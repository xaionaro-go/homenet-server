package api

import (
	"fmt"

	"github.com/xaionaro-go/errors"

	"github.com/xaionaro-go/homenet-server/models"
)

type setNegotiationMessageAnswer struct {
	answerCommon
	Result models.NegotiationMessage
}

func (api *api) SetNegotiationMessage(networkID, peerIDTo, peerIDFrom string, msg *models.NegotiationMessage) (int, *models.NegotiationMessage, error) {
	if len(networkID) == 0 {
		return 0, nil, networkIDCannotBeEmpty.Wrap()
	}

	var answer setNegotiationMessageAnswer
	statusCode, err := api.PUT(&answer, fmt.Sprintf("%s/negotiationMessage/%s/%s", networkID, peerIDTo, peerIDFrom), api.messageToReader(msg))
	return statusCode, &answer.Result, errors.Wrap(err)
}
