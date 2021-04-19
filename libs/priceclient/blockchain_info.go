package priceclient

import (
	"encoding/json"
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
	return &client{inner: libs.NewHttpClient(logger)}
}

func (c *client) CurrentPrice() (*Result, error) {
	r := &Result{}
	u := "https://blockchain.info/ticker"
	err := c.inner.Do(u, "GET", nil, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}
