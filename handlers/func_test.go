package handlers

import (
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"testing"
)

func TestFuncHandler_Handle(t *testing.T) {
	res, err := NewFunc(func(d amqp.Delivery) (u uint8, err error) {
		return 1, errors.New("foo")
	}).Handle(amqp.Delivery{})

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
