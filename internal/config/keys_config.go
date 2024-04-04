package config

import (
	"errors"
	"slices"
	"time"

	"gopkg.in/yaml.v3"
)

type KeysConfig struct {
	Node   *yaml.Node
	Parent Config
}

type KeyEntry struct {
	Name       string
	Private    string
	Public     string
	Passphrase string
	CreatedAt  time.Time
}

func (k KeyEntry) toNode() *yaml.Node {
	return &yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "private"},
			{Kind: yaml.ScalarNode, Value: k.Private},
			{Kind: yaml.ScalarNode, Value: "public"},
			{Kind: yaml.ScalarNode, Value: k.Public},
			{Kind: yaml.ScalarNode, Value: "passphrase"},
			{Kind: yaml.ScalarNode, Value: k.Passphrase},
			{Kind: yaml.ScalarNode, Value: "created"},
			{Kind: yaml.ScalarNode, Value: k.CreatedAt.Format(time.RFC3339)},
		},
	}
}

func mapKeyEntry(entry *configEntry) (KeyEntry, error) {
	key := KeyEntry{Name: entry.KeyNode.Value}

	node := entry.ValueNode
	private := findEntry(node, "private")
	if private == nil {
		return key, errors.New("No private key found")
	}
	key.Private = private.ValueNode.Value

	public := findEntry(node, "public")
	if public == nil {
		return key, errors.New("No public key found")
	}
	key.Public = public.ValueNode.Value

	passphrase := findEntry(node, "passphrase")
	if passphrase != nil {
		key.Passphrase = passphrase.ValueNode.Value
	}

	created := findEntry(node, "created")
	if created == nil {
		return key, errors.New("No create date found")
	}
	createdAt, err := time.Parse(time.RFC3339, created.ValueNode.Value)
	if err != nil {
		return key, errors.New("Created date has wrong format")
	}
	key.CreatedAt = createdAt

	return key, nil
}

func (c *KeysConfig) Add(entry KeyEntry) error {
	keyNode := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: entry.Name,
	}

	if c.Node.Kind != yaml.MappingNode {
		c.Node.Kind = yaml.MappingNode
	}

	c.Node.Content = append(c.Node.Content, keyNode, entry.toNode())

	return c.Parent.Write()
}

func (c *KeysConfig) Remove(name string) error {
	entry := findEntry(c.Node, name)
	if entry == nil {
		return errors.New("Not found")
	}

	c.Node.Content = slices.DeleteFunc(c.Node.Content, func(x *yaml.Node) bool {
		return x.Line == entry.KeyNode.Line || x.Line == entry.ValueNode.Line
	})

	return c.Parent.Write()
}

func (c *KeysConfig) All() []KeyEntry {
	var keys []KeyEntry

	if c.Node.Kind == yaml.MappingNode {
		keysEntries := intoEntries(c.Node)

		for _, keyEntry := range keysEntries {
			key, err := mapKeyEntry(keyEntry)
			if err != nil {
				continue
			}

			keys = append(keys, key)
		}
	}

	return keys
}

func (c *KeysConfig) GetByName(name string) (KeyEntry, error) {
	entry := findEntry(c.Node, name)
	if entry == nil {
		return KeyEntry{}, errors.New("Not found")
	}

	return mapKeyEntry(entry)
}
