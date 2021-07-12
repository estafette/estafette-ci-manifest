package manifest

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v2"
)

func TestReadManifestFromFile(t *testing.T) {

	t.Run("ReturnsErrorForManifestWithUnknownSections", func(t *testing.T) {

		// act
		_, err := ReadManifestFromFile(GetDefaultManifestPreferences(), "test-non-strict-manifest.yaml", true)

		assert.NotNil(t, err)
	})

	t.Run("ReturnsManifestWithoutErrors", func(t *testing.T) {

		// act
		_, err := ReadManifestFromFile(GetDefaultManifestPreferences(), "test-manifest.yaml", true)

		assert.Nil(t, err)
	})

	t.Run("ReturnsManifestWithArchivedTrueWithoutErrors", func(t *testing.T) {

		// act
		_, err := ReadManifestFromFile(GetDefaultManifestPreferences(), "test-manifest-archived.yaml", true)

		assert.Nil(t, err)
	})

	t.Run("ReturnsManifestWithMappedLabels", func(t *testing.T) {

		// act
		manifest, err := ReadManifestFromFile(GetDefaultManifestPreferences(), "test-manifest.yaml", true)

		assert.Nil(t, err)
		assert.Equal(t, "estafette-ci-builder", manifest.Labels["app"])
		assert.Equal(t, "estafette-team", manifest.Labels["team"])
		assert.Equal(t, "golang", manifest.Labels["language"])
	})

	t.Run("ReturnsManifestWithMappedBuilderTrack", func(t *testing.T) {

		// act
		manifest, err := ReadManifestFromFile(GetDefaultManifestPreferences(), "test-manifest.yaml", true)

		assert.Nil(t, err)
		assert.Equal(t, "nanoserver-1809-stable", manifest.Builder.Track)
	})

	t.Run("ReturnsManifestWithBuilderTrackDefaultWindowsServerCore1809", func(t *testing.T) {

		// act
		manifest, err := ReadManifest(GetDefaultManifestPreferences(), `
stages:
  hi:
    image: alpine
    commands:
    - echo 'hi'`, true)

		assert.Nil(t, err)
		assert.Equal(t, "stable", manifest.Builder.Track)
	})

	t.Run("ReturnsManifestWithMappedBuilderOperatingSystem", func(t *testing.T) {

		// act
		manifest, err := ReadManifestFromFile(GetDefaultManifestPreferences(), "test-manifest.yaml", true)

		assert.Nil(t, err)
		assert.Equal(t, OperatingSystemWindows, manifest.Builder.OperatingSystem)
	})

	t.Run("ReturnsManifestWithBuilderOperatingSystemDefaultLinux", func(t *testing.T) {

		// act
		manifest, err := ReadManifest(GetDefaultManifestPreferences(), `
stages:
  hi:
    image: alpine
    commands:
    - echo 'hi'`, true)

		assert.Nil(t, err)
		assert.Equal(t, OperatingSystemLinux, manifest.Builder.OperatingSystem)
	})

	t.Run("ReturnsManifestWithBuilder", func(t *testing.T) {

		// act
		manifest, err := ReadManifestFromFile(GetDefaultManifestPreferences(), "test-manifest.yaml", true)

		assert.Nil(t, err)
		assert.Equal(t, "nanoserver-1809-stable", manifest.Builder.Track)
		assert.Equal(t, OperatingSystemWindows, manifest.Builder.OperatingSystem)
	})

	t.Run("ReturnsManifestWithBuilderForReleaseIfNotOverridden", func(t *testing.T) {

		// act
		manifest, err := ReadManifestFromFile(GetDefaultManifestPreferences(), "test-manifest.yaml", true)

		assert.Nil(t, err)
		assert.NotNil(t, manifest.Releases[1].Builder)
		assert.Equal(t, "nanoserver-1809-stable", manifest.Releases[1].Builder.Track)
		assert.Equal(t, OperatingSystemWindows, manifest.Releases[1].Builder.OperatingSystem)
	})

	t.Run("ReturnsManifestWithReleaseBuilderIfOverridden", func(t *testing.T) {

		// act
		manifest, err := ReadManifestFromFile(GetDefaultManifestPreferences(), "test-manifest.yaml", true)

		assert.Nil(t, err)
		assert.NotNil(t, manifest.Releases[2].Builder)
		assert.Equal(t, "stable", manifest.Releases[2].Builder.Track)
		assert.Equal(t, OperatingSystemLinux, manifest.Releases[2].Builder.OperatingSystem)
	})

	t.Run("ReturnsManifestWithMappedOrderedStagesInSameOrderAsInTheManifest", func(t *testing.T) {

		// act
		manifest, err := ReadManifestFromFile(GetDefaultManifestPreferences(), "test-manifest.yaml", true)

		assert.Nil(t, err)

		assert.Equal(t, 7, len(manifest.Stages))

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
		assert.Equal(t, "status == 'succeeded' && branch == 'main'", manifest.Stages[3].When)

		assert.Equal(t, "slack-notify", manifest.Stages[4].Name)
		assert.Equal(t, "docker:17.03.0-ce", manifest.Stages[4].ContainerImage)
		assert.Equal(t, "curl -X POST --data-urlencode 'payload={\"channel\": \"#build-status\", \"username\": \"estafette-ci-builder\", \"text\": \"Build ${ESTAFETTE_BUILD_VERSION} for ${ESTAFETTE_LABEL_APP} has failed!\"}' ${ESTAFETTE_SLACK_WEBHOOK}", manifest.Stages[4].Commands[0])
		assert.Equal(t, "status == 'failed' || branch == 'main'", manifest.Stages[4].When)

		assert.Equal(t, "some value with spaces", manifest.Stages[4].EnvVars["SOME_ENVIRONMENT_VAR"])
		assert.Equal(t, "value1", manifest.Stages[4].CustomProperties["unknownProperty1"])
		assert.Equal(t, "value2", manifest.Stages[4].CustomProperties["unknownProperty2"])

		// support custom properties of more complex type, so []string can be used
		assert.Equal(t, "supported1", manifest.Stages[4].CustomProperties["unknownProperty3"].([]interface{})[0].(string))
		assert.Equal(t, "supported2", manifest.Stages[4].CustomProperties["unknownProperty3"].([]interface{})[1].(string))

		// support custom properties of more even more complex type, so another map[string]interface can be used
		assert.Equal(t, "extensions", manifest.Stages[5].CustomProperties["container"].(map[string]interface{})["repository"].(string))
		assert.Equal(t, "gke", manifest.Stages[5].CustomProperties["container"].(map[string]interface{})["name"].(string))
		assert.Equal(t, "alpha", manifest.Stages[5].CustomProperties["container"].(map[string]interface{})["tag"].(string))

		_, unknownPropertyExist := manifest.Stages[4].CustomProperties["unknownProperty3"]
		assert.True(t, unknownPropertyExist)

		_, reservedPropertyForGolangNameExist := manifest.Stages[4].CustomProperties["ContainerImage"]
		assert.False(t, reservedPropertyForGolangNameExist)

		_, reservedPropertyForYamlNameExist := manifest.Stages[4].CustomProperties["image"]
		assert.False(t, reservedPropertyForYamlNameExist)
	})

	t.Run("ReturnsWorkDirDefaultIfMissing", func(t *testing.T) {

		// act
		manifest, err := ReadManifestFromFile(GetDefaultManifestPreferences(), "test-manifest.yaml", true)

		assert.Nil(t, err)

		assert.Equal(t, "/go/src/github.com/estafette/estafette-ci-builder", manifest.Stages[0].WorkingDirectory)
	})

	t.Run("ReturnsWorkDirIfSet", func(t *testing.T) {

		// act
		manifest, err := ReadManifestFromFile(GetDefaultManifestPreferences(), "test-manifest.yaml", true)

		assert.Nil(t, err)

		assert.Equal(t, "C:/estafette-work", manifest.Stages[1].WorkingDirectory)
	})

	t.Run("ReturnsShellDefaultIfMissing", func(t *testing.T) {

		// act
		manifest, err := ReadManifestFromFile(GetDefaultManifestPreferences(), "test-manifest.yaml", true)

		assert.Nil(t, err)

		assert.Equal(t, "powershell", manifest.Stages[0].Shell)
	})

	t.Run("ReturnsShellIfSet", func(t *testing.T) {

		// act
		manifest, err := ReadManifestFromFile(GetDefaultManifestPreferences(), "test-manifest.yaml", true)

		assert.Nil(t, err)

		assert.Equal(t, "/bin/bash", manifest.Stages[1].Shell)
	})

	t.Run("ReturnsWhenIfSet", func(t *testing.T) {

		// act
		manifest, err := ReadManifestFromFile(GetDefaultManifestPreferences(), "test-manifest.yaml", true)

		assert.Nil(t, err)

		assert.Equal(t, "status == 'succeeded' && branch == 'main'", manifest.Stages[3].When)
	})

	t.Run("ReturnsWhenDefaultIfMissing", func(t *testing.T) {

		// act
		manifest, err := ReadManifestFromFile(GetDefaultManifestPreferences(), "test-manifest.yaml", true)

		assert.Nil(t, err)

		assert.Equal(t, "status == 'succeeded'", manifest.Stages[0].When)
	})

	t.Run("ReturnsManifestWithSemverVersion", func(t *testing.T) {

		// act
		manifest, err := ReadManifestFromFile(GetDefaultManifestPreferences(), "test-manifest.yaml", true)

		assert.Nil(t, err)
		assert.Equal(t, 1, manifest.Version.SemVer.Major)
		assert.Equal(t, 2, manifest.Version.SemVer.Minor)
		assert.Equal(t, "{{auto}}", manifest.Version.SemVer.Patch)
		assert.Equal(t, "{{branch}}", manifest.Version.SemVer.LabelTemplate)
		assert.Equal(t, "main", manifest.Version.SemVer.ReleaseBranch.Values[0])
	})

	t.Run("ReturnsManifestWithGlobalEnvVars", func(t *testing.T) {

		// act
		manifest, err := ReadManifestFromFile(GetDefaultManifestPreferences(), "test-manifest.yaml", true)

		assert.Nil(t, err)
		assert.Equal(t, "Greetings", manifest.GlobalEnvVars["VAR_A"])
		assert.Equal(t, "World", manifest.GlobalEnvVars["VAR_B"])
	})

	t.Run("ReturnsManifestWithMappedOrderedReleasesInSameOrderAsInTheManifest", func(t *testing.T) {

		// act
		manifest, err := ReadManifestFromFile(GetDefaultManifestPreferences(), "test-manifest.yaml", true)

		assert.Nil(t, err)

		if assert.Equal(t, 6, len(manifest.Releases)) {

			assert.Equal(t, "docker-hub", manifest.Releases[0].Name)
			assert.False(t, *manifest.Releases[0].CloneRepository)
			assert.Equal(t, "push-image", manifest.Releases[0].Stages[0].Name)
			assert.Equal(t, "extensions/push-to-docker-registry:dev", manifest.Releases[0].Stages[0].ContainerImage)

			assert.Equal(t, "beta", manifest.Releases[1].Name)
			assert.False(t, *manifest.Releases[1].CloneRepository)
			assert.Equal(t, "tag-container-image", manifest.Releases[1].Stages[0].Name)
			assert.Equal(t, "extensions/docker:stable", manifest.Releases[1].Stages[0].ContainerImage)
			assert.Equal(t, 1, len(manifest.Releases[1].Stages[0].CustomProperties["tags"].([]interface{})))
			assert.Equal(t, "beta", manifest.Releases[1].Stages[0].CustomProperties["tags"].([]interface{})[0])

			assert.Equal(t, "development", manifest.Releases[2].Name)
			assert.NotNil(t, manifest.Releases[2].Builder)
			assert.Equal(t, "stable", manifest.Releases[2].Builder.Track)
			assert.Equal(t, OperatingSystemLinux, manifest.Releases[2].Builder.OperatingSystem)
			assert.False(t, *manifest.Releases[2].CloneRepository)
			assert.Equal(t, "deploy", manifest.Releases[2].Stages[0].Name)
			assert.Equal(t, "extensions/deploy-to-kubernetes-engine:dev", manifest.Releases[2].Stages[0].ContainerImage)

			assert.Equal(t, "staging", manifest.Releases[3].Name)
			assert.False(t, *manifest.Releases[3].CloneRepository)
			assert.Equal(t, "deploy", manifest.Releases[3].Stages[0].Name)
			assert.Equal(t, "extensions/gke:beta", manifest.Releases[3].Stages[0].ContainerImage)
			assert.Equal(t, 600, manifest.Releases[3].Stages[0].CustomProperties["volumemounts"].([]interface{})[0].(map[string]interface{})["volume"].(map[string]interface{})["secret"].(map[string]interface{})["items"].([]interface{})[0].(map[string]interface{})["mode"])
			assert.Equal(t, true, manifest.Releases[3].Stages[0].CustomProperties["volumemounts"].([]interface{})[0].(map[string]interface{})["volume"].(map[string]interface{})["secret"].(map[string]interface{})["items"].([]interface{})[0].(map[string]interface{})["enabled"])

			assert.Equal(t, "production", manifest.Releases[4].Name)
			assert.True(t, *manifest.Releases[4].CloneRepository)
			assert.Equal(t, "deploy", manifest.Releases[4].Stages[0].Name)
			assert.Equal(t, "extensions/deploy-to-kubernetes-engine:stable", manifest.Releases[4].Stages[0].ContainerImage)
			assert.Equal(t, "create-release-notes", manifest.Releases[4].Stages[1].Name)
			assert.Equal(t, "extensions/create-release-notes-from-changelog:stable", manifest.Releases[4].Stages[1].ContainerImage)

			assert.Equal(t, "tooling", manifest.Releases[5].Name)
			assert.False(t, *manifest.Releases[5].CloneRepository)
		}
	})

	t.Run("ReturnsReleaseTargetWithActions", func(t *testing.T) {

		// act
		manifest, err := ReadManifestFromFile(GetDefaultManifestPreferences(), "test-manifest.yaml", true)

		assert.Nil(t, err)

		if assert.Equal(t, 6, len(manifest.Releases)) {

			assert.Equal(t, "production", manifest.Releases[4].Name)
			assert.Equal(t, 3, len(manifest.Releases[4].Actions))
			assert.Equal(t, "deploy-canary", manifest.Releases[4].Actions[0].Name)
			assert.Equal(t, "rollback-canary", manifest.Releases[4].Actions[1].Name)
			assert.Equal(t, "deploy-stable", manifest.Releases[4].Actions[2].Name)
		}
	})

	t.Run("ReturnsManifestWithReleasesUsingReleaseTemplates", func(t *testing.T) {

		// act
		manifest, err := ReadManifestFromFile(GetDefaultManifestPreferences(), "test-manifest-with-template.yaml", true)

		assert.Nil(t, err)

		if assert.Equal(t, 2, len(manifest.Releases)) {

			assert.Equal(t, "development", manifest.Releases[0].Name)
			assert.NotNil(t, manifest.Releases[0].Builder)
			assert.Equal(t, "stable", manifest.Releases[0].Builder.Track)
			assert.Equal(t, OperatingSystemLinux, manifest.Releases[0].Builder.OperatingSystem)
			assert.False(t, *manifest.Releases[0].CloneRepository)
			assert.Equal(t, "deploy", manifest.Releases[0].Stages[0].Name)
			assert.Equal(t, "extensions/deploy-to-kubernetes-engine:dev", manifest.Releases[0].Stages[0].ContainerImage)

			assert.Equal(t, "staging", manifest.Releases[1].Name)
			assert.NotNil(t, manifest.Releases[1].Builder)
			assert.Equal(t, "stable", manifest.Releases[1].Builder.Track)
			assert.Equal(t, OperatingSystemLinux, manifest.Releases[1].Builder.OperatingSystem)
			assert.True(t, *manifest.Releases[1].CloneRepository)
			assert.Equal(t, "deploy", manifest.Releases[1].Stages[0].Name)
			assert.Equal(t, "extensions/deploy-to-kubernetes-engine:dev", manifest.Releases[1].Stages[0].ContainerImage)
		}
	})

	t.Run("UmarshallingManifestWithDeprecatedPipelinesVerbStillWorks", func(t *testing.T) {

		input := `
pipelines:
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

		// act
		manifest, err := ReadManifest(GetDefaultManifestPreferences(), input, true)

		assert.Nil(t, err)
		assert.Equal(t, 2, len(manifest.Stages))
	})

	t.Run("ReturnsManifestWithTriggersWithoutErrors", func(t *testing.T) {

		// act
		_, err := ReadManifestFromFile(GetDefaultManifestPreferences(), "test-manifest-with-triggers.yaml", true)

		assert.Nil(t, err)
	})

	t.Run("ReturnsTriggers", func(t *testing.T) {

		// act
		manifest, err := ReadManifestFromFile(GetDefaultManifestPreferences(), "test-manifest-with-triggers.yaml", true)

		assert.Nil(t, err)

		if assert.Equal(t, 4, len(manifest.Triggers)) {
			assert.Equal(t, "github.com/estafette/estafette-ci-manifest", manifest.Triggers[0].Pipeline.Name)
			assert.Equal(t, "github.com/estafette/estafette-ci-builder", manifest.Triggers[1].Git.Repository)
			assert.Equal(t, "golang", manifest.Triggers[2].Docker.Image)
			assert.Equal(t, "topic-name", manifest.Triggers[3].PubSub.Topic)
		}
	})

	t.Run("ReturnsReleaseTargetWithTriggers", func(t *testing.T) {

		// act
		manifest, err := ReadManifestFromFile(GetDefaultManifestPreferences(), "test-manifest-with-triggers.yaml", true)

		assert.Nil(t, err)

		if assert.Equal(t, 2, len(manifest.Releases)) {
			assert.Equal(t, "development", manifest.Releases[0].Name)
			assert.Equal(t, 3, len(manifest.Releases[0].Triggers))
			assert.Equal(t, "github.com/estafette/estafette-ci-builder", manifest.Releases[0].Triggers[0].Pipeline.Name)
			assert.Equal(t, "0 10 */1 * *", manifest.Releases[0].Triggers[1].Cron.Schedule)
			assert.Equal(t, "topic-name", manifest.Releases[0].Triggers[2].PubSub.Topic)
		}
	})

	t.Run("ReturnsManifestWithNestedParallelStages", func(t *testing.T) {

		// act
		manifest, err := ReadManifestFromFile(GetDefaultManifestPreferences(), "test-manifest.yaml", true)

		assert.Nil(t, err)

		assert.Equal(t, "parallel-stages-group", manifest.Stages[6].Name)
		assert.Equal(t, 2, len(manifest.Stages[6].ParallelStages))
		assert.Equal(t, "stageA", manifest.Stages[6].ParallelStages[0].Name)
		assert.Equal(t, "stageB", manifest.Stages[6].ParallelStages[1].Name)

		assert.Equal(t, "staging", manifest.Releases[3].Name)
		assert.Equal(t, 2, len(manifest.Releases[3].Stages[1].ParallelStages))
		assert.Equal(t, "stageA", manifest.Releases[3].Stages[1].ParallelStages[0].Name)
		assert.Equal(t, "stageB", manifest.Releases[3].Stages[1].ParallelStages[1].Name)
	})

	t.Run("ReturnsManifestWithStageServices", func(t *testing.T) {

		// act
		manifest, err := ReadManifestFromFile(GetDefaultManifestPreferences(), "test-manifest.yaml", true)

		assert.Nil(t, err)

		assert.Equal(t, "test-alpha-version", manifest.Stages[5].Name)
		assert.Equal(t, 2, len(manifest.Stages[5].Services))

		assert.Equal(t, "kubernetes", manifest.Stages[5].Services[0].Name)
		assert.Equal(t, "bsycorp/kind:latest-1.15", manifest.Stages[5].Services[0].ContainerImage)
		assert.Equal(t, 1, len(manifest.Stages[5].Services[0].EnvVars))
		assert.Equal(t, "some value with spaces", manifest.Stages[5].Services[0].EnvVars["SOME_ENVIRONMENT_VAR"])
		assert.NotNil(t, manifest.Stages[5].Services[0].Readiness)
		assert.Equal(t, "/kubernetes-ready", manifest.Stages[5].Services[0].Readiness.Path)
		assert.Equal(t, 60, manifest.Stages[5].Services[0].Readiness.TimeoutSeconds)

		assert.Equal(t, "database", manifest.Stages[5].Services[1].Name)
		assert.Equal(t, "cockroachdb/cockroach:v19.1.5", manifest.Stages[5].Services[1].ContainerImage)
		assert.Equal(t, "cockroachdb start --insecure --listen-addr=localhost", manifest.Stages[5].Services[1].Commands[0])
	})

	t.Run("ReturnsManifestWithMappedOrderedBotsInSameOrderAsInTheManifest", func(t *testing.T) {

		// act
		manifest, err := ReadManifestFromFile(GetDefaultManifestPreferences(), "test-manifest-with-bots.yaml", true)

		assert.Nil(t, err)

		if assert.Equal(t, 3, len(manifest.Bots)) {

			assert.Equal(t, "pr-bot", manifest.Bots[0].Name)
			assert.Equal(t, 1, len(manifest.Bots[0].Triggers))
			assert.Equal(t, 9, len(manifest.Bots[0].Triggers[0].Bitbucket.Events))
			assert.Equal(t, "self", manifest.Bots[0].Triggers[0].Bitbucket.Repository)
			assert.Equal(t, "pr-bot", manifest.Bots[0].Triggers[0].BotAction.Bot)
			assert.Equal(t, "welcome", manifest.Bots[0].Stages[0].Name)
			assert.Equal(t, "extensions/github-pr-bot:stable", manifest.Bots[0].Stages[0].ContainerImage)

			assert.Equal(t, "any-bot", manifest.Bots[1].Name)
			assert.Equal(t, 2, len(manifest.Bots[1].Triggers))
			assert.Equal(t, 20, len(manifest.Bots[1].Triggers[0].Bitbucket.Events))
			assert.Equal(t, "self", manifest.Bots[1].Triggers[0].Bitbucket.Repository)
			assert.Equal(t, "any-bot", manifest.Bots[1].Triggers[0].BotAction.Bot)
			assert.Equal(t, 34, len(manifest.Bots[1].Triggers[1].Github.Events))
			assert.Equal(t, "self", manifest.Bots[1].Triggers[1].Github.Repository)
			assert.Equal(t, "any-bot", manifest.Bots[1].Triggers[1].BotAction.Bot)
			assert.Equal(t, "backup", manifest.Bots[1].Stages[0].Name)
			assert.Equal(t, "extensions/backup-event:stable", manifest.Bots[1].Stages[0].ContainerImage)

			assert.Equal(t, "stale-issues-bot", manifest.Bots[2].Name)
			assert.Equal(t, 1, len(manifest.Bots[2].Triggers))
			assert.Equal(t, "0 10 * * *", manifest.Bots[2].Triggers[0].Cron.Schedule)
			assert.Equal(t, "stale-issues-bot", manifest.Bots[2].Triggers[0].BotAction.Bot)
			assert.Equal(t, true, *manifest.Bots[2].CloneRepository)
			assert.Equal(t, "welcome", manifest.Bots[2].Stages[0].Name)
			assert.Equal(t, "extensions/github-stale-issue-bot:stable", manifest.Bots[2].Stages[0].ContainerImage)
		}
	})
}

func TestVersion(t *testing.T) {

	t.Run("ReturnsSemverVersionByDefaultIfNoOtherVersionTypeIsSet", func(t *testing.T) {

		// act
		manifest, err := ReadManifest(GetDefaultManifestPreferences(), `
stages:
  git-clone:
    image: extensions/git-clone`, true)

		if assert.Nil(t, err) {
			assert.Nil(t, manifest.Version.Custom)
			assert.NotNil(t, manifest.Version.SemVer)
			assert.Equal(t, 0, manifest.Version.SemVer.Major)
		}
	})

	t.Run("ReturnsCustomVersionWithLabelTemplateDefaultingToRevisionPlaceholder", func(t *testing.T) {

		// act
		manifest, err := ReadManifest(GetDefaultManifestPreferences(), `
version:
  custom:
    labelTemplate: ''

stages:
  git-clone:
    image: extensions/git-clone`, true)

		if assert.Nil(t, err) {
			assert.Nil(t, manifest.Version.SemVer)
			assert.NotNil(t, manifest.Version.Custom)
			assert.Equal(t, "{{revision}}", manifest.Version.Custom.LabelTemplate)
		}
	})

	t.Run("ReturnsSemverVersionIfSemverIsSet", func(t *testing.T) {

		// act
		manifest, err := ReadManifest(GetDefaultManifestPreferences(), `
version:
  semver:
    major: 1
    minor: 2
    patch: '{{auto}}'
    labelTemplate: '{{branch}}'
    releaseBranch: main

stages:
  git-clone:
    image: extensions/git-clone`, true)

		if assert.Nil(t, err) {
			assert.Nil(t, manifest.Version.Custom)
			assert.NotNil(t, manifest.Version.SemVer)
			assert.Equal(t, 1, manifest.Version.SemVer.Major)
		}
	})

	t.Run("ReturnsSemverVersionWithMajorDefaultingToZero", func(t *testing.T) {

		// act
		manifest, err := ReadManifest(GetDefaultManifestPreferences(), `
version:
  semver:
    minor: 2

stages:
  git-clone:
    image: extensions/git-clone`, true)

		if assert.Nil(t, err) {
			assert.Nil(t, manifest.Version.Custom)
			assert.NotNil(t, manifest.Version.SemVer)
			assert.Equal(t, 0, manifest.Version.SemVer.Major)
		}
	})

	t.Run("ReturnsSemverVersionWithMinorDefaultingToZero", func(t *testing.T) {

		// act
		manifest, err := ReadManifest(GetDefaultManifestPreferences(), `
version:
  semver:
    major: 1

stages:
  git-clone:
    image: extensions/git-clone`, true)

		if assert.Nil(t, err) {
			assert.Nil(t, manifest.Version.Custom)
			assert.NotNil(t, manifest.Version.SemVer)
			assert.Equal(t, 0, manifest.Version.SemVer.Minor)
		}
	})

	t.Run("ReturnsSemverVersionWithPatchDefaultingToAutoPlaceholder", func(t *testing.T) {

		// act
		manifest, err := ReadManifest(GetDefaultManifestPreferences(), `
version:
  semver:
    major: 1
    minor: 2

stages:
  git-clone:
    image: extensions/git-clone`, true)

		if assert.Nil(t, err) {
			assert.Nil(t, manifest.Version.Custom)
			assert.NotNil(t, manifest.Version.SemVer)
			assert.Equal(t, "{{auto}}", manifest.Version.SemVer.Patch)
		}
	})

	t.Run("ReturnsSemverVersionWithLabelTemplateDefaultingToBranchPlaceholder", func(t *testing.T) {

		// act
		manifest, err := ReadManifest(GetDefaultManifestPreferences(), `
version:
  semver:
    major: 1
    minor: 2

stages:
  git-clone:
    image: extensions/git-clone`, true)

		if assert.Nil(t, err) {
			assert.Nil(t, manifest.Version.Custom)
			assert.NotNil(t, manifest.Version.SemVer)
			assert.Equal(t, "{{branch}}", manifest.Version.SemVer.LabelTemplate)
		}
	})

	t.Run("ReturnsSemverVersionWithReleaseBranchDefaultingToMasterAndMAin", func(t *testing.T) {

		// act
		manifest, err := ReadManifest(GetDefaultManifestPreferences(), `
version:
  semver:
    major: 1
    minor: 2

stages:
  git-clone:
    image: extensions/git-clone`, true)

		if assert.Nil(t, err) {
			assert.Nil(t, manifest.Version.Custom)
			assert.NotNil(t, manifest.Version.SemVer)
			assert.Equal(t, 2, len(manifest.Version.SemVer.ReleaseBranch.Values))
			assert.Equal(t, "master", manifest.Version.SemVer.ReleaseBranch.Values[0])
			assert.Equal(t, "main", manifest.Version.SemVer.ReleaseBranch.Values[1])
		}
	})
}

func TestManifestToJsonMarshalling(t *testing.T) {

	t.Run("ReturnsStagesAsStages", func(t *testing.T) {

		manifest := EstafetteManifest{
			Stages: []*EstafetteStage{
				&EstafetteStage{
					Name: "build",
				},
			},
		}

		// act
		data, err := json.Marshal(manifest)

		if assert.Nil(t, err) {
			assert.False(t, strings.Contains(string(data), "Pipelines"))
			assert.True(t, strings.Contains(string(data), "\"Stages\""))
		}
	})
}

func TestManifestToYamlMarshalling(t *testing.T) {
	t.Run("UnmarshallingThenMarshallingReturnsTheSameFile", func(t *testing.T) {

		var manifest EstafetteManifest

		input := `builder:
  track: stable
  os: windows
labels:
  app: estafette-ci-builder
  language: golang
  team: estafette-team
version:
  semver:
    major: 0
    minor: 0
    patch: '{{auto}}'
    labelTemplate: '{{branch}}'
    releaseBranch: main
env:
  VAR_A: Greetings
  VAR_B: World
stages:
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
releases:
  staging:
    stages:
      deploy:
        image: extensions/deploy-to-kubernetes-engine:stable
        shell: /bin/sh
        workDir: /estafette-work
        when: status == 'succeeded'
  production:
    stages:
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
		err := yaml.Unmarshal([]byte(input), &manifest)
		assert.Nil(t, err)

		// act
		output, err := yaml.Marshal(manifest)

		assert.Nil(t, err)
		assert.Equal(t, input, string(output))
	})

	t.Run("UnmarshallingWithArchivedTrueThenMarshallingReturnsTheSameFile", func(t *testing.T) {

		var manifest EstafetteManifest

		input := `archived: true
builder:
  track: stable
  os: windows
labels:
  app: estafette-ci-builder
  language: golang
  team: estafette-team
version:
  semver:
    major: 0
    minor: 0
    patch: '{{auto}}'
    labelTemplate: '{{branch}}'
    releaseBranch: main
env:
  VAR_A: Greetings
  VAR_B: World
stages:
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
releases:
  staging:
    builder:
      track: stable
      os: windows
    stages:
      deploy:
        image: extensions/deploy-to-kubernetes-engine:stable
        shell: /bin/sh
        workDir: /estafette-work
        when: status == 'succeeded'
  production:
    builder:
      track: stable
      os: windows
    stages:
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
		err := yaml.Unmarshal([]byte(input), &manifest)
		assert.Nil(t, err)

		// act
		output, err := yaml.Marshal(manifest)

		assert.Nil(t, err)
		assert.Equal(t, input, string(output))
	})
}

func TestManifestToJSONMarshalling(t *testing.T) {

	t.Run("UnmarshallingYamlThenMarshalToJSONForNestedCustomProperties", func(t *testing.T) {

		var manifest EstafetteManifest

		input := `builder:
  track: stable
  os: windows
labels:
  app: estafette-ci-builder
  language: golang
  team: estafette-team
version:
  semver:
    major: 0
    minor: 0
    patch: '{{auto}}'
    labelTemplate: '{{branch}}'
    releaseBranch: main
stages:
  test-alpha-version:
    image: extensions/gke:${ESTAFETTE_BUILD_VERSION}
    credentials: gke-tooling
    app: gke
    namespace: estafette
    visibility: private
    container:
      repository: extensions
      name: gke
      tag: alpha
    cpu:
      request: 100m
      limit: 100m
    memory:
      request: 256Mi
      limit: 256Mi
    dryrun: true
`
		err := yaml.Unmarshal([]byte(input), &manifest)
		assert.Nil(t, err)
		manifest.SetDefaults(*GetDefaultManifestPreferences())

		// act
		output, err := json.Marshal(manifest)

		if assert.Nil(t, err) {
			assert.Equal(t, "{\"Archived\":false,\"Builder\":{\"Track\":\"stable\",\"OperatingSystem\":\"windows\",\"StorageMedium\":\"\",\"BuilderType\":\"docker\"},\"Labels\":{\"app\":\"estafette-ci-builder\",\"language\":\"golang\",\"team\":\"estafette-team\"},\"Version\":{\"SemVer\":{\"Major\":0,\"Minor\":0,\"Patch\":\"{{auto}}\",\"LabelTemplate\":\"{{branch}}\",\"ReleaseBranch\":\"main\"}},\"GlobalEnvVars\":null,\"Triggers\":null,\"Stages\":[{\"Name\":\"test-alpha-version\",\"ContainerImage\":\"extensions/gke:${ESTAFETTE_BUILD_VERSION}\",\"Shell\":\"powershell\",\"WorkingDirectory\":\"C:/estafette-work\",\"When\":\"status == 'succeeded'\",\"CustomProperties\":{\"app\":\"gke\",\"container\":{\"name\":\"gke\",\"repository\":\"extensions\",\"tag\":\"alpha\"},\"cpu\":{\"limit\":\"100m\",\"request\":\"100m\"},\"credentials\":\"gke-tooling\",\"dryrun\":true,\"memory\":{\"limit\":\"256Mi\",\"request\":\"256Mi\"},\"namespace\":\"estafette\",\"visibility\":\"private\"}}],\"Releases\":null,\"ReleaseTemplates\":null,\"Bots\":null}", string(output))
		}
	})
}

func TestGetAllTriggers(t *testing.T) {
	t.Run("ReturnsEmptyArrayIfNoTriggersAreDefined", func(t *testing.T) {

		manifest := EstafetteManifest{}

		// act
		triggers := manifest.GetAllTriggers("github.com", "estafette", "estafette-ci-manifest")

		assert.Equal(t, 0, len(triggers))
	})

	t.Run("ReturnsReleaseTriggersIfNoBuildTriggersAreDefined", func(t *testing.T) {

		manifest := EstafetteManifest{
			Stages: []*EstafetteStage{
				{
					Name: "build",
				},
			},
			Releases: []*EstafetteRelease{
				{
					Name: "tooling",
					Triggers: []*EstafetteTrigger{
						{
							Pipeline:      &EstafettePipelineTrigger{},
							ReleaseAction: &EstafetteTriggerReleaseAction{},
						},
					},
				},
			},
		}

		// act
		triggers := manifest.GetAllTriggers("github.com", "estafette", "estafette-ci-manifest")

		assert.Equal(t, 1, len(triggers))
	})

	t.Run("ReturnsBuildAndReleaseTriggersIfBothAreDefined", func(t *testing.T) {

		manifest := EstafetteManifest{
			Stages: []*EstafetteStage{
				{
					Name: "build",
				},
			},
			Triggers: []*EstafetteTrigger{
				{
					Pipeline:    &EstafettePipelineTrigger{},
					BuildAction: &EstafetteTriggerBuildAction{},
				},
			},
			Releases: []*EstafetteRelease{
				{
					Name: "tooling",
					Triggers: []*EstafetteTrigger{
						{
							Pipeline:      &EstafettePipelineTrigger{},
							ReleaseAction: &EstafetteTriggerReleaseAction{},
						},
					},
				},
			},
		}

		// act
		triggers := manifest.GetAllTriggers("github.com", "estafette", "estafette-ci-manifest")

		assert.Equal(t, 2, len(triggers))
	})

	t.Run("ReplacesPipelineNameWithActualPipelineNameIfValueIsSelf", func(t *testing.T) {

		manifest := EstafetteManifest{
			Stages: []*EstafetteStage{
				{
					Name: "build",
				},
			},
			Triggers: []*EstafetteTrigger{
				{
					Pipeline: &EstafettePipelineTrigger{
						Name: "self",
					},
					BuildAction: &EstafetteTriggerBuildAction{},
				},
			},
			Releases: []*EstafetteRelease{
				{
					Name: "tooling",
					Triggers: []*EstafetteTrigger{
						{
							Pipeline: &EstafettePipelineTrigger{
								Name: "self",
							},
							ReleaseAction: &EstafetteTriggerReleaseAction{},
						},
						{
							Release: &EstafetteReleaseTrigger{
								Name:   "self",
								Target: "tooling",
							},
							ReleaseAction: &EstafetteTriggerReleaseAction{},
						},
					},
				},
			},
		}

		// act
		triggers := manifest.GetAllTriggers("github.com", "estafette", "estafette-ci-manifest")

		assert.Equal(t, "github.com/estafette/estafette-ci-manifest", triggers[0].Pipeline.Name)
		assert.Equal(t, "github.com/estafette/estafette-ci-manifest", triggers[1].Pipeline.Name)
		assert.Equal(t, "github.com/estafette/estafette-ci-manifest", triggers[2].Release.Name)
	})

	t.Run("DoesNotReplacePipelineNameWithActualPipelineNameIfValueIsNotSelf", func(t *testing.T) {

		manifest := EstafetteManifest{
			Stages: []*EstafetteStage{
				{
					Name: "build",
				},
			},
			Triggers: []*EstafetteTrigger{
				{
					Pipeline: &EstafettePipelineTrigger{
						Name: "github.com/estafette/estafette-ci-contracts",
					},
					BuildAction: &EstafetteTriggerBuildAction{},
				},
			},
			Releases: []*EstafetteRelease{
				{
					Name: "tooling",
					Triggers: []*EstafetteTrigger{
						{
							Pipeline: &EstafettePipelineTrigger{
								Name: "github.com/estafette/estafette-ci-crypt",
							},
							ReleaseAction: &EstafetteTriggerReleaseAction{},
						},
						{
							Release: &EstafetteReleaseTrigger{
								Name:   "github.com/estaftte/estafette-ci-api",
								Target: "tooling",
							},
							ReleaseAction: &EstafetteTriggerReleaseAction{},
						},
					},
				},
			},
		}

		// act
		triggers := manifest.GetAllTriggers("github.com", "estafette", "estafette-ci-manifest")

		assert.Equal(t, "github.com/estafette/estafette-ci-contracts", triggers[0].Pipeline.Name)
		assert.Equal(t, "github.com/estafette/estafette-ci-crypt", triggers[1].Pipeline.Name)
		assert.Equal(t, "github.com/estaftte/estafette-ci-api", triggers[2].Release.Name)
	})
}

func TestValidate(t *testing.T) {
	t.Run("ReturnsErrorIfTrackIsNotDevBetaOrStable", func(t *testing.T) {

		manifest := EstafetteManifest{
			Builder: EstafetteBuilder{
				Track: "nightly",
			},
			Stages: []*EstafetteStage{
				{},
			},
		}
		manifest.SetDefaults(*GetDefaultManifestPreferences())

		// act
		err := manifest.Validate(*GetDefaultManifestPreferences())

		assert.NotNil(t, err)
	})

	t.Run("ReturnsNoErrorIfTrackIsDev", func(t *testing.T) {

		manifest := EstafetteManifest{
			Builder: EstafetteBuilder{
				Track: "dev",
			},
			Stages: []*EstafetteStage{
				{
					ContainerImage: "docker",
				},
			},
		}
		manifest.SetDefaults(*GetDefaultManifestPreferences())

		// act
		err := manifest.Validate(*GetDefaultManifestPreferences())

		assert.Nil(t, err)
	})

	t.Run("ReturnsNoErrorIfTrackIsBeta", func(t *testing.T) {

		manifest := EstafetteManifest{
			Builder: EstafetteBuilder{
				Track: "beta",
			},
			Stages: []*EstafetteStage{
				{
					ContainerImage: "docker",
				},
			},
		}
		manifest.SetDefaults(*GetDefaultManifestPreferences())

		// act
		err := manifest.Validate(*GetDefaultManifestPreferences())

		assert.Nil(t, err)
	})

	t.Run("ReturnsNoErrorIfTrackIsStable", func(t *testing.T) {

		manifest := EstafetteManifest{
			Builder: EstafetteBuilder{
				Track: "stable",
			},
			Stages: []*EstafetteStage{
				{
					ContainerImage: "docker",
				},
			},
		}
		manifest.SetDefaults(*GetDefaultManifestPreferences())

		// act
		err := manifest.Validate(*GetDefaultManifestPreferences())

		assert.Nil(t, err)
	})
}

func TestDeepCopy(t *testing.T) {

	t.Run("ReturnsCopyOfManifestAndAllStages", func(t *testing.T) {

		manifest, err := ReadManifestFromFile(GetDefaultManifestPreferences(), "test-manifest.yaml", true)
		assert.Nil(t, err)
		assert.Equal(t, 7, len(manifest.Stages))

		// act
		copiedManifest := manifest.DeepCopy()

		assert.Equal(t, 7, len(copiedManifest.Stages))
		assert.NotSame(t, manifest.Stages[0], copiedManifest.Stages[0])
	})
}
