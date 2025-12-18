package config

import (
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v2"
)

func Bind(v interface{}, path string) error {
	filename, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	cleanedDst := filepath.Clean(filename)
	yamlFile, err := os.ReadFile(cleanedDst)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(yamlFile, v); err != nil {
		return err
	}
	return nil
}
