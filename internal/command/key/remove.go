package key

import (
	"errors"
	"fmt"

	"github.com/bridgex-eu/drs/internal/command"
	"github.com/urfave/cli/v2"
)

func NewRmCmd(drsCli *command.Cli) *cli.Command {
	return &cli.Command{
		Name:      "rm",
		Usage:     "Remove a key",
		Args:      true,
		ArgsUsage: "NAME",
		Action: func(ctx *cli.Context) error {
			name := ctx.Args().First()
			return runRemove(drsCli, name)
		},
	}
}

func runRemove(drsCli *command.Cli, name string) error {
	if name == "" {
		return errors.New("Name cannot be empty")
	}

	if err := drsCli.Config.Keys().Remove(name); err != nil {
		return fmt.Errorf("Failed to remove key: %w", err)
	}

	fmt.Fprintln(drsCli.Out, "Key removed from this computer.")

	return nil
}
