package priceclient

import (
	"encoding/json"
	"errors"
	"github.com/adigunhammedolalekan/cashtroops/libs"
	"github.com/sirupsen/logrus"
)

type Result struct {
	USD struct {
		Last json.Number `json:"last"`
		Buy  json.Number `json:"buy"`
		Sell json.Number `json:"sell"`
	} `json:"usd"`
}

func (r *Result) FloatAmount() float64 {
	if value, err := r.USD.Last.Float64(); err == nil {
		return value
	}
	return 0
}

type Client interface {
	CurrentPrice() (*Result, error)
}

type client struct {
	inner *libs.HttpClient
}

func New(logger *logrus.Logger) Client {
	return &client{inner: libs.NewHttpClient(logger, "")}
}

func (c *client) CurrentPrice() (*Result, error) {
	// retry up to 5 times
	count := 0
	for {
		if count > 5 {
			break
		}
		r := &Result{}
		u := "https://blockchain.info/ticker"
		err := c.inner.Do(u, "GET", nil, r)
		if err != nil {
			count += 1
			continue
		}
		return r, nil
	}
	return nil, errors.New("failed to get BTC price")
}
