package handlers

import (
	"fmt"
	amqpRecipient "github.com/elegant-bro/amqp-recipient"
	"github.com/streadway/amqp"
)

type MapHandler struct {
	handlersMap map[string]amqpRecipient.JobHandler
	keyFn       func(delivery *amqp.Delivery) (string, error)
}

func NewMapHandler(handlersMap map[string]amqpRecipient.JobHandler, keyFn func(delivery *amqp.Delivery) (string, error)) *MapHandler {
	return &MapHandler{handlersMap: handlersMap, keyFn: keyFn}
}

func (m *MapHandler) Handle(d *amqp.Delivery) (uint8, error) {
	key, err := m.keyFn(d)
	if nil != err {
		return 1, fmt.Errorf("fail to compute handler key: %v", err)
	}

	handler, ok := m.handlersMap[key]
	if ok {
		return handler.Handle(d)
	}

	return 1, fmt.Errorf("handler not found for key '%s'", key)
}
