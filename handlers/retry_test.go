package handlers

import (
	"errors"
	"github.com/streadway/amqp"
	"testing"
)

func stubXDeathDelivery(count int64) *amqp.Delivery {
	return &amqp.Delivery{
		Headers: amqp.Table{
			"x-death": []interface{}{
				amqp.Table{
					"count":    count,
					"exchange": "test",
				},
			},
		},
	}
}

func TestRetryHandler_HandleNoErr(t *testing.T) {
	res, err := NewRetry(
		NewFunc(func(d *amqp.Delivery) (u uint8, err error) {
			return 1, nil
		}),
		5,
	).Handle(&amqp.Delivery{})

	if res != 1 {
		t.Errorf("Handler result is %d; 1 expected", res)
	}

	if nil != err {
		t.Errorf("Handler error is %s; nil expected", err.Error())
	}
}

func TestRetryHandler_HandleErr(t *testing.T) {
	res, err := NewRetry(
		NewFunc(func(d *amqp.Delivery) (u uint8, err error) {
			return 1, errors.New("foo")
		}),
		5,
	).Handle(stubXDeathDelivery(3))

	if res != 2 {
		t.Errorf("Handler result is %d; 2 expected", res)
	}

	if nil != err {
		t.Errorf("Handler error is %s; nil expected", err.Error())
	}
}

func TestRetryHandler_HandleErrRetriesReached(t *testing.T) {
	res, err := NewRetry(
		NewFunc(func(d *amqp.Delivery) (u uint8, err error) {
			return 1, errors.New("foo")
		}),
		5,
	).Handle(stubXDeathDelivery(5))

	if res != 0 {
		t.Errorf("Handler result is %d; 0 expected", res)
	}

	if nil == err {
		t.Error("Handler error is nil; foo expected")
		return
	}

	if "foo" != err.Error() {
		t.Errorf("Handler error is %s; foo expected", err.Error())
	}
}

func TestRetryHandler_HandleNoHeaders(t *testing.T) {
	res, err := NewRetry(
		NewFunc(func(d *amqp.Delivery) (u uint8, err error) {
			return 1, errors.New("foo")
		}),
		5,
	).Handle(&amqp.Delivery{})

	if res != 2 {
		t.Errorf("Handler result is %d; 2 expected", res)
	}

	if nil != err {
		t.Errorf("Handler error is %s; nil expected", err.Error())
	}
}
