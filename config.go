package main

import (
	"github.com/spf13/pflag"
)

// A config holds all the flags that have been supplied to the stack command.
type config struct {
	// Flags to pass to influx command
	activeConfig string
	configsPath  string
	host         string
	org          string
	orgID        string
	skipVerify   bool
	token        string

	// Flags used internally
	help      bool
	directory string
	influxCmd string
}

// flagSet generates a flagSet to use in parsing the flags.
func (cfg *config) flagSet() *pflag.FlagSet {
	fs := pflag.NewFlagSet("config", pflag.ContinueOnError)
	fs.StringVarP(&cfg.activeConfig, "active-config", "c", "", "Config name to use for the command. Maps to env var $INFLUX_ACTIVE_CONFIG.")
	fs.StringVar(&cfg.configsPath, "configs-path", "", "Path to the influx cli configurations. Maps to env var $INFLUX_CONFIGS_PATH.")
	fs.StringVar(&cfg.host, "host", "", "HTTP address of InfluxDB. Maps to env var $INFLUX_HOST.")
	fs.StringVarP(&cfg.org, "org", "o", "", "The name of the organization. Maps to env var $INFLUX_ORG.")
	fs.StringVar(&cfg.orgID, "org-id", "", "The ID of the organization. Maps to env var $INFLUX_ORG_ID.")
	fs.BoolVar(&cfg.skipVerify, "skip-verify", false, "Skip TLS certificate chain and host name verification.")
	fs.StringVarP(&cfg.token, "token", "t", "", "Authentication token. Maps to env var $INFLUX_TOKEN.")
	fs.BoolVarP(&cfg.help, "help", "h", false, "Display help for this command.")
	fs.StringVarP(&cfg.directory, "directory", "d", "templates", "Directory to read and write templates from/to.")
	fs.StringVar(&cfg.influxCmd, "influx-cmd", "influx", "Command to call the influx cli, if it not in your path.")
	return fs
}

func (cfg config) generateArgs() []string {
	var args []string
	for name, value := range map[string]string{
		"active-config": cfg.activeConfig,
		"configs-path":  cfg.configsPath,
		"host":          cfg.host,
		"org":           cfg.org,
		"org-id":        cfg.orgID,
		"token":         cfg.token,
	} {
		if value != "" {
			args = append(args, "--"+name, value)
		}
	}

	if cfg.skipVerify {
		args = append(args, "--skip-verify")
	}

	return args
}
