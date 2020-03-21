package main

import (
	amqpRecipient "github.com/elegant-bro/amqp-recipient"
	"github.com/streadway/amqp"
	"log"
)

type Api struct{}

func (a *Api) Send( /*some args*/ ) error {
	// some api call

	return nil // or send error
}

type sendHandler struct {
	api *Api
}

// implementation of github.com/elegant-bro/amqp_recipient/JobHandler
func (s sendHandler) Handle(d *amqp.Delivery) (uint8, error) {
	err := s.api.Send( /*convert delivery to args*/ )
	if nil != err {
		return amqpRecipient.HandlerReject, err
	}

	return amqpRecipient.HandlerAck, nil
}

func rabbitConnFromDsn(dsn string) *amqp.Connection {
	conn, _ := amqp.Dial(dsn)
	// don't ignore error in real project
	return conn
}

func main() {
	// first we create the recipient
	r := amqpRecipient.NewDefaultAmqpRecipient(
		"user.registered", // queue name
		16,                // prefetch count
		rabbitConnFromDsn("rabbit dsn"),
		&sendHandler{&Api{}},
		func(d *amqp.Delivery, err error) {
			log.Printf("Message with id '%s' fails: %v", d.MessageId, err)
		},
	)

	// then we subscribe it
	job, err := r.Subscribe()
	// Internally each subscription fetch channel from the connection and call channel.Consume() method
	// you can call recipient.Subscribe as many times as consumers you need for the queue
	if nil != err {
		log.Fatalf("%s: %v", "can't subscribe to user.registered", err)
	}

	// the last we run the job as goroutine
	go job.Run()

	// prevent exit
	<-make(chan bool)
}
