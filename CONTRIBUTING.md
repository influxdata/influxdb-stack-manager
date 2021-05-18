# Contributing to influxdb-stack-manager

## Bug reports

The influxdb-stack-manager works as a thin layer over the top on the influx cli tool. Before filing an issue
here, please make sure that it is an issue with the stack manager rather than the influx cli tool. This tool
is responsible for:
  * Transforming a template outputted by the influx cli command into multiple templates/flux files.
  * Constructing a single template from the separate templates/flux files.
  * Calling the influx cli tool apprpriately.

If you think there is an issue with the template transformation logic, please include a minimal set of files
that illustrate the problem, along with what your expected output should be. The `influxdb-stack-manager split`
and `influxdb-stack-manager unite` commands should help generate these.

If you think there is a problem with how the tool is calling the influx cli, please use the `--dry-run` option
to generate what is being called, and include that along with the flags you are supplying to the stack-manager.

The easier it is for us to reproduce the problem, the easier it is for us to fix it.
If you have never written a bug report before, or if you want to brush up on your bug reporting skills, we recommend reading [Simon Tatham's essay "How to Report Bugs Effectively."](http://www.chiark.greenend.org.uk/~sgtatham/bugs.html)


## Feature requests
We really like to receive feature requests as it helps us prioritize our work.
Please be clear about your requirements and goals, help us to understand what you would like to see added to the stack-manager
with examples and the reasons why it is important to you.
If you find your feature request already exists as a Github issue please indicate your support for that feature by using the "thumbs up" reaction.


## Submitting a pull request
To submit a pull request you should fork the stack-manager repository, and make your change on a feature branch of your fork.
Then generate a pull request from your branch against *main* of the stack-manager repository.
Include in your pull request details of your change -- the why *and* the how -- as well as the testing you performed.
Also, be sure to run the test suite with your change in place.
Changes that cause tests to fail cannot be merged.

There will usually be some back and forth as we finalize the change, but once that completes it may be merged.


## Security Vulnerability Reporting
InfluxData takes security and our users' trust very seriously.
If you believe you have found a security issue in any of our open source projects, please responsibly disclose it by contacting security@influxdata.com.
More details about security vulnerability reporting, including our GPG key, [can be found here](https://www.influxdata.com/how-to-report-security-vulnerabilities/).


## Signing the CLA

If you are going to be contributing back to InfluxDB please take a second to sign our CLA, which can be found [on our website](https://influxdata.com/community/cla/).


## Building from Source

### Installing Go

influxdb-stack-manager requires Go 1.16.

### Building

The stack-manager command can be built by running `go build .` from the root of the repository.

### Testing

To run tests, run `go test .` from the root of the repository.
