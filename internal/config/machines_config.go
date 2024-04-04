package config

import (
	"errors"
	"net"
	"slices"
	"time"

	"gopkg.in/yaml.v3"
)

type MachinesConfig struct {
	Node   *yaml.Node
	Parent Config
}

type MachineEntry struct {
	Name      string
	Host      net.IP
	User      string
	Key       string
	CreatedAt time.Time
}

func (m MachineEntry) toNode() *yaml.Node {
	return &yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "host"},
			{Kind: yaml.ScalarNode, Value: m.Host.String()},
			{Kind: yaml.ScalarNode, Value: "user"},
			{Kind: yaml.ScalarNode, Value: m.User},
			{Kind: yaml.ScalarNode, Value: "key"},
			{Kind: yaml.ScalarNode, Value: m.Key},
			{Kind: yaml.ScalarNode, Value: "created"},
			{Kind: yaml.ScalarNode, Value: m.CreatedAt.Format(time.RFC3339)},
		},
	}
}

func (m *MachineEntry) fromNode(node *yaml.Node) error {
	host := findEntry(node, "host")
	if host == nil {
		return errors.New("No host found")
	}
	m.Host = net.ParseIP(host.ValueNode.Value)

	user := findEntry(node, "user")
	if user == nil {
		return errors.New("No ssh user found")
	}
	m.User = user.ValueNode.Value

	key := findEntry(node, "key")
	if key != nil {
		m.Key = key.ValueNode.Value
	}

	created := findEntry(node, "created")
	if created == nil {
		return errors.New("No create date found")
	}
	createdAt, err := time.Parse(time.RFC3339, created.ValueNode.Value)
	if err != nil {
		return errors.New("Created date has wrong format")
	}
	m.CreatedAt = createdAt

	return nil
}

func (c *MachinesConfig) Add(entry MachineEntry) error {
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

func (c *MachinesConfig) Remove(nameOrHost string) error {
	machine := c.GetByNameOrHost(nameOrHost)
	if machine == nil {
		return errors.New("Not found")
	}

	entry := findEntry(c.Node, machine.Name)
	if entry == nil {
		return errors.New("Config entry not found")
	}

	c.Node.Content = slices.DeleteFunc(c.Node.Content, func(x *yaml.Node) bool {
		return x.Line == entry.KeyNode.Line || x.Line == entry.ValueNode.Line
	})

	return c.Parent.Write()
}

func (c *MachinesConfig) All() []MachineEntry {
	var machines []MachineEntry

	if c.Node.Kind == yaml.MappingNode {
		machineEntries := intoEntries(c.Node)

		for _, machineNode := range machineEntries {
			machineEntry := MachineEntry{Name: machineNode.KeyNode.Value}
			machineEntry.fromNode(machineNode.ValueNode)
			machines = append(machines, machineEntry)
		}
	}

	return machines
}

func (c *MachinesConfig) GetByNameOrHost(nameOrHost string) *MachineEntry {
	if c.Node.Kind != yaml.MappingNode {
		return nil
	}

	// try to find by name
	entry := findEntry(c.Node, nameOrHost)
	if entry != nil {
		machineEntry := MachineEntry{Name: entry.KeyNode.Value}
		machineEntry.fromNode(entry.ValueNode)
		return &machineEntry
	}

	// try to find by host
	for _, m := range c.All() {
		if m.Host.String() == nameOrHost {
			return &m
		}
	}

	return nil
}
