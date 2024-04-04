package config

import (
	"gopkg.in/yaml.v3"
)

type configEntry struct {
	KeyNode   *yaml.Node
	ValueNode *yaml.Node
}

// findEntry searches for and returns the first entry with the specified key.
// If no entry is found, it returns nil.
func findEntry(node *yaml.Node, key string) *configEntry {
	entries := intoEntries(node)

	for _, entry := range entries {
		if entry.KeyNode.Value == key {
			return entry
		}
	}

	return nil
}

func removeEntry(node *yaml.Node, key string) {
	if node == nil {
		return
	}

	newContent := []*yaml.Node{}

	content := node.Content
	for i := 0; i < len(content); i++ {
		if content[i].Value == key {
			i++ // skip the next node which is this key's value
		} else {
			newContent = append(newContent, content[i])
		}
	}

	node.Content = newContent
}

func intoEntries(node *yaml.Node) []*configEntry {
	var entries []*configEntry

	if node.Kind != yaml.MappingNode {
		// If the node is not a mapping node, return an empty slice
		return entries
	}

	// Content slice goes [key1, value1, key2, value2, ...]
	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		valueNode := node.Content[i+1]
		entry := &configEntry{
			KeyNode:   keyNode,
			ValueNode: valueNode,
		}
		entries = append(entries, entry)
	}

	return entries
}
