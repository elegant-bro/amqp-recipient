package handlers

import (
	rcp "github.com/elegant-bro/amqp-recipient"
	amqp "github.com/rabbitmq/amqp091-go"
)

func NewSkipAck(origin rcp.JobHandler, skip func(d amqp.Delivery) bool) rcp.JobHandler {
	return NewSkip(origin, skip, rcp.HandlerAck)
}

func NewSkip(origin rcp.JobHandler, skip func(d amqp.Delivery) bool, res uint8) rcp.JobHandler {
	return &WrapHandler{
		origin: origin,
		wrapper: func(d amqp.Delivery, wrapped rcp.JobHandler) (uint8, error) {
			if skip(d) {
				return res, nil
			}

			return wrapped.Handle(d)
		},
	}
}
