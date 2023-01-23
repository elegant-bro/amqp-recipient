package jobs

import (
	"context"
	amqpRecipient "github.com/elegant-bro/amqp-recipient"
	amqp "github.com/rabbitmq/amqp091-go"
)

type CancelableAmqpJob struct {
	deliveries <-chan amqp.Delivery
	handler    amqpRecipient.JobHandler
	onFail     amqpRecipient.OnHandlerFails
	onCancel   func()
	ctx        context.Context
}

func NewCancelableAmqpJob(ctx context.Context, deliveries <-chan amqp.Delivery, handler amqpRecipient.JobHandler, onFail amqpRecipient.OnHandlerFails) *CancelableAmqpJob {
	return &CancelableAmqpJob{deliveries: deliveries, handler: handler, onFail: onFail, ctx: ctx}
}

func (job *CancelableAmqpJob) Run() {
	for {
		select {
		case d, ok := <-job.deliveries:
			if ok {
				handleDelivery(d, job.handler, job.onFail)
			}

		case <-job.ctx.Done():
			if nil != job.onCancel {
				job.onCancel()
			}
			return
		}
	}
}

func (job *CancelableAmqpJob) OnCancel(fn func()) {
	job.onCancel = fn
}
