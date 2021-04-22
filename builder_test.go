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

	t.Run("DefaultsOperatingSystemToLinuxIfMissingIfNotPresent", func(t *testing.T) {

		var builder EstafetteBuilder

		// act
		err := yaml.Unmarshal([]byte(`
`), &builder)
		builder.SetDefaults(*GetDefaultManifestPreferences())

		assert.Nil(t, err)
		assert.Equal(t, OperatingSystemLinux, builder.OperatingSystem)
	})

	t.Run("DefaultsTrackToStableIfMissingIfNotPresent", func(t *testing.T) {

		var builder EstafetteBuilder

		// act
		err := yaml.Unmarshal([]byte(` 
`), &builder)
		builder.SetDefaults(*GetDefaultManifestPreferences())

		assert.Nil(t, err)
		assert.Equal(t, "stable", builder.Track)
	})

	t.Run("DefaultsStorageMediumToDefaultIfMissingIfNotPresent", func(t *testing.T) {

		var builder EstafetteBuilder

		// act
		err := yaml.Unmarshal([]byte(` 
`), &builder)
		builder.SetDefaults(*GetDefaultManifestPreferences())

		assert.Nil(t, err)
		assert.Equal(t, StorageMediumDefault, builder.StorageMedium)
	})

	t.Run("KeepsStorageMediumToIfSet", func(t *testing.T) {

		var builder EstafetteBuilder

		// act
		err := yaml.Unmarshal([]byte(`
medium: memory`), &builder)
		builder.SetDefaults(*GetDefaultManifestPreferences())

		assert.Nil(t, err)
		assert.Equal(t, StorageMediumMemory, builder.StorageMedium)
	})

}
