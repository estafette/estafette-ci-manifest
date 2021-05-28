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
		assert.Equal(t, "master", version.SemVer.ReleaseBranch.Values[0])
	})

	t.Run("ReturnsSemverVersionWithMultipleReleaseBranchesForArray", func(t *testing.T) {

		var version EstafetteVersion

		// act
		err := yaml.Unmarshal([]byte(`
semver:
  releaseBranch: 
  - master
  - production`), &version)

		assert.Nil(t, err)
		assert.Nil(t, version.Custom)
		assert.NotNil(t, version.SemVer)
		assert.Equal(t, "master", version.SemVer.ReleaseBranch.Values[0])
		assert.Equal(t, "production", version.SemVer.ReleaseBranch.Values[1])
	})
}

func TestCustomVersion(t *testing.T) {

	t.Run("ReturnsLabelTemplateAsIsWhenItHasNoPlaceholders", func(t *testing.T) {

		version := EstafetteCustomVersion{
			LabelTemplate: "whateveryoulike",
		}
		params := EstafetteVersionParams{
			AutoIncrement: 5,
			Branch:        "release",
			Revision:      "219aae19153da2b20ac1d88e2fd68e0b20274be2",
		}

		// act
		versionString := version.Version(params)

		assert.Equal(t, "whateveryoulike", versionString)
	})

	t.Run("ReturnsLabelTemplateWithAutoPlaceholderReplaced", func(t *testing.T) {

		version := EstafetteCustomVersion{
			LabelTemplate: "{{auto}}",
		}
		params := EstafetteVersionParams{
			AutoIncrement: 5,
			Branch:        "release",
			Revision:      "219aae19153da2b20ac1d88e2fd68e0b20274be2",
		}

		// act
		versionString := version.Version(params)

		assert.Equal(t, "5", versionString)
	})

	t.Run("ReturnsLabelTemplateWithBranchPlaceholderReplaced", func(t *testing.T) {

		version := EstafetteCustomVersion{
			LabelTemplate: "{{branch}}",
		}
		params := EstafetteVersionParams{
			AutoIncrement: 5,
			Branch:        "release",
			Revision:      "219aae19153da2b20ac1d88e2fd68e0b20274be2",
		}

		// act
		versionString := version.Version(params)

		assert.Equal(t, "release", versionString)
	})

	t.Run("ReturnsLabelTemplateWithRevisionPlaceholderReplaced", func(t *testing.T) {

		version := EstafetteCustomVersion{
			LabelTemplate: "{{revision}}",
		}
		params := EstafetteVersionParams{
			AutoIncrement: 5,
			Branch:        "release",
			Revision:      "219aae19153da2b20ac1d88e2fd68e0b20274be2",
		}

		// act
		versionString := version.Version(params)

		assert.Equal(t, "219aae19153da2b20ac1d88e2fd68e0b20274be2", versionString)
	})
}

func TestSemverVersion(t *testing.T) {

	t.Run("ReturnsMajorDotMinorDotPatchDashLabelTemplateAsIsWhenItHasNoPlaceholders", func(t *testing.T) {

		version := EstafetteSemverVersion{
			Major:         5,
			Minor:         3,
			Patch:         "6",
			LabelTemplate: "whateveryoulike",
			ReleaseBranch: StringOrStringArray{Values: []string{"alpha"}},
		}
		params := EstafetteVersionParams{
			AutoIncrement: 5,
			Branch:        "release",
			Revision:      "219aae19153da2b20ac1d88e2fd68e0b20274be2",
		}

		// act
		versionString := version.Version(params)

		assert.Equal(t, "5.3.6-whateveryoulike", versionString)
	})

	t.Run("ReturnsSemverWithAutoPlaceholderInPatchReplaced", func(t *testing.T) {

		version := EstafetteSemverVersion{
			Major:         5,
			Minor:         3,
			Patch:         "{{auto}}",
			LabelTemplate: "whateveryoulike",
			ReleaseBranch: StringOrStringArray{Values: []string{"alpha"}},
		}
		params := EstafetteVersionParams{
			AutoIncrement: 16,
			Branch:        "release",
			Revision:      "219aae19153da2b20ac1d88e2fd68e0b20274be2",
		}

		// act
		versionString := version.Version(params)

		assert.Equal(t, "5.3.16-whateveryoulike", versionString)
	})

	t.Run("ReturnsSemverWithBranchPlaceholderInLabelReplaced", func(t *testing.T) {

		version := EstafetteSemverVersion{
			Major:         5,
			Minor:         3,
			Patch:         "6",
			LabelTemplate: "{{branch}}",
			ReleaseBranch: StringOrStringArray{Values: []string{"release"}},
		}
		params := EstafetteVersionParams{
			AutoIncrement: 16,
			Branch:        "alpha",
			Revision:      "219aae19153da2b20ac1d88e2fd68e0b20274be2",
		}

		// act
		versionString := version.Version(params)

		assert.Equal(t, "5.3.6-alpha", versionString)
	})

	t.Run("ReturnsSemverWithLabelWithoutSlashIfBranchHasSlash", func(t *testing.T) {

		version := EstafetteSemverVersion{
			Major:         5,
			Minor:         3,
			Patch:         "6",
			LabelTemplate: "{{branch}}",
			ReleaseBranch: StringOrStringArray{Values: []string{"release"}},
		}
		params := EstafetteVersionParams{
			AutoIncrement: 16,
			Branch:        "feature/foo",
			Revision:      "219aae19153da2b20ac1d88e2fd68e0b20274be2",
		}

		// act
		versionString := version.Version(params)

		assert.Equal(t, "5.3.6-feature-foo", versionString)
	})

	t.Run("ReturnsSemverWithDNSSafeLabel", func(t *testing.T) {

		version := EstafetteSemverVersion{
			Major:         5,
			Minor:         3,
			Patch:         "6",
			LabelTemplate: "{{branch}}",
			ReleaseBranch: StringOrStringArray{Values: []string{"release"}},
		}
		params := EstafetteVersionParams{
			AutoIncrement: 16,
			Branch:        "HOTFIX_production_deploy",
			Revision:      "219aae19153da2b20ac1d88e2fd68e0b20274be2",
		}

		// act
		versionString := version.Version(params)

		assert.Equal(t, "5.3.6-hotfix-production-deploy", versionString)
	})

	t.Run("ReturnsSemverWithDNSSafeLabelForBranchStartingWithNumber", func(t *testing.T) {

		version := EstafetteSemverVersion{
			Major:         5,
			Minor:         3,
			Patch:         "6",
			LabelTemplate: "{{revision}}",
			ReleaseBranch: StringOrStringArray{Values: []string{"release"}},
		}
		params := EstafetteVersionParams{
			AutoIncrement: 16,
			Branch:        "19.1.0",
			Revision:      "219aae19153da2b20ac1d88e2fd68e0b20274be2",
		}

		// act
		versionString := version.Version(params)

		assert.Equal(t, "5.3.6-revision-219aae19153da2b20ac1d88e2fd68e0b20274be2", versionString)
	})

	t.Run("ReturnsSemverWithDNSSafeLabelForRevisionStartingWithNumber", func(t *testing.T) {

		version := EstafetteSemverVersion{
			Major:         5,
			Minor:         3,
			Patch:         "6",
			LabelTemplate: "{{branch}}",
			ReleaseBranch: StringOrStringArray{Values: []string{"release"}},
		}
		params := EstafetteVersionParams{
			AutoIncrement: 16,
			Branch:        "19.1.0",
			Revision:      "219aae19153da2b20ac1d88e2fd68e0b20274be2",
		}

		// act
		versionString := version.Version(params)

		assert.Equal(t, "5.3.6-branch-19-1-0", versionString)
	})

	t.Run("ReturnsSemverWithDNSSafeLabelStartingWithALetter", func(t *testing.T) {

		version := EstafetteSemverVersion{
			Major:         5,
			Minor:         3,
			Patch:         "6",
			LabelTemplate: "{{branch}}",
			ReleaseBranch: StringOrStringArray{Values: []string{"release"}},
		}
		params := EstafetteVersionParams{
			AutoIncrement: 16,
			Branch:        "123-HOTFIX_production_deploy",
			Revision:      "219aae19153da2b20ac1d88e2fd68e0b20274be2",
		}

		// act
		versionString := version.Version(params)

		assert.Equal(t, "5.3.6-branch-123-hotfix-production-deploy", versionString)
	})

	t.Run("ReturnsSemverWithDNSSafeLabelOfMax63Characters", func(t *testing.T) {

		version := EstafetteSemverVersion{
			Major:         5,
			Minor:         3,
			Patch:         "6",
			LabelTemplate: "{{branch}}",
			ReleaseBranch: StringOrStringArray{Values: []string{"release"}},
		}
		params := EstafetteVersionParams{
			AutoIncrement: 16,
			Branch:        "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz123456789012",
			Revision:      "219aae19153da2b20ac1d88e2fd68e0b20274be2",
		}

		// act
		versionString := version.Version(params)

		assert.Equal(t, "5.3.6-abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz12345678901", versionString)
	})

	t.Run("ReturnsSemverWithLabelHyphensTrimmedAtEnd", func(t *testing.T) {

		version := EstafetteSemverVersion{
			Major:         5,
			Minor:         3,
			Patch:         "6",
			LabelTemplate: "{{branch}}",
			ReleaseBranch: StringOrStringArray{Values: []string{"release"}},
		}
		params := EstafetteVersionParams{
			AutoIncrement: 16,
			Branch:        "-some-branch-surrounded-by-hyphens-",
			Revision:      "219aae19153da2b20ac1d88e2fd68e0b20274be2",
		}

		// act
		versionString := version.Version(params)

		assert.Equal(t, "5.3.6-some-branch-surrounded-by-hyphens", versionString)
	})

	t.Run("ReturnsSemverWithoutLabelIfBranchMatchesReleaseBranch", func(t *testing.T) {

		version := EstafetteSemverVersion{
			Major:         5,
			Minor:         3,
			Patch:         "6",
			LabelTemplate: "{{branch}}",
			ReleaseBranch: StringOrStringArray{Values: []string{"release"}},
		}
		params := EstafetteVersionParams{
			AutoIncrement: 16,
			Branch:        "release",
			Revision:      "219aae19153da2b20ac1d88e2fd68e0b20274be2",
		}

		// act
		versionString := version.Version(params)

		assert.Equal(t, "5.3.6", versionString)
	})

	t.Run("ReturnsSemverWithoutLabelIfBranchMatchesReleaseBranchRegularExpression", func(t *testing.T) {

		version := EstafetteSemverVersion{
			Major:         5,
			Minor:         3,
			Patch:         "6",
			LabelTemplate: "{{branch}}",
			ReleaseBranch: StringOrStringArray{Values: []string{"release/.+"}},
		}
		params := EstafetteVersionParams{
			AutoIncrement: 16,
			Branch:        "release/abc",
			Revision:      "219aae19153da2b20ac1d88e2fd68e0b20274be2",
		}

		// act
		versionString := version.Version(params)

		assert.Equal(t, "5.3.6", versionString)
	})

	t.Run("ReturnsSemverWithLabelIfBranchDoesNotMatchReleaseBranchRegularExpression", func(t *testing.T) {

		version := EstafetteSemverVersion{
			Major:         5,
			Minor:         3,
			Patch:         "6",
			LabelTemplate: "{{branch}}",
			ReleaseBranch: StringOrStringArray{Values: []string{"release/[0-9]+"}},
		}
		params := EstafetteVersionParams{
			AutoIncrement: 16,
			Branch:        "release/abc",
			Revision:      "219aae19153da2b20ac1d88e2fd68e0b20274be2",
		}

		// act
		versionString := version.Version(params)

		assert.Equal(t, "5.3.6-release-abc", versionString)
	})
}
