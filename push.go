package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

const pushUsage = `Push local changes to a stack in influxdb

Usage:
  influxdb-stack-manager push <stack-id> [flags]

Flags:
`

func push(args []string) error {
	var cfg config
	var force bool
	var dataFile string

	fs := cfg.flagSet()
	fs.BoolVar(&force, "force", false, "TTY input, if template will have destructive changes, proceed if set true.")
	fs.StringVar(&dataFile, "data-file", "", "Data file to use for injected data in templates")
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

	tmpFile, err := writeTemplateToFile(cfg.directory, dataFile)
	if err != nil {
		return err
	}

	args = []string{"apply", "--stack-id", fs.Arg(0), "-f", tmpFile}
	args = append(args, cfg.generateArgs()...)
	if cfg.dryRun {
		log.Println("Dry run - calling:")
		log.Println(cfg.influxCmd, strings.Join(args, " "))
		log.Printf("Tempfile %q will not be removed automatically\n", tmpFile)
		return nil
	}
	defer os.Remove(tmpFile)

	cmd := exec.Command(cfg.influxCmd, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func writeTemplateToFile(dir, dataDir string) (string, error) {
	f, err := os.CreateTemp("", "*.yml")
	if err != nil {
		return "", fmt.Errorf("Error: unable to create temp file: %v", err)
	}
	defer f.Close()

	if err := uniteTemplate(dir, f, dataDir); err != nil {
		return "", fmt.Errorf("Error: unable to unite templates: %v", err)
	}

	return f.Name(), nil
}
