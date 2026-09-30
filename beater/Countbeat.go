package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/v7/libbeat/beat"
	conf "github.com/elastic/elastic-agent-libs/config"
	"github.com/elastic/elastic-agent-libs/logp"
	"github.com/elastic/elastic-agent-libs/mapstr"

	"github.com/jimmino/countbeat/config"
)

// countbeat configuration.
type countbeat struct {
	done   chan struct{}
	config config.Config
	logger *logp.Logger
}

// New creates an instance of countbeat.
func New(b *beat.Beat, cfg *conf.C) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &countbeat{
		done:   make(chan struct{}),
		config: c,
		logger: b.Info.Logger.Named("countbeat"),
	}
	return bt, nil
}

// Run starts countbeat.
func (bt *countbeat) Run(b *beat.Beat) error {
	bt.logger.Info("countbeat is running! Hit CTRL-C to stop it.")

	client, err := b.Publisher.Connect()
	if err != nil {
		return err
	}
	defer client.Close()

	ticker := time.NewTicker(bt.config.Period)
	defer ticker.Stop()
	counter := 1
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		event := beat.Event{
			Timestamp: time.Now(),
			Fields: mapstr.M{
				"type":    b.Info.Name,
				"counter": counter,
			},
		}
		client.Publish(event)
		bt.logger.Info("Event sent")
		counter++
	}
}

// Stop stops countbeat.
func (bt *countbeat) Stop() {
	close(bt.done)
}
