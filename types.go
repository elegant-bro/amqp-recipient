package amqp_recipient

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Recipient interface {
	Subscribe() (Job, error)
}

type Sender interface {
	Send(p amqp.Publishing) error
}

type Job interface {
	Run()
}

type JobHandler interface {
	Handle(d amqp.Delivery) (uint8, error)
}

type HandledIds interface {
	Has(key string) (bool, error)
	Save(key string, fn func() (uint8, error)) (uint8, error)
}

type OnHandlerFails func(d amqp.Delivery, err error)

type MapOfHandlers map[string]JobHandler

type Wrapper func(d amqp.Delivery, wrapped JobHandler) (uint8, error)
