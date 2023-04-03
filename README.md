# boss_apu_exporter

BlueOcean APU exporter for prometheus, using `boss-smi` binary to gather metrics.

## Introduction

This is a simple exporter that uses `boss-smi(.exe)` binary to collect, parse and export metrics.
This makes it possible to run it on Windows and get GPU metrics while gaming - no Docker or Linux required.

This project is based on [utkuozdemir/nvidia_gpu_exporter](https://github.com/utkuozdemir/nvidia_gpu_exporter).

**If you are a gamer who's into monitoring, you are in for a treat.**

## Highlights

- Will work on any system that has `boss-smi(.exe)?` binary - Windows, Linux, MacOS... No C bindings required
- Doesn't even need to run the monitored machine: can be configured to execute `boss-smi` command remotely
- No need for a Docker or Kubernetes environment
- Auto-discovery of the metric fields `boss-smi` can expose (future-compatible)
- Comes with its own [Grafana dashboard](https://grafana.com/grafana/dashboards/14574)

## Visualization

You can use the official [Grafana dashboard](https://grafana.com/grafana/dashboards/14574)
to see your GPU metrics in a nicely visualized way.

Here's how it looks like:
![Grafana dashboard](https://github.com/dijiujun/boss_apu_exporter/blob/master/grafana/dashboard.png)

## Installation

See [INSTALL.md](INSTALL.md) for details.

## Configuration

See [CONFIGURE.md](CONFIGURE.md) for details.

## Metrics

See [METRICS.md](METRICS.md) for details.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for details.
