package jobs

import (
	"github.com/elegant-bro/amqp-recipient"
	"github.com/streadway/amqp"
)

type AmqpJob struct {
	deliveries <-chan amqp.Delivery
	handler    amqp_recipient.JobHandler
	onFail     amqp_recipient.OnHandlerFails
}

func NewAmqpJob(deliveries <-chan amqp.Delivery, handler amqp_recipient.JobHandler, onFail amqp_recipient.OnHandlerFails) *AmqpJob {
	return &AmqpJob{deliveries: deliveries, handler: handler, onFail: onFail}
}

func (job *AmqpJob) Run() {
	for d := range job.deliveries {
		handleDelivery(d, job.handler, job.onFail)
	}
}
