package lib

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestYAMLFileLoad(t *testing.T) {

	Convey("Returns an error when the file can't be found", t, func() {
		result, err := NewYAMLFileLoader("", false).Load()
		So(result, ShouldBeEmpty)
		So(err, ShouldNotBeNil)
	})

	Convey("Loads a YAML file", t, func() {
		result, err := NewYAMLFileLoader("../test/test.yaml", false).Load()
		So(err, ShouldBeNil)
		So(result, ShouldResemble, map[string]interface{}{
			"string":  "woohoo",
			"boolean": true,
			"integer": 10,
			"float":   3.5,
			"array":   []interface{}{"woohoo", true, 10, 3.5},
			"object": map[string]interface{}{
				"string":  "woohoo",
				"boolean": true,
				"integer": 10,
				"float":   3.5,
			},
		})
	})
}

func TestParseYAML(t *testing.T) {
	loader := NewYAMLFileLoader("", false)

	Convey("Returns an error when parsing invalid YAML", t, func() {
		result, err := loader.parseYAML([]byte("@"))
		So(result, ShouldBeEmpty)
		So(err, ShouldNotBeNil)
	})

	Convey("Parses valid YAML into a map", t, func() {

		Convey("Parses YAML without sub objects", func() {
			result, err := loader.parseYAML([]byte(`{"a": "b"}`))
			So(result, ShouldResemble, map[string]interface{}{"a": "b"})
			So(err, ShouldBeNil)
		})

		Convey("Parses YAML with sub objects", func() {
			result, err := loader.parseYAML([]byte(`{"a": { "b": "c" } }`))
			So(result, ShouldResemble, map[string]interface{}{"a": map[string]interface{}{"b": "c"}})
			So(err, ShouldBeNil)
		})

		Convey("Returns the original map when duration parsing is disabled", func() {
			result, err := loader.parseYAML([]byte(`{"a": "3s"}`))
			So(result, ShouldResemble, map[string]interface{}{"a": "3s"})
			So(err, ShouldBeNil)
		})

		Convey("Returns a modified map when duration parsing is enabled", func() {
			l := NewYAMLFileLoader("", true)
			result, err := l.parseYAML([]byte(`{"a": "3s"}`))
			So(result, ShouldResemble, map[string]interface{}{"a": 3 * time.Second})
			So(err, ShouldBeNil)
		})
	})
}
