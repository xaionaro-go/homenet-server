package api

import (
	"fmt"

	"github.com/xaionaro-go/errors"

	"github.com/xaionaro-go/homenet-server/models"
)

type registerNetAnswer struct {
	answerCommon
	Result models.NetworkT
}

func (api *api) RegisterNet(networkID string) (int, *models.NetworkT, error) {
	if len(networkID) == 0 {
		return 0, nil, networkIDCannotBeEmpty.Wrap()
	}

	var answer registerNetAnswer
	statusCode, err := api.PUT(&answer, fmt.Sprintf("%s", networkID))
	return statusCode, &answer.Result, errors.Wrap(err)
}
