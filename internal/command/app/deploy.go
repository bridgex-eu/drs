package app

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bridgex-eu/drs/internal/command"
	"github.com/bridgex-eu/drs/internal/command/dockerhelper"
	"github.com/urfave/cli/v2"
)

func NewDeployCmd(drsCli *command.Cli) *cli.Command {
	return &cli.Command{
		Name:      "deploy",
		Usage:     "Deploy a Docker image on a machine",
		Args:      true,
		ArgsUsage: "IMAGE [ARGS]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "name",
				Usage: "Container name",
			},
			&cli.StringFlag{
				Name:     "machine",
				Usage:    "Machine name",
				Aliases:  []string{"m"},
				Required: true,
			},
			&cli.StringSliceFlag{
				Name:    "port",
				Usage:   "Port mapping [domain]:machine_port:container_port",
				Aliases: []string{"p"},
			},
			&cli.StringSliceFlag{
				Name:    "volume",
				Usage:   "Volume mapping",
				Aliases: []string{"v"},
			},
			&cli.StringSliceFlag{
				Name:    "env",
				Usage:   "Environment variable",
				Aliases: []string{"e"},
			},
			&cli.BoolFlag{
				Name:    "yes",
				Usage:   "Automatic yes to prompts; assume \"yes\" as answer to all prompts and run non-interactively.",
				Aliases: []string{"y", "assume-yes"},
			},
			&cli.BoolFlag{
				Name:    "remote",
				Usage:   "Pull the Docker image directly onto a remote machine.",
				Aliases: []string{"r"},
			},
		},
		Action: func(ctx *cli.Context) error {
			name := ctx.String("name")
			machine := ctx.String("machine")
			ports := ctx.StringSlice("port")
			volumes := ctx.StringSlice("volume")
			envs := ctx.StringSlice("env")
			yes := ctx.Bool("yes")
			remote := ctx.Bool("remote")

			image := ctx.Args().First()
			otherArgs := ctx.Args().Tail()

			return runDeploy(drsCli, image, name, machine, ports, volumes, envs, yes, remote, otherArgs...)
		},
	}
}

func runDeploy(drsCli *command.Cli, image, name, machine string, ports, volumes, envs []string, yes, remote bool, args ...string) error {
	if image == "" {
		return errors.New("Image cannot be an empty string.")
	}

	if name == "" {
		name = command.GenerateRandomName()
	}

	client, err := dockerhelper.MachineWithDocker(drsCli, machine)
	if err != nil {
		return err
	}
	defer client.Close()

	if !remote {
		size, err := dockerhelper.LocalImageSize(image)
		if err != nil {
			return err
		}

		const maxSize = 100 * 1024 * 1024 // 100 MB in bytes
		if size > maxSize {
			fmt.Fprintf(drsCli.Out, "!! Your Docker image is over 100 MB (it's %.2f MB). You might find it faster to deploy if you use the '--remote' flag, allowing you to pull the image directly to a remote machine.\n", float64(size)/float64(1024*1024))
		}

		fmt.Fprintln(drsCli.Out, "Sending the Docker image...")
		if err := dockerhelper.SendDockerImage(client, image); err != nil {
			return fmt.Errorf("Failed to send the Docker image: %w", err)
		}
	}

	var traefikLabels []string
	var adaptedPorts []string
	for _, port := range ports {
		dockerPort, labels, err := parsePortFlag(port)
		if err != nil {
			return fmt.Errorf("Failed to parse port flag '%s': %w", port, err)
		}
		adaptedPorts = append(adaptedPorts, dockerPort)
		traefikLabels = append(traefikLabels, labels...)
	}

	var old string

	if dockerhelper.IsContainerExist(client, name) {
		if !yes {
			if confirm, err := command.PromptForConfirmation(drsCli.In, drsCli.Out, "An app with this name already exists. Do you want to replace it?"); err != nil || !confirm {
				return err
			}
		}

		fmt.Fprintln(drsCli.Out, "Stopping the existing app...")
		if err := dockerhelper.StopContainer(client, name); err != nil {
			return fmt.Errorf("Failed to stop old container: %w", err)
		}

		old = name + "-old"

		if err := dockerhelper.RenameContainer(client, name, old); err != nil {
			return fmt.Errorf("Failed to rename old container: %w", err)
		}
	}

	fmt.Fprintln(drsCli.Out, "Running a new app...")
	if err := dockerhelper.RunContainer(client, image, name, adaptedPorts, volumes, traefikLabels, envs, args...); err != nil {
		fmt.Fprintln(drsCli.Out, "Failed to run the new container, attempting to restart the previous one.")

		if startErr := dockerhelper.StartContainer(client, old); startErr != nil {
			return fmt.Errorf("Failed to restart old container: %w", startErr)
		}

		fmt.Fprintln(drsCli.Out, "Container started.")

		if renameErr := dockerhelper.RenameContainer(client, old, name); renameErr != nil {
			return fmt.Errorf("Failed to rename old container back: %w", renameErr)
		}

		return fmt.Errorf("Failed to run image: %w", err)
	}

	if old != "" {
		if err := dockerhelper.RemoveContainer(client, old, false); err != nil {
			return fmt.Errorf("Failed to remove old container: %w", err)
		}
	}

	fmt.Fprintf(drsCli.Out, "The app %s now runs on your machine.\n", name)

	return nil
}

// Parse [domain.com]:serverport:containerport
func parsePortFlag(portFlag string) (dockerPort string, traefikLabels []string, err error) {
	parts := strings.Split(portFlag, ":")

	if len(parts) < 2 || len(parts) > 3 {
		err = fmt.Errorf("invalid port flag format")
		return
	}

	if len(parts) == 3 {
		domain := parts[0]
		hostPort := parts[1]
		containerPort := parts[2]

		// 127.0.0.1 makes container accessible only from internal network
		dockerPort = "127.0.0.1:" + hostPort + ":" + containerPort
		traefikLabels = dockerhelper.GetTraefikRouteLabel(domain, containerPort)
	} else {
		dockerPort = "127.0.0.1:" + parts[0] + ":" + parts[1]
		traefikLabels = nil
	}

	return
}
