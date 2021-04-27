package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
)

const pullUsage = `Pull a template from a stack in influx db and split it.

Warning: This is a destructive operation, the destination directory will be
cleared if it already exists.

Usage:
  influxdb-stack-manager pull <stack-id> [flags]

Flags:
`

func pull(args []string) error {
	var cfg config
	fs := cfg.flagSet()
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("Error: %v\nSee 'influxdb-stack-manager pull -h' for help", err)
	}

	if cfg.help {
		log.Println(pullUsage + fs.FlagUsages())
		return nil
	}

	if fs.NArg() < 1 {
		return errors.New("Error: required arg missing: stack-id\nSee 'influxdb-stack-manager pull -h' for help")
	}

	cmd := exec.Command(cfg.influxCmd, append([]string{"export", "stack", fs.Arg(0)}, cfg.generateArgs()...)...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return errors.New(out.String())
	}

	if err := splitTemplate(cfg.directory, &out); err != nil {
		return fmt.Errorf("Error: couldn't split template: %v\nPlease report this as an issue", err)
	}
	return nil
}
