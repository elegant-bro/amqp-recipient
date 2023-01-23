package handlers

import (
	"errors"
	amqpRecipient "github.com/elegant-bro/amqp-recipient"
	amqp "github.com/rabbitmq/amqp091-go"
	"testing"
)

func TestMapHandler_Handle(t *testing.T) {
	res, err := NewDefaultMap(
		map[string]amqpRecipient.JobHandler{
			"foo": NewFunc(func(d amqp.Delivery) (u uint8, err error) {
				return 0, errors.New("foo call")
			}),
			"bar": NewFunc(func(d amqp.Delivery) (u uint8, err error) {
				return 1, nil
			}),
			"baz": NewFunc(func(d amqp.Delivery) (u uint8, err error) {
				return 0, errors.New("baz call")
			}),
		},
		func(delivery amqp.Delivery) (s string, err error) {
			return "bar", nil
		},
	).Handle(amqp.Delivery{})

	if res != 1 {
		t.Errorf("Handler result is %d; 1 expected", res)
	}

	if nil != err {
		t.Errorf("Handler error is %s; nil expected", err.Error())
	}
}

func TestMapHandler_HandleEmptyMap(t *testing.T) {
	res, err := NewDefaultMap(
		map[string]amqpRecipient.JobHandler{
			"foo": NewFunc(func(d amqp.Delivery) (u uint8, err error) {
				return 1, nil
			}),
		},
		func(delivery amqp.Delivery) (s string, err error) {
			return "", errors.New("bar")
		},
	).Handle(amqp.Delivery{})

	if res != 0 {
		t.Errorf("Handler result is %d; 0 expected", res)
	}

	if nil == err {
		t.Error("Handler error expected")
		return
	}

	if "fail to compute handler key: bar" != err.Error() {
		t.Error("Unexpected error message")
	}
}

func TestMapHandler_HandleKeyFuncFails(t *testing.T) {
	res, err := NewMap(
		map[string]amqpRecipient.JobHandler{},
		func(delivery amqp.Delivery) (s string, err error) {
			return "bar", nil
		},
		2,
	).Handle(amqp.Delivery{})

	if nil == err {
		t.Errorf("Handler error expected")
		return
	}

	if "handler not found: key 'bar'" != err.Error() {
		t.Error("Unexpected error message")
	}

	if res != 2 {
		t.Errorf("Handler result is %d; 2 expected", res)
	}
}
