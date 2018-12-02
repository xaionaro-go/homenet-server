package api

import (
	"fmt"

	"github.com/xaionaro-go/errors"

	"github.com/xaionaro-go/homenet-server/models"
)

type getNegotiationMessageAnswer struct {
	answerCommon
	Result models.NegotiationMessageT
}

func (api *api) GetNegotiationMessage(networkID, peerIDTo, peerIDFrom string) (int, models.NegotiationMessageT, error) {
	if len(networkID) == 0 {
		return 0, models.NegotiationMessageT{}, networkIDCannotBeEmpty.Wrap()
	}

	var answer getNegotiationMessageAnswer
	statusCode, err := api.GET(&answer, fmt.Sprintf("%s/negotiationMessage/%s/%s", networkID, peerIDTo, peerIDFrom))
	return statusCode, answer.Result, errors.Wrap(err)
}
