package key

import (
	"fmt"

	"github.com/bridgex-eu/drs/internal/command"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"golang.design/x/clipboard"
)

func NewPrivateCmd(drsCli *command.Cli) *cli.Command {
	return &cli.Command{
		Name:      "private copy",
		Usage:     "Copy private key to clipboard",
		Args:      true,
		ArgsUsage: "KEY",
		Action: func(ctx *cli.Context) error {
			name := ctx.Args().First()
			return runPrivateCopy(drsCli, name)
		},
	}
}

func runPrivateCopy(drsCli *command.Cli, name string) error {
	if name == "" {
		return errors.New("Name cannot be empty")
	}

	key, err := drsCli.Config.Keys().GetByName(name)
	if err != nil {
		return fmt.Errorf("Failed to retrieve key: %w", err)
	}

	err = clipboard.Init()
	if err != nil {
		return fmt.Errorf("Failed to copy key to clipboard: %w", err)
	}

	clipboard.Write(clipboard.FmtText, []byte(key.Private))

	fmt.Fprintln(drsCli.Out, "Private key copied to your clipboard.")

	return nil
}
