package subscriptions_router

import (
	"github.com/goccy/go-json"
	"github.com/nats-io/nats.go"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Nats *nats.Conn
}

func NewNatsSubscription(opts Opts) (*RouterNats, error) {
	return &RouterNats{
		nc: opts.Nats,
	}, nil
}

type RouterNats struct {
	nc *nats.Conn
}

var _ Router = &RouterNats{}

type WsRouterNatsSubscription struct {
	subs      []*nats.Subscription
	dataChann chan []byte
}

func (c *WsRouterNatsSubscription) Unsubscribe() error {
	for _, sub := range c.subs {
		if err := sub.Unsubscribe(); err != nil {
			return err
		}
	}

	return nil
}

func (c *WsRouterNatsSubscription) GetChannel() chan []byte {
	return c.dataChann
}

func (c *RouterNats) Subscribe(keys []string) (Subscription, error) {
	ch := make(chan []byte, 1)
	subs := make([]*nats.Subscription, 0, len(keys))

	for _, key := range keys {
		sub, err := c.nc.Subscribe(
			key,
			func(msg *nats.Msg) {
				ch <- msg.Data
			},
		)

		if err != nil {
			return nil, err
		}

		subs = append(subs, sub)
	}

	return &WsRouterNatsSubscription{
		subs:      subs,
		dataChann: ch,
	}, nil
}

func (c *RouterNats) Publish(key string, data any) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return c.nc.Publish(key, dataBytes)
}
