package manifest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadManifestFromFile(t *testing.T) {

	t.Run("ReturnsManifestWithoutErrors", func(t *testing.T) {

		// act
		_, err := ReadManifestFromFile("test-manifest.yaml")

		assert.Nil(t, err)
	})

	t.Run("ReturnsManifestWithMappedLabels", func(t *testing.T) {

		// act
		manifest, err := ReadManifestFromFile("test-manifest.yaml")

		assert.Nil(t, err)
		assert.Equal(t, "estafette-ci-builder", manifest.Labels["app"])
		assert.Equal(t, "estafette-team", manifest.Labels["team"])
		assert.Equal(t, "golang", manifest.Labels["language"])
	})

	t.Run("ReturnsManifestWithMappedBuilderTrack", func(t *testing.T) {

		// act
		manifest, err := ReadManifestFromFile("test-manifest.yaml")

		assert.Nil(t, err)
		assert.Equal(t, "dev", manifest.Builder.Track)
	})

	t.Run("ReturnsManifestWithBuilderTrackDefaultStable", func(t *testing.T) {

		// act
		manifest, err := ReadManifest("")

		assert.Nil(t, err)
		assert.Equal(t, "stable", manifest.Builder.Track)
	})

	t.Run("ReturnsManifestWithMappedOrderedPipelinesInSameOrderAsInTheManifest", func(t *testing.T) {

		// act
		manifest, err := ReadManifestFromFile("test-manifest.yaml")

		assert.Nil(t, err)

		assert.Equal(t, 5, len(manifest.Pipelines))

		assert.Equal(t, "build", manifest.Pipelines[0].Name)
		assert.Equal(t, "golang:1.8.0-alpine", manifest.Pipelines[0].ContainerImage)
		assert.Equal(t, "/go/src/github.com/estafette/estafette-ci-builder", manifest.Pipelines[0].WorkingDirectory)
		assert.Equal(t, "go test -v ./...", manifest.Pipelines[0].Commands[0])
		assert.Equal(t, "CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./publish/estafette-ci-builder .", manifest.Pipelines[0].Commands[1])

		assert.Equal(t, "bake", manifest.Pipelines[1].Name)
		assert.Equal(t, "docker:17.03.0-ce", manifest.Pipelines[1].ContainerImage)
		assert.Equal(t, "cp Dockerfile ./publish", manifest.Pipelines[1].Commands[0])
		assert.Equal(t, "docker build -t estafette-ci-builder ./publish", manifest.Pipelines[1].Commands[1])

		assert.Equal(t, "set-build-status", manifest.Pipelines[2].Name)
		assert.Equal(t, "extensions/github-status:0.0.2", manifest.Pipelines[2].ContainerImage)
		assert.Equal(t, 0, len(manifest.Pipelines[2].Commands))
		assert.Equal(t, "server == 'estafette'", manifest.Pipelines[2].When)

		assert.Equal(t, "push-to-docker-hub", manifest.Pipelines[3].Name)
		assert.Equal(t, "docker:17.03.0-ce", manifest.Pipelines[3].ContainerImage)
		assert.Equal(t, "docker login --username=${ESTAFETTE_DOCKER_HUB_USERNAME} --password='${ESTAFETTE_DOCKER_HUB_PASSWORD}'", manifest.Pipelines[3].Commands[0])
		assert.Equal(t, "docker push estafette/${ESTAFETTE_LABEL_APP}:${ESTAFETTE_BUILD_VERSION}", manifest.Pipelines[3].Commands[1])
		assert.Equal(t, "status == 'succeeded' && branch == 'master'", manifest.Pipelines[3].When)

		assert.Equal(t, "slack-notify", manifest.Pipelines[4].Name)
		assert.Equal(t, "docker:17.03.0-ce", manifest.Pipelines[4].ContainerImage)
		assert.Equal(t, "curl -X POST --data-urlencode 'payload={\"channel\": \"#build-status\", \"username\": \"estafette-ci-builder\", \"text\": \"Build ${ESTAFETTE_BUILD_VERSION} for ${ESTAFETTE_LABEL_APP} has failed!\"}' ${ESTAFETTE_SLACK_WEBHOOK}", manifest.Pipelines[4].Commands[0])
		assert.Equal(t, "status == 'failed' || branch == 'master'", manifest.Pipelines[4].When)

		assert.Equal(t, "some value with spaces", manifest.Pipelines[4].EnvVars["SOME_ENVIRONMENT_VAR"])
		assert.Equal(t, "value1", manifest.Pipelines[4].CustomProperties["unknownProperty1"])
		assert.Equal(t, "value2", manifest.Pipelines[4].CustomProperties["unknownProperty2"])

		// support custom properties of more complex type, so []string can be used
		assert.Equal(t, "supported1", manifest.Pipelines[4].CustomProperties["unknownProperty3"].([]interface{})[0].(string))
		assert.Equal(t, "supported2", manifest.Pipelines[4].CustomProperties["unknownProperty3"].([]interface{})[1].(string))

		_, unknownPropertyExist := manifest.Pipelines[4].CustomProperties["unknownProperty3"]
		assert.True(t, unknownPropertyExist)

		_, reservedPropertyForGolangNameExist := manifest.Pipelines[4].CustomProperties["ContainerImage"]
		assert.False(t, reservedPropertyForGolangNameExist)

		_, reservedPropertyForYamlNameExist := manifest.Pipelines[4].CustomProperties["image"]
		assert.False(t, reservedPropertyForYamlNameExist)
	})

	t.Run("ReturnsWorkDirDefaultIfMissing", func(t *testing.T) {

		// act
		manifest, err := ReadManifestFromFile("test-manifest.yaml")

		assert.Nil(t, err)

		assert.Equal(t, "/go/src/github.com/estafette/estafette-ci-builder", manifest.Pipelines[0].WorkingDirectory)
	})

	t.Run("ReturnsWorkDirIfSet", func(t *testing.T) {

		// act
		manifest, err := ReadManifestFromFile("test-manifest.yaml")

		assert.Nil(t, err)

		assert.Equal(t, "/estafette-work", manifest.Pipelines[1].WorkingDirectory)
	})

	t.Run("ReturnsShellDefaultIfMissing", func(t *testing.T) {

		// act
		manifest, err := ReadManifestFromFile("test-manifest.yaml")

		assert.Nil(t, err)

		assert.Equal(t, "/bin/sh", manifest.Pipelines[0].Shell)
	})

	t.Run("ReturnsShellIfSet", func(t *testing.T) {

		// act
		manifest, err := ReadManifestFromFile("test-manifest.yaml")

		assert.Nil(t, err)

		assert.Equal(t, "/bin/bash", manifest.Pipelines[1].Shell)
	})

	t.Run("ReturnsWhenIfSet", func(t *testing.T) {

		// act
		manifest, err := ReadManifestFromFile("test-manifest.yaml")

		assert.Nil(t, err)

		assert.Equal(t, "status == 'succeeded' && branch == 'master'", manifest.Pipelines[3].When)
	})

	t.Run("ReturnsWhenDefaultIfMissing", func(t *testing.T) {

		// act
		manifest, err := ReadManifestFromFile("test-manifest.yaml")

		assert.Nil(t, err)

		assert.Equal(t, "status == 'succeeded'", manifest.Pipelines[0].When)
	})

	t.Run("ReturnsManifestWithSemverVersion", func(t *testing.T) {

		// act
		manifest, err := ReadManifestFromFile("test-manifest.yaml")

		assert.Nil(t, err)
		assert.Equal(t, 1, manifest.Version.SemVer.Major)
		assert.Equal(t, 2, manifest.Version.SemVer.Minor)
		assert.Equal(t, "{{auto}}", manifest.Version.SemVer.Patch)
		assert.Equal(t, "{{branch}}", manifest.Version.SemVer.LabelTemplate)
		assert.Equal(t, "master", manifest.Version.SemVer.ReleaseBranch)
	})
}

func TestVersion(t *testing.T) {

	t.Run("ReturnsCustomVersionByDefaultIfNoOtherVersionTypeIsSet", func(t *testing.T) {

		// act
		manifest, err := ReadManifest("")

		assert.Nil(t, err)
		assert.Nil(t, manifest.Version.SemVer)
		assert.NotNil(t, manifest.Version.Custom)
		assert.Equal(t, "{{revision}}", manifest.Version.Custom.LabelTemplate)
	})

	t.Run("ReturnsCustomVersionWithLabelTemplateDefaultingToRevisionPlaceholder", func(t *testing.T) {

		// act
		manifest, err := ReadManifest(`
version:
  custom:
    labelTemplate: ''`)

		assert.Nil(t, err)
		assert.Nil(t, manifest.Version.SemVer)
		assert.NotNil(t, manifest.Version.Custom)
		assert.Equal(t, "{{revision}}", manifest.Version.Custom.LabelTemplate)
	})

	t.Run("ReturnsSemverVersionIfSemverIsSet", func(t *testing.T) {

		// act
		manifest, err := ReadManifest(`
version:
  semver:
    major: 1
    minor: 2
    patch: '{{auto}}'
    labelTemplate: '{{branch}}'
    releaseBranch: master`)

		assert.Nil(t, err)
		assert.Nil(t, manifest.Version.Custom)
		assert.NotNil(t, manifest.Version.SemVer)
		assert.Equal(t, 1, manifest.Version.SemVer.Major)
	})

	t.Run("ReturnsSemverVersionWithPatchDefaultingToAutoPlaceholder", func(t *testing.T) {

		// act
		manifest, err := ReadManifest(`
version:
  semver:
    major: 1
    minor: 2`)

		assert.Nil(t, err)
		assert.Nil(t, manifest.Version.Custom)
		assert.NotNil(t, manifest.Version.SemVer)
		assert.Equal(t, "{{auto}}", manifest.Version.SemVer.Patch)
	})

	t.Run("ReturnsSemverVersionWithLabelTemplateDefaultingToBranchPlaceholder", func(t *testing.T) {

		// act
		manifest, err := ReadManifest(`
version:
  semver:
    major: 1
    minor: 2`)

		assert.Nil(t, err)
		assert.Nil(t, manifest.Version.Custom)
		assert.NotNil(t, manifest.Version.SemVer)
		assert.Equal(t, "{{branch}}", manifest.Version.SemVer.LabelTemplate)
	})

	t.Run("ReturnsSemverVersionWithReleaseBranchDefaultingToMaster", func(t *testing.T) {

		// act
		manifest, err := ReadManifest(`
version:
  semver:
    major: 1
    minor: 2`)

		assert.Nil(t, err)
		assert.Nil(t, manifest.Version.Custom)
		assert.NotNil(t, manifest.Version.SemVer)
		assert.Equal(t, "master", manifest.Version.SemVer.ReleaseBranch)
	})
}

func TestCustomVersion(t *testing.T) {

	t.Run("ReturnsLabelTemplateAsIsWhenItHasNoPlaceholders", func(t *testing.T) {

		version := EstafetteCustomVersion{
			LabelTemplate: "whateveryoulike",
		}
		params := EstafetteVersionParams{
			auto:     5,
			branch:   "release",
			revision: "219aae19153da2b20ac1d88e2fd68e0b20274be2",
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
			auto:     5,
			branch:   "release",
			revision: "219aae19153da2b20ac1d88e2fd68e0b20274be2",
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
			auto:     5,
			branch:   "release",
			revision: "219aae19153da2b20ac1d88e2fd68e0b20274be2",
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
			auto:     5,
			branch:   "release",
			revision: "219aae19153da2b20ac1d88e2fd68e0b20274be2",
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
			ReleaseBranch: "alpha",
		}
		params := EstafetteVersionParams{
			auto:     5,
			branch:   "release",
			revision: "219aae19153da2b20ac1d88e2fd68e0b20274be2",
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
			ReleaseBranch: "alpha",
		}
		params := EstafetteVersionParams{
			auto:     16,
			branch:   "release",
			revision: "219aae19153da2b20ac1d88e2fd68e0b20274be2",
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
			ReleaseBranch: "release",
		}
		params := EstafetteVersionParams{
			auto:     16,
			branch:   "alpha",
			revision: "219aae19153da2b20ac1d88e2fd68e0b20274be2",
		}

		// act
		versionString := version.Version(params)

		assert.Equal(t, "5.3.6-alpha", versionString)
	})

	t.Run("ReturnsSemverWithoutLabelIfBranchMatchesReleaseBranch", func(t *testing.T) {

		version := EstafetteSemverVersion{
			Major:         5,
			Minor:         3,
			Patch:         "6",
			LabelTemplate: "{{branch}}",
			ReleaseBranch: "release",
		}
		params := EstafetteVersionParams{
			auto:     16,
			branch:   "release",
			revision: "219aae19153da2b20ac1d88e2fd68e0b20274be2",
		}

		// act
		versionString := version.Version(params)

		assert.Equal(t, "5.3.6", versionString)
	})
}

func TestGetReservedPropertyNames(t *testing.T) {

	t.Run("ReturnsListWithPropertyNamesAndYamlNames", func(t *testing.T) {

		// act
		names := getReservedPropertyNames()

		// yaml names
		assert.True(t, isReservedPopertyName(names, "image"))
		assert.True(t, isReservedPopertyName(names, "shell"))
		assert.True(t, isReservedPopertyName(names, "workDir"))
		assert.True(t, isReservedPopertyName(names, "commands"))
		assert.True(t, isReservedPopertyName(names, "when"))
		assert.True(t, isReservedPopertyName(names, "env"))

		// property names
		assert.True(t, isReservedPopertyName(names, "Name"))
		assert.True(t, isReservedPopertyName(names, "ContainerImage"))
		assert.True(t, isReservedPopertyName(names, "Shell"))
		assert.True(t, isReservedPopertyName(names, "WorkingDirectory"))
		assert.True(t, isReservedPopertyName(names, "Commands"))
		assert.True(t, isReservedPopertyName(names, "When"))
		assert.True(t, isReservedPopertyName(names, "EnvVars"))
		assert.True(t, isReservedPopertyName(names, "CustomProperties"))
	})
}
