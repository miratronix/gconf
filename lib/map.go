package lib

// MapLoader defines a loader that loads configurations from a map
type MapLoader struct {
	values map[string]interface{}
}

// NewMapLoader creates a new map loader
func NewMapLoader(stringMap map[string]interface{}) *MapLoader {
	return &MapLoader{
		values: stringMap,
	}
}

// Load returns the underlying map
func (loader *MapLoader) Load() (map[string]interface{}, error) {
	return loader.values, nil
}
