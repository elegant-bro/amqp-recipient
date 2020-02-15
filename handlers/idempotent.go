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

	has, err := h.ids.Has(messageId)
	if nil != err {
		return recipient.HandlerReject, newInternalError(err, "can't check message id has been handled")
	}

	if has {
		return recipient.HandlerAck, nil
	}

	return h.ids.Save(messageId, func() (uint8, error) {
		return h.origin.Handle(d)
	})
}
