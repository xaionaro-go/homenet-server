package api

import (
	"fmt"

	"github.com/xaionaro-go/errors"

	"github.com/xaionaro-go/homenet-server/models"
)

type setNegotiationMessageAnswer struct {
	answerCommon
	Result models.NegotiationMessageT
}

func (api *api) SetNegotiationMessage(networkID, peerIDTo, peerIDFrom string, msg *models.NegotiationMessageT) (int, *models.NegotiationMessageT, error) {
	if len(networkID) == 0 {
		return 0, nil, networkIDCannotBeEmpty.Wrap()
	}

	var answer setNegotiationMessageAnswer
	statusCode, err := api.PUT(&answer, fmt.Sprintf("%s/negotiationMessage/%s/%s", networkID, peerIDTo, peerIDFrom), api.messageToReader(msg))
	return statusCode, &answer.Result, errors.Wrap(err)
}
