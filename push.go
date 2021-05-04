package main

import (
	"log"
	"os"
	"os/exec"
)

const pushUsage = `Push local changes to a stack in influxdb

Usage:
  influxdb-stack-manager push $STACK_ID [flags]

Flags:`

func push(args []string) {
	var cfg config
	fs := cfg.flagSet()

	var force bool
	fs.BoolVar(&force, "force", false, "TTY input, if template will have destructive changes, proceed if set true.")
	if err := fs.Parse(args); err != nil {
		log.Fatalf("Error: %v\nSee 'influxdb-stack-manager push -h' for help\n", err)
	}

	if cfg.help {
		log.Println(pushUsage + fs.FlagUsages())
		return
	}

	if fs.NArg() < 1 {
		log.Fatal("Error: Required arg missing: stack-id\nSee 'influxdb-stack-manager push -h' for help")
	}

	f, err := os.CreateTemp("", "*.yml")
	if err != nil {
		log.Fatalf("Unable to create temp file: %v", err)
	}
	defer f.Close()

	if err := uniteTemplate(cfg.directory, f); err != nil {
		log.Fatalf("Unable to unite templates: %v", err)
	}

	args = []string{"apply", "--stack-id", fs.Arg(0), "-f", f.Name()}
	args = append(args, cfg.generateArgs()...)
	cmd := exec.Command(cfg.influxCmd, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}
}
