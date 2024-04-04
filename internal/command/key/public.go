package key

import (
	"fmt"

	"github.com/bridgex-eu/drs/internal/command"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"golang.design/x/clipboard"
)

func NewPublicCmd(drsCli *command.Cli) *cli.Command {
	return &cli.Command{
		Name:      "public copy",
		Usage:     "Copy public key to clipboard",
		Args:      true,
		ArgsUsage: "KEY",
		Action: func(ctx *cli.Context) error {
			name := ctx.Args().First()
			return runPublicCopy(drsCli, name)
		},
	}
}

func runPublicCopy(drsCli *command.Cli, name string) error {
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

	clipboard.Write(clipboard.FmtText, []byte(key.Public))

	fmt.Fprintln(drsCli.Out, "Public key copied to your clipboard.")

	return nil
}
