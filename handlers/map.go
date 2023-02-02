package handlers

import (
	"fmt"
	rcp "github.com/elegant-bro/amqp-recipient"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MapHandler struct {
	handlersMap   rcp.MapOfHandlers
	keyFn         func(delivery amqp.Delivery) (string, error)
	onErrorResult uint8
}

func NewDefaultMap(handlersMap rcp.MapOfHandlers, keyFn func(delivery amqp.Delivery) (string, error)) *MapHandler {
	return NewMap(handlersMap, keyFn, rcp.HandlerAck)
}

func NewMap(handlersMap rcp.MapOfHandlers, keyFn func(delivery amqp.Delivery) (string, error), onErrorResult uint8) *MapHandler {
	return &MapHandler{handlersMap: handlersMap, keyFn: keyFn, onErrorResult: onErrorResult}
}

func (m *MapHandler) Handle(d amqp.Delivery) (uint8, error) {
	key, err := m.keyFn(d)
	if nil != err {
		return m.onErrorResult, newInternalError(err, "fail to compute handler key")
	}

	handler, ok := m.handlersMap[key]
	if ok {
		return handler.Handle(d)
	}

	return m.onErrorResult, newInternalError(fmt.Errorf("key '%s'", key), "handler not found")
}
