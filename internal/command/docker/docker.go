package docker

import (
	"github.com/bridgex-eu/drs/internal/command"
	"github.com/bridgex-eu/drs/internal/command/dockerhelper"
	"github.com/urfave/cli/v2"
)

func NewDockerCmd(drsCli *command.Cli) *cli.Command {
	return &cli.Command{
		Name:    "docker",
		Usage:   "Run Docker command on machine",
		Args:    true,
		Aliases: []string{"do"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "machine",
				Usage:    "Machine name",
				Aliases:  []string{"m"},
				Required: true,
			},
		},
		Action: func(ctx *cli.Context) error {
			machine := ctx.String("machine")
			args := ctx.Args()
			return runDocker(drsCli, machine, args.Slice()...)
		},
	}
}

func runDocker(drsCli *command.Cli, machine string, args ...string) error {
	client, err := drsCli.MachineClient(machine)
	if err != nil {
		return err
	}
	defer client.Close()

	return dockerhelper.ExecuteDockerCmd(client, drsCli.In, drsCli.Out, drsCli.Err, args...)
}
