package profile

import (
	"fmt"

	"github.com/bridgex-eu/drs/internal/command"
	"github.com/urfave/cli/v2"
)

func NewEmailCmd(drsCli *command.Cli) *cli.Command {
	return &cli.Command{
		Name:  "email",
		Usage: "Manage your email address",
		Subcommands: []*cli.Command{
			{
				Name:  "set",
				Usage: "Set email address",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "email",
						Usage:   "New email value",
						Aliases: []string{"e"},
					},
				},
				Action: func(cli *cli.Context) error {
					email := cli.String("email")
					return setEmail(drsCli, email)
				},
			},
		},
	}
}

func setEmail(drsCli *command.Cli, email string) error {
	if email == "" {
		var err error
		email, err = command.Prompt(drsCli.In, drsCli.Out, "Enter email address", "")
		if err != nil {
			return err
		}
	}

	if err := drsCli.Config.SetEmail(email); err != nil {
		return err
	}

	fmt.Fprintln(drsCli.Out, "Email changed.")
	return nil
}
