package amqp_recipient

import "github.com/streadway/amqp"

type FuncSender struct {
	fn func(p amqp.Publishing) error
}

func NewFuncSender(fn func(p amqp.Publishing) error) *FuncSender {
	return &FuncSender{fn: fn}
}

func (f FuncSender) Send(p amqp.Publishing) error {
	return f.fn(p)
}
