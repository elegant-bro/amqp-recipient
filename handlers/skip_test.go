package handlers

import (
	"errors"
	recipient "github.com/elegant-bro/amqp-recipient"
	"github.com/streadway/amqp"
	"testing"
)

func TestSkipHandler_HandleSkipped(t *testing.T) {
	res, err := NewSkipHandler(
		NewFunc(func(d amqp.Delivery) (uint8, error) {
			return recipient.HandlerReject, errors.New("foo")
		}),
		func(d amqp.Delivery) bool {
			return true
		},
		recipient.HandlerAck,
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
	res, err := NewSkipHandler(
		NewFunc(func(d amqp.Delivery) (uint8, error) {
			return recipient.HandlerReject, errors.New("foo")
		}),
		func(d amqp.Delivery) bool {
			return false
		},
		recipient.HandlerAck,
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
	_, _ = NewSkipHandler(
		NewFunc(func(d amqp.Delivery) (uint8, error) {
			return recipient.HandlerReject, errors.New("foo")
		}),
		func(d amqp.Delivery) bool {
			skipCall++
			if d.MessageId != "82c90db9-501a-4334-afe6-f203198eb04f" {
				t.Errorf("MessageId is %s; '82c90db9-501a-4334-afe6-f203198eb04f' expected", d.MessageId)
			}
			return true
		},
		recipient.HandlerAck,
	).Handle(
		amqp.Delivery{
			MessageId: "82c90db9-501a-4334-afe6-f203198eb04f",
		},
	)

	if skipCall != 1 {
		t.Errorf("Skip called %d times; 1 expected", skipCall)
	}
}

func TestSkipAckHandler(t *testing.T) {
	h := SkipAckHandler(
		NewFunc(func(d amqp.Delivery) (uint8, error) {
			return recipient.HandlerReject, errors.New("foo")
		}),
		func(d amqp.Delivery) bool {
			return true
		},
	)

	if h.res != recipient.HandlerAck {
		t.Errorf("Handler skip result is %d; 0 expected", h.res)
	}
}
