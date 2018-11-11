package api

import (
	"fmt"

	"github.com/xaionaro-go/homenet-server/models"
)

type getNetAnswer struct {
	answerCommon
	Result models.NetworkT
}

func (api *api) GetNet(networkID string) (int, models.NetworkT, error) {
	var answer getNetAnswer
	statusCode, err := api.GET(&answer, fmt.Sprintf("%s", networkID))
	return statusCode, answer.Result, err
}
