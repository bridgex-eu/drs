package key

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/bridgex-eu/drs/internal/command"
	"github.com/bridgex-eu/drs/internal/config"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/ssh"
)

func NewAddCmd(drsCli *command.Cli) *cli.Command {
	return &cli.Command{
		Name:      "add",
		Usage:     "Add a new ssh key from file",
		Args:      true,
		ArgsUsage: "FILE_PATH",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "name",
				Usage: "Key name",
			},
			&cli.StringFlag{
				Name:    "passphrase",
				Usage:   "Key passphrase",
				Aliases: []string{"p"},
			},
		},
		Action: func(ctx *cli.Context) error {
			path := ctx.Args().First()
			name := ctx.String("name")
			passphrase := ctx.String("passphrase")

			return runAdd(drsCli, path, name, passphrase)
		},
	}
}

func runAdd(drsCli *command.Cli, path, name, passphrase string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Failed to read file: %w", err)
	}

	var privateKey ssh.Signer

	if passphrase == "" {
		privateKey, err = ssh.ParsePrivateKey(data)
	} else {
		privateKey, err = ssh.ParsePrivateKeyWithPassphrase(data, []byte(passphrase))
	}

	if err != nil {
		return fmt.Errorf("Failed to parse ssh key: %w", err)
	}

	publicKey := ssh.MarshalAuthorizedKey(privateKey.PublicKey())

	if name == "" {
		name = command.GenerateRandomName()
	}

	_, err = drsCli.Config.Keys().GetByName(name)
	if err == nil {
		return errors.New("Key with this name arleady exist")
	}

	if err := drsCli.Config.Keys().Add(
		config.KeyEntry{
			Name:       name,
			Private:    string(data),
			Public:     string(publicKey),
			Passphrase: passphrase,
			CreatedAt:  time.Now(),
		},
	); err != nil {
		return err
	}

	fmt.Fprintf(drsCli.Out, "Key %s added.\n", name)

	return nil
}
