package config

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"

	"gopkg.in/yaml.v3"
)

const (
	drs_CONFIG_DIR = "drs_CONFIG_DIR"
)

func ConfigDir() (string, error) {
	var path string
	if a := os.Getenv(drs_CONFIG_DIR); a != "" {
		path = a
	} else {
		b, err := os.UserConfigDir()
		if err != nil {
			return "", fmt.Errorf("Failed to retrieve config dir path: %w", err)
		}

		path = filepath.Join(b, "drs")
	}

	if !dirExists(path) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return "", fmt.Errorf("Failed to create config dir: %w", err)
		}

	}

	return path, nil
}

func dirExists(path string) bool {
	f, err := os.Stat(path)
	return err == nil && f.IsDir()
}

func fileExists(path string) bool {
	f, err := os.Stat(path)
	return err == nil && !f.IsDir()
}

func ConfigFile() (string, error) {
	path, err := ConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(path, "default.yml"), nil
}

func ParseDefaultConfig() (Config, error) {
	path, err := ConfigFile()
	if err != nil {
		return nil, err
	}

	return parseConfig(path)
}

func readConfigFile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, pathError(err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func WriteConfigFile(filename string, data []byte) error {
	err := os.MkdirAll(filepath.Dir(filename), 0771)
	if err != nil {
		return pathError(err)
	}

	cfgFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600) // cargo coded from setup
	if err != nil {
		return err
	}
	defer cfgFile.Close()

	_, err = cfgFile.Write(data)
	return err
}

func parseConfigFile(filename string) ([]byte, *yaml.Node, error) {
	data, err := readConfigFile(filename)
	if err != nil {
		return nil, nil, err
	}

	root, err := parseConfigData(data)
	if err != nil {
		return nil, nil, err
	}
	return data, root, err
}

func parseConfigData(data []byte) (*yaml.Node, error) {
	var root yaml.Node
	err := yaml.Unmarshal(data, &root)
	if err != nil {
		return nil, err
	}

	if len(root.Content) == 0 {
		return &yaml.Node{
			Kind:    yaml.DocumentNode,
			Content: []*yaml.Node{{Kind: yaml.MappingNode}},
		}, nil
	}
	if root.Content[0].Kind != yaml.MappingNode {
		return &root, fmt.Errorf("Expected a top level map")
	}
	return &root, nil
}

func parseConfig(filename string) (Config, error) {
	_, root, err := parseConfigFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			root = DefaultConfigRoot()
		} else {
			return nil, fmt.Errorf("Failed to parse config file: %w", err)
		}
	}

	return NewConfig(root), nil
}

func pathError(err error) error {
	var pathError *os.PathError
	if errors.As(err, &pathError) && errors.Is(pathError.Err, syscall.ENOTDIR) {
		if p := findRegularFile(pathError.Path); p != "" {
			return fmt.Errorf("remove or rename regular file `%s` (must be a directory)", p)
		}

	}
	return err
}

func findRegularFile(p string) string {
	for {
		if s, err := os.Stat(p); err == nil && s.Mode().IsRegular() {
			return p
		}
		newPath := filepath.Dir(p)
		if newPath == p || newPath == "/" || newPath == "." {
			break
		}
		p = newPath
	}
	return ""
}
