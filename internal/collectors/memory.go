package collectors

import (
	"context"
	"log"
	"time"

	"github.com/Bitlatte/metrics-agent/internal/types"
	"github.com/shirou/gopsutil/v4/mem"
)

type MemoryCollector struct {
	*BaseCollector
}

func NewMemoryCollector(interval time.Duration) *MemoryCollector {
	return &MemoryCollector{
		BaseCollector: &BaseCollector{
			device:   "memory",
			interval: interval,
			done:     make(chan struct{}),
		},
	}
}

func (c *MemoryCollector) Collect() ([]types.Metric, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	return []types.Metric{
		{
			Type:      types.MetricMemoryUsage,
			ValueType: types.ValueTypeInt,
			IntValue:  Uint64ToIntPtr(v.Used),
			Unit:      "bytes",
			Labels:    types.Labels{"type": "virtual"},
			Timestamp: time.Now().Unix(),
		},
		{
			Type:      types.MetricMemoryTotal,
			ValueType: types.ValueTypeInt,
			IntValue:  Uint64ToIntPtr(v.Total),
			Unit:      "bytes",
			Labels:    types.Labels{"type": "virtual"},
			Timestamp: time.Now().Unix(),
		},
		{
			Type:      types.MetricMemoryFree,
			ValueType: types.ValueTypeInt,
			IntValue:  Uint64ToIntPtr(v.Free),
			Unit:      "bytes",
			Labels:    types.Labels{"type": "virtual"},
			Timestamp: time.Now().Unix(),
		},
	}, nil
}

func (c *MemoryCollector) Start(ctx context.Context, metrics chan<- []types.Metric) error {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			collected, err := c.Collect()
			if err != nil {
				log.Printf("Error collecting %v Metrics: %v", c.device, err)
				continue
			}
			metrics <- collected
		}
	}
}

// Helper to convert uint64 to float64 pointer
func Uint64ToIntPtr(value uint64) *int {
	i := int(value)
	return &i
}
