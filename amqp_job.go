package amqp_recipient

import (
	"github.com/streadway/amqp"
)

type AmqpJob struct {
	deliveries <-chan amqp.Delivery
	handler    JobHandler
	onFail     OnHandlerFails
}

func NewAmqpJob(deliveries <-chan amqp.Delivery, handler JobHandler, onFail OnHandlerFails) *AmqpJob {
	return &AmqpJob{deliveries: deliveries, handler: handler, onFail: onFail}
}

func (job *AmqpJob) Run() {
	for d := range job.deliveries {
		result, err := job.handler.Handle(&d)

		if HandlerAck == result {
			_ = d.Ack(false)
		} else if HandlerRequeue == result {
			_ = d.Reject(true)
		} else if result >= HandlerReject {
			_ = d.Reject(false)
		}

		if nil != err {
			job.onFail(&d, err)
		}
	}
}
