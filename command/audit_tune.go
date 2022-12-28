package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*AuditTuneCommand)(nil)
	_ cli.CommandAutocomplete = (*AuditTuneCommand)(nil)
)

type AuditTuneCommand struct {
	*BaseCommand

	flagDescription string
	flagOptions     map[string]string
}

func (c *AuditTuneCommand) Synopsis() string {
	return "Tunes an audit device configuration"
}

func (c *AuditTuneCommand) Help() string {
	helpText := `
Usage: vault audit tune [options] PATH

  Tunes the configuration options for the audit device at the given PATH. The
  argument corresponds to the PATH where the audit device is enabled, not the
  TYPE!

  Tune an option of the audit method mounted as file/:

      $ vault audit tune -options=log_raw=true file/

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *AuditTuneCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)

	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:   flagNameDescription,
		Target: &c.flagDescription,
		Usage: "Human-friendly description of this audit device. This overrides " +
			"the current stored value, if any.",
	})

	f.StringMapVar(&StringMapVar{
		Name:       "options",
		Target:     &c.flagOptions,
		Completion: complete.PredictAnything,
		Usage: "Key-value pair provided as key=value for the audit device options. " +
			"This can be specified multiple times.",
	})

	return set
}

func (c *AuditTuneCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultAudits()
}

func (c *AuditTuneCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *AuditTuneCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	switch {
	case len(args) < 1:
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected 1, got %d)", len(args)))
		return 1
	case len(args) > 1:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 1, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	auditConfigInput := api.AuditConfigInput{
		Options: c.flagOptions,
	}

	// Set these values only if they are provided in the CLI
	f.Visit(func(fl *flag.Flag) {
		if fl.Name == flagNameDescription {
			auditConfigInput.Description = &c.flagDescription
		}
	})

	// Append a trailing slash to indicate it's a path in the output
	mountPath := ensureTrailingSlash(sanitizePath(args[0]))

	if err := client.Sys().TuneAudit(mountPath, auditConfigInput); err != nil {
		c.UI.Error(fmt.Sprintf("Error tuning audit device %s: %s", mountPath, err))
		return 2
	}

	c.UI.Output(fmt.Sprintf("Success! Tuned the audit device at: %s", mountPath))
	return 0
}
