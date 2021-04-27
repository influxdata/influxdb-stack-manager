package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
)

const pushUsage = `Push local changes to a stack in influxdb

Usage:
  influxdb-stack-manager push <stack-id> [flags]

Flags:`

func push(args []string) error {
	var cfg config
	var force bool

	fs := cfg.flagSet()
	fs.BoolVar(&force, "force", false, "TTY input, if template will have destructive changes, proceed if set true.")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("Error: %v\nSee 'influxdb-stack-manager push -h' for help", err)
	}

	if cfg.help {
		log.Println(pushUsage + fs.FlagUsages())
		return nil
	}

	if fs.NArg() < 1 {
		return errors.New("Error: required arg missing: stack-id\nSee 'influxdb-stack-manager push -h' for help")
	}

	f, err := os.CreateTemp("", "*.yml")
	if err != nil {
		return fmt.Errorf("Error: unable to create temp file: %v", err)
	}
	defer f.Close()

	if err := uniteTemplate(cfg.directory, f); err != nil {
		return fmt.Errorf("Error: unable to unite templates: %v", err)
	}

	args = []string{"apply", "--stack-id", fs.Arg(0), "-f", f.Name()}
	args = append(args, cfg.generateArgs()...)
	cmd := exec.Command(cfg.influxCmd, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
