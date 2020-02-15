package amqp_recipient

import (
	"fmt"
	"github.com/streadway/amqp"
)

type AmqpSender struct {
	Conn       *amqp.Connection
	Exchange   string
	RoutingKey string
}

func NewAmqpSender(conn *amqp.Connection, exchange string, routingKey string) *AmqpSender {
	return &AmqpSender{Conn: conn, Exchange: exchange, RoutingKey: routingKey}
}

func (a AmqpSender) Send(p amqp.Publishing) error {
	ch, err := a.Conn.Channel()
	if nil != err {
		return err
	}
	defer func() {
		fmt.Println(ch.Close())
	}()

	return ch.Publish(a.Exchange, a.RoutingKey, false, false, p)
}
