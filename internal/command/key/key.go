package key

import (
	"github.com/bridgex-eu/drs/internal/command"
	"github.com/urfave/cli/v2"
)

func NewCmd(drsCli *command.Cli) *cli.Command {
	return &cli.Command{
		Name:  "key",
		Usage: "Manage SSH keys",
		Subcommands: []*cli.Command{
			NewAddCmd(drsCli),
			NewListCmd(drsCli),
			NewRmCmd(drsCli),
			NewGenerateCmd(drsCli),
			NewPrivateCmd(drsCli),
			NewPublicCmd(drsCli),
		},
	}
}
