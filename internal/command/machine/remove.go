package machine

import (
	"errors"
	"fmt"

	"github.com/bridgex-eu/drs/internal/command"
	"github.com/urfave/cli/v2"
)

func NewRmCmd(drsCli *command.Cli) *cli.Command {
	return &cli.Command{
		Name:      "rm",
		Usage:     "Remove a machine",
		Args:      true,
		ArgsUsage: "MACHINE",
		Action: func(ctx *cli.Context) error {
			machine := ctx.Args().First()
			return runRemove(drsCli, machine)
		},
	}
}

func runRemove(drsCli *command.Cli, machine string) error {
	if machine == "" {
		return errors.New("Name cannot be empty")
	}

	if err := drsCli.Config.Machines().Remove(machine); err != nil {
		return fmt.Errorf("Failed to remove machine: %w", err)
	}

	fmt.Fprintln(drsCli.Out, "Machine removed from this computer.")

	return nil
}
