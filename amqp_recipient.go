package amqp_recipient

import "github.com/streadway/amqp"

type AmqpRecipient struct {
	Queue      string
	Prefetch   int
	Connection *amqp.Connection
	Handler    JobHandler
}

func (recipient *AmqpRecipient) Subscribe() (Job, error) {
	ch, err := recipient.Connection.Channel()
	if nil != err {
		return nil, err
	}

	err = ch.Qos(recipient.Prefetch, 0, false)
	if nil != err {
		return nil, err
	}

	deliveries, err := ch.Consume(
		recipient.Queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if nil != err {
		return nil, err
	}

	return &AmqpJob{
		Deliveries: deliveries,
		Handler:    recipient.Handler,
	}, nil
}
