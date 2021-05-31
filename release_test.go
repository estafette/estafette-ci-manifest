package manifest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v2"
)

func TestUnmarshalRelease(t *testing.T) {
	t.Run("ReturnsUnmarshaledRelease", func(t *testing.T) {

		var release EstafetteRelease

		// act
		err := yaml.Unmarshal([]byte(`
stages:
  deploy:
    image: extensions/deploy-to-kubernetes-engine:stable

  create-release-notes:
    image: extensions/create-release-notes-from-changelog:stable
`), &release)

		assert.Nil(t, err)
		assert.Equal(t, 2, len(release.Stages))
		assert.Equal(t, "deploy", release.Stages[0].Name)
		assert.Equal(t, "extensions/deploy-to-kubernetes-engine:stable", release.Stages[0].ContainerImage)
		assert.Equal(t, "create-release-notes", release.Stages[1].Name)
		assert.Equal(t, "extensions/create-release-notes-from-changelog:stable", release.Stages[1].ContainerImage)
	})
}

func TestReleaseToYamlMarshalling(t *testing.T) {
	t.Run("UnmarshallingThenMarshallingReturnsTheSameFile", func(t *testing.T) {

		var release EstafetteRelease

		input := `stages:
  deploy:
    image: extensions/deploy-to-kubernetes-engine:stable
    shell: /bin/sh
    workDir: /estafette-work
    when: status == 'succeeded'
  create-release-notes:
    image: extensions/create-release-notes-from-changelog:stable
    shell: /bin/sh
    workDir: /estafette-work
    when: status == 'succeeded'
`
		err := yaml.Unmarshal([]byte(input), &release)
		assert.Nil(t, err)

		// act
		output, err := yaml.Marshal(release)

		assert.Nil(t, err)
		assert.Equal(t, input, string(output))
	})
}
