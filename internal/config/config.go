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
		CPU    CPUCollectorConfig  `yaml:"cpu"`
		Memory CollectorConfig     `yaml:"memory"`
		Disk   DiskCollectorConfig `yaml:"disk"`
	}
}

type CPUCollectorConfig struct {
	CollectorConfig `yaml:",inline"`
	IncludeTemps    bool `yaml:"include_temps"`
}

type DiskCollectorConfig struct {
	CollectorConfig `yaml:",inline"`
	IgnorePaths     []string `yaml:"ignore_paths"`
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

	if c.Collection.Collectors.Memory.Enabled {
		if err := validateInterval("Memory", c.Collection.Collectors.Memory.Interval); err != nil {
			return err
		}
	}

	if c.Collection.Collectors.Disk.Enabled {
		if err := validateInterval("Disk", c.Collection.Collectors.Disk.Interval); err != nil {
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
				CPU    CPUCollectorConfig  `yaml:"cpu"`
				Memory CollectorConfig     `yaml:"memory"`
				Disk   DiskCollectorConfig `yaml:"disk"`
			}{
				CPU: CPUCollectorConfig{
					CollectorConfig: CollectorConfig{
						Enabled:  false,
						Interval: time.Second * 60,
					},
					IncludeTemps: false,
				},
				Memory: CollectorConfig{
					Enabled:  false,
					Interval: time.Second * 10,
				},
				Disk: DiskCollectorConfig{
					CollectorConfig: CollectorConfig{
						Enabled:  false,
						Interval: time.Second * 60,
					},
					IgnorePaths: []string{"/proc", "/sys", "/dev"},
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
