package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type api struct {
	urlRoot      string
	passwordHash string
}

type answer interface {
	GetStatus() string
	GetErrorDescription() string
}

type answerCommon struct {
	Status           string
	ErrorDescription string
}

func (answer *answerCommon) GetStatus() string {
	return answer.Status
}
func (answer *answerCommon) GetErrorDescription() string {
	return answer.ErrorDescription
}

func New(urlRoot, passwordHash string) *api {
	return &api{
		urlRoot:      urlRoot,
		passwordHash: passwordHash,
	}
}

func (api *api) query(result answer, method, uri string) (int, error) {
	req, err := http.NewRequest(method, api.urlRoot+uri, nil)
	if err != nil {
		return 0, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return resp.StatusCode, err
	}
	return resp.StatusCode, api.wrapResultError(result.GetErrorDescription())
}

func (api *api) GET(result answer, uri string) (int, error) {
	return api.query(result, `GET`, uri)
}

func (api *api) PUT(result answer, uri string) (int, error) {
	return api.query(result, `PUT`, uri)
}

func (api *api) DELETE(result answer, uri string) (int, error) {
	return api.query(result, `DELETE`, uri)
}

func (api *api) wrapResultError(errorDescription string) error {
	if errorDescription == "" {
		return nil
	}
	return errors.New(errorDescription)
}
