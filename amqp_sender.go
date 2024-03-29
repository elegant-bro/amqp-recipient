package amqp_recipient

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
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
		if err := ch.Close(); nil != err {
			fmt.Println("closing AmqpSender channel failed: ", err)
		}
	}()

	return ch.Publish(a.Exchange, a.RoutingKey, false, false, p)
}
