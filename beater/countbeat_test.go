package beater

import (
	"testing"
	"time"

	"github.com/elastic/beats/v7/libbeat/beat"
	pubtest "github.com/elastic/beats/v7/libbeat/publisher/testing"
	conf "github.com/elastic/elastic-agent-libs/config"
	"github.com/elastic/elastic-agent-libs/logp/logptest"
)

func newTestBeat(t *testing.T, client beat.Client) *beat.Beat {
	return &beat.Beat{
		Info:      beat.Info{Name: "countbeat", Logger: logptest.NewTestingLogger(t, "")},
		Publisher: pubtest.PublisherWithClient(client),
	}
}

func newTestBeater(t *testing.T, b *beat.Beat, settings map[string]any) beat.Beater {
	bt, err := New(b, conf.MustNewConfigFrom(settings))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return bt
}

// runAsync starts Run and returns a channel that receives its result.
func runAsync(bt beat.Beater, b *beat.Beat) <-chan error {
	errCh := make(chan error, 1)
	go func() { errCh <- bt.Run(b) }()
	return errCh
}

func waitForRun(t *testing.T, errCh <-chan error) {
	t.Helper()
	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("Run returned error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Run did not return after Stop")
	}
}

func TestRunPublishesIncrementingCounter(t *testing.T) {
	client := pubtest.NewChanClient(10)
	b := newTestBeat(t, client)
	bt := newTestBeater(t, b, map[string]any{"period": "10ms"})

	errCh := runAsync(bt, b)
	for want := 1; want <= 3; want++ {
		select {
		case event := <-client.Channel:
			if got := event.Fields["counter"]; got != want {
				t.Errorf("counter = %v, want %d", got, want)
			}
			if got := event.Fields["type"]; got != "countbeat" {
				t.Errorf("type = %v, want countbeat", got)
			}
		case <-time.After(5 * time.Second):
			t.Fatalf("timed out waiting for event %d", want)
		}
	}

	bt.Stop()
	waitForRun(t, errCh)
}

func TestStopBeforeRun(t *testing.T) {
	client := pubtest.NewChanClient(10)
	b := newTestBeat(t, client)
	bt := newTestBeater(t, b, map[string]any{"period": "1h"})

	// Stop used to dereference the client, which is only set once Run connects.
	bt.Stop()
	waitForRun(t, runAsync(bt, b))
}

func TestNewConfig(t *testing.T) {
	b := newTestBeat(t, pubtest.NewChanClient(1))

	bt := newTestBeater(t, b, map[string]any{})
	if got := bt.(*countbeat).config.Period; got != time.Second {
		t.Errorf("default period = %v, want 1s", got)
	}

	bt = newTestBeater(t, b, map[string]any{"period": "250ms"})
	if got := bt.(*countbeat).config.Period; got != 250*time.Millisecond {
		t.Errorf("period = %v, want 250ms", got)
	}

	if _, err := New(b, conf.MustNewConfigFrom(map[string]any{"period": "soon"})); err == nil {
		t.Error("New accepted an invalid period")
	}
}
