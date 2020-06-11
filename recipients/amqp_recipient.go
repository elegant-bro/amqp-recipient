package recipients

import (
	"github.com/elegant-bro/amqp-recipient"
	"github.com/elegant-bro/amqp-recipient/jobs"
	"github.com/streadway/amqp"
)

type AmqpRecipient struct {
	queue          string
	prefetch       int
	conn           *amqp.Connection
	handler        amqp_recipient.JobHandler
	onFail         amqp_recipient.OnHandlerFails
	consumeOptions ConsumeOptions
}

func NewDefaultAmqpRecipient(
	queue string,
	prefetch int,
	connection *amqp.Connection,
	handler amqp_recipient.JobHandler,
	onFail amqp_recipient.OnHandlerFails,
) *AmqpRecipient {
	return NewAmqpRecipient(queue, prefetch, connection, handler, onFail, ConsumeOptions{})
}

func NewAmqpRecipient(
	queue string,
	prefetch int,
	connection *amqp.Connection,
	handler amqp_recipient.JobHandler,
	onFail amqp_recipient.OnHandlerFails,
	opt ConsumeOptions,
) *AmqpRecipient {
	return &AmqpRecipient{
		queue:          queue,
		prefetch:       prefetch,
		conn:           connection,
		handler:        handler,
		onFail:         onFail,
		consumeOptions: opt,
	}
}

func (recipient *AmqpRecipient) Subscribe() (amqp_recipient.Job, error) {
	dd := NewAmqpDeliveries(recipient.conn, recipient.queue, recipient.prefetch, recipient.consumeOptions)
	deliveries, err := dd.All()
	if nil != err {
		return nil, err
	}

	return jobs.NewAmqpJob(deliveries, recipient.handler, recipient.onFail), nil
}
