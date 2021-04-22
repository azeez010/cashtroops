package libs

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpClient struct {
	inner       *http.Client
	bearerToken string
	logger      *logrus.Logger
	debug       bool
}

func NewHttpClient(logger *logrus.Logger, bearerToken string) *HttpClient {
	return &HttpClient{logger: logger, bearerToken: bearerToken, inner: &http.Client{Timeout: 60 * time.Second}}
}

func (o *HttpClient) Do(url, method string, reqBody, resp interface{}) error {
	b := &bytes.Buffer{}
	if reqBody != nil {
		if err := json.NewEncoder(b).Encode(reqBody); err != nil {
			return err
		}
	}

	request, err := http.NewRequest(method, url, b)
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", "application/json")
	if o.bearerToken != "" {
		request.Header.Add("Authorization", "Bearer "+o.bearerToken)
	}

	res, err := o.inner.Do(request)
	if err != nil {
		return err
	}
	defer func(res *http.Response) {
		err := res.Body.Close()
		if err != nil {
			o.logger.WithError(err).Error("error closing response body")
		}
	}(res)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if o.debug {
		o.logger.WithField("api_call_response", string(body)).Info("response body")
	}
	if res.StatusCode < 101 || res.StatusCode > 299 {
		o.logger.WithError(errors.New(string(body))).Error("API call error")
		return errors.New("failed to complete api call")
	}
	if resp != nil {
		if err := json.Unmarshal(body, resp); err != nil {
			return err
		}
	}
	return nil
}
