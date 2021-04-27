package main

import (
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
		log.Println(usage)
		return
	}

	var err error
	switch args[0] {
	case "pull":
		err = pull(args[1:])

	case "push":
		err = push(args[1:])

	case "split":
		err = split(args[1:])

	case "unite":
		err = unite(args[1:])

	default:
		log.Println(usage)
	}

	if err != nil {
		log.Fatal(err)
	}
}
