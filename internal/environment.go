package internal

import (
	"os"
	"strings"
)

// EnvironmentLoader defines a loader that loads configurations from environment variables
type EnvironmentLoader struct {
	lowerCase bool
	prefix    string
	separator string
}

// NewEnvironmentLoader creates a new environment loader
func NewEnvironmentLoader(lowerCase bool, separator string, prefix string) *EnvironmentLoader {
	return &EnvironmentLoader{
		lowerCase: lowerCase,
		prefix:    prefix,
		separator: separator,
	}
}

// Load loads environment variables
func (loader *EnvironmentLoader) Load() (map[string]interface{}, error) {
	return loader.parseEnvironment(os.Environ())
}

// parseEnvironment parses environment variables into a configuration map
func (loader *EnvironmentLoader) parseEnvironment(environmentData []string) (map[string]interface{}, error) {
	config := map[string]interface{}{}

	for _, environmentLine := range environmentData {

		// Split the env entry on =
		keyValue := strings.SplitN(environmentLine, "=", 2)

		// If there was no equals, ignore this line
		if keyValue == nil || len(keyValue) < 2 {
			continue
		}

		// If we have a configured prefix and the key doesn't match it, ignore this line
		if len(loader.prefix) > 0 && !strings.HasPrefix(keyValue[0], loader.prefix) {
			continue
		}

		// Trim the prefix off the key and trim the separator if it's there as a prefix (that would result in an empty key)
		trimmedKey := strings.TrimPrefix(keyValue[0], loader.prefix)
		trimmedKey = strings.TrimPrefix(trimmedKey, loader.separator)

		// Ignore keys that are empty after trimming
		if len(trimmedKey) == 0 {
			continue
		}

		// Lowercase the key if that option is enabled
		if loader.lowerCase {
			trimmedKey = strings.ToLower(trimmedKey)
		}

		// Separate it on the separator if required
		separatedKeys := []string{trimmedKey}
		if len(loader.separator) > 0 {
			separatedKeys = strings.Split(trimmedKey, loader.separator)
		}

		// Set the nested value in the map
		value := parseString(keyValue[1])
		_, err := set(config, separatedKeys, value)
		if err != nil {
			return config, err
		}
	}

	return config, nil
}
