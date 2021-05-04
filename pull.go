package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

const pullUsage = `Pull a template from a stack in influx db and split it.

Usage:
  influxdb-stack-manager pull $STACK_ID [flags]

Flags:
%v

All other flags will be passed directly to the influx cli command.
`

func pull(args []string) {
	cfg := config{}
	fs := cfg.flagSet()
	if err := fs.Parse(args); err != nil {
		fmt.Printf("Error: %v\nSee 'influxdb-stack-manager pull -h' for help\n", err)
		return
	}

	if cfg.help {
		fmt.Println(pullUsage + fs.FlagUsages())
		return
	}

	if fs.NArg() < 1 {
		fmt.Println("Error: Required arg missing: stack-id\nSee 'influxdb-stack-manager pull -h' for help")
		return
	}

	cmd := exec.Command(cfg.influxCmd, append([]string{"export", "stack", fs.Arg(0)}, cfg.generateArgs()...)...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println(out.String())
		return
	}

	if err := splitTemplate(cfg.directory, &out); err != nil {
		fmt.Println("error splitting template")
	}
}
