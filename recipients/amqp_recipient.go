package recipients

import (
	"fmt"
	"github.com/elegant-bro/amqp-recipient"
	"github.com/elegant-bro/amqp-recipient/jobs"
	amqp "github.com/rabbitmq/amqp091-go"
	"math/rand"
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
	ch, err := recipient.conn.Channel()
	if nil != err {
		return nil, err
	}

	err = ch.Qos(recipient.prefetch, 0, false)
	if nil != err {
		return nil, err
	}

	tag := fmt.Sprintf("%d", rand.Int())
	deliveries, err := ch.Consume(
		recipient.queue,
		tag,
		recipient.consumeOptions.AutoAck,
		recipient.consumeOptions.Exclusive,
		recipient.consumeOptions.NoLocal,
		recipient.consumeOptions.NoWait,
		recipient.consumeOptions.Args,
	)
	if nil != err {
		return nil, err
	}

	return jobs.NewAmqpJob(deliveries, recipient.handler, recipient.onFail), nil
}
