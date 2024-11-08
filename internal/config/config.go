package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Collection CollectionConfig `yaml:"collection"`
}

type CollectionConfig struct {
	BatchSize  uint8         `yaml:"batch_size"`
	Interval   time.Duration `yaml:"interval"`
	Collectors struct {
		CPU    CPUCollectorConfig `yaml:"cpu"`
		Memory CollectorConfig    `yaml:"memory"`
	}
}

type CPUCollectorConfig struct {
	CollectorConfig `yaml:",inline"`
	IncludeTemps    bool `yaml:"include_temps"`
}

type CollectorConfig struct {
	Enabled  bool          `yaml:"enabled"`
	Interval time.Duration `yaml:"interval"`
}

func (c *Config) Validate() error {
	if c.Collection.BatchSize <= 0 || c.Collection.BatchSize > 255 {
		return fmt.Errorf("batch size must be between 1 and 255")
	}

	if c.Collection.Interval <= 0 {
		return fmt.Errorf("collection interval must be positive")
	}

	// Validate collector intervals
	validateInterval := func(name string, interval time.Duration) error {
		if interval <= 0 {
			return fmt.Errorf("%s collector interval must be positive", name)
		}
		return nil
	}

	if c.Collection.Collectors.CPU.Enabled {
		if err := validateInterval("CPU", c.Collection.Collectors.CPU.Interval); err != nil {
			return err
		}
	}

	return nil
}

func Load() (*Config, error) {
	config := &Config{
		Collection: CollectionConfig{
			BatchSize: 100,
			Interval:  time.Second * 15,
			Collectors: struct {
				CPU    CPUCollectorConfig `yaml:"cpu"`
				Memory CollectorConfig    `yaml:"memory"`
			}{
				CPU: CPUCollectorConfig{
					CollectorConfig: CollectorConfig{
						Enabled:  true,
						Interval: time.Second * 60,
					},
					IncludeTemps: true,
				},
				Memory: CollectorConfig{
					Enabled:  true,
					Interval: time.Second * 10,
				},
			},
		},
	}

	defaultPaths := []string{
		"agent.yaml",
		os.Getenv("METRICS_AGENT_CONFIG"),
	}

	for _, path := range defaultPaths {
		if path == "" {
			continue
		}

		if data, err := os.ReadFile(path); err == nil {
			if err := yaml.Unmarshal(data, config); err != nil {
				return nil, err
			}
			break
		}
	}

	// Validate the configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %v", err)
	}

	return config, nil
}
