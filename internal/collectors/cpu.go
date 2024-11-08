package collectors

import (
	"context"
	"log"
	"time"

	"github.com/Bitlatte/metrics-agent/internal/types"
	"github.com/shirou/gopsutil/v4/cpu"
)

type CPUCollector struct {
	*BaseCollector
}

func NewCPUCollector(interval time.Duration) *CPUCollector {
	return &CPUCollector{
		BaseCollector: &BaseCollector{
			device:   "cpu",
			interval: interval,
			done:     make(chan struct{}),
		},
	}
}

func (c *CPUCollector) Collect() ([]types.Metric, error) {
	// Collect CPU Metrics
	var metrics []types.Metric

	percentages, err := cpu.Percent(0, false)
	if err != nil {
		return nil, err
	}
	usage := percentages[0]
	metrics = append(metrics, types.Metric{
		Type:       types.MetricCPUUsage,
		ValueType:  types.ValueTypeFloat,
		FloatValue: &usage,
		Unit:       "percent",
		Labels:     types.Labels{"type": "total"},
		Timestamp:  time.Now().Unix(),
	})

	info, err := cpu.Info()
	if err != nil {
		return nil, err
	}
	model := info[0].ModelName
	metrics = append(metrics, types.Metric{
		Type:        types.MetricCPUModel,
		ValueType:   types.ValueTypeString,
		StringValue: &model,
		Unit:        "nil",
		Timestamp:   time.Now().Unix(),
	})

	physical, err := cpu.Counts(false)
	if err != nil {
		return nil, err
	}
	metrics = append(metrics, types.Metric{
		Type:      types.MetricCPUCountPhysical,
		ValueType: types.ValueTypeInt,
		IntValue:  &physical,
		Unit:      "cores",
		Labels:    types.Labels{"type": "physical"},
		Timestamp: time.Now().Unix(),
	})

	logical, err := cpu.Counts(true)
	if err != nil {
		return nil, err
	}
	metrics = append(metrics, types.Metric{
		Type:      types.MetricCPUCountLogical,
		ValueType: types.ValueTypeInt,
		IntValue:  &logical,
		Unit:      "cores",
		Labels:    types.Labels{"type": "logical"},
		Timestamp: time.Now().Unix(),
	})

	return metrics, nil
}

func (c *CPUCollector) Start(ctx context.Context, metrics chan<- []types.Metric) error {
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
