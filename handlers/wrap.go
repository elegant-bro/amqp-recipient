package handlers

import (
	rcp "github.com/elegant-bro/amqp-recipient"
	"github.com/rabbitmq/amqp091-go"
)

type WrapHandler struct {
	origin  rcp.JobHandler
	wrapper rcp.Wrapper
}

func NewWrapHandler(origin rcp.JobHandler, wrapper rcp.Wrapper) rcp.JobHandler {
	return &WrapHandler{origin: origin, wrapper: wrapper}
}

func (w *WrapHandler) Handle(d amqp091.Delivery) (uint8, error) {
	return w.wrapper(d, w.origin)
}
