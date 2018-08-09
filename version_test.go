package manifest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v2"
)

func TestUnmarshalVersion(t *testing.T) {

	t.Run("ReturnsSemverVersionByDefaultIfNoOtherVersionTypeIsSet", func(t *testing.T) {

		var version EstafetteVersion

		// act
		err := yaml.Unmarshal([]byte(``), &version)
		version.SetDefaults()

		assert.Nil(t, err)
		assert.Nil(t, version.Custom)
		assert.NotNil(t, version.SemVer)
		assert.Equal(t, 0, version.SemVer.Major)
	})

	t.Run("ReturnsCustomVersionWithLabelTemplateDefaultingToRevisionPlaceholder", func(t *testing.T) {

		var version EstafetteVersion

		// act
		err := yaml.Unmarshal([]byte(`
custom:
  labelTemplate: ''`), &version)

		assert.Nil(t, err)
		assert.Nil(t, version.SemVer)
		assert.NotNil(t, version.Custom)
		assert.Equal(t, "{{revision}}", version.Custom.LabelTemplate)
	})

	t.Run("ReturnsSemverVersionIfSemverIsSet", func(t *testing.T) {

		var version EstafetteVersion

		// act
		err := yaml.Unmarshal([]byte(`
semver:
  major: 1
  minor: 2
  patch: '{{auto}}'
  labelTemplate: '{{branch}}'
  releaseBranch: master`), &version)

		assert.Nil(t, err)
		assert.Nil(t, version.Custom)
		assert.NotNil(t, version.SemVer)
		assert.Equal(t, 1, version.SemVer.Major)
	})

	t.Run("ReturnsSemverVersionWithMajorDefaultingToZero", func(t *testing.T) {

		var version EstafetteVersion

		// act
		err := yaml.Unmarshal([]byte(`
semver:
  minor: 2`), &version)

		assert.Nil(t, err)
		assert.Nil(t, version.Custom)
		assert.NotNil(t, version.SemVer)
		assert.Equal(t, 0, version.SemVer.Major)
	})

	t.Run("ReturnsSemverVersionWithMinorDefaultingToZero", func(t *testing.T) {

		var version EstafetteVersion

		// act
		err := yaml.Unmarshal([]byte(`
  semver:
    major: 1`), &version)

		assert.Nil(t, err)
		assert.Nil(t, version.Custom)
		assert.NotNil(t, version.SemVer)
		assert.Equal(t, 0, version.SemVer.Minor)
	})

	t.Run("ReturnsSemverVersionWithPatchDefaultingToAutoPlaceholder", func(t *testing.T) {

		var version EstafetteVersion

		// act
		err := yaml.Unmarshal([]byte(`
semver:
  major: 1
  minor: 2`), &version)

		assert.Nil(t, err)
		assert.Nil(t, version.Custom)
		assert.NotNil(t, version.SemVer)
		assert.Equal(t, "{{auto}}", version.SemVer.Patch)
	})

	t.Run("ReturnsSemverVersionWithLabelTemplateDefaultingToBranchPlaceholder", func(t *testing.T) {

		var version EstafetteVersion

		// act
		err := yaml.Unmarshal([]byte(`
semver:
  major: 1
  minor: 2`), &version)

		assert.Nil(t, err)
		assert.Nil(t, version.Custom)
		assert.NotNil(t, version.SemVer)
		assert.Equal(t, "{{branch}}", version.SemVer.LabelTemplate)
	})

	t.Run("ReturnsSemverVersionWithReleaseBranchDefaultingToMaster", func(t *testing.T) {

		var version EstafetteVersion

		// act
		err := yaml.Unmarshal([]byte(`
semver:
  major: 1
  minor: 2`), &version)

		assert.Nil(t, err)
		assert.Nil(t, version.Custom)
		assert.NotNil(t, version.SemVer)
		assert.Equal(t, "master", version.SemVer.ReleaseBranch)
	})

}
