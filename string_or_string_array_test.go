package manifest

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v2"
)

func TestYAMLUnmarshalStringOrStringArray(t *testing.T) {
	t.Run("ReturnsUnmarshaledStringOrStringArrayWithoutValuesForEmptyValue", func(t *testing.T) {

		var stringOrStringArray StringOrStringArray

		// act
		err := yaml.Unmarshal([]byte(``), &stringOrStringArray)

		assert.Nil(t, err)
		assert.Equal(t, 0, len(stringOrStringArray.Values))
	})

	t.Run("ReturnsUnmarshaledStringOrStringArrayWithSingleValuesForSingleValue", func(t *testing.T) {

		var stringOrStringArray StringOrStringArray

		// act
		err := yaml.Unmarshal([]byte(`
singlevalue`), &stringOrStringArray)

		assert.Nil(t, err)
		assert.Equal(t, 1, len(stringOrStringArray.Values))
		assert.Equal(t, "singlevalue", stringOrStringArray.Values[0])
	})

	t.Run("ReturnsUnmarshaledStringOrStringArrayWithMultipleValuesForMultipleValues", func(t *testing.T) {

		var stringOrStringArray StringOrStringArray

		// act
		err := yaml.Unmarshal([]byte(`
- value1
- value2`), &stringOrStringArray)

		assert.Nil(t, err)
		assert.Equal(t, 2, len(stringOrStringArray.Values))
		assert.Equal(t, "value1", stringOrStringArray.Values[0])
		assert.Equal(t, "value2", stringOrStringArray.Values[1])
	})
}

func TestYAMLMarshalStringOrStringArray(t *testing.T) {
	t.Run("ReturnsYamlStringForStringOrStringArrayWithNoValue", func(t *testing.T) {

		stringOrStringArray := StringOrStringArray{}

		// act
		output, err := yaml.Marshal(stringOrStringArray)

		assert.Nil(t, err)
		assert.Equal(t, "\"\"\n", string(output))
	})

	t.Run("ReturnsUnmarshaledStringOrStringArrayWithSingleValuesForSingleValue", func(t *testing.T) {

		stringOrStringArray := StringOrStringArray{
			Values: []string{
				"singlevalue",
			},
		}

		// act
		output, err := yaml.Marshal(stringOrStringArray)

		assert.Nil(t, err)
		assert.Equal(t, "singlevalue\n", string(output))
	})

	t.Run("ReturnsUnmarshaledStringOrStringArrayWithMultipleValuesForMultipleValues", func(t *testing.T) {

		stringOrStringArray := StringOrStringArray{
			Values: []string{
				"value1",
				"value2",
			},
		}

		// act
		output, err := yaml.Marshal(stringOrStringArray)

		assert.Nil(t, err)
		assert.Equal(t, "- value1\n- value2\n", string(output))
	})
}

func TestJSONMarshalStringOrStringArray(t *testing.T) {
	t.Run("ReturnsJSONStringForStringOrStringArrayWithNoValue", func(t *testing.T) {

		stringOrStringArray := StringOrStringArray{}

		// act
		output, err := json.Marshal(stringOrStringArray)

		assert.Nil(t, err)
		assert.Equal(t, "\"\"", string(output))
	})

	t.Run("ReturnsUnmarshaledStringOrStringArrayWithSingleValuesForSingleValue", func(t *testing.T) {

		stringOrStringArray := StringOrStringArray{
			Values: []string{
				"singlevalue",
			},
		}

		// act
		output, err := json.Marshal(stringOrStringArray)

		assert.Nil(t, err)
		assert.Equal(t, "\"singlevalue\"", string(output))
	})

	t.Run("ReturnsUnmarshaledStringOrStringArrayWithMultipleValuesForMultipleValues", func(t *testing.T) {

		stringOrStringArray := StringOrStringArray{
			Values: []string{
				"value1",
				"value2",
			},
		}

		// act
		output, err := json.Marshal(stringOrStringArray)

		assert.Nil(t, err)
		assert.Equal(t, "[\"value1\",\"value2\"]", string(output))
	})
}

func TestJSONUnmarshalStringOrStringArray(t *testing.T) {
	t.Run("ReturnsUnmarshaledStringOrStringArrayWithoutValuesForEmptyValue", func(t *testing.T) {

		var stringOrStringArray StringOrStringArray

		// act
		err := json.Unmarshal([]byte("\"\""), &stringOrStringArray)

		assert.Nil(t, err)
		assert.Equal(t, 0, len(stringOrStringArray.Values))
	})

	t.Run("ReturnsUnmarshaledStringOrStringArrayWithSingleValuesForSingleValue", func(t *testing.T) {

		var stringOrStringArray StringOrStringArray

		// act
		err := json.Unmarshal([]byte("\"singlevalue\""), &stringOrStringArray)

		assert.Nil(t, err)
		assert.Equal(t, 1, len(stringOrStringArray.Values))
		assert.Equal(t, "singlevalue", stringOrStringArray.Values[0])
	})

	t.Run("ReturnsUnmarshaledStringOrStringArrayWithMultipleValuesForMultipleValues", func(t *testing.T) {

		var stringOrStringArray StringOrStringArray

		// act
		err := json.Unmarshal([]byte("[\"value1\",\"value2\"]"), &stringOrStringArray)

		assert.Nil(t, err)
		assert.Equal(t, 2, len(stringOrStringArray.Values))
		assert.Equal(t, "value1", stringOrStringArray.Values[0])
		assert.Equal(t, "value2", stringOrStringArray.Values[1])
	})
}

func TestContains(t *testing.T) {
	t.Run("ReturnsTrueIfValueMatchesOneOfTheValues", func(t *testing.T) {

		stringOrStringArray := StringOrStringArray{
			Values: []string{
				"value1",
				"value2",
			},
		}

		// act
		result := stringOrStringArray.Contains("value2")

		assert.True(t, result)
	})

	t.Run("ReturnsFalseIfValueMatchesNoneOfTheValues", func(t *testing.T) {

		stringOrStringArray := StringOrStringArray{
			Values: []string{
				"value1",
				"value2",
			},
		}

		// act
		result := stringOrStringArray.Contains("value3")

		assert.False(t, result)
	})

	t.Run("ReturnsFalseIfStringOrStringArrayHasNoValues", func(t *testing.T) {

		stringOrStringArray := StringOrStringArray{}

		// act
		result := stringOrStringArray.Contains("value3")

		assert.False(t, result)
	})
}
