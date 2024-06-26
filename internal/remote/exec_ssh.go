package remote

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type sshExecutor struct {
	client  *ssh.Client
	session *ssh.Session
}

func SshExecutor(client *ssh.Client) Executor {
	return &sshExecutor{
		client: client,
	}
}

func (e *sshExecutor) Addr() net.Addr {
	return e.client.RemoteAddr()
}

func (e *sshExecutor) Start(cmd string, in io.Reader, out, stderr io.Writer) error {
	if e.session != nil {
		return errors.New("Failed to start cmd: command already stared")
	}
	session, err := e.client.NewSession()
	if err != nil {
		return fmt.Errorf("Failed to create SSH session: %w", err)
	}
	session.Stdin = in
	session.Stdout = out
	session.Stderr = stderr

	e.session = session
	if err := e.session.Start(cmd); err != nil {
		return fmt.Errorf("Failed to start SSH session: %w", err)
	}

	return nil
}

func (e *sshExecutor) Wait() error {
	if e.session == nil {
		return errors.New("Failed to wait command: command not started")
	}

	if err := e.session.Wait(); err != nil {
		if exitErr, ok := err.(*ssh.ExitError); ok {
			return &ExitError{
				Status:  exitErr.ExitStatus(),
				Content: exitErr.String(),
			}
		}

		return fmt.Errorf("Failed to wait SSH session: %w", err)
	}

	return nil
}

func (e *sshExecutor) Close() error {
	if e.session == nil {
		return errors.New("Failed to wait command: command not started")
	}

	return e.Close()
}

func SshClient(host, user, private, passphrase string) (*ssh.Client, error) {
	var sshAuth ssh.AuthMethod
	var err error

	if private != "" {
		sshAuth, err = authorizeWithKey(private, passphrase)
	} else {
		sshAuth, err = authorizeWithSSHAgent()
	}
	if err != nil {
		return nil, err
	}

	// Set up SSH client configuration
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			sshAuth,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 5,
	}

	return ssh.Dial("tcp", host+":22", config)
}

func authorizeWithKey(key, passphrase string) (ssh.AuthMethod, error) {
	var signer ssh.Signer
	var err error

	if passphrase != "" {
		signer, err = ssh.ParsePrivateKeyWithPassphrase([]byte(key), []byte(passphrase))
	} else {
		signer, err = ssh.ParsePrivateKey([]byte(key))
	}
	if err != nil {
		return nil, err
	}

	return ssh.PublicKeys(signer), nil
}

func authorizeWithSSHAgent() (ssh.AuthMethod, error) {
	conn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to ssh-agent: %w", err)
	}
	defer conn.Close()

	sshAgent := agent.NewClient(conn)
	return ssh.PublicKeysCallback(sshAgent.Signers), nil
}
