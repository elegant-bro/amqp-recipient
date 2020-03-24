package recipients

import "testing"

func TestNewDefaultAmqpRecipient(t *testing.T) {
	r := NewDefaultAmqpRecipient(
		"foo",
		8,
		nil,
		nil,
		nil,
	)

	if "foo" != r.queue {
		t.Errorf("queue name expected to be 'foo', '%s' given", r.queue)
	}

	if 8 != r.prefetch {
		t.Errorf("prefetch expected to be 8, '%d' given", r.prefetch)
	}

	if nil != r.consumeOptions.Args {
		t.Error("Options Args expected to be null")
	}

	if r.consumeOptions.NoWait {
		t.Error("Options NoWait expected to be false")
	}

	if r.consumeOptions.NoLocal {
		t.Error("Options NoLocal expected to be false")
	}

	if r.consumeOptions.Exclusive {
		t.Error("Options Exclusive expected to be false")
	}

	if r.consumeOptions.AutoAck {
		t.Error("Options AutoAck expected to be false")
	}
}
