package internal

import (
	"gopkg.in/yaml.v3"
	"os"
)

// YAMLFileLoader defines a loader that loads configurations from a YAML file
type YAMLFileLoader struct {
	filePath       string
	parseDurations bool
}

// NewYAMLFileLoader creates a new YAML file loader
func NewYAMLFileLoader(filePath string, parseDurations bool) *YAMLFileLoader {
	return &YAMLFileLoader{
		filePath:       filePath,
		parseDurations: parseDurations,
	}
}

// Load loads a YAML file
func (loader *YAMLFileLoader) Load() (map[string]interface{}, error) {
	file, err := os.ReadFile(loader.filePath)
	if err != nil {
		return map[string]interface{}{}, err
	}

	return loader.parseYAML(file)
}

// parseYAML parses yaml into a configuration map
func (loader *YAMLFileLoader) parseYAML(bytes []byte) (map[string]interface{}, error) {
	config := map[string]interface{}{}
	err := yaml.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}

	// If we were configured to parse durations, do that
	if loader.parseDurations {
		return convertDurationStrings(config), nil
	}

	return config, nil
}
