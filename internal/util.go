package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// has checks if the supplied map contains the supplied key
func has(m map[string]interface{}, key string) bool {
	_, keyExists := m[key]
	return keyExists
}

// set sets the value of a nested key in the supplied map
func set(m map[string]interface{}, keys []string, value interface{}) (map[string]interface{}, error) {

	// If we're not adding interface{} more keys, return this map
	if len(keys) == 0 {
		return m, nil
	}

	key := keys[0]
	keyExists := has(m, key)

	// Last key, just write it in and return
	if len(keys) == 1 {

		// Key that we're trying to set already exists
		if keyExists {
			return m, fmt.Errorf("configuration option '%s' already present", key)
		}

		m[key] = value
		return m, nil
	}

	// Initialize a new map. We'll put this in the parent map if there isn't already a key there
	castValue := map[string]interface{}{}

	// The key already exists but if it's a map we can still go into it
	if keyExists {
		var castSuccessfully bool
		castValue, castSuccessfully = m[key].(map[string]interface{})

		if !castSuccessfully {
			return m, fmt.Errorf("configuration option '%s' already present and not a map", key)
		}
	}

	// Recurse and add the next nested value
	submap, err := set(castValue, keys[1:], value)
	if err != nil {
		return m, err
	}

	m[key] = submap
	return m, nil
}

// get gets the value of a nested key in the supplied map
func get(m map[string]interface{}, keys []string) (interface{}, error) {

	key := keys[0]
	keyExists := has(m, key)

	if !keyExists {
		return nil, fmt.Errorf("key '%s' was not found", key)
	}

	// If this is the last key, return the value
	if len(keys) == 1 {
		return m[key], nil
	}

	// Not the last key in the chain, make sure the next key is a map
	mapValue, castMapValue := m[key].(map[string]interface{})
	if !castMapValue {
		return nil, fmt.Errorf("key '%s' is not a map that can contain sub keys", key)
	}

	return get(mapValue, keys[1:])
}

// merge merges two maps recursively
func merge(map1 map[string]interface{}, map2 map[string]interface{}) map[string]interface{} {

	for key, value := range map2 {

		// If we don't have the key in map 1, just take the whole thing
		if !has(map1, key) {
			map1[key] = value
			continue
		}

		// We have the key in map 1 and map 2, let's see if it's a map in both so we can merge those
		map1Value, castMap1Value := map1[key].(map[string]interface{})
		map2Value, castMap2Value := map2[key].(map[string]interface{})

		// If we failed to cast one of these to a map then we can't merge them. Just ignore the key
		if !castMap1Value || !castMap2Value {
			continue
		}

		// Both of them are maps, keep merging
		map1[key] = merge(map1Value, map2Value)
	}

	return map1
}

// parseString parses a string into a variety of types
func parseString(value string) interface{} {

	// Check if it's an int
	intValue, err := strconv.ParseInt(value, 10, 0)
	if err == nil {
		return int(intValue) // Cast into an integer (parse int returns a int64)
	}

	// Check if it's a float
	floatValue, err := strconv.ParseFloat(value, 64)
	if err == nil {
		return floatValue
	}

	// Check if it's a bool
	boolValue, err := strconv.ParseBool(value)
	if err == nil {
		return boolValue
	}

	// Convert to bytes so we can do JSON checks
	bytes := []byte(value)

	// Check if it's a JSON object
	var jsonObject map[string]interface{}
	jsonObjectError := json.Unmarshal(bytes, &jsonObject)
	if jsonObjectError == nil {
		return jsonObject
	}

	// Check if it's a JSON array
	var jsonArray []interface{}
	jsonArrayError := json.Unmarshal(bytes, &jsonArray)
	if jsonArrayError == nil {
		return jsonArray
	}

	// Finally, try it as a duration
	return parseDurationString(value)
}

// parseDurationString attempts to parse a string into a duration, returning the original value if parsing failed
func parseDurationString(value string) interface{} {
	durationValue, err := time.ParseDuration(value)
	if err == nil {
		return durationValue
	}
	return value
}

// convertDurationStrings converts all the string durations in the supplied map to duration values
func convertDurationStrings(m map[string]interface{}) map[string]interface{} {
	for key, value := range m {

		// If the value is a string, apply duration parsing
		stringValue, castStringValue := value.(string)
		if castStringValue {
			m[key] = parseDurationString(stringValue)
			continue
		}

		// If the value is not a map, keep going
		mapValue, castMapValue := value.(map[string]interface{})
		if !castMapValue {
			continue
		}

		// Value is a map, recurse
		m[key] = convertDurationStrings(mapValue)
	}

	return m
}

// splitKey splits the supplied key into an array
func splitKey(key string) []string {
	return strings.Split(key, ":")
}

// cast casts the supplied object into the supplied type
func cast[T any](obj interface{}) (T, error) {
	value, success := obj.(T)
	if !success {
		return value, fmt.Errorf("failed to cast value to %T", *new(T))
	}
	return value, nil
}

// castIntegerSlice casts a value to an integer slice
func castIntegerSlice(obj interface{}) ([]int, error) {

	// Try a regular cast
	value, success := obj.([]int)
	if success {
		return value, nil
	}

	// Try to cast it into a float slice (the default when reading JSON)
	floatSlice, success := obj.([]float64)
	if !success {
		return nil, errors.New("failed to cast value to integer slice")
	}

	// Managed to convert to a float slice, convert to an integer slice
	intSlice := make([]int, len(floatSlice))
	for i, v := range floatSlice {
		intSlice[i] = int(v)

		if float64(intSlice[i]) != v {
			return nil, errors.New("failed to cast value to integer slice")
		}
	}

	return intSlice, nil
}

// castInteger casts a value to an integer
func castInteger(obj interface{}) (int, error) {

	// Try a regular integer cast
	value, success := obj.(int)
	if success {
		return value, nil
	}

	// Try to cast it into a float (the default when reading JSON) and then convert to int
	floatValue, cast := obj.(float64)
	if !cast {
		return 0, errors.New("failed to cast value to integer")
	}

	if float64(int(floatValue)) != floatValue {
		return 0, errors.New("failed to cast value to integer")
	}

	return int(floatValue), nil
}
