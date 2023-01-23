package handlers

import (
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"testing"
)

type stubIds struct {
	has bool
	err error
}

func (s stubIds) Has(_ string) (bool, error) {
	return s.has, s.err
}

func (s stubIds) Save(_ string, fn func() (uint8, error)) (uint8, error) {
	return fn()
}

func TestIdempotentHandler_HandleEmptyMsgId(t *testing.T) {
	res, err := NewIdempotent(
		NewFunc(func(d amqp.Delivery) (u uint8, err error) {
			return 1, errors.New("baz")
		}),
		nil,
	).Handle(amqp.Delivery{})

	if res != 1 {
		t.Errorf("Handler result is %d; 1 expected", res)
	}

	if nil == err {
		t.Error("Handler error is nil; baz expected")
		return
	}

	if "baz" != err.Error() {
		t.Errorf("Handler error is %s; baz expected", err.Error())
	}
}

func TestIdempotentHandler_HandleHasMsg(t *testing.T) {
	res, err := NewIdempotent(nil, stubIds{has: true}).
		Handle(amqp.Delivery{MessageId: "some_id"})

	if res != 0 {
		t.Errorf("Handler result is %d; 0 expected", res)
	}

	if nil != err {
		t.Errorf("Handler error is %s; nil expected", err.Error())
	}
}

func TestIdempotentHandler_Handle(t *testing.T) {
	res, err := NewIdempotent(
		NewFunc(func(d amqp.Delivery) (u uint8, err error) {
			return 1, errors.New("foo")
		}),
		stubIds{has: false},
	).Handle(amqp.Delivery{MessageId: "some_id"})

	if res != 1 {
		t.Errorf("Handler result is %d; 1 expected", res)
	}

	if nil == err {
		t.Error("Handler error is nil; foo expected")
		return
	}

	if "foo" != err.Error() {
		t.Errorf("Handler error is %s; foo expected", err.Error())
	}
}

func TestIdempotentHandler_HandleIdsFails(t *testing.T) {
	res, err := NewIdempotent(
		NewFunc(func(d amqp.Delivery) (u uint8, err error) {
			return 1, errors.New("foo")
		}),
		stubIds{has: false, err: errors.New("bar")},
	).Handle(amqp.Delivery{MessageId: "some_id"})

	if res != 2 {
		t.Errorf("Handler result is %d; 1 expected", res)
	}

	if nil == err {
		t.Error("Handler error is nil; bar expected")
		return
	}

	if "can't check message id has been handled: bar" != err.Error() {
		t.Errorf("Handler error is %s; %s expected", err.Error(), "can't check message id has been handled: bar")
	}
}
