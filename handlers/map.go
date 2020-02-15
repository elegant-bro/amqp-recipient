package handlers

import (
	"fmt"
	amqpRecipient "github.com/elegant-bro/amqp-recipient"
	"github.com/streadway/amqp"
)

type MapHandler struct {
	handlersMap   map[string]amqpRecipient.JobHandler
	keyFn         func(delivery *amqp.Delivery) (string, error)
	onErrorResult uint8
}

func NewDefaultMap(handlersMap map[string]amqpRecipient.JobHandler, keyFn func(delivery *amqp.Delivery) (string, error)) *MapHandler {
	return NewMap(handlersMap, keyFn, amqpRecipient.HandlerAck)
}

func NewMap(handlersMap map[string]amqpRecipient.JobHandler, keyFn func(delivery *amqp.Delivery) (string, error), onErrorResult uint8) *MapHandler {
	return &MapHandler{handlersMap: handlersMap, keyFn: keyFn, onErrorResult: onErrorResult}
}

func (m *MapHandler) Handle(d *amqp.Delivery) (uint8, error) {
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
