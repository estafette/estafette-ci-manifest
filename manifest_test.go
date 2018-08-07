package manifest

import (
	"encoding/json"
	"strings"
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

	t.Run("ReturnsManifestWithMappedOrderedStagesInSameOrderAsInTheManifest", func(t *testing.T) {

		// act
		manifest, err := ReadManifestFromFile("test-manifest.yaml")

		assert.Nil(t, err)

		assert.Equal(t, 5, len(manifest.Stages))

		assert.Equal(t, "build", manifest.Stages[0].Name)
		assert.Equal(t, "golang:1.8.0-alpine", manifest.Stages[0].ContainerImage)
		assert.Equal(t, "/go/src/github.com/estafette/estafette-ci-builder", manifest.Stages[0].WorkingDirectory)
		assert.Equal(t, "go test -v ./...", manifest.Stages[0].Commands[0])
		assert.Equal(t, "CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./publish/estafette-ci-builder .", manifest.Stages[0].Commands[1])

		assert.Equal(t, "bake", manifest.Stages[1].Name)
		assert.Equal(t, "docker:17.03.0-ce", manifest.Stages[1].ContainerImage)
		assert.Equal(t, "cp Dockerfile ./publish", manifest.Stages[1].Commands[0])
		assert.Equal(t, "docker build -t estafette-ci-builder ./publish", manifest.Stages[1].Commands[1])

		assert.Equal(t, "set-build-status", manifest.Stages[2].Name)
		assert.Equal(t, "extensions/github-status:0.0.2", manifest.Stages[2].ContainerImage)
		assert.Equal(t, 0, len(manifest.Stages[2].Commands))
		assert.Equal(t, "server == 'estafette'", manifest.Stages[2].When)

		assert.Equal(t, "push-to-docker-hub", manifest.Stages[3].Name)
		assert.Equal(t, "docker:17.03.0-ce", manifest.Stages[3].ContainerImage)
		assert.Equal(t, "docker login --username=${ESTAFETTE_DOCKER_HUB_USERNAME} --password='${ESTAFETTE_DOCKER_HUB_PASSWORD}'", manifest.Stages[3].Commands[0])
		assert.Equal(t, "docker push estafette/${ESTAFETTE_LABEL_APP}:${ESTAFETTE_BUILD_VERSION}", manifest.Stages[3].Commands[1])
		assert.Equal(t, "status == 'succeeded' && branch == 'master'", manifest.Stages[3].When)

		assert.Equal(t, "slack-notify", manifest.Stages[4].Name)
		assert.Equal(t, "docker:17.03.0-ce", manifest.Stages[4].ContainerImage)
		assert.Equal(t, "curl -X POST --data-urlencode 'payload={\"channel\": \"#build-status\", \"username\": \"estafette-ci-builder\", \"text\": \"Build ${ESTAFETTE_BUILD_VERSION} for ${ESTAFETTE_LABEL_APP} has failed!\"}' ${ESTAFETTE_SLACK_WEBHOOK}", manifest.Stages[4].Commands[0])
		assert.Equal(t, "status == 'failed' || branch == 'master'", manifest.Stages[4].When)

		assert.Equal(t, "some value with spaces", manifest.Stages[4].EnvVars["SOME_ENVIRONMENT_VAR"])
		assert.Equal(t, "value1", manifest.Stages[4].CustomProperties["unknownProperty1"])
		assert.Equal(t, "value2", manifest.Stages[4].CustomProperties["unknownProperty2"])

		// support custom properties of more complex type, so []string can be used
		// assert.Equal(t, "supported1", manifest.Stages[4].CustomProperties["unknownProperty3"].([]interface{})[0].(string))
		// assert.Equal(t, "supported2", manifest.Stages[4].CustomProperties["unknownProperty3"].([]interface{})[1].(string))

		_, unknownPropertyExist := manifest.Stages[4].CustomProperties["unknownProperty3"]
		assert.True(t, unknownPropertyExist)

		_, reservedPropertyForGolangNameExist := manifest.Stages[4].CustomProperties["ContainerImage"]
		assert.False(t, reservedPropertyForGolangNameExist)

		_, reservedPropertyForYamlNameExist := manifest.Stages[4].CustomProperties["image"]
		assert.False(t, reservedPropertyForYamlNameExist)
	})

	t.Run("ReturnsWorkDirDefaultIfMissing", func(t *testing.T) {

		// act
		manifest, err := ReadManifestFromFile("test-manifest.yaml")

		assert.Nil(t, err)

		assert.Equal(t, "/go/src/github.com/estafette/estafette-ci-builder", manifest.Stages[0].WorkingDirectory)
	})

	t.Run("ReturnsWorkDirIfSet", func(t *testing.T) {

		// act
		manifest, err := ReadManifestFromFile("test-manifest.yaml")

		assert.Nil(t, err)

		assert.Equal(t, "/estafette-work", manifest.Stages[1].WorkingDirectory)
	})

	t.Run("ReturnsShellDefaultIfMissing", func(t *testing.T) {

		// act
		manifest, err := ReadManifestFromFile("test-manifest.yaml")

		assert.Nil(t, err)

		assert.Equal(t, "/bin/sh", manifest.Stages[0].Shell)
	})

	t.Run("ReturnsShellIfSet", func(t *testing.T) {

		// act
		manifest, err := ReadManifestFromFile("test-manifest.yaml")

		assert.Nil(t, err)

		assert.Equal(t, "/bin/bash", manifest.Stages[1].Shell)
	})

	t.Run("ReturnsWhenIfSet", func(t *testing.T) {

		// act
		manifest, err := ReadManifestFromFile("test-manifest.yaml")

		assert.Nil(t, err)

		assert.Equal(t, "status == 'succeeded' && branch == 'master'", manifest.Stages[3].When)
	})

	t.Run("ReturnsWhenDefaultIfMissing", func(t *testing.T) {

		// act
		manifest, err := ReadManifestFromFile("test-manifest.yaml")

		assert.Nil(t, err)

		assert.Equal(t, "status == 'succeeded'", manifest.Stages[0].When)
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

	t.Run("ReturnsManifestWithGlobalEnvVars", func(t *testing.T) {

		// act
		manifest, err := ReadManifestFromFile("test-manifest.yaml")

		assert.Nil(t, err)
		assert.Equal(t, "Greetings", manifest.GlobalEnvVars["VAR_A"])
		assert.Equal(t, "World", manifest.GlobalEnvVars["VAR_B"])
	})

	t.Run("ReturnsManifestWithMappedOrderedReleasesInSameOrderAsInTheManifest", func(t *testing.T) {

		// act
		manifest, err := ReadManifestFromFile("test-manifest.yaml")

		assert.Nil(t, err)

		assert.Equal(t, 6, len(manifest.Releases))

		assert.Equal(t, "docker-hub", manifest.Releases[0].Name)
		assert.Equal(t, "push-image", manifest.Releases[0].Stages[0].Name)
		assert.Equal(t, "extensions/push-to-docker-registry:dev", manifest.Releases[0].Stages[0].ContainerImage)

		assert.Equal(t, "beta", manifest.Releases[1].Name)
		assert.Equal(t, "tag-image", manifest.Releases[1].Stages[0].Name)
		assert.Equal(t, "extensions/tag-container-image:dev", manifest.Releases[1].Stages[0].ContainerImage)

		assert.Equal(t, "development", manifest.Releases[2].Name)
		assert.Equal(t, "deploy", manifest.Releases[2].Stages[0].Name)
		assert.Equal(t, "extensions/deploy-to-kubernetes-engine:dev", manifest.Releases[2].Stages[0].ContainerImage)

		assert.Equal(t, "staging", manifest.Releases[3].Name)
		assert.Equal(t, "deploy", manifest.Releases[3].Stages[0].Name)
		assert.Equal(t, "extensions/deploy-to-kubernetes-engine:beta", manifest.Releases[3].Stages[0].ContainerImage)

		assert.Equal(t, "production", manifest.Releases[4].Name)
		assert.Equal(t, "deploy", manifest.Releases[4].Stages[0].Name)
		assert.Equal(t, "extensions/deploy-to-kubernetes-engine:stable", manifest.Releases[4].Stages[0].ContainerImage)
		assert.Equal(t, "create-release-notes", manifest.Releases[4].Stages[1].Name)
		assert.Equal(t, "extensions/create-release-notes-from-changelog:stable", manifest.Releases[4].Stages[1].ContainerImage)

		assert.Equal(t, "tooling", manifest.Releases[5].Name)
	})
}

func TestVersion(t *testing.T) {

	t.Run("ReturnsSemverVersionByDefaultIfNoOtherVersionTypeIsSet", func(t *testing.T) {

		// act
		manifest, err := ReadManifest("")

		assert.Nil(t, err)
		assert.Nil(t, manifest.Version.Custom)
		assert.NotNil(t, manifest.Version.SemVer)
		assert.Equal(t, 0, manifest.Version.SemVer.Major)
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

	t.Run("ReturnsSemverVersionWithMajorDefaultingToZero", func(t *testing.T) {

		// act
		manifest, err := ReadManifest(`
version:
  semver:
    minor: 2`)

		assert.Nil(t, err)
		assert.Nil(t, manifest.Version.Custom)
		assert.NotNil(t, manifest.Version.SemVer)
		assert.Equal(t, 0, manifest.Version.SemVer.Major)
	})

	t.Run("ReturnsSemverVersionWithMinorDefaultingToZero", func(t *testing.T) {

		// act
		manifest, err := ReadManifest(`
version:
  semver:
    major: 1`)

		assert.Nil(t, err)
		assert.Nil(t, manifest.Version.Custom)
		assert.NotNil(t, manifest.Version.SemVer)
		assert.Equal(t, 0, manifest.Version.SemVer.Minor)
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
			ReleaseBranch: "alpha",
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
			ReleaseBranch: "alpha",
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
			ReleaseBranch: "release",
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

	t.Run("ReturnsSemverWithoutLabelIfBranchMatchesReleaseBranch", func(t *testing.T) {

		version := EstafetteSemverVersion{
			Major:         5,
			Minor:         3,
			Patch:         "6",
			LabelTemplate: "{{branch}}",
			ReleaseBranch: "release",
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
}

func TestJsonMarshalling(t *testing.T) {

	t.Run("ReturnsStagesAsPipelines", func(t *testing.T) {

		manifest := EstafetteManifest{
			Stages: []*EstafetteStage{
				&EstafetteStage{
					Name: "build",
				},
			},
		}

		// act
		data, err := json.Marshal(manifest)

		assert.Nil(t, err)
		assert.True(t, strings.Contains(string(data), "Pipelines"))
		assert.False(t, strings.Contains(string(data), "Stages"))
	})
}
