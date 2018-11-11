package api

import (
	"fmt"

	"github.com/xaionaro-go/homenet-server/models"
)

type getPeersAnswer struct {
	answerCommon
	Result models.Peers
}

func (api *api) GetPeers(networkID string) (int, models.Peers, error) {
	var answer getPeersAnswer
	statusCode, err := api.GET(&answer, fmt.Sprintf("%s/peers", networkID))
	return statusCode, answer.Result, err
}
