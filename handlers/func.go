package handlers

import "github.com/streadway/amqp"

type FuncHandler struct {
	fn func(d amqp.Delivery) (uint8, error)
}

func NewFunc(fn func(d amqp.Delivery) (uint8, error)) *FuncHandler {
	return &FuncHandler{fn: fn}
}

func (f FuncHandler) Handle(d amqp.Delivery) (uint8, error) {
	return f.fn(d)
}
