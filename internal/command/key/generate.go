package key

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/bridgex-eu/drs/internal/command"
	"github.com/bridgex-eu/drs/internal/config"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/ssh"
)

func NewGenerateCmd(drsCli *command.Cli) *cli.Command {
	return &cli.Command{
		Name:  "generate",
		Usage: "Generate a new SSH key",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "name",
				Usage: "Key name",
			},
		},
		Action: func(ctx *cli.Context) error {
			name := ctx.String("name")

			return runGenerate(drsCli, name)
		},
	}
}

func runGenerate(drsCli *command.Cli, name string) error {
	if name == "" {
		name = command.GenerateRandomName()
	}

	private, public, err := generateSSHKeyPair()
	if err != nil {
		return fmt.Errorf("Failed to generate ssh key: %w", err)
	}

	_, err = drsCli.Config.Keys().GetByName(name)
	if err == nil {
		return errors.New("Key with this name arleady exist")
	}

	if err = drsCli.Config.Keys().Add(
		config.KeyEntry{
			Name:      name,
			Private:   string(private),
			Public:    string(public),
			CreatedAt: time.Now(),
		},
	); err != nil {
		return err
	}

	fmt.Fprintf(drsCli.Out, "Key %s added.\n", name)

	return nil
}

func generateSSHKeyPair() (privateKey, publicKey string, err error) {
	privateRSAKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}

	// Generate the private key PEM block.
	privatePEMBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateRSAKey)}

	// Encode the private key to PEM format.
	privateKeyBytes := pem.EncodeToMemory(privatePEMBlock)
	privateKey = string(privateKeyBytes)

	// Generate the public key for the private key.
	publicRSAKey, err := ssh.NewPublicKey(&privateRSAKey.PublicKey)
	if err != nil {
		return "", "", err
	}

	// Encode the public key to the authorized_keys format.
	publicKeyBytes := ssh.MarshalAuthorizedKey(publicRSAKey)
	publicKey = string(publicKeyBytes)

	return privateKey, publicKey, nil
}
