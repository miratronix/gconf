package gconf

import (
	"gconf/lib"
	"sync"
)

var configSingleton *lib.Config
var once sync.Once

// New creates a new configuration structure
func New() *lib.Config {
	return lib.NewConfig()
}

// Instance returns a singleton configuration structure instance
func Instance() *lib.Config {
	once.Do(func() {
		configSingleton = lib.NewConfig()
	})
	return configSingleton
}

// Arguments creates a new command line argument loader
func Arguments(separator string, prefix string) *lib.ArgumentLoader {
	return lib.NewArgumentLoader(separator, prefix)
}

// Environment creates a new environment variable loader
func Environment(lowerCase bool, separator string, prefix string) *lib.EnvironmentLoader {
	return lib.NewEnvironmentLoader(lowerCase, separator, prefix)
}

// JSONFile creates a new JSON file loader
func JSONFile(filePath string, parseDurations bool) *lib.JSONFileLoader {
	return lib.NewJSONFileLoader(filePath, parseDurations)
}

// YAMLFile creates a new YAML file loader
func YAMLFile(filePath string, parseDurations bool) *lib.YAMLFileLoader {
	return lib.NewYAMLFileLoader(filePath, parseDurations)
}

// Map creates a new map laoder
func Map(stringMap map[string]interface{}) *lib.MapLoader {
	return lib.NewMapLoader(stringMap)
}
