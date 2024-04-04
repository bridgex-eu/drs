package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Config interface {
	Write() error
	Machines() *MachinesConfig
	Keys() *KeysConfig
	Email() string
	SetEmail(email string) error
}

func DefaultConfigRoot() *yaml.Node {
	return &yaml.Node{
		Kind: yaml.DocumentNode,
		Content: []*yaml.Node{
			{
				Kind: yaml.MappingNode,
				Content: []*yaml.Node{
					{
						Kind:  yaml.ScalarNode,
						Value: "version",
					},
					{
						Kind:  yaml.ScalarNode,
						Value: "1",
					},
				},
			},
			{
				Kind: yaml.MappingNode,
				Content: []*yaml.Node{
					{
						HeadComment: "List of machines to deploy.",
						Kind:        yaml.ScalarNode,
						Value:       "machines",
					},
					{
						Kind:  yaml.ScalarNode,
						Value: "",
					},
				},
			},
			{
				Kind: yaml.MappingNode,
				Content: []*yaml.Node{
					{
						HeadComment: "List of your SSH keys.",
						Kind:        yaml.ScalarNode,
						Value:       "keys",
					},
					{
						Kind:  yaml.ScalarNode,
						Value: "",
					},
				},
			},
			{
				Kind: yaml.MappingNode,
				Content: []*yaml.Node{
					{
						HeadComment: "List of your SSH keys.",
						Kind:        yaml.ScalarNode,
						Value:       "email",
					},
					{
						Kind:  yaml.ScalarNode,
						Value: "",
					},
				},
			},
		},
	}
}

type fileConfig struct {
	documentNode *yaml.Node
}

func NewConfig(documentNode *yaml.Node) Config {
	return &fileConfig{
		documentNode: documentNode,
	}
}

func (c *fileConfig) Write() error {
	data, err := yaml.Marshal(c.documentNode)
	if err != nil {
		return err
	}

	filename, err := ConfigFile()
	if err != nil {
		return err
	}

	if err := WriteConfigFile(filename, data); err != nil {
		return fmt.Errorf("Failed to save config file: %w", err)
	}

	return nil
}

func (c *fileConfig) root() *yaml.Node {
	return c.documentNode.Content[0]
}

func (c *fileConfig) Machines() *MachinesConfig {
	machines := findEntry(c.root(), "machines")

	if machines == nil {
		keyNode := &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: "machines",
		}
		valueNode := &yaml.Node{
			Kind:    yaml.MappingNode,
			Content: []*yaml.Node{},
		}

		c.root().Content = append(c.root().Content, keyNode, valueNode)

		machines = &configEntry{
			KeyNode:   keyNode,
			ValueNode: valueNode,
		}
	}

	return &MachinesConfig{
		Parent: c,
		Node:   machines.ValueNode,
	}
}

func (c *fileConfig) Keys() *KeysConfig {
	keys := findEntry(c.root(), "keys")

	if keys == nil {
		keyNode := &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: "keys",
		}
		valueNode := &yaml.Node{
			Kind:    yaml.MappingNode,
			Content: []*yaml.Node{},
		}

		c.root().Content = append(c.root().Content, keyNode, valueNode)

		keys = &configEntry{
			KeyNode:   keyNode,
			ValueNode: valueNode,
		}
	}

	return &KeysConfig{
		Parent: c,
		Node:   keys.ValueNode,
	}
}

func (c *fileConfig) getEmailNode() *configEntry {
	email := findEntry(c.root(), "email")

	if email == nil {
		keyNode := &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: "email",
		}
		valueNode := &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: "",
		}

		c.root().Content = append(c.root().Content, keyNode, valueNode)

		email = &configEntry{
			KeyNode:   keyNode,
			ValueNode: valueNode,
		}
	}

	return email
}

func (c *fileConfig) Email() string {
	email := c.getEmailNode()

	return email.ValueNode.Value
}

func (c *fileConfig) SetEmail(email string) error {
	emailNode := c.getEmailNode()
	emailNode.ValueNode.Value = email

	return c.Write()
}
