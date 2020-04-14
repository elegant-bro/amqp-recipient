package amqp_recipient

import "sync"

type Run struct {
	entries []RunEntry
	onFail  func(err error)
}

func NewRun(entries []RunEntry, onFail func(err error)) *Run {
	return &Run{entries: entries, onFail: onFail}
}

func (run *Run) All() {
	var jobs []Job
	for _, entry := range run.entries {
		for i := 0; i < entry.consumers; i++ {
			job, err := entry.recipient.Subscribe()
			if nil != err {
				run.onFail(err)
			} else {
				jobs = append(jobs, job)
			}
		}
	}

	var wg sync.WaitGroup
	wg.Add(len(jobs))
	for _, readyJob := range jobs {
		go func(j Job) {
			defer wg.Done()
			j.Run()
		}(readyJob)
	}
	wg.Wait()
}

func (run *Run) AllAsync() <-chan error {
	done := make(chan error)
	go func() {
		defer close(done)
		run.All()
	}()

	return done
}

type RunEntry struct {
	recipient Recipient
	consumers int
}

func OneRunEntry(recipient Recipient) RunEntry {
	return NewRunEntry(recipient, 1)
}

func NewRunEntry(recipient Recipient, consumers int) RunEntry {
	return RunEntry{recipient: recipient, consumers: consumers}
}
