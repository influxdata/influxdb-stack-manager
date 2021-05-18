# InfluxDB Stack Manager

This tool is used to help manage influxdb stacks (dashboards, tasks etc).
It builds on top of the existing influx cli tool, transforming a single
template into multiple, more manageable, templates and extracting any flux
queries to make them easier to work with.

## Installation

This tool can be built using [go](https://golang.org/) by running:
`go install github.com/influxdata/influxdb-stack-manager`.

The influx cli tool will also need to be installed, and can be found
[here](https://github.com/influxdata/influxdb).


## Usage

To fetch the template belonging to a stack, run:

```
influxdb-stack-manager pull <stack-id>
```

Please be aware though, that this will destructively update the template
directory (which can be specified using the `--directory` or `-d` argument).

To apply any changes you've made to a stack, run:

```
influxdb-stack-manager push <stack-id>
```

Help can be found on by supplying an `-h` or `--help` argument to any command.

