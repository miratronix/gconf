package internal

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParseArguments(t *testing.T) {

	Convey("Returns an empty map if there are no arguments", t, func() {
		loader := NewArgumentLoader("", "")
		result, err := loader.parseArguments([]string{})
		So(result, ShouldBeEmpty)
		So(err, ShouldBeNil)
	})

	Convey("Parses arguments without a separator or prefix", t, func() {
		loader := NewArgumentLoader("", "")

		Convey("Ignores arguments not starting with '-' or '--'", func() {
			result, err := loader.parseArguments([]string{"string=testing"})
			So(result, ShouldBeEmpty)
			So(err, ShouldBeNil)
		})

		Convey("Ignores arguments without an '='", func() {
			result, err := loader.parseArguments([]string{"testing"})
			So(result, ShouldBeEmpty)
			So(err, ShouldBeNil)
		})

		Convey("Parses arguments starting with '-'", func() {
			result, err := loader.parseArguments([]string{"-string=testing"})
			So(result, ShouldResemble, map[string]interface{}{"string": "testing"})
			So(err, ShouldBeNil)
		})

		Convey("Parses arguments starting with '--'", func() {
			result, err := loader.parseArguments([]string{"--string=testing"})
			So(result, ShouldResemble, map[string]interface{}{"string": "testing"})
			So(err, ShouldBeNil)
		})
	})

	Convey("Parses arguments with a separator", t, func() {
		loader := NewArgumentLoader("__", "")

		Convey("Nests the configuration using the separator", func() {
			result, err := loader.parseArguments([]string{"--map__string=testing"})
			So(result, ShouldResemble, map[string]interface{}{"map": map[string]interface{}{"string": "testing"}})
			So(err, ShouldBeNil)
		})

		Convey("Fails when a nested option would override another option", func() {
			result, err := loader.parseArguments([]string{"--test=testing", "--test__stuff=things"})
			So(result, ShouldResemble, map[string]interface{}{"test": "testing"})
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Parses arguments with a prefix", t, func() {
		loader := NewArgumentLoader("", "TEST")

		Convey("Strips the prefix from the argument", func() {
			result, err := loader.parseArguments([]string{"--TESTing=testing"})
			So(result, ShouldResemble, map[string]interface{}{"ing": "testing"})
			So(err, ShouldBeNil)
		})

		Convey("Ignores arguments without the prefix", func() {
			result, err := loader.parseArguments([]string{"--notTesting=testing"})
			So(result, ShouldBeEmpty)
			So(err, ShouldBeNil)
		})
	})
}
