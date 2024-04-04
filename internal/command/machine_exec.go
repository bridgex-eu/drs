package command

import (
	"github.com/bridgex-eu/drs/internal/remote"
	"golang.org/x/crypto/ssh"
)

type MachineClient struct {
	client *ssh.Client
}

func (c *MachineClient) Command(name string, args ...string) *remote.Cmd {
	exec := remote.SshExecutor(c.client)

	return remote.Command(exec, name, args...)
}

func (c *MachineClient) Close() error {
	return c.client.Close()
}
