# InfluxDB Stack Manager

The [influx cli tool](https://github.com/influxdata/influx-cli) provides
the ability to manage dashboards, tasks etc. in a collection called a "stack".
It provides the ability to export a template for all of these objects, and
apply a modified version to update our stack.

This tool builds on top of the influx cli to make working with that template
file a better, and easier experience. It transforms the single large yaml file,
with inlined flux code like this:

![Full template](/screenshots/full-template.png?raw=true)


Into a well-organised set of individual templates, with the flux code extracted
so we can work with it easily, like this:

![Separated templates](/screenshots/separated-templates.png?raw=true)



## Installation

This tool can be built using [go](https://golang.org/) by running:
`go install github.com/influxdata/influxdb-stack-manager`.

The influx cli tool will also need to be installed, and can be found
[here](https://github.com/influxdata/influxdb).


## Basic Usage

To first create a stack to manage, you will need to use the influx cli
tool itself. To create a stack, run:

```bash
influx stacks init -n MyStackName -d MyStackDescription
```

To add resources to the stack, you will need to know their IDs which
can be found either from the influxdb UI, or by listing them, and then
updating the stack with the new resources:

```bash
influx dashboards
influx task list

influx stacks update --stack-id $STACK_ID \
    --addResource=Dashboard=$DASHBOARD_ID
    --addResource=Task=$TASK_ID
```

With the stack created, you can use the influxdb-stack-manager to actually
fetch your templates and push up any changes.

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


## Templating

If you would like to use the same templates for multiple stacks (in the same or
different influxdb clusters), you may want to inject data into the templates or
flux code. To do this:

### 1. Add fields to your templates/flux queries

The templates and flux queries can be amended to use injected fields using the
[golang templating syntax](https://pkg.go.dev/text/template). For example, if
I wanted to inject a set of thresholds which I would alert on, I could inject
the field `{{ .Thresholds.CPU }}`.


### 2. Create a data file with the values to inject

The data to be injected can be specified either in json or yaml format, for example:

```json
{
  "Thresholds": {
    "CPU": 75.0,
    "Memory": 60.0
  }
}
```

```yaml
Thresholds:
  CPU: 75.0
  Memory: 60.0
```

### 3. Specify the data file when pushing

When pushing, add the `--data-file` flag:

```
influxdb-stack-manager push <stack-id> --data-file "data/cluster-1.yml"
```


## TODO

 - [ ] Provide release binaries
 - [ ] Provide docker images
 - [ ] Allow running queries with injected data

