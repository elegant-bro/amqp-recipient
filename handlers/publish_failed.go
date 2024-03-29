package handlers

import (
	rcp "github.com/elegant-bro/amqp-recipient"
	amqp "github.com/rabbitmq/amqp091-go"
)

type PublishFailedHandler struct {
	Origin      rcp.JobHandler
	Sender      rcp.Sender
	MakeHeaders func(d amqp.Delivery, err error) map[string]string
}

func NewPublishFailed(origin rcp.JobHandler, sender rcp.Sender, makeHeaders func(d amqp.Delivery, err error) map[string]string) rcp.JobHandler {
	return &PublishFailedHandler{
		Origin:      origin,
		Sender:      sender,
		MakeHeaders: makeHeaders,
	}
}

func NewPublishFailedStd(origin rcp.JobHandler, connection *amqp.Connection, exchange string, routingKey string) rcp.JobHandler {
	return NewPublishFailed(
		origin,
		rcp.NewAmqpSender(connection, exchange, routingKey),
		defaultMakeHeaders,
	)
}

func (f *PublishFailedHandler) Handle(d amqp.Delivery) (uint8, error) {
	res, err := f.Origin.Handle(d)
	if nil != err {
		var newHeaders amqp.Table

		if nil == d.Headers {
			newHeaders = amqp.Table{}
		} else {
			newHeaders = d.Headers
		}

		for k, v := range f.MakeHeaders(d, err) {
			newHeaders["x-"+k] = v
		}

		pubErr := f.Sender.Send(amqp.Publishing{
			Headers:         newHeaders,
			ContentType:     d.ContentType,
			ContentEncoding: d.ContentEncoding,
			DeliveryMode:    d.DeliveryMode,
			Priority:        d.Priority,
			CorrelationId:   d.CorrelationId,
			ReplyTo:         d.ReplyTo,
			Expiration:      d.Expiration,
			MessageId:       d.MessageId,
			Timestamp:       d.Timestamp,
			Type:            d.Type,
			UserId:          d.UserId,
			AppId:           d.AppId,
			Body:            d.Body,
		})
		if nil != pubErr {
			return rcp.HandlerRequeue, newInternalError(pubErr, "sending failed message fails")
		}

		return rcp.HandlerAck, err
	}

	return res, nil
}

func defaultMakeHeaders(d amqp.Delivery, err error) map[string]string {
	return map[string]string{
		"exception-message":  err.Error(),
		"origin-exchange":    d.Exchange,
		"origin-routing-key": d.RoutingKey,
	}
}
