package amqp_recipient

import "github.com/streadway/amqp"

type AmqpRecipient struct {
	queue          string
	prefetch       int
	conn           *amqp.Connection
	handler        JobHandler
	onFail         OnHandlerFails
	consumeOptions ConsumeOptions
}

func NewDefaultAmqpRecipient(queue string, prefetch int, connection *amqp.Connection, handler JobHandler, onFail OnHandlerFails) *AmqpRecipient {
	return NewAmqpRecipient(queue, prefetch, connection, handler, onFail, ConsumeOptions{})
}

func NewAmqpRecipient(queue string, prefetch int, connection *amqp.Connection, handler JobHandler, onFail OnHandlerFails, opt ConsumeOptions) *AmqpRecipient {
	return &AmqpRecipient{
		queue:          queue,
		prefetch:       prefetch,
		conn:           connection,
		handler:        handler,
		onFail:         onFail,
		consumeOptions: opt,
	}
}

func (recipient *AmqpRecipient) Subscribe() (Job, error) {
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

	return NewAmqpJob(deliveries, recipient.handler, recipient.onFail), nil
}

type ConsumeOptions struct {
	Consumer  string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqp.Table
}
