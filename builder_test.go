package manifest

import (
	"testing"

	yaml "github.com/buildkite/yaml"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshalBuilder(t *testing.T) {
	t.Run("ReturnsUnmarshaledBuilder", func(t *testing.T) {

		var builder EstafetteBuilder

		// act
		err := yaml.Unmarshal([]byte(`
track: dev`), &builder)

		assert.Nil(t, err)
		assert.Equal(t, "dev", builder.Track)
	})

	t.Run("DefaultsTrackToStableIfMissingIfNotPresent", func(t *testing.T) {

		var builder EstafetteBuilder

		// act
		err := yaml.Unmarshal([]byte(` 
`), &builder)
		builder.setDefaults()

		assert.Nil(t, err)
		assert.Equal(t, "stable", builder.Track)
	})
}
