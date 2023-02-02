package handlers

import (
	rcp "github.com/elegant-bro/amqp-recipient"
	amqp "github.com/rabbitmq/amqp091-go"
)

type FuncHandler struct {
	fn func(d amqp.Delivery) (uint8, error)
}

func NewFunc(fn func(d amqp.Delivery) (uint8, error)) rcp.JobHandler {
	return &FuncHandler{fn: fn}
}

func (f FuncHandler) Handle(d amqp.Delivery) (uint8, error) {
	return f.fn(d)
}
