package collectors

import (
	"context"
	"log"
	"time"

	"github.com/Bitlatte/metrics-agent/internal/types"
	"github.com/shirou/gopsutil/v4/disk"
)

type DiskCollector struct {
	*BaseCollector
	ignorePaths []string
}

func NewDiskCollector(interval time.Duration, ignorePaths []string) *DiskCollector {
	return &DiskCollector{
		BaseCollector: &BaseCollector{
			device:   "disk",
			interval: interval,
			done:     make(chan struct{}),
		},
		ignorePaths: ignorePaths,
	}
}

func (c *DiskCollector) Collect() ([]types.Metric, error) {
	var metrics []types.Metric

	partitions, err := disk.Partitions(true)
	if err != nil {
		return nil, err
	}

	for _, partition := range partitions {
		if c.shouldIgnorePath(partition.Mountpoint) {
			continue
		}

		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			return nil, err
		}

		var used int = int(usage.Used)
		var free int = int(usage.Free)
		var total int = int(usage.Total)

		metrics = append(metrics, []types.Metric{
			{
				Type:       types.MetricDiskUsage,
				ValueType:  types.ValueTypeFloat,
				FloatValue: &usage.UsedPercent,
				Unit:       "percent",
				Labels: types.Labels{
					"device":     partition.Device,
					"mountpoint": partition.Mountpoint,
					"fstype":     partition.Fstype,
				},
				Timestamp: time.Now().Unix(),
			},
			{
				Type:      types.MetricDiskUsed,
				ValueType: types.ValueTypeInt,
				IntValue:  &used,
				Unit:      "bytes",
				Labels: types.Labels{
					"device":     partition.Device,
					"mountpoint": partition.Mountpoint,
					"fstype":     partition.Fstype,
				},
				Timestamp: time.Now().Unix(),
			},
			{
				Type:      types.MetricDiskFree,
				ValueType: types.ValueTypeInt,
				IntValue:  &free,
				Unit:      "bytes",
				Labels: types.Labels{
					"device":     partition.Device,
					"mountpoint": partition.Mountpoint,
					"fstype":     partition.Fstype,
				},
			},
			{
				Type:      types.MetricDiskTotal,
				ValueType: types.ValueTypeInt,
				IntValue:  &total,
				Unit:      "bytes",
				Labels: types.Labels{
					"device":     partition.Device,
					"mountpoint": partition.Mountpoint,
					"fstype":     partition.Fstype,
				},
			},
		}...)
	}

	return metrics, nil
}

func (c *DiskCollector) Start(ctx context.Context, metrics chan<- []types.Metric) error {
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

func (c *DiskCollector) shouldIgnorePath(path string) bool {
	if len(c.ignorePaths) == 0 {
		return false
	}
	for _, ignorePath := range c.ignorePaths {
		if path == ignorePath {
			return true
		}
	}
	return false
}
