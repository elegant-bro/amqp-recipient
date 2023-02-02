package handlers

import (
	"errors"
	rcp "github.com/elegant-bro/amqp-recipient"
	amqp "github.com/rabbitmq/amqp091-go"
	"testing"
)

func TestSkipHandler_HandleSkipped(t *testing.T) {
	res, err := NewSkip(
		NewFunc(func(d amqp.Delivery) (uint8, error) {
			return rcp.HandlerReject, errors.New("foo")
		}),
		func(d amqp.Delivery) bool {
			return true
		},
		rcp.HandlerAck,
	).Handle(
		amqp.Delivery{},
	)

	if res != 0 {
		t.Errorf("Handler result is %d; 0 expected", res)
	}

	if nil != err {
		t.Errorf("Handler error is %s; nil expected", err.Error())
	}
}

func TestSkipHandler_HandleNotSkipped(t *testing.T) {
	res, err := NewSkip(
		NewFunc(func(d amqp.Delivery) (uint8, error) {
			return rcp.HandlerReject, errors.New("foo")
		}),
		func(d amqp.Delivery) bool {
			return false
		},
		rcp.HandlerAck,
	).Handle(
		amqp.Delivery{},
	)

	if res != 2 {
		t.Errorf("Handler result is %d; 2 expected", res)
	}

	if nil == err {
		t.Error("Handler error is nil; 'foo' expected")
		return
	}

	if err.Error() != "foo" {
		t.Errorf("Error is %s; 'foo' expected", err.Error())
	}
}

func TestSkipHandler_HandleSkipCalled(t *testing.T) {
	skipCall := 0
	_, _ = NewSkipAck(
		NewFunc(func(d amqp.Delivery) (uint8, error) {
			return rcp.HandlerReject, errors.New("foo")
		}),
		func(d amqp.Delivery) bool {
			skipCall++
			if d.MessageId != "82c90db9-501a-4334-afe6-f203198eb04f" {
				t.Errorf("MessageId is %s; '82c90db9-501a-4334-afe6-f203198eb04f' expected", d.MessageId)
			}
			return true
		},
	).Handle(
		amqp.Delivery{
			MessageId: "82c90db9-501a-4334-afe6-f203198eb04f",
		},
	)

	if skipCall != 1 {
		t.Errorf("Skip called %d times; 1 expected", skipCall)
	}
}
