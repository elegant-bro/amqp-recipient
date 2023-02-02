package handlers

import (
	rcp "github.com/elegant-bro/amqp-recipient"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RetryHandler struct {
	origin     rcp.JobHandler
	maxRetries int64
}

func NewRetry(origin rcp.JobHandler, maxRetries int64) rcp.JobHandler {
	return &RetryHandler{origin: origin, maxRetries: maxRetries}
}

func (r *RetryHandler) Handle(d amqp.Delivery) (res uint8, err error) {
	res, err = r.origin.Handle(d)
	if nil != err {
		if xDeath(d.Headers) >= r.maxRetries {
			return rcp.HandlerAck, err
		}

		return rcp.HandlerReject, nil
	}

	return
}

func xDeath(headers amqp.Table) int64 {
	if xDeathHeader, ok := headers["x-death"]; ok {
		if unboxed, ok := xDeathHeader.([]interface{}); ok {
			if len(unboxed) > 0 {
				if xDeathUnboxed, ok := unboxed[0].(amqp.Table); ok {
					if count, ok := xDeathUnboxed["count"]; ok {
						return count.(int64)
					}
				}
			}
		}
	}

	return 0
}
