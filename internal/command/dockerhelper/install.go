package dockerhelper

import (
	"errors"
	"fmt"

	"github.com/bridgex-eu/drs/internal/command"
)

// MachineWithDocker initialize the machine client and verify that all necessary components are installed
func MachineWithDocker(drsCli *command.Cli, machine string) (*command.MachineClient, error) {
	client, err := drsCli.MachineClient(machine)
	if err != nil {
		return nil, err
	}

	if !isDockerInstalled(client) {
		fmt.Fprintln(drsCli.Out, "Missing docker on machine. Installing...")

		if !isSuperUser(client) {
			return nil, errors.New("Docker installation failed: superuser access is required.")
		}

		if err := installDocker(client); err != nil {
			return nil, fmt.Errorf("Docker installation failed: %w", err)
		}

		if err := CreateDockerNetwork(client, drsNetwork); err != nil {
			return nil, fmt.Errorf("Failed to create Docker network: %w", err)
		}
	}

	if !IsNetworkExist(client, drsNetwork) {
		fmt.Fprintln(drsCli.Out, "Missing docker network. Creating...")

		if err := CreateDockerNetwork(client, drsNetwork); err != nil {
			return nil, fmt.Errorf("Failed to create Docker network: %w", err)
		}
	}

	if !IsContainerExist(client, traefikName) {
		fmt.Fprintln(drsCli.Out, "Traefik is missing. Attempting to start it...")
		if err := runTraefik(client, drsCli); err != nil {
			return nil, fmt.Errorf("Unable to start Traefik: %w", err)
		}
	}

	return client, err
}
