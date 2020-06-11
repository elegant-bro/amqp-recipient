package recipients

import (
	"context"
	"github.com/elegant-bro/amqp-recipient"
	"github.com/elegant-bro/amqp-recipient/jobs"
	"github.com/streadway/amqp"
	"log"
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
	dd := NewAmqpDeliveries(recipient.conn, recipient.queue, recipient.prefetch, recipient.consumeOptions)
	deliveries, err := dd.All()
	if nil != err {
		return nil, err
	}

	go func() {
		<-recipient.ctx.Done()
		if err := dd.Close(); nil != err {
			log.Print(err)
		}

	}()

	return jobs.NewCancelableAmqpJob(recipient.ctx, deliveries, recipient.handler, recipient.onFail), nil
}
