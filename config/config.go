package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	configPath string
}

type BluePrintConfig struct {
	TemplatePath    string   `json:"template_path"`
	IgnoreFilePaths []string `json:"ignore_file_paths"`
}

var bluePrintConfig BluePrintConfig

func NewConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, ".config", "blueprint", "config.json")

	c := &Config{
		configPath: configPath,
	}

	if err := c.initializeConfig(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Config) GetConfig() (BluePrintConfig, error) {
	configFile, err := os.ReadFile(c.configPath)
	if err != nil {
		return BluePrintConfig{}, fmt.Errorf("failed to read the config file: %w", err)
	}

	if err = json.Unmarshal(configFile, &bluePrintConfig); err != nil {
		return BluePrintConfig{}, fmt.Errorf("failed to unmarshal the config file: %w", err)
	}

	if bluePrintConfig.TemplatePath == "" {
		return BluePrintConfig{}, fmt.Errorf("template path is not especified in config path: %s", c.configPath)
	}
	return bluePrintConfig, nil
}

func (c *Config) initializeConfig() error {
	if err := c.createConfigFolder(); err != nil {
		return err
	}
	return c.createConfigFile()
}

func (c *Config) createConfigFile() error {
	if c.filePathExists(c.configPath) {
		return nil
	}

	baseBluePrintConfigBytes, err := json.MarshalIndent(bluePrintConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal base config: %w", err)
	}

	err = os.WriteFile(c.configPath, baseBluePrintConfigBytes, 0600)
	if err != nil {
		return fmt.Errorf("failed to create the config file: %w", err)
	}

	return nil
}

func (c *Config) createConfigFolder() error {
	dir := filepath.Dir(c.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create the config folder: %w", err)
	}
	return nil
}

func (c *Config) filePathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
