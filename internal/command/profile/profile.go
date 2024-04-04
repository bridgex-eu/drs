package profile

import (
	"github.com/bridgex-eu/drs/internal/command"
	"github.com/urfave/cli/v2"
)

func NewCmd(drsCli *command.Cli) *cli.Command {
	return &cli.Command{
		Name:  "profile",
		Usage: "Manage your profile",
		Subcommands: []*cli.Command{
			NewEmailCmd(drsCli),
		},
	}
}
