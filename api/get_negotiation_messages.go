package api

import (
	"fmt"

	"github.com/xaionaro-go/errors"

	"github.com/xaionaro-go/homenet-server/models"
)

type getNegotiationMessagesAnswer struct {
	answerCommon
	Result map[string]models.NegotiationMessage
}

func (api *api) GetNegotiationMessages(networkID, peerIDTo string) (int, map[string]models.NegotiationMessage, error) {
	if len(networkID) == 0 {
		return 0, nil, networkIDCannotBeEmpty.Wrap()
	}

	var answer getNegotiationMessagesAnswer
	statusCode, err := api.GET(&answer, fmt.Sprintf("%s/negotiationMessage/%s", networkID, peerIDTo))
	return statusCode, answer.Result, errors.Wrap(err)
}
