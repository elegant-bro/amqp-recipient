package recipients

import (
	"fmt"
	"github.com/streadway/amqp"
	"math/rand"
)

type AmqpDeliveries struct {
	conn     *amqp.Connection
	queue    string
	prefetch int
	options  ConsumeOptions
	ch       *amqp.Channel
}

func NewAmqpDeliveries(conn *amqp.Connection, queue string, prefetch int, options ConsumeOptions) *AmqpDeliveries {
	return &AmqpDeliveries{conn: conn, queue: queue, prefetch: prefetch, options: options}
}

func (a AmqpDeliveries) All() (<-chan amqp.Delivery, error) {
	ch, err := a.conn.Channel()
	if nil != err {
		return nil, err
	}
	a.ch = ch

	err = a.ch.Qos(a.prefetch, 0, false)
	if nil != err {
		return nil, err
	}

	tag := fmt.Sprintf("%d", rand.Int())
	return a.ch.Consume(
		a.queue,
		tag,
		a.options.AutoAck,
		a.options.Exclusive,
		a.options.NoLocal,
		a.options.NoWait,
		a.options.Args,
	)
}

func (a AmqpDeliveries) Close() error {
	return a.ch.Close()
}
