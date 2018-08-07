package manifest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v2"
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
		//builder.SetDefaults()

		assert.Nil(t, err)
		assert.Equal(t, "stable", builder.Track)
	})
}
