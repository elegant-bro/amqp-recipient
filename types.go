package amqp_recipient

import "github.com/streadway/amqp"

type Recipient interface {
	Subscribe() (Job, error)
}

type Job interface {
	Run()
}

type JobHandler interface {
	Handle(d *amqp.Delivery) (uint8, error)
}

type HandledIds interface {
	Has(key string) (bool, error)
	Save(key string, fn func() (uint8, error)) (uint8, error)
}