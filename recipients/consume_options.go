package recipients

import amqp "github.com/rabbitmq/amqp091-go"

type ConsumeOptions struct {
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqp.Table
}
