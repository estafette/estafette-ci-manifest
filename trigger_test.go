package manifest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetDefaults(t *testing.T) {

	t.Run("DefaultsFilterBranchToMasterForEventPipelineBuildStarted", func(t *testing.T) {

		trigger := EstafetteTrigger{
			Event: "pipeline-build-started",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, "master", trigger.Filter.Branch)
	})

	t.Run("DefaultsFilterBranchToMasterForEventPipelineBuildFinished", func(t *testing.T) {

		trigger := EstafetteTrigger{
			Event: "pipeline-build-finished",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, "master", trigger.Filter.Branch)
	})

	t.Run("DefaultsFilterStatusToSucceededForEventPipelineBuildFinished", func(t *testing.T) {

		trigger := EstafetteTrigger{
			Event: "pipeline-build-finished",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, "succeeded", trigger.Filter.Status)
	})

	t.Run("DefaultsFilterBranchToAnyForEventPipelineReleaseStarted", func(t *testing.T) {

		trigger := EstafetteTrigger{
			Event: "pipeline-release-started",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, ".+", trigger.Filter.Branch)
	})

	t.Run("DefaultsFilterBranchToAnyForEventPipelineReleaseFinished", func(t *testing.T) {

		trigger := EstafetteTrigger{
			Event: "pipeline-release-finished",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, ".+", trigger.Filter.Branch)
	})

	t.Run("DefaultsFilterStatusToSucceededForEventPipelineReleaseFinished", func(t *testing.T) {

		trigger := EstafetteTrigger{
			Event: "pipeline-release-finished",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, "succeeded", trigger.Filter.Status)
	})

	t.Run("DefaultsFilterPipelineToThisPipelineToAnyForEventPipelineReleaseStarted", func(t *testing.T) {

		trigger := EstafetteTrigger{
			Event: "pipeline-release-started",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, "this", trigger.Filter.Pipeline)
	})

	t.Run("DefaultsFilterPipelineToThisPipelineToAnyForEventPipelineReleaseFinished", func(t *testing.T) {

		trigger := EstafetteTrigger{
			Event: "pipeline-release-finished",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, "this", trigger.Filter.Pipeline)
	})

	t.Run("DefaultsThenBranchToMasterForEventPipelineBuildStarted", func(t *testing.T) {

		trigger := EstafetteTrigger{
			Event: "pipeline-build-started",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, "master", trigger.Run.Branch)
	})

	t.Run("DefaultsThenBranchToMasterForEventPipelineBuildFinished", func(t *testing.T) {

		trigger := EstafetteTrigger{
			Event: "pipeline-build-finished",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, "master", trigger.Run.Branch)
	})

	t.Run("DefaultsThenBranchToAnyForEventPipelineReleaseStarted", func(t *testing.T) {

		trigger := EstafetteTrigger{
			Event: "pipeline-release-started",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, ".+", trigger.Run.Branch)
	})

	t.Run("DefaultsThenBranchToAnyForEventPipelineReleaseFinished", func(t *testing.T) {

		trigger := EstafetteTrigger{
			Event: "pipeline-release-finished",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, ".+", trigger.Run.Branch)
	})
}
