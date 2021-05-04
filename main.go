package main

import (
	"fmt"
	"log"
	"os"
)

const usage = `InfluxDB Stack Manager

This tool is designed to make managing stack templates simpler.

Usage:
  influxdb-stack-manager [command] [flags]

Available Commands:
  pull		Fetch a stack template from influxdb and split it.
  push		Apply templates changes to a stack in influxdb.
  split		Split a local template file.
  unite		Unite a set of parsed templates to generate a local template file.

Flags:
  -h,		Help for the influx command
`

func main() {
	log.SetFlags(0)
	run(os.Args[1:])
}

func run(args []string) {
	if len(args) == 0 {
		fmt.Println(usage)
		return
	}

	switch args[1] {
	case "pull":
		pull(args[2:])

	case "push":
		push(args[2:])

	case "split":
		split(args[2:])

	case "unite":
		unite(args[2:])

	default:
		fmt.Println(usage)
	}

}
