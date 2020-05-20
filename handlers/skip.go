package handlers

import (
	recipient "github.com/elegant-bro/amqp-recipient"
	"github.com/streadway/amqp"
)

type SkipHandler struct {
	origin recipient.JobHandler
	skip   func(d amqp.Delivery) bool
	res    uint8
}

func SkipAckHandler(origin recipient.JobHandler, skip func(d amqp.Delivery) bool) *SkipHandler {
	return NewSkipHandler(origin, skip, recipient.HandlerAck)
}

func NewSkipHandler(origin recipient.JobHandler, skip func(d amqp.Delivery) bool, res uint8) *SkipHandler {
	return &SkipHandler{origin: origin, skip: skip, res: res}
}

func (s SkipHandler) Handle(d amqp.Delivery) (uint8, error) {
	if s.skip(d) {
		return s.res, nil
	}

	return s.origin.Handle(d)
}
