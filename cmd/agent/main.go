package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/Bitlatte/metrics-agent/internal/collectors"
	"github.com/Bitlatte/metrics-agent/internal/config"
	"github.com/Bitlatte/metrics-agent/internal/types"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Optimize GC for lower latency
	debug.SetGCPercent(20)

	// Add memory limiting
	var memLimit int64 = 20000000 // 20 Megabytes
	debug.SetMemoryLimit(memLimit)

	// Load Configuration from file
	config, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	manager, err := initializeCollectors(config)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	if err := manager.Start(ctx); err != nil {
		log.Fatalf("Failed to start collectors: %v", err)
	}

	// Setup signal handling for graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	// Start Collection Loop
	go func() {
		ticker := time.NewTicker(config.Collection.Interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				metrics := manager.GetMetrics()
				if len(metrics) > 0 {
					// For now, just log the metrics
					for _, metric := range metrics {
						log.Printf("Metric - Type: %s, Value: %v, Unit: %s, Labels: %v\n",
							metric.Type,
							getMetricValue(metric),
							metric.Unit,
							metric.Labels)
					}
				}
			}
		}
	}()

	log.Println("Metrics Agent started successfully")

	// Wait for shutdown signal
	<-shutdown
	log.Println("Shutting down gracefully...")
	cancel()

	time.Sleep(time.Second * 2)
}

// Helper function to get the appropriate value based on the metric's value type
func getMetricValue(m types.Metric) interface{} {
	switch m.ValueType {
	case types.ValueTypeFloat:
		if m.FloatValue != nil {
			return *m.FloatValue
		}
	case types.ValueTypeString:
		if m.StringValue != nil {
			return *m.StringValue
		}
	case types.ValueTypeBool:
		if m.BoolValue != nil {
			return *m.BoolValue
		}
	case types.ValueTypeInt:
		if m.IntValue != nil {
			return *m.IntValue
		}
	}
	return "<no value>"
}

func initializeCollectors(config *config.Config) (*collectors.Manager, error) {
	manager, err := collectors.NewManager(config.Collection)
	if err != nil {
		return nil, err
	}

	if config.Collection.Collectors.CPU.Enabled {
		log.Printf("Registering CPU collector with interval %v", config.Collection.Collectors.CPU.Interval)
		collector := collectors.NewCPUCollector(config.Collection.Collectors.CPU.Interval)

		// Add debug logging
		if collector == nil {
			log.Printf("Warning: CPU collector was not created properly")
		}

		err := manager.RegisterCollector("cpu", collector)
		if err != nil {
			log.Printf("Failed to register CPU collector: %v", err)
			return nil, err
		}
		log.Printf("CPU collector registered successfully")
	}

	if config.Collection.Collectors.Memory.Enabled {
		log.Printf("Registering Memory collector with interval %v", config.Collection.Collectors.Memory.Interval)
		collector := collectors.NewMemoryCollector(config.Collection.Collectors.Memory.Interval)

		if collector == nil {
			log.Printf("Warning: Memory collector was not created properly")
		}

		err := manager.RegisterCollector("memory", collector)
		if err != nil {
			log.Printf("Failed to register Memory collector: %v", err)
			return nil, err
		}
		log.Printf("Memory collector registered successfully")
	}

	return manager, nil
}
