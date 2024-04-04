package machine

import (
	"github.com/bridgex-eu/drs/internal/command"
	"github.com/urfave/cli/v2"
)

func NewCmd(drsCli *command.Cli) *cli.Command {
	return &cli.Command{
		Name:  "machine",
		Usage: "Manage your machines",
		Subcommands: []*cli.Command{
			NewAddCmd(drsCli),
			NewListCmd(drsCli),
			NewRmCmd(drsCli),
			NewExecuteCmd(drsCli),
		},
	}
}
