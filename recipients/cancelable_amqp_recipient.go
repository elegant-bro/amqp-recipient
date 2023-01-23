package recipients

import (
	"context"
	"fmt"
	amqp_recipient "github.com/elegant-bro/amqp-recipient"
	"github.com/elegant-bro/amqp-recipient/jobs"
	amqp "github.com/rabbitmq/amqp091-go"
	"math/rand"
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
		return nil, fmt.Errorf("recipient for %s queue conn.Channel fails: %w", recipient.queue, err)
	}

	err = ch.Qos(recipient.prefetch, 0, false)
	if nil != err {
		return nil, fmt.Errorf("recipient for %s queue channel.Qos fails: %w", recipient.queue, err)
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
		return nil, fmt.Errorf("consuming %s queue, tag %s fails: %w", recipient.queue, tag, err)
	}

	go func() {
		<-recipient.ctx.Done()
		if !recipient.conn.IsClosed() {
			if err := ch.Cancel(tag, false); nil != err {
				fmt.Println("CancelableAmqpRecipient channel canceling failed: ", err)
			}
		}
	}()

	return jobs.NewCancelableAmqpJob(recipient.ctx, deliveries, recipient.handler, recipient.onFail), nil
}
