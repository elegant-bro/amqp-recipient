package handlers

import (
	rcp "github.com/elegant-bro/amqp-recipient"
	amqp "github.com/rabbitmq/amqp091-go"
)

type LoggableHandler struct {
	origin   rcp.JobHandler
	writeLog rcp.OnHandlerFails
}

func NewLoggable(origin rcp.JobHandler, writeLog rcp.OnHandlerFails) rcp.JobHandler {
	return &LoggableHandler{origin: origin, writeLog: writeLog}
}

func (l *LoggableHandler) Handle(d amqp.Delivery) (res uint8, err error) {
	res, err = l.origin.Handle(d)
	if nil != err {
		l.writeLog(d, err)
	}
	return
}
