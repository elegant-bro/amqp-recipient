package amqp_recipient

import "github.com/streadway/amqp"

type StubRecipient struct {
	job Job
	err error
}

func NoErrorStubRecipient(job Job) *StubRecipient {
	return NewStubRecipient(job, nil)
}

func NewStubRecipient(job Job, err error) *StubRecipient {
	return &StubRecipient{job: job, err: err}
}

func (s StubRecipient) Subscribe() (Job, error) {
	return s.job, s.err
}

type StubJob struct {
	fn func()
}

func NewStubJob(fn func()) *StubJob {
	return &StubJob{fn: fn}
}

func (s StubJob) Run() {
	s.fn()
}

type StubJobHandler struct {
	fn func(d amqp.Delivery) (uint8, error)
}

func (s StubJobHandler) Handle(d amqp.Delivery) (uint8, error) {
	return s.fn(d)
}

func NewStubJobHandler(fn func(d amqp.Delivery) (uint8, error)) *StubJobHandler {
	return &StubJobHandler{fn}
}

type StubSender struct {
	fn func(p amqp.Publishing) error
}

func NewStubSender(fn func(p amqp.Publishing) error) *StubSender {
	return &StubSender{fn: fn}
}

func (f StubSender) Send(p amqp.Publishing) error {
	return f.fn(p)
}
