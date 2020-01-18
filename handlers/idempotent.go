package handlers

import (
	recipient "github.com/elegant-bro/amqp-recipient"
	"github.com/streadway/amqp"
)

type IdempotentHandler struct {
	origin recipient.JobHandler
	ids    recipient.HandledIds
}

func NewIdempotent(origin recipient.JobHandler, ids recipient.HandledIds) *IdempotentHandler {
	return &IdempotentHandler{origin: origin, ids: ids}
}

func (h *IdempotentHandler) Handle(d *amqp.Delivery) (uint8, error) {
	messageId := d.MessageId
	if len(messageId) == 0 {
		return h.origin.Handle(d)
	}

	if has, err := h.ids.Has(messageId); nil != err || has {
		return 0, err
	}

	return h.ids.Save(messageId, func() (uint8, error) {
		return h.origin.Handle(d)
	})
}
