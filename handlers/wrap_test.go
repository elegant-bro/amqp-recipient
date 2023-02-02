package handlers

import (
	"errors"
	rcp "github.com/elegant-bro/amqp-recipient"
	amqp "github.com/rabbitmq/amqp091-go"
	"testing"
)

func TestWrapHandler_Handle(t *testing.T) {
	origCall := 0
	w := NewWrapHandler(
		NewFunc(func(d amqp.Delivery) (uint8, error) {
			origCall++
			return 1, errors.New("bar")
		}),
		func(d amqp.Delivery, wrapped rcp.JobHandler) (uint8, error) {
			if "foo" != d.MessageId {
				t.Errorf("delivery message id is %s; foo expected", d.MessageId)
			}
			return wrapped.Handle(d)
		},
	)

	res, err := w.Handle(amqp.Delivery{MessageId: "foo"})
	if origCall != 1 {
		t.Errorf("origin was called %d times; 1 expected", origCall)
	}

	if res != 1 {
		t.Errorf("Handler result is %d; 1 expected", res)
	}

	if "bar" != err.Error() {
		t.Errorf("Error message is %s; bar expected", err.Error())
	}
}
