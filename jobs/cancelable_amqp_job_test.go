package jobs

import (
	"context"
	amqpRecipient "github.com/elegant-bro/amqp-recipient"
	"github.com/streadway/amqp"
	"sync"
	"testing"
	"time"
)

func TestCancelableAmqpJob_Run(t *testing.T) {

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(5)
	job := NewCancelableAmqpJob(
		ctx,
		queue(5),
		amqpRecipient.NewStubJobHandler(func(d amqp.Delivery) (u uint8, err error) {
			wg.Done()
			return amqpRecipient.HandlerDoNothing, nil
		}),
		func(d amqp.Delivery, err error) {

		},
	)

	go job.Run()
	wg.Wait()
}

func TestCancelableAmqpJob_RunDone(t *testing.T) {

	infiniteQueue := func() <-chan amqp.Delivery {
		deliveries := make(chan amqp.Delivery)
		go func() {
			defer close(deliveries)
			for {
				deliveries <- amqp.Delivery{}
			}
		}()

		return deliveries
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	call := 0
	terminated := make(chan bool)
	job := NewCancelableAmqpJob(
		ctx,
		infiniteQueue(),
		amqpRecipient.NewStubJobHandler(func(d amqp.Delivery) (u uint8, err error) {
			call++
			return amqpRecipient.HandlerDoNothing, nil
		}),
		func(d amqp.Delivery, err error) {

		},
	)

	job.OnCancel(func() {
		close(terminated)
	})

	go job.Run()

	<-terminated
	if call < 1 {
		t.Error("Handler was not called")
	}

}

func TestCancelableAmqpJob_RunEnsureDoneLongTask(t *testing.T) {

	infiniteQueue := func() <-chan amqp.Delivery {
		deliveries := make(chan amqp.Delivery)
		go func() {
			defer close(deliveries)
			for {
				deliveries <- amqp.Delivery{}
			}
		}()

		return deliveries
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	call := 0
	terminated := make(chan bool)
	job := NewCancelableAmqpJob(
		ctx,
		infiniteQueue(),
		amqpRecipient.NewStubJobHandler(func(d amqp.Delivery) (u uint8, err error) {
			time.Sleep(50 * time.Millisecond) // long blocking task
			call++
			return amqpRecipient.HandlerDoNothing, nil
		}),
		func(d amqp.Delivery, err error) {

		},
	)

	job.OnCancel(func() {
		close(terminated)
	})

	go job.Run()

	<-terminated
	if call < 1 {
		t.Error("Handler was not called")
	}

}
