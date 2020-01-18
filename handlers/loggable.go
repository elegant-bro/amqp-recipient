package handlers

import (
	recipient "github.com/elegant-bro/amqp-recipient"
	"github.com/streadway/amqp"
)

type LoggableHandler struct {
	origin   recipient.JobHandler
	writeLog func(err error)
}

func NewLoggable(origin recipient.JobHandler, writeLog func(err error)) *LoggableHandler {
	return &LoggableHandler{origin: origin, writeLog: writeLog}
}

func (l *LoggableHandler) Handle(d *amqp.Delivery) (res uint8, err error) {
	res, err = l.origin.Handle(d)
	if nil != err {
		l.writeLog(err)
	}
	return
}
