package recipients

import (
	"context"
	"github.com/elegant-bro/amqp-recipient"
	"github.com/elegant-bro/amqp-recipient/jobs"
	"github.com/streadway/amqp"
)

type CancelableAmqpRecipient struct {
	queue          string
	prefetch       int
	conn           *amqp.Connection
	handler        amqp_recipient.JobHandler
	onFail         amqp_recipient.OnHandlerFails
	consumeOptions ConsumeOptions
	ctx            context.Context
}

func NewDefaultCancelableAmqpRecipient(
	ctx context.Context,
	queue string,
	prefetch int,
	connection *amqp.Connection,
	handler amqp_recipient.JobHandler,
	onFail amqp_recipient.OnHandlerFails,
) *CancelableAmqpRecipient {
	return NewCancelableAmqpRecipient(ctx, queue, prefetch, connection, handler, onFail, ConsumeOptions{})
}

func NewCancelableAmqpRecipient(
	ctx context.Context,
	queue string,
	prefetch int,
	connection *amqp.Connection,
	handler amqp_recipient.JobHandler,
	onFail amqp_recipient.OnHandlerFails,
	opt ConsumeOptions,
) *CancelableAmqpRecipient {
	return &CancelableAmqpRecipient{
		queue:          queue,
		prefetch:       prefetch,
		conn:           connection,
		handler:        handler,
		onFail:         onFail,
		consumeOptions: opt,
		ctx:            ctx,
	}
}

func (recipient *CancelableAmqpRecipient) Subscribe() (amqp_recipient.Job, error) {
	ch, err := recipient.conn.Channel()
	if nil != err {
		return nil, err
	}

	err = ch.Qos(recipient.prefetch, 0, false)
	if nil != err {
		return nil, err
	}

	deliveries, err := ch.Consume(
		recipient.queue,
		recipient.consumeOptions.Consumer,
		recipient.consumeOptions.AutoAck,
		recipient.consumeOptions.Exclusive,
		recipient.consumeOptions.NoLocal,
		recipient.consumeOptions.NoWait,
		recipient.consumeOptions.Args,
	)
	if nil != err {
		return nil, err
	}

	return jobs.NewCancelableAmqpJob(recipient.ctx, deliveries, recipient.handler, recipient.onFail), nil
}
