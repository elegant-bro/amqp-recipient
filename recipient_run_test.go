package amqp_recipient

import (
	"errors"
	"sync"
	"testing"
)

func TestNewRunEntry(t *testing.T) {
	entry := NewRunEntry(NewStubRecipient(nil, nil), 2)

	if nil == entry.recipient {
		t.Errorf("expected not nil recipient")
	}

	if 2 != entry.consumers {
		t.Errorf("expected consumers is 2, %d given", entry.consumers)
	}
}

func TestOneRunEntry(t *testing.T) {
	entry := OneRunEntry(NewStubRecipient(nil, nil))

	if nil == entry.recipient {
		t.Errorf("expected not nil recipient")
	}

	if 1 != entry.consumers {
		t.Errorf("expected consumers is 1, %d given", entry.consumers)
	}
}

func TestRun_All(t *testing.T) {
	var wg sync.WaitGroup
	runCalled := 0
	run := NewRun(
		[]RunEntry{
			OneRunEntry(
				NoErrorStubRecipient(NewStubJob(func() {
					runCalled++
					wg.Done()
				})),
			),
			NewRunEntry(
				NoErrorStubRecipient(NewStubJob(func() {
					runCalled++
					wg.Done()
				})), 2,
			),
		},
		func(err error) {
			t.Error("unexpected call")
		},
	)

	wg.Add(3)
	run.All()
	wg.Wait()

	if 3 != runCalled {
		t.Errorf("expected job run called 3 times, %d happens", runCalled)
	}
}

func TestRun_AllFail(t *testing.T) {
	var wg sync.WaitGroup
	runCalled := 0
	run := NewRun(
		[]RunEntry{
			OneRunEntry(
				NewStubRecipient(NewStubJob(func() {
					t.Error("unexpected call")
				}), errors.New("foo")),
			),
			NewRunEntry(
				NoErrorStubRecipient(NewStubJob(func() {
					runCalled++
					wg.Done()
				})), 2,
			),
		},
		func(err error) {
			runCalled++
			wg.Done()
		},
	)

	wg.Add(3)
	run.All()
	wg.Wait()

	if 3 != runCalled {
		t.Errorf("expected job run called 3 times, %d happens", runCalled)
	}
}
