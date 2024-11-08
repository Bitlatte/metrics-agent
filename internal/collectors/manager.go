package collectors

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/Bitlatte/metrics-agent/internal/config"
	"github.com/Bitlatte/metrics-agent/internal/types"
)

type Manager struct {
	collectors  map[string]Collector
	metrics     []types.Metric
	metricsChan chan []types.Metric
	mu          sync.RWMutex
	wg          sync.WaitGroup
}

func NewManager(config config.CollectionConfig) (*Manager, error) {
	if config.BatchSize <= 0 {
		return nil, fmt.Errorf("batch size must be between 0 and 255")
	}

	m := &Manager{
		collectors:  make(map[string]Collector),
		metricsChan: make(chan []types.Metric, 100),
	}

	return m, nil
}

func (m *Manager) RegisterCollector(name string, collector Collector) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.collectors[name]; exists {
		return errors.New(fmt.Sprintf("collector '%s' already registered", name))
	}

	m.collectors[name] = collector
	return nil
}

func (m *Manager) GetMetrics() []types.Metric {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var metrics []types.Metric
	for {
		select {
		case m := <-m.metricsChan:
			metrics = append(metrics, m...)
		default:
			return metrics
		}
	}
}

func (m *Manager) Start(ctx context.Context) error {
	if len(m.collectors) == 0 {
		return errors.New("no collectors registered")
	}

	for name, collector := range m.collectors {
		m.wg.Add(1)
		go func(name string, c Collector) error {
			defer m.wg.Done()
			if err := c.Start(ctx, m.metricsChan); err != nil {
				msg := fmt.Sprintf("Error starting collector '%s': %v", name, err)
				return errors.New(msg)
			}
			return nil
		}(name, collector)
	}
	return nil
}

func (m *Manager) Stop() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var errs []error
	for name, collector := range m.collectors {
		if err := collector.Stop(); err != nil {
			errs = append(errs, fmt.Errorf("failed to stop collector '%s': %v", name, err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors stopping collectors: %v", errs)
	}
	return nil
}
