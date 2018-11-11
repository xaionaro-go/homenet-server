package api

import (
	"fmt"

	"github.com/xaionaro-go/homenet-server/models"
)

type registerNetAnswer struct {
	answerCommon
	Result models.NetworkT
}

func (api *api) RegisterNet(networkID string) (int, models.NetworkT, error) {
	var answer registerNetAnswer
	statusCode, err := api.PUT(&answer, fmt.Sprintf("%s", networkID))
	return statusCode, answer.Result, err
}
