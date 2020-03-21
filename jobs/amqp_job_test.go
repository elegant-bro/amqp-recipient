package jobs

import (
	"github.com/elegant-bro/amqp-recipient"
	"github.com/streadway/amqp"
	"sync"
	"testing"
)

func TestAmqpJob_Run(t *testing.T) {

	queue := func(count int) <-chan amqp.Delivery {
		deliveries := make(chan amqp.Delivery)
		go func() {
			defer close(deliveries)
			for i := 0; i < count; i++ {
				deliveries <- amqp.Delivery{}
			}
		}()

		return deliveries
	}

	var wg sync.WaitGroup
	wg.Add(5)
	job := NewAmqpJob(
		queue(5),
		amqp_recipient.NewStubJobHandler(func(d *amqp.Delivery) (u uint8, err error) {
			wg.Done()
			return amqp_recipient.HandlerDoNothing, nil
		}),
		func(d *amqp.Delivery, err error) {

		},
	)

	go job.Run()
	wg.Wait()

}
