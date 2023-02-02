package handlers

import (
	rcp "github.com/elegant-bro/amqp-recipient"
	amqp "github.com/rabbitmq/amqp091-go"
)

func NewIdempotent(origin rcp.JobHandler, ids rcp.HandledIds) rcp.JobHandler {
	return &WrapHandler{
		origin: origin,
		wrapper: func(d amqp.Delivery, wrapped rcp.JobHandler) (uint8, error) {
			messageId := d.MessageId
			if len(messageId) == 0 {
				return wrapped.Handle(d)
			}

			has, err := ids.Has(messageId)
			if nil != err {
				return rcp.HandlerReject, newInternalError(err, "can't check message id has been handled")
			}

			if has {
				return rcp.HandlerAck, nil
			}

			return ids.Save(messageId, func() (uint8, error) {
				return wrapped.Handle(d)
			})
		},
	}
}
