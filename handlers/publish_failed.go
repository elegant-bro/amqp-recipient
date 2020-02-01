package handlers

import (
	amqpRecipient "github.com/elegant-bro/amqp-recipient"
	"github.com/streadway/amqp"
)

type PublishFailedHandler struct {
	Origin      amqpRecipient.JobHandler
	Sender      amqpRecipient.Sender
	MakeHeaders func(d *amqp.Delivery, err error) map[string]string
}

func NewPublishFailed(origin amqpRecipient.JobHandler, sender amqpRecipient.Sender, makeHeaders func(d *amqp.Delivery, err error) map[string]string) *PublishFailedHandler {
	return &PublishFailedHandler{
		Origin:      origin,
		Sender:      sender,
		MakeHeaders: makeHeaders,
	}
}

func NewPublishFailedStd(origin amqpRecipient.JobHandler, connection *amqp.Connection, exchange string, routingKey string) *PublishFailedHandler {
	return NewPublishFailed(
		origin,
		amqpRecipient.NewAmqpSender(connection, exchange, routingKey),
		func(d *amqp.Delivery, err error) map[string]string {
			return map[string]string{
				"exception-message":  err.Error(),
				"origin-exchange":    d.Exchange,
				"origin-routing-key": d.RoutingKey,
			}
		},
	)
}

func (f PublishFailedHandler) Handle(d *amqp.Delivery) (uint8, error) {
	res, err := f.Origin.Handle(d)
	if nil != err {
		newHeaders := d.Headers
		for k, v := range f.MakeHeaders(d, err) {
			newHeaders["x-"+k] = v
		}

		err := f.Sender.Send(amqp.Publishing{
			Headers:         d.Headers,
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
		if nil != err {
			return 1, err
		}
	}

	return res, nil
}
