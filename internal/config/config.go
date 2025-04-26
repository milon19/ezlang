package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type FileConfig struct {
	Path string `yaml:"path"`
	Lang string `yaml:"lang"`
}

type Config struct {
	Files []FileConfig `yaml:"files"`
}

func Load(path string) (*Config, error) {
	if path == "" {
		path = ".ezlang.yaml"
	}
	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("error reading config file: %v\n", err)
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		fmt.Printf("Error unmarshalling YAML data: %v\n", err)
		return nil, err
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

func (c *Config) Validate() error {
	if len(c.Files) == 0 {
		return fmt.Errorf("no files specified in configuration")
	}
	for _, file := range c.Files {
		if file.Path == "" {
			return fmt.Errorf("file path is empty")
		}

		if file.Lang == "" {
			return fmt.Errorf("language is empty for file %s", file.Path)
		}

		if _, err := os.Stat(file.Path); os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", file.Path)
		}
	}
	return nil
}
