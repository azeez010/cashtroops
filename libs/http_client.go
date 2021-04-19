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
	inner  *http.Client
	logger *logrus.Logger
}

func NewHttpClient(logger *logrus.Logger) *HttpClient {
	return &HttpClient{logger: logger, inner: &http.Client{Timeout: 60 * time.Second}}
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

	res, err := o.inner.Do(request)
	if err != nil {
		return err
	}
	defer func(res *http.Response) {
		err := res.Body.Close()
		o.logger.WithError(err).Error("error closing response body")
	}(res)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode < 101 || res.StatusCode > 299 {
		o.logger.WithError(errors.New(string(body))).Error("BC API error")
		return errors.New("failed to complete api call")
	}
	if resp != nil {
		if err := json.Unmarshal(body, resp); err != nil {
			return err
		}
	}
	return nil
}
