# Metrics Agent

Metrics Agent is a lightweight, open-source system monitoring agent written in Go. It collects key metrics from the host system and exposes them for easy consumption by monitoring and observability platforms.

## Features

- Collects CPU, memory, and disk usage metrics
- Configurable collection intervals per metric type  
- Batches metrics to reduce network overhead
- Exposes metrics in a structured format for easy ingestion
- Lightweight and efficient, with minimal impact on host resources

## Roadmap

- [ ] Add support for collecting network metrics
- [ ] Provide pre-built binaries for major operating systems
- [ ] Develop a plugin system for easy extensibility 
- [ ] Integrate with popular monitoring systems out-of-the-box (e.g. Prometheus, InfluxDB)
- [ ] Enhance configuration options and reloading
- [ ] Improve test coverage and performance benchmarking
- [ ] Create comprehensive documentation and getting started guides

## Getting Started

### Prerequisites

- Go 1.23.2 or higher

### Installation

1. Clone the repository:
   ```
   git clone https://github.com/Bitlatte/metrics-agent.git
   ```

2. Build the agent:
   ```
   cd metrics-agent
   go build -o metrics-agent ./cmd/agent
   ```

3. Configure the agent by editing the `agent.yaml` file.

4. Run the agent:
   ```
   ./metrics-agent
   ```

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.