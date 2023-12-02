package lib

import (
	"encoding/json"
	"os"
)

// JSONFileLoader defines a loader that loads configurations from a JSON file
type JSONFileLoader struct {
	filePath       string
	parseDurations bool
}

// NewJSONFileLoader creates a new JSON file loader
func NewJSONFileLoader(filePath string, parseDurations bool) *JSONFileLoader {
	return &JSONFileLoader{
		filePath:       filePath,
		parseDurations: parseDurations,
	}
}

// Load loads a JSON file
func (loader *JSONFileLoader) Load() (map[string]interface{}, error) {
	file, err := os.ReadFile(loader.filePath)
	if err != nil {
		return map[string]interface{}{}, err
	}

	return loader.parseJSON(file)
}

// parseJSON parses json into a configuration map
func (loader *JSONFileLoader) parseJSON(bytes []byte) (map[string]interface{}, error) {
	config := map[string]interface{}{}
	err := json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}

	// If we were configured to parse durations, do that
	if loader.parseDurations {
		return convertDurationStrings(config), nil
	}

	return config, nil
}
