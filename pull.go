package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
)

const pullUsage = `Pull a template from a stack in influx db and split it.

Usage:
  influxdb-stack-manager pull $STACK_ID [flags]

Flags:
`

func pull(args []string) {
	var cfg config
	fs := cfg.flagSet()
	if err := fs.Parse(args); err != nil {
		log.Fatalf("Error: %v\nSee 'influxdb-stack-manager pull -h' for help", err)
	}

	if cfg.help {
		log.Println(pullUsage + fs.FlagUsages())
		return
	}

	if fs.NArg() < 1 {
		log.Fatalf("Error: Required arg missing: stack-id\nSee 'influxdb-stack-manager pull -h' for help")
	}

	cmd := exec.Command(cfg.influxCmd, append([]string{"export", "stack", fs.Arg(0)}, cfg.generateArgs()...)...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(out.String())
	}

	if err := splitTemplate(cfg.directory, &out); err != nil {
		log.Fatalf("Error splitting template: %v\nPlease report this as an issue.", err)
	}
}
