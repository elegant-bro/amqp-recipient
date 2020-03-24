package jobs

import (
	amqpRecipient "github.com/elegant-bro/amqp-recipient"
	"github.com/streadway/amqp"
)

func handleDelivery(d amqp.Delivery, handler amqpRecipient.JobHandler, onFail amqpRecipient.OnHandlerFails) {
	result, err := handler.Handle(d)

	if amqpRecipient.HandlerAck == result {
		_ = d.Ack(false)
	} else if amqpRecipient.HandlerRequeue == result {
		_ = d.Reject(true)
	} else if amqpRecipient.HandlerReject == result {
		_ = d.Reject(false)
	}

	if nil != err {
		onFail(d, err)
	}
}
