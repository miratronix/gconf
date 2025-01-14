package internal

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParseEnvironment(t *testing.T) {

	Convey("Parses environment data without lower casing, a separator, or a prefix", t, func() {
		loader := NewEnvironmentLoader(false, "", "")
		result, err := loader.parseEnvironment([]string{"TEST=abc"})
		So(result, ShouldResemble, map[string]interface{}{"TEST": "abc"})
		So(err, ShouldBeNil)
	})

	Convey("Lower cases environment keys when enabled", t, func() {
		loader := NewEnvironmentLoader(true, "", "")
		result, err := loader.parseEnvironment([]string{"TEST=abc"})
		So(result, ShouldResemble, map[string]interface{}{"test": "abc"})
		So(err, ShouldBeNil)
	})

	Convey("Returns an error when a environment variable is defined twice", t, func() {
		loader := NewEnvironmentLoader(false, "", "")
		result, err := loader.parseEnvironment([]string{"test=abc", "test=def"})
		So(result, ShouldResemble, map[string]interface{}{"test": "abc"})
		So(err, ShouldNotBeNil)
	})

	Convey("Parses environment variables using the configured separator", t, func() {
		loader := NewEnvironmentLoader(false, "__", "")

		Convey("Doesn't nest the key when the environment variable doesn't contain the separator", func() {
			result, err := loader.parseEnvironment([]string{"test=abc"})
			So(result, ShouldResemble, map[string]interface{}{"test": "abc"})
			So(err, ShouldBeNil)
		})

		Convey("Nests the key when the environment variable does contain the separator", func() {
			result, err := loader.parseEnvironment([]string{"test__ing=abc"})
			So(result, ShouldResemble, map[string]interface{}{"test": map[string]interface{}{"ing": "abc"}})
			So(err, ShouldBeNil)
		})

		Convey("Doesn't nest the key when it starts with the separator", func() {
			result, err := loader.parseEnvironment([]string{"__test=abc"})
			So(result, ShouldResemble, map[string]interface{}{"test": "abc"})
			So(err, ShouldBeNil)
		})

		Convey("Returns an error when overriding an existing environment variable with a nested one", func() {
			result, err := loader.parseEnvironment([]string{"test=abc", "test__ing=abc"})
			So(result, ShouldResemble, map[string]interface{}{"test": "abc"})
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Parses environment variables using the configured prefix", t, func() {
		loader := NewEnvironmentLoader(false, "", "prefix")

		Convey("Ignores environment variables without the prefix", func() {
			result, err := loader.parseEnvironment([]string{"test=abc"})
			So(result, ShouldBeEmpty)
			So(err, ShouldBeNil)
		})

		Convey("Removes the prefix from environment variables", func() {
			result, err := loader.parseEnvironment([]string{"prefixtest=abc"})
			So(result, ShouldResemble, map[string]interface{}{"test": "abc"})
			So(err, ShouldBeNil)
		})

		Convey("Doesn't lowercase the prefix when lower casing is enabled", func() {
			loader := NewEnvironmentLoader(true, "", "PREFIX")
			result, err := loader.parseEnvironment([]string{"prefixtest=abc"})
			So(result, ShouldBeEmpty)
			So(err, ShouldBeNil)
		})
	})

	Convey("Allows an equal sign in the value", t, func() {
		value := "a+b=c"
		loader := NewEnvironmentLoader(false, "", "")
		result, err := loader.parseEnvironment([]string{fmt.Sprintf("test=%v", value)})
		So(result["test"], ShouldEqual, value)
		So(err, ShouldBeNil)
	})
}
