package amqp_recipient

import "github.com/streadway/amqp"

type AmqpJob struct {
	Deliveries <-chan amqp.Delivery
	Handler    JobHandler
}

func (job *AmqpJob) Run() {
	for d := range job.Deliveries {
		result, _ := job.Handler.Handle(&d)
		if 0 == result {
			_ = d.Ack(false)
		} else if 1 == result {
			_ = d.Reject(true)
		} else {
			_ = d.Reject(false)
		}
	}
}