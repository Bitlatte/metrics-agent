package collectors

import (
	"context"
	"time"

	"github.com/Bitlatte/metrics-agent/internal/types"
)

type Collector interface {
	Collect() ([]types.Metric, error)
	Start(context.Context, chan<- []types.Metric) error
	Stop() error
	GetInterval() time.Duration
}

type BaseCollector struct {
	device   string
	interval time.Duration
	done     chan struct{}
}

func NewBaseCollector(device string, interval time.Duration) BaseCollector {
	return BaseCollector{
		device:   device,
		interval: interval,
		done:     make(chan struct{}),
	}
}

func (c *BaseCollector) GetInterval() time.Duration {
	return c.interval
}

func (c *BaseCollector) Stop() error {
	close(c.done)
	return nil
}
