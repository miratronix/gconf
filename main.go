package gconf

import (
	"github.com/miratronix/gconf/internal"
	"sync"
)

var configSingleton *internal.Config
var once sync.Once

// New creates a new configuration structure
func New() *internal.Config {
	return internal.NewConfig()
}

// Instance returns a singleton configuration structure instance
func Instance() *internal.Config {
	once.Do(func() {
		configSingleton = internal.NewConfig()
	})
	return configSingleton
}

// Arguments creates a new command line argument loader
func Arguments(separator string, prefix string) *internal.ArgumentLoader {
	return internal.NewArgumentLoader(separator, prefix)
}

// Environment creates a new environment variable loader
func Environment(lowerCase bool, separator string, prefix string) *internal.EnvironmentLoader {
	return internal.NewEnvironmentLoader(lowerCase, separator, prefix)
}

// JSONFile creates a new JSON file loader
func JSONFile(filePath string, parseDurations bool) *internal.JSONFileLoader {
	return internal.NewJSONFileLoader(filePath, parseDurations)
}

// YAMLFile creates a new YAML file loader
func YAMLFile(filePath string, parseDurations bool) *internal.YAMLFileLoader {
	return internal.NewYAMLFileLoader(filePath, parseDurations)
}

// Map creates a new map laoder
func Map(stringMap map[string]interface{}) *internal.MapLoader {
	return internal.NewMapLoader(stringMap)
}
