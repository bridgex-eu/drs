package machine

import (
	"errors"
	"fmt"

	"github.com/bridgex-eu/drs/internal/command"
	"github.com/bridgex-eu/drs/internal/remote"
	"github.com/urfave/cli/v2"
)

func NewExecuteCmd(drsCli *command.Cli) *cli.Command {
	return &cli.Command{
		Name:        "exec",
		Description: "Execute shell command on machine",
		Args:        true,
		ArgsUsage:   "MACHINE \"COMMAND\"",
		Action: func(ctx *cli.Context) error {
			machine := ctx.Args().First()
			cmd := ctx.Args().Get(1)

			return runExecute(drsCli, machine, cmd)
		},
	}
}

func runExecute(drsCli *command.Cli, machine string, cmd string) error {
	if machine == "" {
		return errors.New("Name is required")
	}

	client, err := drsCli.MachineClient(machine)
	if err != nil {
		return err
	}
	defer client.Close()

	output, err := client.Command(cmd).CombinedOutput()
	if err != nil {
		if _, ok := err.(*remote.ExitError); !ok {
			return err
		}
	}

	fmt.Fprintln(drsCli.Out, output)

	return nil
}
