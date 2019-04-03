package manifest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEstafettePipelineTriggerFires(t *testing.T) {
	t.Run("ReturnsTrueIfEventStatusNameAndBranchMatch", func(t *testing.T) {

		event := EstafettePipelineEvent{
			RepoSource: "github.com",
			RepoOwner:  "estafette",
			RepoName:   "estafette-ci-api",
			Branch:     "master",
			Status:     "succeeded",
			Event:      "finished",
		}

		trigger := EstafettePipelineTrigger{
			Event:  "finished",
			Status: "succeeded",
			Name:   "github.com/estafette/estafette-ci-api",
			Branch: "master",
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
			Branch: "!~ master",
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
			Branch:     "master",
			Status:     "succeeded",
			Event:      "finished",
		}

		trigger := EstafettePipelineTrigger{
			Event:  "finished",
			Status: "succeeded",
			Name:   "github.com/estafette/estafette-ci-api",
			Branch: "!~ master",
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
			Branch:     "master",
			Status:     "succeeded",
			Event:      "finished",
		}

		trigger := EstafettePipelineTrigger{
			Event:  "started",
			Status: "",
			Name:   "github.com/estafette/estafette-ci-api",
			Branch: "!= master",
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
			Branch:     "master",
			Status:     "succeeded",
			Event:      "finished",
		}

		trigger := EstafettePipelineTrigger{
			Event:  "finished",
			Status: "failed",
			Name:   "github.com/estafette/estafette-ci-api",
			Branch: "master",
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
			Branch:     "master",
			Status:     "succeeded",
			Event:      "finished",
		}

		trigger := EstafettePipelineTrigger{
			Event:  "finished",
			Status: "succeeded",
			Name:   "github.com/estafette/estafette-ci-builder",
			Branch: "master",
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
			Branch:     "master",
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
