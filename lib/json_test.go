package lib

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestJSONFileLoad(t *testing.T) {

	Convey("Returns an error when the file can't be found", t, func() {
		result, err := NewJSONFileLoader("", false).Load()
		So(result, ShouldBeEmpty)
		So(err, ShouldNotBeNil)
	})

	Convey("Loads a JSON file", t, func() {
		result, err := NewJSONFileLoader("../test/test.json", false).Load()
		So(err, ShouldBeNil)
		So(result, ShouldResemble, map[string]interface{}{
			"string":  "woohoo",
			"boolean": true,
			"integer": 10.0,
			"float":   3.5,
			"array":   []interface{}{"woohoo", true, float64(10), 3.5},
			"object": map[string]interface{}{
				"string":  "woohoo",
				"boolean": true,
				"integer": float64(10),
				"float":   3.5,
			},
		})
	})
}

func TestParseJSON(t *testing.T) {
	loader := NewJSONFileLoader("", false)

	Convey("Returns an error when parsing invalid JSON", t, func() {
		result, err := loader.parseJSON([]byte(""))
		So(result, ShouldBeEmpty)
		So(err, ShouldNotBeNil)
	})

	Convey("Parses valid JSON into a map", t, func() {

		Convey("Parses JSON without sub objects", func() {
			result, err := loader.parseJSON([]byte(`{"a": "b"}`))
			So(result, ShouldResemble, map[string]interface{}{"a": "b"})
			So(err, ShouldBeNil)
		})

		Convey("Parses JSON with sub objects", func() {
			result, err := loader.parseJSON([]byte(`{"a": { "b": "c" } }`))
			So(result, ShouldResemble, map[string]interface{}{"a": map[string]interface{}{"b": "c"}})
			So(err, ShouldBeNil)
		})

		Convey("Returns the original map when duration parsing is disabled", func() {
			result, err := loader.parseJSON([]byte(`{"a": "3s"}`))
			So(result, ShouldResemble, map[string]interface{}{"a": "3s"})
			So(err, ShouldBeNil)
		})

		Convey("Returns a modified map when duration parsing is enabled", func() {
			l := NewJSONFileLoader("", true)
			result, err := l.parseJSON([]byte(`{"a": "3s"}`))
			So(result, ShouldResemble, map[string]interface{}{"a": 3 * time.Second})
			So(err, ShouldBeNil)
		})
	})
}
