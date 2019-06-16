package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/xaionaro-go/errors"

	"github.com/xaionaro-go/homenet-server/helpers"
)

var (
	networkIDCannotBeEmpty = errors.InvalidArguments.New("networkID cannot be empty")
	peerIDCannotBeEmpty    = errors.InvalidArguments.New("peerID cannot be empty")
)

type loggers struct {
	debug logger
}

type api struct {
	urlRoot      string
	passwordHash []byte

	logger loggers
}

type answer interface {
	GetStatus() string
	GetErrorDescription() string
}

type answerCommon struct {
	Status           string
	ErrorDescription string
}

type logger interface {
	Printf(format string, v ...interface{})
	Print(...interface{})
}

func (answer *answerCommon) GetStatus() string {
	return answer.Status
}
func (answer *answerCommon) GetErrorDescription() string {
	return answer.ErrorDescription
}

func New(urlRoot, passwordHash string, options ...Option) *api {
	result := &api{
		urlRoot:      urlRoot,
		passwordHash: helpers.Hash([]byte(passwordHash)),
	}

	for _, optI := range options {
		switch opt := optI.(type) {
		case *optSetLoggerDebug:
			result.logger.debug = opt.GetLogger()
		default:
			panic(fmt.Errorf("Unknown option: %T", opt))
		}
	}

	return result
}

func (api *api) ifDebug(fn func(logger)) {
	if api.logger.debug == nil {
		return
	}
	fn(api.logger.debug)
}

func (api *api) query(result answer, method, uri string, bodyReader io.Reader, options ...map[string]interface{}) (resultStatusCode int, resultErr error) {
	var body []byte
	defer func() { resultErr = errors.Wrap(resultErr, `JSON-message:"`+string(body)+`"`) }()

	v := url.Values{}
	if len(options) >= 1 {
		for paramName, paramValue := range options[0] {
			v.Add(paramName, fmt.Sprintf("%v", paramValue))
		}
	}
	req, err := http.NewRequest(method, api.urlRoot+uri+"?"+v.Encode(), bodyReader)
	if err != nil {
		return 0, err
	}
	req.Header.Set("X-Homenet-Accesshash", base64.StdEncoding.EncodeToString(api.passwordHash))
	client := &http.Client{}
	api.ifDebug(func(log logger) {
		if dump, err := httputil.DumpRequestOut(req, true); err == nil {
			log.Printf("API-request: %v", string(dump))
		}
	})
	resp, err := client.Do(req)
	api.ifDebug(func(log logger) {
		if resp == nil {
			return
		}
		if dump, err := httputil.DumpResponse(resp, true); err == nil {
			log.Printf("API-response: %v", string(dump))
		}
	})
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return resp.StatusCode, err
	}
	return resp.StatusCode, api.wrapResultError(result.GetErrorDescription())
}

func (api *api) GET(result answer, uri string, options ...map[string]interface{}) (int, error) {
	return api.query(result, `GET`, uri, nil, options...)
}

func (api *api) PUT(result answer, uri string, body io.Reader, options ...map[string]interface{}) (int, error) {
	return api.query(result, `PUT`, uri, body, options...)
}

func (api *api) DELETE(result answer, uri string, options ...map[string]interface{}) (int, error) {
	return api.query(result, `DELETE`, uri, nil, options...)
}

func (api *api) wrapResultError(errorDescription string, args ...interface{}) error {
	if errorDescription == "" {
		return nil
	}
	return errors.New(errorDescription, args...)
}
func (api *api) messageToReader(msg interface{}) io.Reader {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return bytes.NewBuffer(b)
}
