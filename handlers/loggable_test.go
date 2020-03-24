package handlers

import (
	"errors"
	"github.com/streadway/amqp"
	"testing"
)

func TestLoggableHandler_Handle(t *testing.T) {
	res, err := NewLoggable(
		NewFunc(func(d amqp.Delivery) (u uint8, err error) {
			if "foo" != string(d.Body) {
				t.Errorf("Message body is %s; foo expected", string(d.Body))
			}
			return 1, nil
		}),
		nil,
	).Handle(amqp.Delivery{Body: []byte("foo")})

	if res != 1 {
		t.Errorf("Handler result is %d; 1 expected", res)
	}

	if nil != err {
		t.Errorf("Handler error is %s; nil expected", err.Error())
	}
}

func TestLoggableHandler_HandleErr(t *testing.T) {
	res, err := NewLoggable(
		NewFunc(func(d amqp.Delivery) (u uint8, err error) {
			return 1, errors.New("bar")
		}),
		func(d amqp.Delivery, err error) {
			if "foo" != string(d.Body) {
				t.Errorf("Message body is %s; foo expected", string(d.Body))
			}

			if "bar" != err.Error() {
				t.Errorf("Handler error is %s; bar expected", err.Error())
			}
		},
	).Handle(amqp.Delivery{Body: []byte("foo")})

	if res != 1 {
		t.Errorf("Handler result is %d; 1 expected", res)
	}

	if nil == err {
		t.Error("Handler error is nil; bar expected")
		return
	}

	if "bar" != err.Error() {
		t.Errorf("Handler error is %s; bar expected", err.Error())
	}
}
