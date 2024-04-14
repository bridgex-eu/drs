package machine

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/bridgex-eu/drs/internal/command"
	"github.com/bridgex-eu/drs/internal/config"
	"github.com/urfave/cli/v2"
)

func NewAddCmd(drsCli *command.Cli) *cli.Command {
	return &cli.Command{
		Name:        "add",
		Usage:       "Add a new machine",
		Description: "This command adds the machine to the drs. Ensure your computer has SSH access to this machine.",
		Args:        true,
		ArgsUsage:   "HOST",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "name",
				Usage: "Machine name",
			},
			&cli.StringFlag{
				Name:    "user",
				Aliases: []string{"u"},
				Value:   "root",
				Usage:   "SSH user name",
			},
			&cli.StringFlag{
				Name:    "key",
				Aliases: []string{"k"},
				Usage:   "SSH key name",
			},
		},
		Action: func(ctx *cli.Context) error {
			host := ctx.Args().First()
			name := ctx.String("name")
			user := ctx.String("user")
			key := ctx.String("key")

			return runAdd(drsCli, host, name, user, key)
		},
	}
}

func runAdd(drsCli *command.Cli, host, name, user, key string) error {
	machines := drsCli.Config.Machines()

	if name == "" {
		name = command.GenerateRandomName()
	}

	if key != "" {
		_, err := drsCli.Config.Keys().GetByName(key)
		if err != nil {
			return fmt.Errorf("Failed to retrieve machine: %w", err)
		}
	}

	if host == "" {
		return errors.New("Host cannot be empty")
	}

	if machines.GetOrDefault(host) != nil || machines.GetOrDefault(name) != nil {
		return errors.New("Machine with this host or name arleady exist")
	}

	hostIp := net.ParseIP(host)
	if hostIp == nil {
		return errors.New("Host must be valid ip address")
	}

	if err := machines.Add(config.MachineEntry{
		Name:      name,
		Host:      hostIp,
		User:      user,
		CreatedAt: time.Now(),
	}); err != nil {
		return fmt.Errorf("Failed to create machine: %w", err)
	}

	fmt.Fprintf(drsCli.Out, "Machine %s added.\n", name)

	return nil
}
