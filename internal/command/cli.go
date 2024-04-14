package command

import (
	"errors"
	"fmt"
	"io"

	"github.com/bridgex-eu/drs/internal/config"
	"github.com/bridgex-eu/drs/internal/remote"
)

type Cli struct {
	Config   config.Config
	Out, Err io.Writer
	In       io.ReadCloser
}

// MachineClient provides a client to interact with a machine
func (c *Cli) MachineClient(machine string) (*MachineClient, error) {
	m := c.Config.Machines().GetOrDefault(machine)
	if m == nil {
		return nil, errors.New("Machine with this name or host not found")
	}

	var private, passphrase string
	if m.Key != "" {
		key, err := c.Config.Keys().GetByName(m.Key)
		if err != nil {
			return nil, fmt.Errorf("Failed to retrieve key: %w", err)
		}

		private = key.Private
		passphrase = key.Passphrase
	}

	client, err := remote.SshClient(m.Host.String(), m.User, private, passphrase)
	if err != nil {
		return nil, fmt.Errorf("Failed to create ssh client: %w", err)
	}

	return &MachineClient{client: client}, nil
}
