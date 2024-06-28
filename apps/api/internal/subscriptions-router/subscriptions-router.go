package subscriptions_router

type Router interface {
	Subscribe(keys []string) (Subscription, error)
	Publish(key string, data any) error
}

type Subscription interface {
	GetChannel() chan []byte
	Unsubscribe() error
}
