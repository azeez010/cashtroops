package bc

import (
	"fmt"
	"github.com/adigunhammedolalekan/cashtroops/libs"
	"github.com/sirupsen/logrus"
)

const (
	bcAddress               = "https://api.blockcypher.com/v1/%s/%s"
	EventTypeTxConfirmation = "confirmed-tx"
)

type Address struct {
	Private string `json:"private"`
	Public  string `json:"public"`
	Address string `json:"address"`
	WIF     string `json:"wif"`
}

type Event struct {
	// {"event": "unconfirmed-tx", "address": "15qx9ug952GWGTNn7Uiv6vode4RcGrRemh", "url": "https://my.domain.com/callbacks/new-tx"}
	ID      string `json:"id"`
	Event   string `json:"event"`
	URL     string `json:"url"`
	Address string `json:"address"`
}

type Client interface {
	GenerateAddress() (*Address, error)
	SetupWebHooks(hooks []Event) error
	ListHooks() ([]Event, error)
	AddHook(event Event) (*Event, error)
}

type client struct {
	token, coin, network string
	httpClient           *libs.HttpClient
	logger               *logrus.Logger
}

func New(token, coin, network string, logger *logrus.Logger) (Client, error) {
	return &client{
		token:      token,
		coin:       coin,
		network:    network,
		httpClient: libs.NewHttpClient(logger, ""),
		logger:     logger,
	}, nil
}

func (o *client) GenerateAddress() (*Address, error) {
	addrUrl := fmt.Sprintf(bcAddress+"/addrs?token=%s", o.coin, o.network, o.token)
	address := &Address{}
	err := o.httpClient.Do(addrUrl, "POST", nil, address)
	if err != nil {
		return nil, err
	}
	return address, nil
}

func (o *client) SetupWebHooks(hooks []Event) error {
	events, err := o.ListHooks()
	if err != nil {
		return err
	}

	if len(events) > 0 {
		if err := o.deleteHooks(events); err != nil {
			return err
		}
	}
	createHooksUrl := fmt.Sprintf(bcAddress+"/hooks?token=%s", o.coin, o.network, o.token)
	for _, nextEvent := range hooks {
		if err := o.httpClient.Do(createHooksUrl, "POST", &nextEvent, &Event{}); err != nil {
			return err
		}
	}
	return nil
}

func (o *client) ListHooks() ([]Event, error) {
	events := make([]Event, 0)
	webhooksUrl := fmt.Sprintf(bcAddress+"/hooks?token=%s", o.coin, o.network, o.token)
	err := o.httpClient.Do(webhooksUrl, "GET", nil, &events)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (o *client) deleteHooks(hooks []Event) error {
	for _, next := range hooks {
		deleteHookUrl := fmt.Sprintf(bcAddress+"/hooks/%s?token=%s", o.coin, o.network, next.ID, o.token)
		if err := o.httpClient.Do(deleteHookUrl, "DELETE", nil, nil); err != nil {
			return err
		}
	}
	return nil
}

func (o *client) AddHook(event Event) (*Event, error) {
	createHooksUrl := fmt.Sprintf(bcAddress+"/hooks?token=%s", o.coin, o.network, o.token)
	e := &Event{}
	err := o.httpClient.Do(createHooksUrl, "POST", event, e)
	if err != nil {
		return nil, err
	}
	return e, nil
}
