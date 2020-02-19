package handlers

import (
	"errors"
	"github.com/elegant-bro/amqp-recipient"
	"github.com/streadway/amqp"
	"testing"
)

func TestPublishFailedHandler_HandleWithError(t *testing.T) {
	senderCall := 0
	h := NewPublishFailed(
		NewFunc(func(d *amqp.Delivery) (u uint8, err error) {
			return 2, errors.New("foo")
		}),
		amqp_recipient.NewStubSender(func(p amqp.Publishing) error {
			senderCall++

			if p.Headers["x-bar"] != "baz" {
				t.Errorf("Header 'x-bar' expected to be 'baz', '%s' given", p.Headers["x-bar"])
			}

			return nil
		}),
		func(d *amqp.Delivery, err error) map[string]string {
			return map[string]string{"bar": "baz"}
		},
	)

	res, err := h.Handle(&amqp.Delivery{})

	if 0 == senderCall {
		t.Error("Sender was not called")
	}

	if nil == err {
		t.Error("Handler error is nil; foo expected")
	}

	if res != 0 {
		t.Errorf("Handler result is %d; 0 expected", res)
	}
}

func TestPublishFailedHandler_HandleWithSendFails(t *testing.T) {
	senderCall := 0
	h := NewPublishFailed(
		NewFunc(func(d *amqp.Delivery) (u uint8, err error) {
			return 2, errors.New("foo")
		}),
		amqp_recipient.NewStubSender(func(p amqp.Publishing) error {
			senderCall++

			if p.Headers["x-bar"] != "baz" {
				t.Errorf("Header 'x-bar' expected to be 'baz', '%s' given", p.Headers["x-bar"])
			}

			if p.Headers["x-origin"] != "hello" {
				t.Errorf("Header 'x-origin' expected to be 'hello', '%s' given", p.Headers["x-origin"])
			}

			return errors.New("bar")
		}),
		func(d *amqp.Delivery, err error) map[string]string {
			return map[string]string{"bar": "baz"}
		},
	)

	res, err := h.Handle(&amqp.Delivery{
		Headers: map[string]interface{}{"x-origin": "hello"},
	})

	if 0 == senderCall {
		t.Error("Sender was not called")
	}

	if nil == err {
		t.Error("Handler error is nil; bar expected")
	}

	if res != 1 {
		t.Errorf("Handler result is %d; 1 expected", res)
	}
}

func TestPublishFailedHandler_HandleWithoutError(t *testing.T) {
	senderCall := 0
	h := NewPublishFailed(
		NewFunc(func(d *amqp.Delivery) (u uint8, err error) {
			return 1, nil
		}),
		amqp_recipient.NewStubSender(func(p amqp.Publishing) error {
			senderCall++
			return nil
		}),
		func(d *amqp.Delivery, err error) map[string]string {
			return map[string]string{"bar": "baz"}
		},
	)

	res, err := h.Handle(&amqp.Delivery{})

	if 0 != senderCall {
		t.Error("Sender call is not expected")
	}

	if nil != err {
		t.Errorf("Handler error is %s; nil expected", err.Error())
	}

	if res != 1 {
		t.Errorf("Handler result is %d; 1 expected", res)
	}
}

func TestDefaultMakeHeaders(t *testing.T) {
	h := defaultMakeHeaders(
		&amqp.Delivery{
			Exchange:   "bar",
			RoutingKey: "baz",
		},
		errors.New("foo"),
	)

	if h["exception-message"] != "foo" {
		t.Errorf("Header 'exception-message' expected to be 'foo', '%s' given", h["exception-message"])
	}

	if h["origin-exchange"] != "bar" {
		t.Errorf("Header 'origin-exchange' expected to be 'bar', '%s' given", h["origin-exchange"])
	}

	if h["origin-routing-key"] != "baz" {
		t.Errorf("Header 'origin-routing-key' expected to be 'baz', '%s' given", h["origin-routing-key"])
	}
}

func TestNewPublishFailedStd(t *testing.T) {
	h := NewPublishFailedStd(
		NewFunc(func(d *amqp.Delivery) (u uint8, err error) {
			return 0, nil
		}),
		&amqp.Connection{},
		"foo",
		"bar",
	)
	if nil == h {
		t.Error("Handler is nil")
	}
}
