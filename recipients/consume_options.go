package recipients

import "github.com/streadway/amqp"

type ConsumeOptions struct {
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqp.Table
}
