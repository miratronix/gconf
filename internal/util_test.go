package internal

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHas(t *testing.T) {

	Convey("Returns false if the map doesn't have the key", t, func() {
		result := has(map[string]interface{}{}, "something")
		So(result, ShouldBeFalse)
	})

	Convey("Returns true if the map has the key", t, func() {
		result := has(map[string]interface{}{"something": "stuff"}, "something")
		So(result, ShouldBeTrue)
	})
}

func TestSet(t *testing.T) {

	Convey("Returns the input map when no keys are specified", t, func() {
		result, err := set(map[string]interface{}{}, []string{}, nil)
		So(result, ShouldBeEmpty)
		So(err, ShouldBeNil)
	})

	Convey("Sets a non-nested key to the specified value", t, func() {
		result, err := set(map[string]interface{}{}, []string{"a"}, "testing")
		So(result, ShouldResemble, map[string]interface{}{"a": "testing"})
		So(err, ShouldBeNil)
	})

	Convey("Sets a nested key to the specified value", t, func() {
		result, err := set(map[string]interface{}{}, []string{"a", "b", "c"}, "testing")
		So(result, ShouldResemble, map[string]interface{}{"a": map[string]interface{}{"b": map[string]interface{}{"c": "testing"}}})
		So(err, ShouldBeNil)
	})

	Convey("Returns an error when a non-nested key is already present", t, func() {
		result, err := set(map[string]interface{}{"a": true}, []string{"a"}, nil)
		So(result, ShouldResemble, map[string]interface{}{"a": true})
		So(err, ShouldNotBeNil)
	})

	Convey("Returns an error when a nested key is already present", t, func() {
		result, err := set(map[string]interface{}{"a": map[string]interface{}{"b": true}}, []string{"a", "b"}, nil)
		So(result, ShouldResemble, map[string]interface{}{"a": map[string]interface{}{"b": true}})
		So(err, ShouldNotBeNil)
	})

	Convey("Returns an error when a nested key is already present as a nested key", t, func() {
		result, err := set(map[string]interface{}{"a": map[string]interface{}{"b": true}}, []string{"a", "b", "c"}, nil)
		So(result, ShouldResemble, map[string]interface{}{"a": map[string]interface{}{"b": true}})
		So(err, ShouldNotBeNil)
	})
}

func TestGet(t *testing.T) {

	Convey("Gets a non-nested key", t, func() {
		result, err := get(map[string]interface{}{"string": "woohoo"}, []string{"string"})
		So(result, ShouldEqual, "woohoo")
		So(err, ShouldBeNil)
	})

	Convey("Gets a nested key", t, func() {
		result, err := get(map[string]interface{}{"map": map[string]interface{}{"string": "woohoo"}}, []string{"map", "string"})
		So(result, ShouldEqual, "woohoo")
		So(err, ShouldBeNil)
	})

	Convey("Returns an error when a non-existent non-nested key is requested", t, func() {
		result, err := get(map[string]interface{}{}, []string{"non-existent"})
		So(result, ShouldBeNil)
		So(err, ShouldNotBeNil)
	})

	Convey("Returns an error when a non-existent nested key is requested", t, func() {
		result, err := get(map[string]interface{}{"string": "woohoo"}, []string{"string", "subString"})
		So(result, ShouldBeNil)
		So(err, ShouldNotBeNil)
	})
}

func TestMerge(t *testing.T) {

	Convey("Merges non-nested keys", t, func() {
		result := merge(map[string]interface{}{"one": 1}, map[string]interface{}{"two": 2})
		So(result, ShouldResemble, map[string]interface{}{"one": 1, "two": 2})
	})

	Convey("Doesn't override existing keys", t, func() {
		result := merge(map[string]interface{}{"one": 1}, map[string]interface{}{"one": 2})
		So(result, ShouldResemble, map[string]interface{}{"one": 1})
	})

	Convey("Merges nested keys", t, func() {
		result := merge(map[string]interface{}{"one": map[string]interface{}{"one": 1}}, map[string]interface{}{"one": map[string]interface{}{"two": 2}})
		So(result, ShouldResemble, map[string]interface{}{"one": map[string]interface{}{"one": 1, "two": 2}})
	})

	Convey("Doesn't override existing nested keys", t, func() {
		result := merge(map[string]interface{}{"one": map[string]interface{}{"one": 1}}, map[string]interface{}{"one": map[string]interface{}{"one": 2}})
		So(result, ShouldResemble, map[string]interface{}{"one": map[string]interface{}{"one": 1}})
	})
}

func TestParseString(t *testing.T) {

	Convey("Parses booleans", t, func() {

		Convey("Parses 'true' into true", func() {
			result := parseString("true")
			So(result, ShouldEqual, true)
		})

		Convey("Parses 'false' into false", func() {
			result := parseString("false")
			So(result, ShouldEqual, false)
		})
	})

	Convey("Parses numbers", t, func() {

		Convey("Parses integers", func() {
			result := parseString("10")
			So(result, ShouldEqual, 10)
		})

		Convey("Parses floats", func() {
			result := parseString("10.5")
			So(result, ShouldEqual, 10.5)
		})

		Convey("Parses floats that end with .0", func() {
			result := parseString("10.0")
			So(result, ShouldEqual, 10.0)
		})
	})

	Convey("Parses JSON arrays", t, func() {
		result := parseString("[1,2,3]")
		So(result, ShouldResemble, []interface{}{float64(1), float64(2), float64(3)})
	})

	Convey("Parses JSON objects", t, func() {
		result := parseString(`{"a": 1, "b":2}`)
		So(result, ShouldResemble, map[string]interface{}{"a": float64(1), "b": float64(2)})
	})

	Convey("Parses durations", t, func() {
		result := parseString("3s")
		So(result, ShouldEqual, 3*time.Second)
	})

	Convey("Leaves every other format untouched", t, func() {
		result := parseString("Hello")
		So(result, ShouldEqual, "Hello")
	})
}

func TestParseDuration(t *testing.T) {

	Convey("Returns a duration when the string matches duration format", t, func() {
		result := parseDurationString("3s")
		So(result, ShouldEqual, 3*time.Second)
	})

	Convey("Returns the original string when the string doesn't match duration format", t, func() {
		result := parseDurationString("definitely not 3s")
		So(result, ShouldEqual, "definitely not 3s")
	})
}

func TestConvertDurationStrings(t *testing.T) {

	Convey("Parses string durations", t, func() {
		m := convertDurationStrings(map[string]interface{}{"a": "3s"})
		So(m, ShouldResemble, map[string]interface{}{"a": 3 * time.Second})
	})

	Convey("Recurses into submaps", t, func() {
		m := convertDurationStrings(map[string]interface{}{"a": map[string]interface{}{"b": "3s"}})
		So(m, ShouldResemble, map[string]interface{}{"a": map[string]interface{}{"b": 3 * time.Second}})
	})
}

func TestCasts(t *testing.T) {

	Convey("castMap", t, func() {

		Convey("Casts an empty interface into a map", func() {
			result, err := castMap(map[string]interface{}{})
			So(result, ShouldResemble, map[string]interface{}{})
			So(err, ShouldBeNil)
		})

		Convey("Returns an error if the cast can't be done", func() {
			result, err := castMap(5)
			So(result, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("castSlice", t, func() {

		Convey("Casts an empty interface into a slice", func() {
			result, err := castSlice([]interface{}{1, "two"})
			So(result, ShouldResemble, []interface{}{1, "two"})
			So(err, ShouldBeNil)
		})

		Convey("Returns an error if the cast can't be done", func() {
			result, err := castSlice(5)
			So(result, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("castStringSlice", t, func() {

		Convey("Casts an empty interface into a string slice", func() {
			result, err := castStringSlice([]string{"one", "two"})
			So(result, ShouldResemble, []string{"one", "two"})
			So(err, ShouldBeNil)
		})

		Convey("Returns an error if the cast can't be done", func() {
			result, err := castStringSlice(5)
			So(result, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("castString", t, func() {

		Convey("Casts an empty interface into a string", func() {
			result, err := castString("one")
			So(result, ShouldEqual, "one")
			So(err, ShouldBeNil)
		})

		Convey("Returns an error if the cast can't be done", func() {
			result, err := castString(5)
			So(result, ShouldBeZeroValue)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("castIntegerSlice", t, func() {

		Convey("Casts an empty interface into a integer slice", func() {
			result, err := castIntegerSlice([]int{1, 2})
			So(result, ShouldResemble, []int{1, 2})
			So(err, ShouldBeNil)
		})

		Convey("Attempts to cast a float slice to an int", func() {
			result, err := castIntegerSlice([]float64{1, 2})
			So(result, ShouldResemble, []int{1, 2})
			So(err, ShouldBeNil)
		})

		Convey("Returns an error when casting a float slice to an int slice would truncate a value", func() {
			result, err := castIntegerSlice([]float64{1.5, 2})
			So(result, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})

		Convey("Returns an error if the cast can't be done", func() {
			result, err := castIntegerSlice(5)
			So(result, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("castInteger", t, func() {

		Convey("Casts an empty interface into a integer", func() {
			result, err := castInteger(1)
			So(result, ShouldEqual, 1)
			So(err, ShouldBeNil)
		})

		Convey("Attempts to cast a float to an int", func() {
			result, err := castInteger(float64(1))
			So(result, ShouldEqual, 1)
			So(err, ShouldBeNil)
		})

		Convey("Returns an error when a converting to a int would truncate the value", func() {
			result, err := castInteger(1.5)
			So(result, ShouldBeZeroValue)
			So(err, ShouldNotBeNil)
		})

		Convey("Returns an error if the cast can't be done", func() {
			result, err := castInteger("Hello")
			So(result, ShouldBeZeroValue)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("castBooleanSlice", t, func() {

		Convey("Casts an empty interface into a boolean slice", func() {
			result, err := castBooleanSlice([]bool{true, false})
			So(result, ShouldResemble, []bool{true, false})
			So(err, ShouldBeNil)
		})

		Convey("Returns an error if the cast can't be done", func() {
			result, err := castBooleanSlice(5)
			So(result, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("castBoolean", t, func() {

		Convey("Casts an empty interface into a boolean", func() {
			result, err := castBoolean(true)
			So(result, ShouldEqual, true)
			So(err, ShouldBeNil)
		})

		Convey("Returns an error if the cast can't be done", func() {
			result, err := castBoolean(5)
			So(result, ShouldBeZeroValue)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("castFloatSlice", t, func() {

		Convey("Casts an empty interface into a float slice", func() {
			result, err := castFloatSlice([]float64{1, 2})
			So(result, ShouldResemble, []float64{1, 2})
			So(err, ShouldBeNil)
		})

		Convey("Returns an error if the cast can't be done", func() {
			result, err := castFloatSlice(5)
			So(result, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("castFloat", t, func() {

		Convey("Casts an empty interface into a float", func() {
			result, err := castFloat(3.3)
			So(result, ShouldEqual, 3.3)
			So(err, ShouldBeNil)
		})

		Convey("Returns an error if the cast can't be done", func() {
			result, err := castFloat("Hello")
			So(result, ShouldBeZeroValue)
			So(err, ShouldNotBeNil)
		})
	})
}
