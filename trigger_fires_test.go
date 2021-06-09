package manifest

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEstafettePipelineTriggerFires(t *testing.T) {
	t.Run("ReturnsTrueIfEventStatusNameAndBranchMatch", func(t *testing.T) {

		event := EstafettePipelineEvent{
			RepoSource: "github.com",
			RepoOwner:  "estafette",
			RepoName:   "estafette-ci-api",
			Branch:     "main",
			Status:     "succeeded",
			Event:      "finished",
		}

		trigger := EstafettePipelineTrigger{
			Event:  "finished",
			Status: "succeeded",
			Name:   "github.com/estafette/estafette-ci-api",
			Branch: "main",
		}

		// act
		fires := trigger.Fires(&event)

		assert.True(t, fires)
	})

	t.Run("ReturnsTrueIfNegativeLookupBranchRegexDoesMatch", func(t *testing.T) {

		event := EstafettePipelineEvent{
			RepoSource: "github.com",
			RepoOwner:  "estafette",
			RepoName:   "estafette-ci-api",
			Branch:     "development",
			Status:     "succeeded",
			Event:      "finished",
		}

		trigger := EstafettePipelineTrigger{
			Event:  "finished",
			Status: "succeeded",
			Name:   "github.com/estafette/estafette-ci-api",
			Branch: "!~ main",
		}

		// act
		fires := trigger.Fires(&event)

		assert.True(t, fires)
	})

	t.Run("ReturnsTrueIfNegativeLookupBranchRegexDoesNotMatch", func(t *testing.T) {

		event := EstafettePipelineEvent{
			RepoSource: "github.com",
			RepoOwner:  "estafette",
			RepoName:   "estafette-ci-api",
			Branch:     "main",
			Status:     "succeeded",
			Event:      "finished",
		}

		trigger := EstafettePipelineTrigger{
			Event:  "finished",
			Status: "succeeded",
			Name:   "github.com/estafette/estafette-ci-api",
			Branch: "!~ main",
		}

		// act
		fires := trigger.Fires(&event)

		assert.False(t, fires)
	})

	t.Run("ReturnsFalseIfEventDoesNotMatch", func(t *testing.T) {

		event := EstafettePipelineEvent{
			RepoSource: "github.com",
			RepoOwner:  "estafette",
			RepoName:   "estafette-ci-api",
			Branch:     "main",
			Status:     "succeeded",
			Event:      "finished",
		}

		trigger := EstafettePipelineTrigger{
			Event:  "started",
			Status: "",
			Name:   "github.com/estafette/estafette-ci-api",
			Branch: "!= main",
		}

		// act
		fires := trigger.Fires(&event)

		assert.False(t, fires)
	})

	t.Run("ReturnsFalseIfStatusDoesNotMatch", func(t *testing.T) {

		event := EstafettePipelineEvent{
			RepoSource: "github.com",
			RepoOwner:  "estafette",
			RepoName:   "estafette-ci-api",
			Branch:     "main",
			Status:     "succeeded",
			Event:      "finished",
		}

		trigger := EstafettePipelineTrigger{
			Event:  "finished",
			Status: "failed",
			Name:   "github.com/estafette/estafette-ci-api",
			Branch: "main",
		}

		// act
		fires := trigger.Fires(&event)

		assert.False(t, fires)
	})

	t.Run("ReturnsFalseIfNameDoesNotMatch", func(t *testing.T) {

		event := EstafettePipelineEvent{
			RepoSource: "github.com",
			RepoOwner:  "estafette",
			RepoName:   "estafette-ci-api",
			Branch:     "main",
			Status:     "succeeded",
			Event:      "finished",
		}

		trigger := EstafettePipelineTrigger{
			Event:  "finished",
			Status: "succeeded",
			Name:   "github.com/estafette/estafette-ci-builder",
			Branch: "main",
		}

		// act
		fires := trigger.Fires(&event)

		assert.False(t, fires)
	})

	t.Run("ReturnsFalseIfBranchDoesNotMatch", func(t *testing.T) {

		event := EstafettePipelineEvent{
			RepoSource: "github.com",
			RepoOwner:  "estafette",
			RepoName:   "estafette-ci-api",
			Branch:     "main",
			Status:     "succeeded",
			Event:      "finished",
		}

		trigger := EstafettePipelineTrigger{
			Event:  "finished",
			Status: "succeeded",
			Name:   "github.com/estafette/estafette-ci-api",
			Branch: "development",
		}

		// act
		fires := trigger.Fires(&event)

		assert.False(t, fires)
	})
}

func TestEstafetteReleaseTriggerFires(t *testing.T) {
	t.Run("ReturnsTrueIfEventStatusNameAndBranchMatch", func(t *testing.T) {

		event := EstafetteReleaseEvent{
			RepoSource: "github.com",
			RepoOwner:  "estafette",
			RepoName:   "estafette-ci-api",
			Target:     "development",
			Status:     "succeeded",
			Event:      "finished",
		}

		trigger := EstafetteReleaseTrigger{
			Event:  "finished",
			Status: "succeeded",
			Name:   "github.com/estafette/estafette-ci-api",
			Target: "development",
		}

		// act
		fires := trigger.Fires(&event)

		assert.True(t, fires)
	})

	t.Run("ReturnsFalseIfEventDoesNotMatch", func(t *testing.T) {

		event := EstafetteReleaseEvent{
			RepoSource: "github.com",
			RepoOwner:  "estafette",
			RepoName:   "estafette-ci-api",
			Target:     "development",
			Status:     "succeeded",
			Event:      "finished",
		}

		trigger := EstafetteReleaseTrigger{
			Event:  "started",
			Status: "",
			Name:   "github.com/estafette/estafette-ci-api",
			Target: "development",
		}

		// act
		fires := trigger.Fires(&event)

		assert.False(t, fires)
	})

	t.Run("ReturnsFalseIfStatusDoesNotMatch", func(t *testing.T) {

		event := EstafetteReleaseEvent{
			RepoSource: "github.com",
			RepoOwner:  "estafette",
			RepoName:   "estafette-ci-api",
			Target:     "development",
			Status:     "succeeded",
			Event:      "finished",
		}

		trigger := EstafetteReleaseTrigger{
			Event:  "finished",
			Status: "failed",
			Name:   "github.com/estafette/estafette-ci-api",
			Target: "development",
		}

		// act
		fires := trigger.Fires(&event)

		assert.False(t, fires)
	})

	t.Run("ReturnsFalseIfNameDoesNotMatch", func(t *testing.T) {

		event := EstafetteReleaseEvent{
			RepoSource: "github.com",
			RepoOwner:  "estafette",
			RepoName:   "estafette-ci-api",
			Target:     "development",
			Status:     "succeeded",
			Event:      "finished",
		}

		trigger := EstafetteReleaseTrigger{
			Event:  "finished",
			Status: "succeeded",
			Name:   "github.com/estafette/estafette-ci-builder",
			Target: "development",
		}

		// act
		fires := trigger.Fires(&event)

		assert.False(t, fires)
	})

	t.Run("ReturnsFalseIfBranchDoesNotMatch", func(t *testing.T) {

		event := EstafetteReleaseEvent{
			RepoSource: "github.com",
			RepoOwner:  "estafette",
			RepoName:   "estafette-ci-api",
			Target:     "development",
			Status:     "succeeded",
			Event:      "finished",
		}

		trigger := EstafetteReleaseTrigger{
			Event:  "finished",
			Status: "succeeded",
			Name:   "github.com/estafette/estafette-ci-api",
			Target: "staging",
		}

		// act
		fires := trigger.Fires(&event)

		assert.False(t, fires)
	})
}

func TestEstafetteCronTriggerFires(t *testing.T) {
	t.Run("ReturnsTrueIfEventTimeMatchesCronSchedule", func(t *testing.T) {

		event := EstafetteCronEvent{
			Time: time.Date(2019, 4, 5, 11, 10, 0, 0, time.UTC),
		}

		trigger := EstafetteCronTrigger{
			Schedule: "*/5 * * * *",
		}

		// act
		fires := trigger.Fires(&event)

		assert.True(t, fires)
	})

	t.Run("ReturnsTrueIfEventTimeMatchesCronSchedule", func(t *testing.T) {

		event := EstafetteCronEvent{
			Time: time.Date(2019, 4, 5, 11, 12, 1, 0, time.UTC),
		}

		trigger := EstafetteCronTrigger{
			Schedule: "*/5 * * * *",
		}

		// act
		fires := trigger.Fires(&event)

		assert.False(t, fires)
	})
}

func TestEstafetteGitTriggerFires(t *testing.T) {
	t.Run("ReturnsTrueIfEventStatusNameAndBranchMatch", func(t *testing.T) {

		event := EstafetteGitEvent{
			Event:      "push",
			Repository: "bitbucket.org/xivart/icarus_to_email_service_trigger",
			Branch:     "main"}

		trigger := EstafetteGitTrigger{
			Event:      "push",
			Repository: "bitbucket.org/xivart/icarus_to_email_service_trigger",
			Branch:     "main",
		}

		// act
		fires := trigger.Fires(&event)

		assert.True(t, fires)
	})
}

func TestEstafetteGithubTriggerFires(t *testing.T) {
	t.Run("ReturnsTrueIfEventIsContainedInTriggerEvents", func(t *testing.T) {

		event := EstafetteGithubEvent{
			Event: "create",
		}

		trigger := EstafetteGithubTrigger{
			Events: []string{
				"commit_comment",
				"create",
				"delete",
				"deployment",
				"deployment_status",
				"fork",
				"gollum",
				"installation",
				"installation_repositories",
			},
		}

		// act
		fires := trigger.Fires(&event)

		assert.True(t, fires)
	})
}

func TestEstafetteBitbucketTriggerFires(t *testing.T) {
	t.Run("ReturnsTrueIfEventIsContainedInTriggerEvents", func(t *testing.T) {

		event := EstafetteGithubEvent{
			Event: "pullrequest:comment_created",
		}

		trigger := EstafetteGithubTrigger{
			Events: []string{
				"pullrequest:fulfilled",
				"pullrequest:rejected",
				"pullrequest:comment_created",
				"pullrequest:comment_updated",
				"pullrequest:comment_deleted",
			},
		}

		// act
		fires := trigger.Fires(&event)

		assert.True(t, fires)
	})
}

func TestEstafettePubsubTriggerFires(t *testing.T) {
	t.Run("ReturnsTrueIfTopicAndProjectMatch", func(t *testing.T) {

		event := EstafettePubSubEvent{
			Project: "my-project",
			Topic:   "my-topic",
		}

		trigger := EstafettePubSubTrigger{
			Project: "my-project",
			Topic:   "my-topic",
		}

		// act
		fires := trigger.Fires(&event)

		assert.True(t, fires)
	})

	t.Run("ReturnsFalseIfProjectDoesNotMatch", func(t *testing.T) {

		event := EstafettePubSubEvent{
			Project: "another-project",
			Topic:   "my-topic",
		}

		trigger := EstafettePubSubTrigger{
			Project: "my-project",
			Topic:   "my-topic",
		}

		// act
		fires := trigger.Fires(&event)

		assert.False(t, fires)
	})

	t.Run("ReturnsFalseIfTopicDoesNotMatch", func(t *testing.T) {

		event := EstafettePubSubEvent{
			Project: "my-project",
			Topic:   "another-topic",
		}

		trigger := EstafettePubSubTrigger{
			Project: "my-project",
			Topic:   "my-topic",
		}

		// act
		fires := trigger.Fires(&event)

		assert.False(t, fires)
	})

	t.Run("ReturnsTrueIfTopicAndProjectMatchAsRegex", func(t *testing.T) {

		event := EstafettePubSubEvent{
			Project: "my-project",
			Topic:   "my-topic",
		}

		trigger := EstafettePubSubTrigger{
			Project: ".+-project",
			Topic:   ".+-topic",
		}

		// act
		fires := trigger.Fires(&event)

		assert.True(t, fires)
	})

	t.Run("ReturnsFalseIfProjectDoesNotMatchAsRegex", func(t *testing.T) {

		event := EstafettePubSubEvent{
			Project: "-project",
			Topic:   "my-topic",
		}

		trigger := EstafettePubSubTrigger{
			Project: ".+-project",
			Topic:   ".+-topic",
		}

		// act
		fires := trigger.Fires(&event)

		assert.False(t, fires)
	})

	t.Run("ReturnsFalseIfTopicDoesNotMatchAsRegex", func(t *testing.T) {

		event := EstafettePubSubEvent{
			Project: "my-project",
			Topic:   "-topic",
		}

		trigger := EstafettePubSubTrigger{
			Project: ".+-project",
			Topic:   ".+-topic",
		}

		// act
		fires := trigger.Fires(&event)

		assert.False(t, fires)
	})
}
